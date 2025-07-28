package models

import (
	"time"
)

// PlayerCharacter represents an in-game character, linked to a PlayerAccount.
type PlayerCharacter struct {
	ID              string    `json:"id"`
	PlayerAccountID string    `json:"player_account_id"`
	Name            string    `json:"name"`
	RaceID          string    `json:"race_id"`
	ProfessionID    string    `json:"profession_id"`
	CurrentRoomID   string    `json:"current_room_id"`
	Health          int       `json:"health"`
	MaxHealth       int       `json:"max_health"`
	Inventory       string    `json:"inventory"` // JSON string
	VisitedRoomIDs  string    `json:"visited_room_ids"` // JSON string
	CreatedAt       time.Time `json:"created_at"`
	LastPlayedAt    time.Time `json:"last_played_at,omitempty"`
}
