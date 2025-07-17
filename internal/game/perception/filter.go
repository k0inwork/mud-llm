package perception

import (
	"fmt"
	"math"

	"mud/internal/dal"
	"mud/internal/game/events"
	"mud/internal/models"
)

// PerceptionFilter service handles the subjective interpretation of ActionEvents.
type PerceptionFilter struct {
	roomDAL      dal.RoomDALInterface
	raceDAL      dal.RaceDALInterface
	professionDAL dal.ProfessionDALInterface
	baseActionSignificance map[string]map[string]float64 // ActionType -> ObserverType -> Score
}

// NewPerceptionFilter creates a new PerceptionFilter.
func NewPerceptionFilter(
	roomDAL dal.RoomDALInterface,
	raceDAL dal.RaceDALInterface,
	professionDAL dal.ProfessionDALInterface,
) *PerceptionFilter {
	return &PerceptionFilter{
		roomDAL:      roomDAL,
		raceDAL:      raceDAL,
		professionDAL: professionDAL,
		baseActionSignificance: map[string]map[string]float64{
			"attack": {"npc": 10.0, "owner": 10.0, "questmaker": 10.0, "player": 10.0},
			"pray": {"npc": 2.0, "owner": 10.0, "questmaker": 2.0, "player": 5.0}, // Significant for owners
			"say": {"npc": 10.0, "owner": 1.0, "questmaker": 1.0, "player": 5.0},   // Significant for npcs
			"talk": {"npc": 10.0, "owner": 1.0, "questmaker": 1.0, "player": 5.0},  // Significant for npcs
			"use_skill": {"npc": 5.0, "owner": 5.0, "questmaker": 5.0, "player": 5.0},
			"magic_action": {"npc": 7.0, "owner": 7.0, "questmaker": 7.0, "player": 7.0},
			"combat_action": {"npc": 8.0, "owner": 8.0, "questmaker": 8.0, "player": 8.0},
			"subterfuge_action": {"npc": 6.0, "owner": 6.0, "questmaker": 6.0, "player": 6.0},
			"strange_magic": {"npc": 3.0, "owner": 3.0, "questmaker": 3.0, "player": 3.0},
			"unclear_action": {"npc": 0.5, "owner": 0.5, "questmaker": 0.5, "player": 0.5},
			"tamper_lock": {"npc": 7.0, "owner": 7.0, "questmaker": 7.0, "player": 7.0},
			"cast_hostile_spell": {"npc": 12.0, "owner": 12.0, "questmaker": 12.0, "player": 12.0},
			"attack_ally": {"npc": 15.0, "owner": 15.0, "questmaker": 15.0, "player": 15.0},
			"healing_magic": {"npc": 4.0, "owner": 4.0, "questmaker": 4.0, "player": 4.0},
			"arcane_weaving": {"npc": 6.0, "owner": 6.0, "questmaker": 6.0, "player": 6.0},
			"disable_trap": {"npc": 5.0, "owner": 5.0, "questmaker": 5.0, "player": 5.0},
			"gather_item": {"npc": 3.0, "owner": 3.0, "questmaker": 3.0, "player": 3.0},
			"deliver_item": {"npc": 2.0, "owner": 2.0, "questmaker": 2.0, "player": 2.0},
			"find_item": {"npc": 3.0, "owner": 3.0, "questmaker": 3.0, "player": 3.0},
			"return_item_to_npc": {"npc": 4.0, "owner": 4.0, "questmaker": 4.0, "player": 4.0},
			"observe_area": {"npc": 1.0, "owner": 1.0, "questmaker": 1.0, "player": 1.0},
			"report_to_npc": {"npc": 2.0, "owner": 2.0, "questmaker": 2.0, "player": 2.0},
			"defeat_dummy": {"npc": 5.0, "owner": 5.0, "questmaker": 5.0, "player": 5.0},
		},
	}
}

