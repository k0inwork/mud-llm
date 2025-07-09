package models

// SkillInfo represents a skill and its associated percentage.
type SkillInfo struct {
	SkillID    string `json:"skill_id"`
	Percentage int    `json:"percentage"`
}
