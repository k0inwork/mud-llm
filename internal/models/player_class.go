package models

// PlayerClass tracks a player's progression in each class they have acquired.
type PlayerClass struct {
	PlayerID   string `json:"player_id"`
	ClassID    string `json:"class_id"`
	Level      int    `json:"level"`
	Experience int    `json:"experience"`
}
