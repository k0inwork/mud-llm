package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"os" // Added for file operations
)

const (
	telnetPort = ":4000"
	webPort    = ":7070"
)

// Room represents a location in the MUD
type Room struct {
	ID          int
	Name        string
	Description string
	Exits       map[string]*Exit // Direction (e.g., "north") -> Exit
}

// Exit represents a connection to another room
type Exit struct {
	Direction    string
	TargetRoomID int
}

// Player represents a connected player
type Player struct {
	ID          int
	Name        string
	Conn        net.Conn
	CurrentRoomID int
}

// ServerStats holds statistics about the MUD server
type ServerStats struct {
	sync.RWMutex
	ConnectedPlayers int
	TotalCommands    int
	Uptime           time.Time
}

var (
	players     = make(map[net.Conn]*Player)
	playerIDGen = 0
	stats       = ServerStats{
		Uptime: time.Now(),
	}
	rooms       = make(map[int]*Room)
	roomsMutex  sync.RWMutex
	playersMutex sync.RWMutex // Added playersMutex
	editorTemplate *template.Template
)

func init() {
	// Initialize rooms
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room1 := &Room{
		ID:          1,
		Name:        "The Starting Room",
		Description: "You are in a dimly lit room. There's a faint smell of dust and old parchment.",
		Exits:       make(map[string]*Exit),
	}
	room2 := &Room{
		ID:          2,
		Name:        "A Dusty Corridor",
		Description: "A narrow, dusty corridor stretches before you. The air is still and heavy.",
		Exits:       make(map[string]*Exit),
	}

	room1.Exits["east"] = &Exit{Direction: "east", TargetRoomID: 2}
	room2.Exits["west"] = &Exit{Direction: "west", TargetRoomID: 1}

	rooms[room1.ID] = room1
	rooms[room2.ID] = room2

	var err error
	editorTemplate, err = template.ParseFiles("templates/editor.html")
	if err != nil {
		log.Fatalf("Error parsing editor template: %v", err)
	}
}

// sendRoomDescription sends the current room's description and exits to the player
func sendRoomDescription(p *Player) {
	roomsMutex.RLock()
	currentRoom, ok := rooms[p.CurrentRoomID]
	roomsMutex.RUnlock()

	if !ok {
		p.Conn.Write([]byte("You are in a void. Something is terribly wrong.\r\n"))
		return
	}

	p.Conn.Write([]byte(fmt.Sprintf("\r\n--- %s ---\r\n", currentRoom.Name)))
	p.Conn.Write([]byte(fmt.Sprintf("%s\r\n", currentRoom.Description)))

	if len(currentRoom.Exits) > 0 {
		p.Conn.Write([]byte("\r\nExits: "))
		exitStrings := []string{}
		for dir := range currentRoom.Exits {
			exitStrings = append(exitStrings, dir)
		}
		p.Conn.Write([]byte(strings.Join(exitStrings, ", ") + "\r\n"))
	} else {
		p.Conn.Write([]byte("There are no obvious exits.\r\n"))
	}
	p.Conn.Write([]byte("\r\n> ")) // Prompt
}

// movePlayer attempts to move the player in the given direction
func movePlayer(p *Player, direction string) {
	roomsMutex.RLock()
	currentRoom, ok := rooms[p.CurrentRoomID]
	roomsMutex.RUnlock()

	if !ok {
		p.Conn.Write([]byte("You are lost in the void.\r\n> "))
		return
	}

	exit, exists := currentRoom.Exits[direction]
	if !exists {
		p.Conn.Write([]byte("You can't go that way.\r\n> "))
		return
	}

	roomsMutex.RLock()
	_, targetRoomExists := rooms[exit.TargetRoomID]
	roomsMutex.RUnlock()

	if !targetRoomExists {
		p.Conn.Write([]byte("That exit leads nowhere. Something is broken.\r\n> "))
		return
		}

	p.CurrentRoomID = exit.TargetRoomID
	sendRoomDescription(p)
}

