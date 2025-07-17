package actionsignificance

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"mud/internal/dal"
	"mud/internal/game/events"
	"mud/internal/game"
	"mud/internal/game/perception"
	"mud/internal/models"
)

// ActionBuffer stores perceived actions for a specific player and entity.
type ActionBuffer struct {
	mu      sync.Mutex
	records []perception.PerceivedActionRecord
}

// ActionSignificanceMonitor manages player actions and triggers LLM interactions.
type ActionSignificanceMonitor struct {
	eventBus            *events.EventBus
	perceptionFilter    game.PerceptionFilterInterface
	npcDAL              dal.NPCDALInterface
	ownerDAL            dal.OwnerDALInterface
	questmakerDAL       dal.QuestmakerDALInterface
	sentientEntityManager game.SentientEntityManagerInterface
	playerEntityBuffers map[string]map[string]*ActionBuffer // playerID -> entityID -> *ActionBuffer
	mu                  sync.RWMutex
}

// NewMonitor creates a new ActionSignificanceMonitor.
func NewMonitor(
	eventBus *events.EventBus,
	perceptionFilter game.PerceptionFilterInterface,
	npcDAL dal.NPCDALInterface,
	ownerDAL dal.OwnerDALInterface,
	questmakerDAL dal.QuestmakerDALInterface,
	sentientEntityManager game.SentientEntityManagerInterface,
) *ActionSignificanceMonitor {
	m := &ActionSignificanceMonitor{
		eventBus:            eventBus,
		perceptionFilter:    perceptionFilter,
		npcDAL:              npcDAL,
		ownerDAL:            ownerDAL,
		questmakerDAL:       questmakerDAL,
		sentientEntityManager: sentientEntityManager,
		playerEntityBuffers: make(map[string]map[string]*ActionBuffer),
	}

	// Subscribe to ActionEvents
	actionEventChannel := make(chan interface{})
	eventBus.Subscribe(events.ActionEventType, actionEventChannel)
	go func() {
		for event := range actionEventChannel {
			if actionEvent, ok := event.(*events.ActionEvent); ok {
				m.HandleActionEvent(actionEvent)
			} else {
				logrus.Errorf("ActionSignificanceMonitor: received unexpected event type on ActionEventType channel: %T", event)
			}
		}
	}()
	return m
}

// HandleActionEvent is the event handler for ActionEvents.
func (m *ActionSignificanceMonitor) HandleActionEvent(actionEvent *events.ActionEvent) {

	// Identify local observers (NPCs, Owners, Questmakers in the same room)
	// This is a simplified approach. A more robust solution would involve spatial indexing.
	// For now, we'll iterate through all NPCs, Owners, Questmakers and check their location.
	// This will be inefficient for large worlds, but sufficient for a prototype.

	// Get all NPCs
	npcs, err := m.npcDAL.GetAllNPCs()
	if err != nil {
		logrus.Errorf("ActionSignificanceMonitor: failed to get all NPCs: %v", err)
		return
	}

	// Get all Owners
	owners, err := m.ownerDAL.GetAllOwners()
	if err != nil {
		logrus.Errorf("ActionSignificanceMonitor: failed to get all Owners: %v", err)
		return
	}

	// Get all Questmakers
	questmakers, err := m.questmakerDAL.GetAllQuestmakers()
	if err != nil {
		logrus.Errorf("ActionSignificanceMonitor: failed to get all Questmakers: %v", err)
		return
	}

	// Combine all potential observers
	var observers []interface{}
	for _, npc := range npcs {
		if npc.CurrentRoomID == actionEvent.Room.ID { // Only consider NPCs in the same room
			observers = append(observers, npc)
		}
	}
	for _, owner := range owners {
		// Owners can monitor locations, races, professions.
		// For location-based owners, check if they monitor the event's room.
		if owner.MonitoredAspect == "location" && owner.AssociatedID == actionEvent.Room.ID {
			observers = append(observers, owner)
		} else if owner.MonitoredAspect == "race" && actionEvent.Player != nil && owner.AssociatedID == actionEvent.Player.RaceID {
			observers = append(observers, owner)
		} else if owner.MonitoredAspect == "profession" && actionEvent.Player != nil && actionEvent.Player.ProfessionID == owner.AssociatedID {
			observers = append(observers, owner)
		}
	}
	for _, questmaker := range questmakers {
		// Questmakers might observe actions relevant to their quests, regardless of location.
		// This requires more sophisticated logic to link actions to quests.
		// For now, we'll skip direct questmaker observation here and rely on Owners for quest-related triggers.
		// However, we still need to include them in the observer list for the perception filter to process them.
		observers = append(observers, questmaker)
	}

	for _, observer := range observers {
		perceivedAction, err := m.perceptionFilter.Filter(actionEvent, observer)
		if err != nil {
			logrus.Errorf("ActionSignificanceMonitor: failed to filter perception for observer %T: %v", observer, err)
			continue
		}

		// Calculate significance score
		// Final Score = (BaseScore + Σ AdditiveBonuses) * Multiplier * Clarity
		// For now, no additive bonuses or multipliers are implemented, so it's just BaseSignificance * Clarity
		significance := perceivedAction.BaseSignificance * perceivedAction.Clarity

		// Store the perceived action and its significance
		m.addPerceivedAction(actionEvent.Player.ID, getObserverID(observer), perceivedAction, significance)

		// Check and trigger reaction
		m.checkAndTriggerReaction(actionEvent.Player.ID, getObserverID(observer), observer)
	}
}

