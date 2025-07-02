package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"mud/internal/presentation"
	"strings"
)

// TelnetServer represents the Telnet server for the MUD.
type TelnetServer struct {
	port    string
	renderer *presentation.TelnetRenderer
}

// NewTelnetServer creates a new TelnetServer.
func NewTelnetServer(port string) *TelnetServer {
	return &TelnetServer{
		port:    port,
		renderer: presentation.NewTelnetRenderer(),
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

		// Process input (for now, just echo it back)
		processedInput := fmt.Sprintf("You typed: %s", strings.TrimSpace(input))
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
