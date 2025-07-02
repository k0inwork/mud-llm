package presentation

import (
	"fmt"
	"strings"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Black  = "\033[30m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Magenta = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	BrightBlack = "\033[90m"
	BrightRed   = "\033[91m"
	BrightGreen = "\033[92m"
	BrightYellow = "\033[93m"
	BrightBlue  = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan  = "\033[96m"
	BrightWhite = "\033[97m"
)

// TelnetRenderer translates semantic messages into ANSI-colored strings.
type TelnetRenderer struct {
	colorMap map[SemanticColorType]string
}

// NewTelnetRenderer creates a new TelnetRenderer with a default color map.
func NewTelnetRenderer() *TelnetRenderer {
	return &TelnetRenderer{
		colorMap: map[SemanticColorType]string{
			ColorDefault:   White,
			ColorHighlight: Yellow,
			ColorSuccess:   Green,
			ColorError:     Red,
			ColorNPC:       Cyan,
			ColorPlayer:    Blue,
			ColorItem:      BrightYellow,
			ColorLore:      Magenta,
			ColorQuest:     BrightBlue,
			ColorOwner:     BrightMagenta,
		},
	}
}

// RenderMessage converts a SemanticMessage into an ANSI-colored string.
func (r *TelnetRenderer) RenderMessage(msg SemanticMessage) string {
	colorCode, ok := r.colorMap[msg.Color]
	if !ok {
		colorCode = r.colorMap[ColorDefault] // Fallback to default color
	}

	var builder strings.Builder
	builder.WriteString(colorCode)
	builder.WriteString(msg.Content)
	builder.WriteString(Reset)
	builder.WriteString("\r\n") // Add carriage return and newline for Telnet clients

	return builder.String()
}

// RenderRawString applies a specific semantic color to a raw string.
func (r *TelnetRenderer) RenderRawString(s string, color SemanticColorType) string {
	colorCode, ok := r.colorMap[color]
	if !ok {
		colorCode = r.colorMap[ColorDefault]
	}
	return fmt.Sprintf("%s%s%s", colorCode, s, Reset)
}
