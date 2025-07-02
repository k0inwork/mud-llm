package models

// NPC represents a Non-Player Character.
type NPC struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	CurrentRoomID    string `json:"current_room_id"`
	Health           int    `json:"health"`
	MaxHealth        int    `json:"max_health"`
	Inventory        string `json:"inventory"`         // JSON array of item IDs and quantities
	OwnerIDs         string `json:"owner_ids"`         // JSON array of Owner IDs this NPC is associated with
	MemoriesAboutPlayers string `json:"memories_about_players"` // JSON object mapping player IDs to arrays of memory strings
	PersonalityPrompt string `json:"personality_prompt"`
	AvailableTools   string `json:"available_tools"`   // JSON array of conceptual tools LLM can call
	BehaviorState    string `json:"behavior_state"`    // JSON object for dynamic behavior state
}
