package models

// Owner represents an LLM-driven world guardian.
type Owner struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	Description          string  `json:"description"`
	MonitoredAspect      string  `json:"monitored_aspect"`
	AssociatedID         string  `json:"associated_id"`
	LLMPromptContext     string  `json:"llm_prompt_context"`
	MemoriesAboutPlayers map[string][]string `json:"memories_about_players"` // Map of player IDs to arrays of private memory strings
	CurrentInfluenceBudget float64 `json:"current_influence_budget"`
	MaxInfluenceBudget   float64 `json:"max_influence_budget"`
	BudgetRegenRate      float64 `json:"budget_regen_rate"`
	AvailableTools       []Tool  `json:"available_tools"`   // Array of conceptual tools LLM can call
	InitiatedQuests      []string `json:"initiated_quests"` // Array of quest IDs this owner can initiate/offer
	ReactionThreshold    int      `json:"reaction_threshold"`
}
