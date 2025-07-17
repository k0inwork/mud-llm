package sentiententitymanager

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"mud/internal/dal"
	"mud/internal/game"
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
}

// NewSentientEntityManager creates a new SentientEntityManager.
func NewSentientEntityManager(
	llmService game.LLMServiceInterface,
	npcDAL dal.NPCDALInterface,
	ownerDAL dal.OwnerDALInterface,
	questmakerDAL dal.QuestmakerDALInterface,
	toolDispatcher game.ToolDispatcherInterface,
	telnetRenderer game.TelnetRendererInterface,
) *SentientEntityManager {
	return &SentientEntityManager{
		llmService:    llmService,
		npcDAL:        npcDAL,
		ownerDAL:      ownerDAL,
		questmakerDAL: questmakerDAL,
		toolDispatcher: toolDispatcher,
		telnetRenderer: telnetRenderer,
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
		owner, err := m.ownerDAL.GetOwnerByID(entityID)
		if err == nil && owner != nil {
			entity = owner
		} else {
			// Attempt to retrieve as Questmaker
			questmaker, err := m.questmakerDAL.GetQuestmakerByID(entityID)
			if err == nil && questmaker != nil {
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

	// 4. Construct prompt (placeholder for now, will use PromptAssembler)
	// For now, just a simple prompt based on the first action
	prompt := fmt.Sprintf("Player %s performed action %s (clarity %.2f). Respond to this.", player.Name, perceivedActions[0].PerceivedAction.PerceivedActionType, perceivedActions[0].PerceivedAction.Clarity)

	// 5. Send to LLM
	llmResponse, err := m.llmService.ProcessAction(context.Background(), entity, player, prompt)
	if err != nil {
		return fmt.Errorf("LLM Service ProcessAction failed for entity %s: %w", entityID, err)
	}

	// 6. Handle LLM Response
	if llmResponse != nil {
		// Send narrative to player (placeholder: need player's connection)
		// For now, just log it.
		logrus.Printf("LLM Narrative for %s: %s", entityID, llmResponse.Narrative)

		// Dispatch tool calls (placeholder: need to implement tool dispatcher logic)
		if len(llmResponse.ToolCalls) > 0 {
			logrus.Printf("LLM Tool Calls for %s: %+v", entityID, llmResponse.ToolCalls)
			// m.toolDispatcher.Dispatch(ctx, player, entity, llmResponse.ToolCalls)
		}
	}
	return nil
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