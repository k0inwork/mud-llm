package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mud/internal/dal"
	"mud/internal/game"
	"mud/internal/game/events"
	"mud/internal/models"
	"mud/internal/presentation"
)

// ConnectionState defines the state of a client connection.
type ConnectionState int

const (
	StateWelcome ConnectionState = iota
	StateLoginUsername
	StateLoginPassword
	StateCreateAccountUsername
	StateCreateAccountPassword
	StateCreateAccountEmail
	StateCharacterSelection
	StateCharacterCreationName
	StateCharacterCreationRace
	StateCharacterCreationProfession
	StateInGame
)

// client represents a connected client.
type client struct {
	conn net.Conn
	writer        *bufio.Writer
	state         ConnectionState
	account       *models.PlayerAccount
	character     *models.PlayerCharacter
	tempUsername  string
	tempPassword  string
}

// TelnetServer represents the Telnet server for the MUD.
type TelnetServer struct {
	listener           net.Listener
	renderer           game.TelnetRendererInterface
	eventBus           *events.EventBus
	dal                *dal.DAL
	llmService         game.LLMServiceInterface
	playerConnections  map[string]*client // Map characterID to client
	connectionsMutex   sync.RWMutex
	Ready              chan bool
}

// NewTelnetServer creates a new TelnetServer.
func NewTelnetServer(listener net.Listener, renderer game.TelnetRendererInterface, eventBus *events.EventBus, dal *dal.DAL, llmService game.LLMServiceInterface) *TelnetServer {
	s := &TelnetServer{
		listener:          listener,
		renderer:          renderer,
		eventBus:          eventBus,
		dal:               dal,
		llmService:        llmService,
		playerConnections: make(map[string]*client),
		Ready:             make(chan bool),
	}

	// Subscribe to PlayerMessageEvent
	playerMessageChannel := make(chan interface{}, 100) // Buffered channel
	s.eventBus.Subscribe(events.PlayerMessageEventType, playerMessageChannel)
	go func() {
		for event := range playerMessageChannel {
			if pm, ok := event.(*events.PlayerMessageEvent); ok {
				s.connectionsMutex.RLock()
				client, found := s.playerConnections[pm.PlayerID]
				s.connectionsMutex.RUnlock()

				if found {
					msg := presentation.SemanticMessage{
						Type:    presentation.NarrativeMessage,
						Content: pm.Content,
						Color:   presentation.ColorDefault,
					}
					s.sendMessage(client, msg)
				} else {
					logrus.Infof("TelnetServer: Player connection not found for PlayerID %s, cannot send message: %s", pm.PlayerID, pm.Content)
				}
			} else {
				logrus.Infof("TelnetServer: Received unexpected event type on PlayerMessageEventType: %T", event)
			}
		}
	}()

	return s
}

// Start begins listening for incoming Telnet connections.
func (s *TelnetServer) Start() {
	defer s.listener.Close()
	logrus.Infof("Telnet server listening on port %d\n", s.listener.Addr().(*net.TCPAddr).Port)
	s.Ready <- true

	for {

		conn, err := s.listener.Accept()
		if err != nil {
			// Check if the error is because the listener was closed.
			if errors.Is(err, net.ErrClosed) {
				logrus.Info("Telnet server listener closed, shutting down.")
				break // Exit the loop gracefully.
			}
			logrus.Infof("Failed to accept connection: %v", err)
			continue
		}
		client := &client{conn: conn, state: StateWelcome, writer: bufio.NewWriter(conn)}
		go s.handleConnection(client)
	}
}

// handleConnection manages a single Telnet client connection.
func (s *TelnetServer) handleConnection(c *client) {
	defer func() {
		if c.character != nil {
			s.connectionsMutex.Lock()
			delete(s.playerConnections, c.character.ID)
			s.connectionsMutex.Unlock()
		}
		logrus.Infof("Client %s disconnected: %v\n", c.conn.RemoteAddr(), c.conn.Close())
	}()

	logrus.Infof("New Telnet connection from %s\n", c.conn.RemoteAddr())
	s.showWelcomeMenu(c)

	reader := bufio.NewReader(c.conn)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			logrus.Infof("Client %s disconnected: %v\n", c.conn.RemoteAddr(), err)
			return
		}
		trimmedInput := strings.TrimSpace(input)
		s.handleInput(c, trimmedInput)
	}
}

