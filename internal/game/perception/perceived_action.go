package perception

import (
	"mud/internal/models"
	"time"
)

// PerceivedAction represents an entity's subjective interpretation of an ActionEvent.
// It is the output of the PerceptionFilter.
type PerceivedAction struct {
	Observer     interface{} // The entity doing the perceiving.
	SourcePlayer *models.Player
	Target       interface{} // The specific target of this perception.

	// PerceivedActionType is a string representing the observer's understanding of the action.
	// e.g., "tamper_lock", "cast_hostile_spell", "attack_ally".
	PerceivedActionType string

	// Clarity indicates how well the action was understood, from 0.0 (not at all) to 1.0 (perfectly).
	Clarity float64

	// ApparentSkillLevel is the observer's guess as to how skillfully the action was performed (1-100).
	ApparentSkillLevel int

	// IsCriminal indicates if the action was perceived as illegal according to the observer's morals/laws.
	IsCriminal bool

	Timestamp time.Time
	BaseSignificance float64 // The base significance score for this perceived action, before clarity is applied.
}

// PerceivedActionRecord stores a perceived action and its calculated significance.
type PerceivedActionRecord struct {
	PerceivedAction *PerceivedAction
	Significance    float64
	Timestamp       time.Time
}
