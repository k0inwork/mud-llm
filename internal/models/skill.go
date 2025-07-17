package models

// Skill defines available skills (active and passive).
type Skill struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`	
	Category           string `json:"category"` // e.g., "combat", "magic", "social", "crafting"
	Description        string `json:"description"`
	Type               string `json:"type"` // e.g., "active", "passive"
	AssociatedClassID  string `json:"associated_class_id"`
	GrantedByEntityType string `json:"granted_by_entity_type"` // e.g., "Questmaker", "Owner"
	GrantedByEntityID  string `json:"granted_by_entity_id"` // ID of the specific LLM entity
	Effects            string `json:"effects"` // JSON array of structured effect objects
	Cost               int    `json:"cost"`
	Cooldown           int    `json:"cooldown"`
	MinClassLevel      int    `json:"min_class_level"`
}