func handleTelnetConnection(conn net.Conn) {
	defer conn.Close()

	// Telnet negotiation: Suppress Go-Ahead and enable local echo (optional)
	conn.Write([]byte{255, 251, 3}) // IAC WILL SUPPRESS_GO_AHEAD
	conn.Write([]byte{255, 252, 1}) // IAC WONT ECHO (client should echo)

	// Simple welcome message
	conn.Write([]byte("Welcome to the GoMUD!\r\n"))
	conn.Write([]byte("What is your name?\r\n> "))

	reader := bufio.NewReader(conn)
	name, err := readLine(reader)
	if err != nil {
		log.Printf("Error reading name: %v", err)
		return
	}


	playerIDGen++
	p := &Player{
		ID:   playerIDGen,
		Name: name,
		Conn: conn,
		CurrentRoomID: 1, // Start in room 1
	}

	playersMutex.Lock()
	players[conn] = p
	playersMutex.Unlock()

	stats.Lock()
	stats.ConnectedPlayers++
	stats.Unlock()

	log.Printf("Player %s (%d) connected from %s", p.Name, p.ID, conn.RemoteAddr())

	conn.Write([]byte(fmt.Sprintf("Welcome, %s! Type 'quit' to exit.\r\n", p.Name)))
	sendRoomDescription(p) // Send initial room description

	for {
		input, err := readLine(reader)
		if err != nil {
			log.Printf("Player %s (%d) disconnected: %v", p.Name, p.ID, err)
			break
		}
		log.Printf("Received raw from %s: %q", p.Name, input)

		stats.Lock()
		stats.TotalCommands++
		stats.Unlock()

		lowerInput := strings.ToLower(input)
		log.Printf("Processed input from %s: %q", p.Name, lowerInput)

		switch lowerInput {
		case "quit":
			conn.Write([]byte("Goodbye!\r\n"))
			break
		case "look":
			sendRoomDescription(p)
		case "north", "n":
			movePlayer(p, "north")
		case "south", "s":
			movePlayer(p, "south")
		case "east", "e":
			movePlayer(p, "east")
		case "west", "w":
			movePlayer(p, "west")
		case "up", "u":
			movePlayer(p, "up")
		case "down", "d":
			movePlayer(p, "down")
		default:
			if strings.HasPrefix(lowerInput, "go ") {
				direction := strings.TrimPrefix(lowerInput, "go ")
				movePlayer(p, direction)
			} else {
				conn.Write([]byte(fmt.Sprintf("Unknown command: %s\r\n", input)))
				conn.Write([]byte("> ")) // Prompt
			}
		}
	}

	playersMutex.Lock()
	delete(players, conn)
	playersMutex.Unlock()

	stats.Lock()
	stats.ConnectedPlayers--
	stats.Unlock()
}

func readLine(reader *bufio.Reader) (string, error) {
	var line []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return "", err
		}

		// Telnet option negotiation stripping (IAC is 255)
		if b == 255 { // IAC
			// Consume the next two bytes (command and option)
			for i := 0; i < 2; i++ {
				_, err := reader.ReadByte()
				if err != nil {
					break // Incomplete sequence
				}
			}
			continue // Skip appending IAC and its parameters
		}

		// Handle newline characters
		if b == '\r' || b == '\n' {
			// Consume any remaining newline characters (e.g., if client sends \r\n)
			for reader.Buffered() > 0 {
				nextByte, err := reader.Peek(1)
				if err != nil {
					break // Cannot peek, break
				}
				if nextByte[0] == '\r' || nextByte[0] == '\n' {
					reader.ReadByte() // Consume it
				} else {
					break
				}
			}
			return strings.TrimSpace(string(line)), nil
		}

		// Only append printable ASCII characters (32-126) and tab (9)
		if (b >= 32 && b <= 126) || b == '\t' {
			line = append(line, b)
		}
		// All other characters (including NULL, other control characters) are discarded
	}
}





func startTelnetServer() {
	listener, err := net.Listen("tcp", telnetPort)
	if err != nil {
		log.Fatalf("Error listening on Telnet port: %v", err)
	}
	defer listener.Close()
	log.Printf("Telnet server listening on %s", telnetPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting Telnet connection: %v", err)
			continue
		}
		go handleTelnetConnection(conn)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>GoMUD Admin</title>
			<style>
				body { font-family: sans-serif; margin: 20px; }
				nav a { margin-right: 15px; }
			</style>
		</head>
		<body>
			<h1>GoMUD Admin Panel</h1>
			<nav>
				<a href="/">Home</a>
				<a href="/stats">Stats</a>
				<a href="/editor">Editor</a>
			</nav>
			<p>Welcome to the GoMUD administration panel.</p>
		</body>
		</html>
	`)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	stats.RLock()
	defer stats.RUnlock()

	uptime := time.Since(stats.Uptime).Round(time.Second)

	fmt.Fprintf(w, "<h1>Server Statistics</h1>")
	fmt.Fprintf(w, "<p>Connected Players: %d</p>", stats.ConnectedPlayers)
	fmt.Fprintf(w, "<p>Total Commands Processed: %d</p>", stats.TotalCommands)
	fmt.Fprintf(w, "<p>Uptime: %s</p>", uptime)
}

func editorHandler(w http.ResponseWriter, r *http.Request) {
	err := editorTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering editor: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveEditorContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	content := r.FormValue("content")
	err := writeEditorContentToFile(content)
	if err != nil {
		http.Error(w, "Error saving content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Content saved successfully!")
}

func writeEditorContentToFile(content string) error {
	filePath := "mud_content.txt"
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func startWebServer() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/editor", editorHandler)
	http.HandleFunc("/save-editor-content", saveEditorContentHandler)

	log.Printf("Web server listening on %s", webPort)
	log.Fatal(http.ListenAndServe(webPort, nil))
}

func main() {
	go startTelnetServer()
	startWebServer()
}