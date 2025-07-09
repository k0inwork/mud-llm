package models

// Tool represents a conceptual tool that an LLM-driven entity can call.
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}
