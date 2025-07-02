package models

// Questmaker represents an LLM-driven quest entity.
type Questmaker struct {
	ID                   string `json:"id"`
	Name                 string `json:name"`
	LLMPromptContext     string `json:"llm_prompt_context"`
	CurrentInfluenceBudget float64 `json:"current_influence_budget"`
	MaxInfluenceBudget   float64 `json:"max_influence_budget"`
	BudgetRegenRate      float64 `json:"budget_regen_rate"`
	MemoriesAboutPlayers string `json:"memories_about_players"` // JSON object mapping player IDs to arrays of private memory strings
	AvailableTools       string `json:"available_tools"`   // JSON array of conceptual tools LLM can call
}