func (s *TelnetServer) handleInput(c *client, input string) {
	switch c.state {
	case StateWelcome:
		s.handleWelcomeInput(c, input)
	case StateLoginUsername:
		s.handleLoginUsername(c, input)
	case StateLoginPassword:
		s.handleLoginPassword(c, input)
	case StateCreateAccountUsername:
		s.handleCreateUsername(c, input)
	case StateCreateAccountPassword:
		s.handleCreatePassword(c, input)
	case StateCreateAccountEmail:
		s.handleCreateEmail(c, input)
	case StateCharacterSelection:
		s.handleCharacterSelection(c, input)
	case StateCharacterCreationName:
		s.handleCharacterCreationName(c, input)
	case StateCharacterCreationRace:
		// s.handleCharacterCreationRace(c, input)
	case StateCharacterCreationProfession:
		// s.handleCharacterCreationProfession(c, input)
	case StateInGame:
		s.handleInGameInput(c, input)
	}
}

func (s *TelnetServer) showWelcomeMenu(c *client) {
	msg := `
Welcome to GoMUD!
1. Login
2. Create Account
`
	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: msg, Color: presentation.ColorSuccess})
}

func (s *TelnetServer) handleWelcomeInput(c *client, input string) {
	switch input {
	case "1":
		c.state = StateLoginUsername
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Enter username: ", Color: presentation.ColorDefault})
	case "2":
		c.state = StateCreateAccountUsername
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Enter desired username: ", Color: presentation.ColorDefault})
	default:
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Invalid option.", Color: presentation.ColorError})
		s.showWelcomeMenu(c)
	}
}

func (s *TelnetServer) handleLoginUsername(c *client, username string) {
	c.tempUsername = username
	c.state = StateLoginPassword
	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Enter password: ", Color: presentation.ColorDefault})
}

func (s *TelnetServer) handleLoginPassword(c *client, password string) {
	account, err := s.dal.PlayerAccountDAL.Authenticate(c.tempUsername, password)
	if err != nil || account == nil {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Invalid username or password.", Color: presentation.ColorError})
		c.state = StateWelcome
		s.showWelcomeMenu(c)
		return
	}
	c.account = account
	s.dal.PlayerAccountDAL.UpdateLastLogin(account.ID)
	s.showCharacterSelection(c)
}

func (s *TelnetServer) handleCreateUsername(c *client, username string) {
	// Basic validation
	if len(username) < 3 {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Username must be at least 3 characters.", Color: presentation.ColorError})
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Enter desired username: ", Color: presentation.ColorDefault})
		return
	}
	// Check if username exists
	existing, err := s.dal.PlayerAccountDAL.GetAccountByUsername(username)
	if err != nil {
		logrus.Infof("Error checking username: %v", err)
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "An error occurred. Please try again.", Color: presentation.ColorError})
		c.state = StateWelcome
		s.showWelcomeMenu(c)
		return
	}
	if existing != nil {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Username already taken.", Color: presentation.ColorError})
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Enter desired username: ", Color: presentation.ColorDefault})
		return
	}
	c.tempUsername = username
	c.state = StateCreateAccountPassword
	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Enter password: ", Color: presentation.ColorDefault})
}

func (s *TelnetServer) handleCreatePassword(c *client, password string) {
	if len(password) < 6 {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Password must be at least 6 characters.", Color: presentation.ColorError})
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Enter password: ", Color: presentation.ColorDefault})
		return
	}
	c.tempPassword = password
	c.state = StateCreateAccountEmail
	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Enter email (optional): ", Color: presentation.ColorDefault})
}

func (s *TelnetServer) handleCreateEmail(c *client, email string) {
	account, err := s.dal.PlayerAccountDAL.CreateAccount(c.tempUsername, c.tempPassword, email)
	if err != nil {
		logrus.Infof("Error creating account: %v", err)
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Failed to create account. Please try again.", Color: presentation.ColorError})
		c.state = StateWelcome
		s.showWelcomeMenu(c)
		return
	}
	c.account = account
	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Account created successfully!", Color: presentation.ColorSuccess})
	s.showCharacterSelection(c)
}

