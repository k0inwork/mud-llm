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
	Inventory        string    `json:"inventory"`         // JSON array of item IDs and quantities
	VisitedRoomIDs   string    `json:"visited_room_ids"`  // JSON array of room IDs
	CreatedAt        time.Time `json:"created_at"`
	LastLoginAt      time.Time `json:"last_login_at"`
	LastLogoutAt     time.Time `json:"last_logout_at"`
}
