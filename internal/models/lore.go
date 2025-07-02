package models

// Lore stores various pieces of lore, categorized by scope.
type Lore struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Scope        string `json:"scope"`        // e.g., "global", "zone", "faction", "item"
	AssociatedID string `json:"associated_id"` // ID of entity/zone/faction if scope is not global
}
