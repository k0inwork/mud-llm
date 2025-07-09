package models

// Profession defines available player professions/classes.
type Profession struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BaseSkills  []SkillInfo `json:"base_skills"` // Array of SkillInfo (skill ID and percentage)
}
