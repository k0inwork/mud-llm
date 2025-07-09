package models

import (
	"time"
)

// Player represents a player character in the MUD.
type Player struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	RaceID           string    `json:"race_id"`
	ProfessionID     string    `json:"profession_id"`
	CurrentRoomID    string    `json:"current_room_id"`
	Health           int       `json:"health"`
	MaxHealth        int       `json:"max_health"`
	Inventory        []string  `json:"inventory"`         // Array of item IDs
	VisitedRoomIDs   map[string]bool `json:"visited_room_ids"`  // Map of room IDs to boolean (for quick lookup)
	CreatedAt        time.Time `json:"created_at"`
	LastLoginAt      time.Time `json:"last_login_at"`
	LastLogoutAt     time.Time `json:"last_logout_at"`
}
