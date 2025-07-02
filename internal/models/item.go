package models

// Item represents a definition of an item type.
type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`       // e.g., "weapon", "armor", "consumable", "quest_item"
	Properties  string `json:"properties"` // JSON object for item-specific properties
}
