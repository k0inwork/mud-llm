package models

import "time"

// PlayerAccount represents a user's account, which can have multiple characters.
type PlayerAccount struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"` // Do not expose hashed password in JSON
	Email          string    `json:"email,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	LastLoginAt    time.Time `json:"last_login_at,omitempty"`
}
