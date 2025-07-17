package events

import (
	"fmt"
	"sync"
)

// EventType is a string alias for event types.
type EventType string

const (
	// ActionEventType represents a player action event.
	ActionEventType EventType = "ActionEvent"
	// Add other event types here as needed
)

// EventBus manages the subscription and publication of events.
type EventBus struct {
	subscribers map[EventType][]chan interface{}
	mu          sync.RWMutex
}

// NewEventBus creates a new EventBus instance.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]chan interface{}),
	}
}

// Subscribe adds a new subscriber channel for a specific event type.
func (bus *EventBus) Subscribe(eventType EventType, handlerChannel chan interface{}) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.subscribers[eventType] = append(bus.subscribers[eventType], handlerChannel)
}

// Publish sends an event to all registered subscribers for that event type.
// This is currently synchronous. For a more advanced system, this could be made asynchronous.
func (bus *EventBus) Publish(eventType EventType, event interface{}) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	if channels, found := bus.subscribers[eventType]; found {
		for _, ch := range channels {
			// In a production system, you might run each subscriber in a separate goroutine
			// with a panic handler to prevent one subscriber from crashing the bus.
			select {
			case ch <- event:
				// Event sent successfully
			default:
				// Non-blocking send, drop event if channel is full
				fmt.Printf("EventBus: Dropping event of type %s, channel is full.\n", eventType)
			}
		}
	}
}