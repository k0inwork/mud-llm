package models

// Profession defines available player professions/classes.
type Profession struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BaseSkills  string `json:"base_skills"` // JSON array of skill IDs granted by profession
}
