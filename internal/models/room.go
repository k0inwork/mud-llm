package models

// Exit represents a single exit from a room.
type Exit struct {
	Direction    string `json:"direction"`
	TargetRoomID string `json:"TargetRoomID"`
	IsLocked     bool   `json:"is_locked"`
	KeyID        string `json:"key_id,omitempty"`
}

// Room represents a game room or location.
type Room struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Exits       string `json:"exits"`     // JSON object mapping directions to room IDs
	OwnerID          string                 `json:"owner_id"` // Optional: ID of the Owner controlling this room
	TerritoryID      string                 `json:"territory_id"` // ID for grouping rooms into a larger territory
	Properties       string                 `json:"properties"` // JSON object for dynamic room properties
	PerceptionBiases map[string]float64     `json:"perception_biases"` // Territorial biases
}