// Filter processes an ActionEvent through an observer's perception layers
// to produce a PerceivedAction.
func (pf *PerceptionFilter) Filter(event *events.ActionEvent, observer interface{}) (*PerceivedAction, error) {
	perceivedAction := &PerceivedAction{
		SourcePlayer: event.Player,
		Timestamp:    event.Timestamp,
		Clarity:      1.0, // Start with perfect clarity
	}

	// Determine observer type and get relevant biases
		var racialBiases map[string]float64
	var professionBiases map[string]float64
	var roomBiases map[string]float64
	var observerType string

	switch obs := observer.(type) {
	case *models.NPC:
		observerType = "npc"
		// Fetch racial biases
		if obs.RaceID != "" {
			race, err := pf.raceDAL.GetRaceByID(obs.RaceID)
			if err != nil {
				return nil, fmt.Errorf("failed to get race for NPC %s: %w", obs.ID, err)
			}
			if race != nil {
				racialBiases = race.PerceptionBiases
			}
		}

		// Fetch profession biases
		if obs.ProfessionID != "" {
			profession, err := pf.professionDAL.GetProfessionByID(obs.ProfessionID)
			if err != nil {
				return nil, fmt.Errorf("failed to get profession for NPC %s: %w", obs.ID, err)
			}
			if profession != nil {
				professionBiases = profession.PerceptionBiases
			}
		}

		// Fetch room and territory info for NPC
		room, err := pf.roomDAL.GetRoomByID(obs.CurrentRoomID)
		if err != nil {
			return nil, fmt.Errorf("failed to get room for NPC %s: %w", obs.ID, err)
		}
		if room != nil {
			roomBiases = room.PerceptionBiases
		}

	case *models.Owner:
		observerType = "owner"
		switch obs.MonitoredAspect {
		case "location":
			room, err := pf.roomDAL.GetRoomByID(obs.AssociatedID)
			if err != nil {
				return nil, fmt.Errorf("failed to get room for Owner %s: %w", obs.ID, err)
			}
			if room != nil {
				roomBiases = room.PerceptionBiases
			}
		case "race":
			race, err := pf.raceDAL.GetRaceByID(obs.AssociatedID)
			if err != nil {
				return nil, fmt.Errorf("failed to get race for Owner %s: %w", obs.ID, err)
			}
			if race != nil {
				racialBiases = race.PerceptionBiases
			}
		case "profession":
			profession, err := pf.professionDAL.GetProfessionByID(obs.AssociatedID)
			if err != nil {
				return nil, fmt.Errorf("failed to get profession for Owner %s: %w", obs.ID, err)
			}
			if profession != nil {
				professionBiases = profession.PerceptionBiases
			}
		}
	case *models.Questmaker:
		observerType = "questmaker"
		// Questmakers are associated with quests, not directly with locations/races/professions.
		// Their perception might be more abstract or tied to quest objectives.
		// For now, we'll give them a neutral bias.
		racialBiases = map[string]float64{}
		professionBiases = map[string]float64{}
		roomBiases = map[string]float64{}
	case *models.Player:
		observerType = "player"
		// Fetch racial biases for player
		if obs.RaceID != "" {
			race, err := pf.raceDAL.GetRaceByID(obs.RaceID)
			if err != nil {
				return nil, fmt.Errorf("failed to get race for Player %s: %w", obs.ID, err)
			}
			if race != nil {
				racialBiases = race.PerceptionBiases
			}
		}

		// Fetch profession biases for player
		if obs.ProfessionID != "" {
			profession, err := pf.professionDAL.GetProfessionByID(obs.ProfessionID)
			if err != nil {
				return nil, fmt.Errorf("failed to get profession for Player %s: %w", obs.ID, err)
			}
			if profession != nil {
				professionBiases = profession.PerceptionBiases
			}
		}

		// Fetch room and territory info for player
		room, err := pf.roomDAL.GetRoomByID(obs.CurrentRoomID)
			if err != nil {
				return nil, fmt.Errorf("failed to get room for Player %s: %w", obs.ID, err)
			}
			if room != nil {
				roomBiases = room.PerceptionBiases
			}
	default:
		return nil, fmt.Errorf("unsupported observer type: %T", observer)
	}

	// Set BaseSignificance based on observer type and action type
	if scores, ok := pf.baseActionSignificance[event.ActionType]; ok {
		if score, ok := scores[observerType]; ok {
			perceivedAction.BaseSignificance = score
		} else {
			perceivedAction.BaseSignificance = 1.0 // Default low significance if observer type not specified
		}
	} else {
		perceivedAction.BaseSignificance = 0.5 // Default very low significance if action type not specified
	}

	// Layer 0: Physical Sensory Check (Placeholder)
	// For now, assume perfect sensory input.
	// In a real implementation, this would check line of sight, distance, etc.

	// Layer 1: Innate & Cultural Bias (Racial and Territorial)
	// Apply racial biases
	if racialBiases != nil {
		// Check for bias by ActionType
		if bias, ok := racialBiases[event.ActionType]; ok {
			perceivedAction.Clarity += bias
		}
		// Check for bias by Skill Category if skill is used
		if event.SkillUsed != nil && event.SkillUsed.Category != "" {
			if bias, ok := racialBiases[event.SkillUsed.Category]; ok {
				perceivedAction.Clarity += bias
			}
		}
	}

	// Apply territorial biases (from room)
	if roomBiases != nil {
		// Check for bias by ActionType
		if bias, ok := roomBiases[event.ActionType]; ok {
			perceivedAction.Clarity += bias
		}
		// Check for bias by Skill Category if skill is used
		if event.SkillUsed != nil && event.SkillUsed.Category != "" {
			if bias, ok := roomBiases[event.SkillUsed.Category]; ok {
				perceivedAction.Clarity += bias
			}
		}
	}

	// Layer 2: Knowledge & Experience (Profession/Class and Skill Proficiency)
	if professionBiases != nil {
		// Check for bias by ActionType
		if bias, ok := professionBiases[event.ActionType]; ok {
			perceivedAction.Clarity += bias
		}
		// Check for bias by Skill Category if skill is used
		if event.SkillUsed != nil && event.SkillUsed.Category != "" {
			if bias, ok := professionBiases[event.SkillUsed.Category]; ok {
				perceivedAction.Clarity += bias
			}
		}
	}

	// Placeholder for Skill Proficiency:
	// If observer has a relevant skill, increase clarity.
	// This would involve checking observer's skills against event.SkillUsed.

	// Layer 3: Explicit Modifiers (Passive Skills & Buffs) (Placeholder)
	// This would involve checking observer's active buffs/debuffs or passive skills.

	// Cap clarity between 0.0 and 1.0
	perceivedAction.Clarity = math.Max(0.0, math.Min(1.0, perceivedAction.Clarity))

	// Determine PerceivedActionType based on Clarity and ActionType/SkillCategory
	perceivedAction.PerceivedActionType = pf.determinePerceivedActionType(event, perceivedAction.Clarity)

	// Placeholder for ApparentSkillLevel and IsCriminal
	perceivedAction.ApparentSkillLevel = int(perceivedAction.Clarity * 100) // Simple mapping for now
	perceivedAction.IsCriminal = false                                     // Default to not criminal

	return perceivedAction, nil
}

