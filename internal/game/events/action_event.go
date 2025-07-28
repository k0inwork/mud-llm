package events

import (
	"mud/internal/models"
	"time"
)

// ActionEvent represents an action taken by a player in the game.
type ActionEvent struct {
	Player     *models.PlayerCharacter
	ActionType string
	Room       *models.Room
	Timestamp  time.Time
	SkillUsed  *models.Skill
	Targets    []interface{} // Can be *models.NPC, *models.Item, etc.
}