func (s *TelnetServer) showCharacterSelection(c *client) {
	characters, err := s.dal.PlayerCharacterDAL.GetCharactersByAccountID(c.account.ID)
	if err != nil {
		logrus.Infof("Error getting characters: %v", err)
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Error loading characters.", Color: presentation.ColorError})
		return
	}

	c.state = StateCharacterSelection
	var msg strings.Builder
	msg.WriteString("\n--- Character Selection ---\n")
	if len(characters) == 0 {
		msg.WriteString("You have no characters.\n")
	} else {
		for i, char := range characters {
			msg.WriteString(fmt.Sprintf("%d. %s (%s %s)\n", i+1, char.Name, char.RaceID, char.ProfessionID))
		}
	}
	msg.WriteString("\nType a number to select a character, or 'new' to create one.\n")
	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: msg.String(), Color: presentation.ColorDefault})
}

func (s *TelnetServer) handleCharacterSelection(c *client, input string) {
	if strings.ToLower(input) == "new" {
		c.state = StateCharacterCreationName
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Enter character name: ", Color: presentation.ColorDefault})
		return
	}

	characters, err := s.dal.PlayerCharacterDAL.GetCharactersByAccountID(c.account.ID)
	if err != nil {
		logrus.Infof("Error getting characters for selection: %v", err)
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Error loading characters for selection.", Color: presentation.ColorError})
		s.showCharacterSelection(c)
		return
	}

	selection, err := strconv.Atoi(input)
	if err != nil || selection < 1 || selection > len(characters) {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Invalid selection. Please enter a number or 'new'.", Color: presentation.ColorError})
		s.showCharacterSelection(c)
		return
	}

	selectedChar := characters[selection-1]
	c.character = selectedChar
	s.enterGame(c)
}

func (s *TelnetServer) handleCharacterCreationName(c *client, name string) {
	// Basic validation
	if len(name) < 3 {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Name must be at least 3 characters.", Color: presentation.ColorError})
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Enter character name: ", Color: presentation.ColorDefault})
		return
	}

	// In a real MUD, you'd check for name uniqueness here.

	character := &models.PlayerCharacter{
		ID:              uuid.New().String(),
		PlayerAccountID: c.account.ID,
		Name:            name,
		RaceID:          "human",     // Placeholder
		ProfessionID:    "warrior",   // Placeholder
		CurrentRoomID:   "bag_end",
		Health:          100,
		MaxHealth:       100,
		Inventory:       "[]",
		VisitedRoomIDs:  "[\"bag_end\"]",
		CreatedAt:       time.Now(),
	}

	err := s.dal.PlayerCharacterDAL.CreateCharacter(character)
	if err != nil {
		logrus.Infof("Error creating character: %v", err)
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Failed to create character.", Color: presentation.ColorError})
		s.showCharacterSelection(c)
		return
	}

	c.character = character
	s.enterGame(c)
}

func (s *TelnetServer) enterGame(c *client) {
	c.state = StateInGame
	s.connectionsMutex.Lock()
	s.playerConnections[c.character.ID] = c
	s.connectionsMutex.Unlock()

	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: fmt.Sprintf("Welcome, %s!", c.character.Name), Color: presentation.ColorSuccess})
	s.renderRoomDescription(c)
}