// determinePerceivedActionType maps clarity and action details to a perceived action type.
func (pf *PerceptionFilter) determinePerceivedActionType(event *events.ActionEvent, clarity float64) string {
	// Check if the action type is explicitly defined in our base significance map
	_, actionTypeDefined := pf.baseActionSignificance[event.ActionType]

	if clarity > 0.9 {
		// Highly clear, use specific action/skill name if available, otherwise ActionType
		if event.SkillUsed != nil && event.SkillUsed.Name != "" {
			return event.SkillUsed.Name // e.g., "Lesser Heal"
		}
		if actionTypeDefined {
			return event.ActionType // e.g., "say"
		}
		return "unclear_action" // Fallback for high clarity but undefined action
	} else if clarity > 0.5 {
		// Moderately clear, use a category if available, otherwise ActionType_general
		if event.SkillUsed != nil && event.SkillUsed.Category != "" {
			return event.SkillUsed.Category + "_action" // e.g., "magic_action", "combat_action"
		}
		if actionTypeDefined {
			return event.ActionType + "_general" // e.g., "say_general"
		}
		return "unclear_action" // Fallback for moderate clarity but undefined action
	} else {
		// Low clarity, vague perception
		if event.SkillUsed != nil && event.SkillUsed.Category != "" {
			return "strange_" + event.SkillUsed.Category // e.g., "strange_magic"
		}
		return "unclear_action" // Default if no specific skill category or action type is known
	}
}