package models

type QuestOwner struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	LLMPromptContext     string `json:"llm_prompt_context"`
	CurrentInfluenceBudget float64 `json:"current_influence_budget"`
	MaxInfluenceBudget   float64 `json:"max_influence_budget"`
	BudgetRegenRate      float64 `json:"budget_regen_rate"`
	AssociatedQuestmakerIDs string `json:"associated_questmaker_ids"` // JSON array of questmaker IDs
}