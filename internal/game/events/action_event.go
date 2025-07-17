package events

import (
	"mud/internal/models"
	"time"
)

// ActionEvent represents the objective ground truth of a player's action.
// It is created by a command handler and published to the event bus for processing.
type ActionEvent struct {
	Player     *models.Player
	ActionType string        // The canonical action type, e.g., "use_skill", "say".
	SkillUsed  *models.Skill // A reference to the skill model, if applicable.
	Targets    []interface{} // A slice of entities (NPCs, Items, Players, etc.) to support AoE.
	Room       *models.Room
	Timestamp  time.Time
	// Metadata can be used for additional context, e.g., the text of a "say" command.
	Metadata map[string]string
}
