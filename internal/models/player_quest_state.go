package models

import "time"

// PlayerQuestState stores the dynamic state of a player's progress on active quests.
type PlayerQuestState struct {
	PlayerID                      string    `json:"player_id"`
	QuestID                       string    `json:"quest_id"`
	CurrentProgress               string    `json:"current_progress"` // JSON object tracking objective progress
	LastActionTimestamp           time.Time `json:"last_action_timestamp"`
	QuestmakerInfluenceAccumulated float64   `json:"questmaker_influence_accumulated"`
	Status                        string    `json:"status"` // e.g., "active", "completed", "failed", "abandoned"
}