func (s *TelnetServer) handleInGameInput(c *client, input string) {
	processedInput := fmt.Sprintf("You typed: %s", input)
	echoMsg := presentation.SemanticMessage{
		Type:    presentation.SystemMessage,
		Content: processedInput,
		Color:   presentation.ColorDefault,
	}
	s.sendMessage(c, echoMsg)

	parts := strings.Fields(input)
	if len(parts) == 0 {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "What would you like to do?", Color: presentation.ColorDefault})
		return
	}

	actionType := parts[0]
	var target string
	if len(parts) > 1 {
		target = strings.Join(parts[1:], " ")
	}

	// Handle movement commands
	switch strings.ToLower(actionType) {
	case "n", "s", "e", "w", "u", "d":
		s.handleMovement(c, actionType)
	case "move":
		if len(parts) > 1 {
			s.handleMovement(c, parts[1])
		} else {
			s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Move where? (n, s, e, w, u, d)", Color: presentation.ColorWarning})
		}
	case "look":
		s.renderRoomDescription(c)
	case "talk":
		if target == "" {
			s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Talk to whom?", Color: presentation.ColorWarning})
			return
		}
		s.handleTalkCommand(c, target)
	case "gather":
		if target == "" {
			s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Gather what?", Color: presentation.ColorWarning})
			return
		}
		s.handleGatherCommand(c, target)
	case "give":
		s.handleGiveCommand(c, input)
	}

	room, err := s.dal.RoomDAL.GetRoomByID(c.character.CurrentRoomID)
	if err != nil || room == nil {
		logrus.Infof("TelnetServer: Room %s not found for character %s or error: %v", c.character.CurrentRoomID, c.character.ID, err)
		// Fallback to a default room
		room, _ = s.dal.RoomDAL.GetRoomByID("bag_end")
		if room == nil {
			logrus.Fatalf("TelnetServer: Default room 'bag_end' not found. Database not seeded?")
		}
	}

	actionEvent := &events.ActionEvent{
		Player:     c.character,
		ActionType: actionType,
		Room:       room,
		Timestamp:  time.Now(),
	}
	s.eventBus.Publish(events.ActionEventType, actionEvent)
}

func (s *TelnetServer) handleMovement(c *client, direction string) {
	room, err := s.dal.RoomDAL.GetRoomByID(c.character.CurrentRoomID)
	if err != nil || room == nil {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "You are in a strange place and cannot move.", Color: presentation.ColorError})
		return
	}

	var exits map[string]models.Exit
	if err := json.Unmarshal([]byte(room.Exits), &exits); err != nil {
		logrus.Infof("Error unmarshaling exits for room %s: %v", room.ID, err)
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "The exits seem to be broken.", Color: presentation.ColorError})
		return
	}

	exit, found := exits[strings.ToLower(direction)]
	if !found {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "You cannot go that way.", Color: presentation.ColorWarning})
		return
	}

	if exit.IsLocked {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "That way is locked.", Color: presentation.ColorWarning})
		return
	}

	// Update player's current room
	c.character.CurrentRoomID = exit.TargetRoomID
	err = s.dal.PlayerCharacterDAL.UpdateCharacter(c.character)
	if err != nil {
		logrus.Infof("Error updating character room: %v", err)
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "An error occurred while moving.", Color: presentation.ColorError})
		return
	}

	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: fmt.Sprintf("You move %s.", direction), Color: presentation.ColorDefault})
	s.renderRoomDescription(c)

	// Publish a "move" action event
	actionEvent := &events.ActionEvent{
		Player:     c.character,
		ActionType: "move",
		Room:       room, // Old room
		Timestamp:  time.Now(),
		Targets:    []interface{}{exit.TargetRoomID}, // Target is the new room ID
	}
	s.eventBus.Publish(events.ActionEventType, actionEvent)
}

func (s *TelnetServer) handleTalkCommand(c *client, target string) {
	// For now, just acknowledge the talk command and publish an event.
	// Actual NPC interaction logic will be handled by the SentientEntityManager.
	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: fmt.Sprintf("You try to talk to %s.", target), Color: presentation.ColorDefault})

	room, err := s.dal.RoomDAL.GetRoomByID(c.character.CurrentRoomID)
	if err != nil || room == nil {
		logrus.Infof("TelnetServer: Room %s not found for character %s or error: %v", c.character.CurrentRoomID, c.character.ID, err)
		return
	}

	actionEvent := &events.ActionEvent{
		Player:     c.character,
		ActionType: "talk",
		Room:       room,
		Timestamp:  time.Now(),
		Targets:    []interface{}{target}, // Target is the NPC name/ID
	}
	s.eventBus.Publish(events.ActionEventType, actionEvent)
}

func (s *TelnetServer) handleGatherCommand(c *client, target string) {
	// For now, just acknowledge the gather command and publish an event.
	// Actual item gathering logic (e.g., checking if item exists in room, adding to inventory)
	// will be handled by a game logic service.
	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: fmt.Sprintf("You attempt to gather %s.", target), Color: presentation.ColorDefault})

	room, err := s.dal.RoomDAL.GetRoomByID(c.character.CurrentRoomID)
	if err != nil || room == nil {
		logrus.Infof("TelnetServer: Room %s not found for character %s or error: %v", c.character.CurrentRoomID, c.character.ID, err)
		return
	}

	actionEvent := &events.ActionEvent{
		Player:     c.character,
		ActionType: "gather_item", // Specific action type for gathering
		Room:       room,
		Timestamp:  time.Now(),
		Targets:    []interface{}{target}, // Target is the item name/ID
	}
	s.eventBus.Publish(events.ActionEventType, actionEvent)
}

