package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	telnetServerAddress = "localhost:4000"
	// ANSI color codes for stripping
	ansiReset  = "\033[0m"
	ansiSuccess = "\033[32m"
	ansiDefault = "\033[37m"
)

// TestTelnetServer_ConnectionAndEcho tests basic connection and echo functionality.
func TestTelnetServer_ConnectionAndEcho(t *testing.T) {
	// Ensure the server is running before starting tests.
	// In a real CI/CD environment, you'd start the server here.
	// For this exercise, we assume the server is already running from `go run main.go &`

	// Give the server a moment to fully start up
	time.Sleep(500 * time.Millisecond)

	conn, err := net.Dial("tcp", telnetServerAddress)
	assert.NoError(t, err, "Failed to connect to Telnet server")
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read welcome message
	welcome, err := readLine(reader)
	assert.NoError(t, err, "Failed to read welcome message")
	assert.Contains(t, stripANSI(welcome), "Welcome to GoMUD! Type 'help' for commands.", "Welcome message mismatch")

	// Send a command
	command := "look"
	_, err = fmt.Fprintf(conn, "%s\n", command)
	assert.NoError(t, err, "Failed to send command")

	// Read echo response
	echo, err := readLine(reader)
	assert.NoError(t, err, "Failed to read echo response")
	assert.Contains(t, stripANSI(echo), fmt.Sprintf("You typed: %s", command), "Echo response mismatch")

	// Send another command
	command2 := "say hello"
	_, err = fmt.Fprintf(conn, "%s\n", command2)
	assert.NoError(t, err, "Failed to send second command")

	// Read echo response for second command
	echo2, err := readLine(reader)
	assert.NoError(t, err, "Failed to read second echo response")
	assert.Contains(t, stripANSI(echo2), fmt.Sprintf("You typed: %s", command2), "Second echo response mismatch")
}

// readLine reads a line from the bufio.Reader, handling Telnet's \r\n.
func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(line, "\r\n"), nil
}

// stripANSI removes ANSI escape codes from a string.
func stripANSI(str string) string {
	// This is a simplified stripper. A more robust one might use regex.
	str = strings.ReplaceAll(str, ansiReset, "")
	str = strings.ReplaceAll(str, ansiSuccess, "")
	str = strings.ReplaceAll(str, ansiDefault, "")
	// Add more as needed
	return str
}


