package models

// Class defines the available player classes and their associated skill trees.
type Class struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	TotalLevels        int    `json:"total_levels"`
	ParentClassID      string `json:"parent_class_id"`
	AssociatedEntityType string `json:"associated_entity_type"` // e.g., "Questmaker", "Owner"
	AssociatedEntityID string `json:"associated_entity_id"` // ID of the specific LLM entity
	LevelUpRewards     string `json:"level_up_rewards"`     // JSON object mapping level to skill choices/unlocks
}