func (m *ActionSignificanceMonitor) addPerceivedAction(playerID string, observerID string, perceivedAction *perception.PerceivedAction, significance float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.playerEntityBuffers[playerID]; !ok {
		m.playerEntityBuffers[playerID] = make(map[string]*ActionBuffer)
	}

	if _, ok := m.playerEntityBuffers[playerID][observerID]; !ok {
		m.playerEntityBuffers[playerID][observerID] = &ActionBuffer{}
	}

	m.playerEntityBuffers[playerID][observerID].mu.Lock()
	defer m.playerEntityBuffers[playerID][observerID].mu.Unlock()

	m.playerEntityBuffers[playerID][observerID].records = append(m.playerEntityBuffers[playerID][observerID].records, perception.PerceivedActionRecord{
		PerceivedAction: perceivedAction,
		Significance:    significance,
		Timestamp:       time.Now(),
	})
}

func (m *ActionSignificanceMonitor) getCumulativeSignificance(playerID, observerID string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	playerBuffers, ok := m.playerEntityBuffers[playerID]
	if !ok {
		return 0.0
	}

	entityBuffer, ok := playerBuffers[observerID]
	if !ok {
		return 0.0
	}

	entityBuffer.mu.Lock()
	defer entityBuffer.mu.Unlock()

	totalSignificance := 0.0
	for _, record := range entityBuffer.records {
		totalSignificance += record.Significance
	}
	return totalSignificance
}

func (m *ActionSignificanceMonitor) clearPerceivedActions(playerID, observerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if playerBuffers, ok := m.playerEntityBuffers[playerID]; ok {
		if entityBuffer, ok := playerBuffers[observerID]; ok {
			entityBuffer.mu.Lock()
			entityBuffer.records = []perception.PerceivedActionRecord{} // Clear the slice
			entityBuffer.mu.Unlock()
		}
	}
}

func (m *ActionSignificanceMonitor) checkAndTriggerReaction(playerID, observerID string, observer interface{}) {
	cumulativeSignificance := m.getCumulativeSignificance(playerID, observerID)

	var reactionThreshold int
	var observerName string
	switch obs := observer.(type) {
	case *models.NPC:
		reactionThreshold = obs.ReactionThreshold
		observerName = obs.Name
	case *models.Owner:
		reactionThreshold = obs.ReactionThreshold
		observerName = obs.Name
	case *models.Questmaker:
		reactionThreshold = obs.ReactionThreshold
		observerName = obs.Name
	default:
		logrus.Errorf("ActionSignificanceMonitor: unsupported observer type for reaction check: %T", observer)
		return
	}

	if cumulativeSignificance >= float64(reactionThreshold) {
		logrus.Infof("ActionSignificanceMonitor: Triggering reaction for %s (ID: %s) to player %s. Cumulative Significance: %.2f, Threshold: %d",
			observerName, observerID, playerID, cumulativeSignificance, reactionThreshold)

		// Trigger reaction via SentientEntityManager
		err := m.sentientEntityManager.TriggerReaction(observer, m.GetBatchedPerceivedActions(playerID, observerID))
		if err != nil {
			logrus.Errorf("ActionSignificanceMonitor: failed to trigger reaction for %s (ID: %s): %v", observerName, observerID, err)
		}

		m.clearPerceivedActions(playerID, observerID) // Clear buffer after triggering
	}
}

// GetBatchedPerceivedActions retrieves and clears batched perceived actions for a specific player and entity.
func (m *ActionSignificanceMonitor) GetBatchedPerceivedActions(playerID, observerID string) []perception.PerceivedActionRecord {
	m.mu.Lock()
	defer m.mu.Unlock()

	var batchedRecords []perception.PerceivedActionRecord
	if playerBuffers, ok := m.playerEntityBuffers[playerID]; ok {
		if entityBuffer, ok := playerBuffers[observerID]; ok {
			entityBuffer.mu.Lock()
			batchedRecords = make([]perception.PerceivedActionRecord, len(entityBuffer.records))
			copy(batchedRecords, entityBuffer.records)
			entityBuffer.records = []perception.PerceivedActionRecord{} // Clear the slice after retrieving
			entityBuffer.mu.Unlock()
		}
	}
	return batchedRecords
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