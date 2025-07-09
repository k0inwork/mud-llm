package models

// NPC represents a Non-Player Character.
type NPC struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	CurrentRoomID    string `json:"current_room_id"`
	Health           int    `json:"health"`
	MaxHealth        int    `json:"max_health"`
	Inventory        []string `json:"inventory"`         // Array of item IDs
	OwnerIDs         []string `json:"owner_ids"`         // Array of Owner IDs this NPC is associated with
	MemoriesAboutPlayers map[string][]string `json:"memories_about_players"` // Map of player IDs to arrays of memory strings
	PersonalityPrompt string `json:"personality_prompt"`
	AvailableTools   []Tool `json:"available_tools"`   // Array of conceptual tools LLM can call
	BehaviorState    string `json:"behavior_state"`    // JSON object for dynamic behavior state
}
