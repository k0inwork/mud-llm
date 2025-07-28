package mocks

import (
	"fmt"
	"mud/internal/presentation"
	"strings"
	"sync"
)

// TestRenderer is a mock implementation of the TelnetRendererInterface for testing.
type TestRenderer struct {
	renderedMessages []string
	mu               sync.Mutex
}

// NewTestRenderer creates a new TestRenderer.
func NewTestRenderer() *TestRenderer {
	return &TestRenderer{
		renderedMessages: []string{},
	}
}

// RenderMessage appends a simple, predictable string representation of a SemanticMessage to the internal slice.
func (r *TestRenderer) RenderMessage(msg presentation.SemanticMessage) string {
	rendered := fmt.Sprintf("[%s] %s\n", msg.Type, msg.Content)
	r.mu.Lock()
	r.renderedMessages = append(r.renderedMessages, rendered)
	r.mu.Unlock()
	return rendered
}

// RenderRawString appends a simple, predictable string representation of a raw string with a color to the internal slice.
func (r *TestRenderer) RenderRawString(s string, color presentation.SemanticColorType) string {
	rendered := fmt.Sprintf("[%s] %s\n", color, s)
	r.mu.Lock()
	r.renderedMessages = append(r.renderedMessages, rendered)
	r.mu.Unlock()
	return rendered
}

// GetRenderedMessages returns all messages rendered so far and clears the internal slice.
func (r *TestRenderer) GetRenderedMessages() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	messages := r.renderedMessages
	r.renderedMessages = []string{} // Clear after reading
	return messages
}

// ClearRenderedMessages clears the internal slice of rendered messages.
func (r *TestRenderer) ClearRenderedMessages() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.renderedMessages = []string{}
}

// ContainsMessage checks if a specific message is present in the rendered messages.
func (r *TestRenderer) ContainsMessage(expected string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, msg := range r.renderedMessages {
		if strings.Contains(msg, expected) {
			return true
		}
	}
	return false
}

// AllMessages returns all messages rendered so far as a single string.
func (r *TestRenderer) AllMessages() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return strings.Join(r.renderedMessages, "")
}

