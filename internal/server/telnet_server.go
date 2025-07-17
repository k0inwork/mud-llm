package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"mud/internal/dal"
	"mud/internal/game/events"
	"mud/internal/models"
	"mud/internal/presentation"
)

// TelnetServer represents the Telnet server for the MUD.
type TelnetServer struct {
	port     string
	renderer *presentation.TelnetRenderer
	eventBus *events.EventBus
	playerDAL *dal.PlayerDAL
	roomDAL  *dal.RoomDAL
}

// NewTelnetServer creates a new TelnetServer.
func NewTelnetServer(port string, eventBus *events.EventBus, dal *dal.DAL) *TelnetServer {
	return &TelnetServer{
		port:     port,
		renderer: presentation.NewTelnetRenderer(),
		eventBus: eventBus,
		playerDAL: dal.PlayerDAL,
		roomDAL:  dal.RoomDAL,
	}
}

// Start begins listening for incoming Telnet connections.
func (s *TelnetServer) Start() {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Fatalf("Failed to start Telnet server on port %s: %v", s.port, err)
	}
	defer listener.Close()
	fmt.Printf("Telnet server listening on port %s\n", s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection manages a single Telnet client connection.
func (s *TelnetServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New Telnet connection from %s\n", conn.RemoteAddr())

	// Send a welcome message
	welcomeMsg := presentation.SemanticMessage{
		Type:    presentation.SystemMessage,
		Content: "Welcome to GoMUD! Type 'help' for commands.",
		Color:   presentation.ColorSuccess,
	}
	s.sendMessage(conn, welcomeMsg)

	reader := bufio.NewReader(conn)
	for {
		// Read input from the client
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected: %v\n", conn.RemoteAddr(), err)
			return
		}

		// Process input
		trimmedInput := strings.TrimSpace(input)
		parts := strings.Fields(trimmedInput)

		playerID := conn.RemoteAddr().String() // Placeholder for actual player ID

		// Fetch player and room information
		player, err := s.playerDAL.GetPlayerByID(playerID)
		if err != nil || player == nil {
			log.Printf("TelnetServer: Player %s not found or error: %v", playerID, err)
			// For now, create a dummy player if not found (e.g., for initial connection)
			player = &models.Player{ID: playerID, Name: "Guest", CurrentRoomID: "bag_end"} // Default to bag_end
			// In a real game, this would involve character creation/login.
		}

		room, err := s.roomDAL.GetRoomByID(player.CurrentRoomID)
		if err != nil || room == nil {
			log.Printf("TelnetServer: Room %s not found for player %s or error: %v", player.CurrentRoomID, playerID, err)
			// Fallback to a default room if player's room is not found
			room, _ = s.roomDAL.GetRoomByID("bag_end")
			if room == nil {
				log.Fatalf("TelnetServer: Default room 'bag_end' not found. Database not seeded?")
			}
		}

		actionType := ""
		// targetID is currently unused in ActionEvent, but kept for potential future use
		// targetID := ""

		if len(parts) > 0 {
			actionType = parts[0]
		}
		// if len(parts) > 1 {
		// 	targetID = parts[1]
		// }

		if actionType != "" {
			// Create and publish ActionEvent
			actionEvent := &events.ActionEvent{
				Player:     player,
				ActionType: actionType,
				Room:       room,
				Timestamp:  time.Now(),
				// SkillUsed and Targets are not parsed from raw input yet.
				SkillUsed:  nil,
				Targets:    nil,
			}
			s.eventBus.Publish(events.ActionEventType, actionEvent)
		}

		processedInput := fmt.Sprintf("You typed: %s", trimmedInput)
		echoMsg := presentation.SemanticMessage{
			Type:    presentation.SystemMessage,
			Content: processedInput,
			Color:   presentation.ColorDefault,
		}
		s.sendMessage(conn, echoMsg)
	}
}

// sendMessage sends a SemanticMessage to a specific connection.
func (s *TelnetServer) sendMessage(conn net.Conn, msg presentation.SemanticMessage) {
	rendered := s.renderer.RenderMessage(msg)
	_, err := conn.Write([]byte(rendered))
	if err != nil {
		log.Printf("Failed to send message to %s: %v\n", conn.RemoteAddr(), err)
	}
}