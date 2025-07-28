package sentiententitymanager

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"mud/internal/dal"
	"mud/internal/game"
	"mud/internal/game/events"
	"mud/internal/game/perception"
	"mud/internal/models"
)

// SentientEntityManager orchestrates AI responses based on significant player actions.
type SentientEntityManager struct {
	llmService    game.LLMServiceInterface
	npcDAL        dal.NPCDALInterface
	ownerDAL      dal.OwnerDALInterface
	questmakerDAL dal.QuestmakerDALInterface
	toolDispatcher game.ToolDispatcherInterface
	telnetRenderer game.TelnetRendererInterface
	eventBus      *events.EventBus
}

// NewSentientEntityManager creates a new SentientEntityManager.
func NewSentientEntityManager(
	llmService game.LLMServiceInterface,
	npcDAL dal.NPCDALInterface,
	ownerDAL dal.OwnerDALInterface,
	questmakerDAL dal.QuestmakerDALInterface,
	toolDispatcher game.ToolDispatcherInterface,
	telnetRenderer game.TelnetRendererInterface,
	eventBus *events.EventBus,
) *SentientEntityManager {
	return &SentientEntityManager{
		llmService:    llmService,
		npcDAL:        npcDAL,
		ownerDAL:      ownerDAL,
		questmakerDAL: questmakerDAL,
		toolDispatcher: toolDispatcher,
		telnetRenderer: telnetRenderer,
		eventBus:      eventBus,
	}
}

func (m *SentientEntityManager) TriggerReaction(observer interface{}, perceivedActions []perception.PerceivedActionRecord) error {
	logrus.Printf("Triggering reaction for entity %s", getObserverID(observer))

	if len(perceivedActions) == 0 {
		return fmt.Errorf("no perceived actions provided for reaction")
	}

	// For now, let's just log the actions. Later, these will be used to build the prompt.
	for _, record := range perceivedActions {
		logrus.Printf("  Perceived Action: %s (Clarity: %.2f, Significance: %.2f)", record.PerceivedAction.PerceivedActionType, record.PerceivedAction.Clarity, record.Significance)
	}

	// 2. Retrieve the entity (NPC, Owner, or Questmaker)
	var entity interface{}
	var err error

	entityID := getObserverID(observer)

	// Attempt to retrieve as NPC
	npc, err := m.npcDAL.GetNPCByID(entityID)
	if err == nil && npc != nil {
		entity = npc
	} else {
		// Attempt to retrieve as Owner
		owner, err2 := m.ownerDAL.GetOwnerByID(entityID)
		if err2 == nil && owner != nil {
			entity = owner
		} else {
			// Attempt to retrieve as Questmaker
			questmaker, err3 := m.questmakerDAL.GetQuestmakerByID(entityID)
			if err3 == nil && questmaker != nil {
				entity = questmaker
			} else {
				return fmt.Errorf("could not find entity with ID %s for triggering reaction: %w", entityID, err)
			}
		}
	}

	if entity == nil {
		return fmt.Errorf("entity %s not found after all DAL lookups.", entityID)
	}

	// 3. Retrieve the player (from the first perceived action)
	player := perceivedActions[0].PerceivedAction.SourcePlayer
	if player == nil {
		return fmt.Errorf("source player not found in perceived action")
	}

	// Determine the entity's reaction threshold
	var reactionThreshold float64
	switch e := entity.(type) {
	case *models.NPC:
		reactionThreshold = float64(e.ReactionThreshold)
	case *models.Owner:
		reactionThreshold = float64(e.ReactionThreshold)
	case *models.Questmaker:
		reactionThreshold = float64(e.ReactionThreshold)
	default:
		return fmt.Errorf("unsupported entity type for reaction threshold: %T", entity)
	}

	// Filter perceived actions based on reaction threshold
	var relevantPerceivedActions []perception.PerceivedActionRecord
	for _, record := range perceivedActions {
		if record.Significance >= reactionThreshold {
			relevantPerceivedActions = append(relevantPerceivedActions, record)
		}
	}

	if len(relevantPerceivedActions) == 0 {
		logrus.Printf("No perceived actions met the reaction threshold (%.2f) for entity %s. No LLM call made.", reactionThreshold, entityID)
		return nil // No reaction needed if no actions meet the threshold
	}

	// 4. Construct prompt (placeholder for now, will use PromptAssembler)
	// For now, just a simple prompt based on the first RELEVANT action
	prompt := fmt.Sprintf("Player %s performed action %s (clarity %.2f). Respond to this.", player.Name, relevantPerceivedActions[0].PerceivedAction.PerceivedActionType, relevantPerceivedActions[0].PerceivedAction.Clarity)

	// 5. Send to LLM
	llmResponse, err := m.llmService.ProcessAction(context.Background(), entity, player, prompt)
	if err != nil {
		return fmt.Errorf("LLM Service ProcessAction failed for entity %s: %w", entityID, err)
	}

	// 6. Handle LLM Response
	if llmResponse != nil {
		// Publish narrative to player
		if llmResponse.Narrative != "" {
			playerMessage := &events.PlayerMessageEvent{
				PlayerID: player.ID,
				Content:  fmt.Sprintf("%s says: %s", getObserverName(observer), llmResponse.Narrative),
			}
			m.eventBus.Publish(events.PlayerMessageEventType, playerMessage)
			logrus.Printf("LLM Narrative for %s: %s", entityID, llmResponse.Narrative)
		}

		// Dispatch tool calls
		if len(llmResponse.ToolCalls) > 0 {
			logrus.Printf("LLM Tool Calls for %s: %+v", entityID, llmResponse.ToolCalls)
			err = m.toolDispatcher.Dispatch(context.Background(), player, entity, llmResponse.ToolCalls)
			if err != nil {
				logrus.Errorf("Failed to dispatch tool calls for entity %s: %v", entityID, err)
			}
		}
	}
	return nil
}

// Helper to get observer name
func getObserverName(observer interface{}) string {
	switch obs := observer.(type) {
	case *models.NPC:
		return obs.Name
	case *models.Owner:
		return obs.Name
	case *models.Questmaker:
		return obs.Name
	default:
		return "Unknown Entity"
	}
}

// Helper to get observer ID
func getObserverID(observer interface{}) string {
	switch obs := observer.(type) {
	case *models.NPC:
		return obs.ID
	case *models.Owner:
		return obs.ID
	case *models.Questmaker:
		return obs.ID
	default:
		return ""
	}
}