func (s *TelnetServer) handleGiveCommand(c *client, input string) {
	// Expected format: "give <item_name> to <npc_name>"
	parts := strings.SplitN(input, " to ", 2)
	if len(parts) != 2 {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Invalid 'give' command format. Use: give <item> to <target>.", Color: presentation.ColorWarning})
		return
	}

	itemPart := strings.TrimSpace(strings.TrimPrefix(parts[0], "give"))
	targetPart := strings.TrimSpace(parts[1])

	if itemPart == "" || targetPart == "" {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "Invalid 'give' command. Specify both item and target.", Color: presentation.ColorWarning})
		return
	}

	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: fmt.Sprintf("You try to give %s to %s.", itemPart, targetPart), Color: presentation.ColorDefault})

	room, err := s.dal.RoomDAL.GetRoomByID(c.character.CurrentRoomID)
	if err != nil || room == nil {
		logrus.Infof("TelnetServer: Room %s not found for character %s or error: %v", c.character.CurrentRoomID, c.character.ID, err)
		return
	}

	actionEvent := &events.ActionEvent{
		Player:     c.character,
		ActionType: "deliver_item", // Specific action type for giving/delivering
		Room:       room,
		Timestamp:  time.Now(),
		Targets:    []interface{}{itemPart, targetPart}, // Target is item and NPC
	}
	s.eventBus.Publish(events.ActionEventType, actionEvent)
}

func (s *TelnetServer) renderRoomDescription(c *client) {
	room, err := s.dal.RoomDAL.GetRoomByID(c.character.CurrentRoomID)
	if err != nil || room == nil {
		s.sendMessage(c, presentation.SemanticMessage{Type: presentation.SystemMessage, Content: "You are in a void. There is nothing to see here.", Color: presentation.ColorError})
		return
	}

	var exits map[string]models.Exit
	json.Unmarshal([]byte(room.Exits), &exits) // Error handling already done in handleMovement

	var exitDescriptions []string
	for dir, exit := range exits {
		status := ""
		if exit.IsLocked {
			status = " (locked)"
		}
		exitDescriptions = append(exitDescriptions, fmt.Sprintf("%s (%s%s)", dir, exit.TargetRoomID, status))
	}

	roomDesc := fmt.Sprintf("\n--- %s ---\n%s\nExits: %s\n",
		room.Name,
		room.Description,
		strings.Join(exitDescriptions, ", "),
	)

	// List NPCs in the room
	npcs, err := s.dal.NpcDAL.GetNPCsByRoom(room.ID)
	if err == nil && len(npcs) > 0 {
		roomDesc += "NPCs present: "
		var npcNames []string
		for _, npc := range npcs {
			npcNames = append(npcNames, npc.Name)
		}
		roomDesc += strings.Join(npcNames, ", ") + "\n"
	}

	s.sendMessage(c, presentation.SemanticMessage{Type: presentation.RoomUpdate, Content: roomDesc, Color: presentation.ColorDefault})
}

// sendMessage sends a SemanticMessage to a specific connection.
func (s *TelnetServer) sendMessage(c *client, msg presentation.SemanticMessage) {
	rendered := s.renderer.RenderMessage(msg)
	logrus.Infof("Attempting to send to %s: %s", c.conn.RemoteAddr(), rendered)
	_, err := c.writer.WriteString(rendered + "\n")
	if err != nil {
		logrus.Infof("Failed to write message to buffer for %s: %v", c.conn.RemoteAddr(), err)
		return
	}
	err = c.writer.Flush()
	if err != nil {
		logrus.Infof("Failed to flush buffer for %s: %v", c.conn.RemoteAddr(), err)
	} else {
		logrus.Infof("Successfully sent message to %s", c.conn.RemoteAddr())
	}
}
