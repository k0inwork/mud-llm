package models

// Quest represents a definition of a quest.
type Quest struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	QuestOwnerID       string            `json:"quest_owner_id"` // ID of the thematic Quest Owner
	QuestmakerID       string            `json:"questmaker_id"`
	InfluencePointsMap map[string]float64 `json:"influence_points_map"` // Map of player actions to influence points granted
	Objectives         string            `json:"objectives"`         // JSON array of quest objectives
	Rewards            string            `json:"rewards"`            // JSON array of rewards
}
