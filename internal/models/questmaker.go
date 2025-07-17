package models

type Questmaker struct {
	ID                   string            `json:"id"`
	Name                 string            `json:"name"`
	LLMPromptContext     string            `json:"llm_prompt_context"`
	CurrentInfluenceBudget float64           `json:"current_influence_budget"`
	MaxInfluenceBudget   float64           `json:"max_influence_budget"`
	BudgetRegenRate      float64           `json:"budget_regen_rate"`
	MemoriesAboutPlayers map[string][]string `json:"memories_about_players"`
	AvailableTools       []Tool            `json:"available_tools"`
	ReactionThreshold    int               `json:"reaction_threshold"`
}
