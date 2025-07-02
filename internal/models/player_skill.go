package models

// PlayerSkill tracks skills learned by players and their current percentage.
type PlayerSkill struct {
	PlayerID          string `json:"player_id"`
	SkillID           string `json:"skill_id"`
	Percentage        int    `json:"percentage"`
	GrantedByEntityType string `json:"granted_by_entity_type"` // e.g., "Questmaker", "Owner"
	GrantedByEntityID string `json:"granted_by_entity_id"` // ID of the specific LLM entity
}
