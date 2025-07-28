package actionsignificance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mud/internal/dal"
	"mud/internal/game/events"
	"mud/internal/game/perception"
	"mud/internal/models"
)

// Mock implementations for testing

type MockNPCDAL struct {
	npcs map[string]*models.NPC
	GetAllNPCsFunc func() ([]*models.NPC, error)
	GetNPCByIDFunc func(id string) (*models.NPC, error)
	GetNPCsByOwnerFunc func(ownerID string) ([]*models.NPC, error)
	UpdateNPCFunc func(npc *models.NPC) error
}

func (m *MockNPCDAL) GetAllNPCs() ([]*models.NPC, error) {
	if m.GetAllNPCsFunc != nil {
		return m.GetAllNPCsFunc()
	}
	var allNPCs []*models.NPC
	for _, npc := range m.npcs {
		allNPCs = append(allNPCs, npc)
	}
	return allNPCs, nil
}

func (m *MockNPCDAL) GetNPCByID(id string) (*models.NPC, error) {
	if m.GetNPCByIDFunc != nil {
		return m.GetNPCByIDFunc(id)
	}
	return m.npcs[id], nil
}

func (m *MockNPCDAL) GetNPCsByOwner(ownerID string) ([]*models.NPC, error) {
	if m.GetNPCsByOwnerFunc != nil {
		return m.GetNPCsByOwnerFunc(ownerID)
	}
	var ownedNPCs []*models.NPC
	for _, npc := range m.npcs {
		for _, oid := range npc.OwnerIDs { // Corrected to iterate through OwnerIDs
			if oid == ownerID {
				ownedNPCs = append(ownedNPCs, npc)
				break // Found, move to next NPC
			}
		}
	}
	return ownedNPCs, nil
}

func (m *MockNPCDAL) UpdateNPC(npc *models.NPC) error {
	if m.UpdateNPCFunc != nil {
		return m.UpdateNPCFunc(npc)
	}
	m.npcs[npc.ID] = npc
	return nil
}

func (m *MockNPCDAL) CreateNPC(npc *models.NPC) error { return nil }
func (m *MockNPCDAL) DeleteNPC(id string) error { return nil }
func (m *MockNPCDAL) GetNPCsByRoom(roomID string) ([]*models.NPC, error) { return nil, nil }
func (m *MockNPCDAL) Cache() dal.CacheInterface { return nil }

type MockOwnerDAL struct {
	owners map[string]*models.Owner
	GetAllOwnersFunc func() ([]*models.Owner, error)
	GetOwnerByIDFunc func(id string) (*models.Owner, error)
	UpdateOwnerFunc func(owner *models.Owner) error
}

func (m *MockOwnerDAL) GetAllOwners() ([]*models.Owner, error) {
	if m.GetAllOwnersFunc != nil {
		return m.GetAllOwnersFunc()
	}
	var allOwners []*models.Owner
	for _, owner := range m.owners {
		allOwners = append(allOwners, owner)
	}
	return allOwners, nil
}

func (m *MockOwnerDAL) GetOwnerByID(id string) (*models.Owner, error) {
	if m.GetOwnerByIDFunc != nil {
		return m.GetOwnerByIDFunc(id)
	}
	return m.owners[id], nil
}

func (m *MockOwnerDAL) UpdateOwner(owner *models.Owner) error {
	if m.UpdateOwnerFunc != nil {
		return m.UpdateOwnerFunc(owner)
	}
	m.owners[owner.ID] = owner
	return nil
}

func (m *MockOwnerDAL) CreateOwner(owner *models.Owner) error { return nil }
func (m *MockOwnerDAL) DeleteOwner(id string) error { return nil }
func (m *MockOwnerDAL) Cache() dal.CacheInterface { return nil }

type MockQuestmakerDAL struct {
	questmakers map[string]*models.Questmaker
	GetAllQuestmakersFunc func() ([]*models.Questmaker, error)
	GetQuestmakerByIDFunc func(id string) (*models.Questmaker, error)
}

func (m *MockQuestmakerDAL) GetAllQuestmakers() ([]*models.Questmaker, error) {
	if m.GetAllQuestmakersFunc != nil {
		return m.GetAllQuestmakersFunc()
	}
	var allQuestmakers []*models.Questmaker
	for _, qm := range m.questmakers {
		allQuestmakers = append(allQuestmakers, qm)
	}
	return allQuestmakers, nil
}

func (m *MockQuestmakerDAL) GetQuestmakerByID(id string) (*models.Questmaker, error) {
	if m.GetQuestmakerByIDFunc != nil {
		return m.GetQuestmakerByIDFunc(id)
	}
	return m.questmakers[id], nil
}

func (m *MockQuestmakerDAL) CreateQuestmaker(questmaker *models.Questmaker) error { return nil }
func (m *MockQuestmakerDAL) UpdateQuestmaker(questmaker *models.Questmaker) error { return nil }
func (m *MockQuestmakerDAL) DeleteQuestmaker(id string) error { return nil }
func (m *MockQuestmakerDAL) Cache() dal.CacheInterface { return nil }

type MockPerceptionFilter struct {
	FilterFunc func(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error)
}

func (m *MockPerceptionFilter) Filter(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error) {
	return m.FilterFunc(event, observer)
}

type MockSentientEntityManager struct {
	TriggerReactionFunc func(observer interface{}, perceivedActions []perception.PerceivedActionRecord) error
}

func (m *MockSentientEntityManager) TriggerReaction(observer interface{}, perceivedActions []perception.PerceivedActionRecord) error {
	return m.TriggerReactionFunc(observer, perceivedActions)
}

// MockDAL for ActionSignificanceMonitor
type MockDAL struct {
	NPCDAL        dal.NPCDALInterface
	OwnerDAL      dal.OwnerDALInterface
	QuestmakerDAL dal.QuestmakerDALInterface
}

func NewMockDAL(npcDAL dal.NPCDALInterface, ownerDAL dal.OwnerDALInterface, questmakerDAL dal.QuestmakerDALInterface) *MockDAL {
	return &MockDAL{
		NPCDAL:        npcDAL,
		OwnerDAL:      ownerDAL,
		QuestmakerDAL: questmakerDAL,
	}
}

func TestActionSignificanceMonitor_HandleActionEvent(t *testing.T) {
	// Setup mock data
	mockNPCs := map[string]*models.NPC{
		"npc1": {
			ID:                "npc1",
			CurrentRoomID:     "room1",
			ReactionThreshold: 10, // Changed to 10 to test cumulative significance
			Name:              "Test NPC 1",
		},
		"npc2": {
			ID:                "npc2",
			CurrentRoomID:     "room2",
			ReactionThreshold: 10,
			Name:              "Test NPC 2",
		},
	}
	mockOwners := map[string]*models.Owner{
		"owner1": {
			ID:                "owner1",
			MonitoredAspect:   "location",
			AssociatedID:      "room1",
			ReactionThreshold: 5, // Changed to 5 to trigger reaction with single event
			Name:              "Test Owner 1",
		},
		"owner_race": {
			ID:                "owner_race",
			MonitoredAspect:   "race",
			AssociatedID:      "human",
			ReactionThreshold: 12,
			Name:              "Race Owner",
		},
	}
	mockQuestmakers := map[string]*models.Questmaker{
		"questmaker1": {
			ID:                "questmaker1",
			ReactionThreshold: 3,
			Name:              "Test Questmaker 1",
		},
	}

	mockNPCDAL := &MockNPCDAL{npcs: mockNPCs}
	mockOwnerDAL := &MockOwnerDAL{owners: mockOwners}
	mockQuestmakerDAL := &MockQuestmakerDAL{questmakers: mockQuestmakers}
	
	// Mock PerceptionFilter to return a controlled PerceivedAction
	mockPerceivedAction := &perception.PerceivedAction{
		PerceivedActionType: "test_action",
		Clarity:             1.0,
		BaseSignificance:    6.0, // Will trigger reaction for npc1 (threshold 5)
	}
	mockPerceptionFilter := &MockPerceptionFilter{
		FilterFunc: func(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error) {
			return mockPerceivedAction, nil
		},
	}

	// Mock SentientEntityManager to capture triggered reactions
	triggeredReactions := make(map[string][]perception.PerceivedActionRecord)
	mockSentientEntityManager := &MockSentientEntityManager{
		TriggerReactionFunc: func(observer interface{}, perceivedActions []perception.PerceivedActionRecord) error {
			observerID := ""
			switch o := observer.(type) {
			case *models.NPC:
				observerID = o.ID
			case *models.Owner:
				observerID = o.ID
			case *models.Questmaker:
				observerID = o.ID
			}
			triggeredReactions[observerID] = perceivedActions
			return nil
		},
	}

	// Setup EventBus and Monitor
	eventBus := events.NewEventBus()
	monitor := NewMonitor(
		eventBus,
		mockPerceptionFilter,
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockSentientEntityManager,
	)

	// Create a player for the action event
	player := &models.PlayerCharacter{
		ID:           "player1",
		Name:         "TestPlayer",
		RaceID:       "human",
		ProfessionID: "warrior",
	}

	// Create an action event in room1
	actionEvent := &events.ActionEvent{
		ActionType: "test_action",
		Player:     player,
		Room:       &models.Room{ID: "room1"},
		Timestamp:  time.Now(),
	}

	// Handle the first action event directly
	monitor.HandleActionEvent(actionEvent)

	// Assertions after the first event
	assert.NotContains(t, triggeredReactions, "npc1", "Reaction should NOT be triggered for npc1 after first event")
	assert.Contains(t, triggeredReactions, "owner1", "Reaction should be triggered for owner1")
	assert.Contains(t, triggeredReactions, "questmaker1", "Reaction should be triggered for questmaker1")

	// Verify the content of triggered reactions for owner1 and questmaker1
	if records, ok := triggeredReactions["owner1"]; ok {
		assert.Len(t, records, 1)
		assert.Equal(t, mockPerceivedAction.PerceivedActionType, records[0].PerceivedAction.PerceivedActionType)
		assert.InDelta(t, mockPerceivedAction.BaseSignificance*mockPerceivedAction.Clarity, records[0].Significance, 0.001)
	}
	if records, ok := triggeredReactions["questmaker1"]; ok {
		assert.Len(t, records, 1)
		assert.Equal(t, mockPerceivedAction.PerceivedActionType, records[0].PerceivedAction.PerceivedActionType)
		assert.InDelta(t, mockPerceivedAction.BaseSignificance*mockPerceivedAction.Clarity, records[0].Significance, 0.001)
	}

	// Clear triggered reactions for the next part of the test
	triggeredReactions = make(map[string][]perception.PerceivedActionRecord)

	// Handle another event directly to test cumulative significance and clearing
	monitor.HandleActionEvent(actionEvent)

	// Assertions after the second event
	assert.Contains(t, triggeredReactions, "npc1", "Reaction should be triggered for npc1 after second event")
	assert.Contains(t, triggeredReactions, "owner1", "Reaction should be triggered for owner1 after second event") // owner1 triggers again as its buffer was cleared
	assert.Contains(t, triggeredReactions, "owner_race", "Reaction should be triggered for owner_race after second event")
	assert.Contains(t, triggeredReactions, "questmaker1", "Reaction should be triggered for questmaker1 after second event") // questmaker1 triggers again

	// Verify the content of triggered reactions for npc1
	if records, ok := triggeredReactions["npc1"]; ok {
		assert.Len(t, records, 2) // Two events accumulated
		assert.Equal(t, mockPerceivedAction.PerceivedActionType, records[0].PerceivedAction.PerceivedActionType)
		assert.InDelta(t, mockPerceivedAction.BaseSignificance*mockPerceivedAction.Clarity, records[0].Significance, 0.001)
		assert.Equal(t, mockPerceivedAction.PerceivedActionType, records[1].PerceivedAction.PerceivedActionType)
		assert.InDelta(t, mockPerceivedAction.BaseSignificance*mockPerceivedAction.Clarity, records[1].Significance, 0.001)

		// Sum of significances should be 12.0
		cumulativeSig := records[0].Significance + records[1].Significance
		assert.InDelta(t, 12.0, cumulativeSig, 0.001, "Cumulative significance for npc1 should be 12.0")
	}

	// Verify the content of triggered reactions for owner1 after second event
	if records, ok := triggeredReactions["owner1"]; ok {
		assert.Len(t, records, 1) // Triggered again, new record
		assert.InDelta(t, mockPerceivedAction.BaseSignificance*mockPerceivedAction.Clarity, records[0].Significance, 0.001)
	}

	// Verify the content of triggered reactions for owner_race after second event
	if records, ok := triggeredReactions["owner_race"]; ok {
		assert.Len(t, records, 2) // Two events accumulated
		assert.InDelta(t, mockPerceivedAction.BaseSignificance*mockPerceivedAction.Clarity, records[0].Significance, 0.001)
		assert.InDelta(t, mockPerceivedAction.BaseSignificance*mockPerceivedAction.Clarity, records[1].Significance, 0.001)
		cumulativeSig := records[0].Significance + records[1].Significance
		assert.InDelta(t, 12.0, cumulativeSig, 0.001, "Cumulative significance for owner_race should be 12.0")
	}

	// Verify the content of triggered reactions for questmaker1 after second event
	if records, ok := triggeredReactions["questmaker1"]; ok {
		assert.Len(t, records, 1) // Triggered again, new record
		assert.InDelta(t, mockPerceivedAction.BaseSignificance*mockPerceivedAction.Clarity, records[0].Significance, 0.001)
	}

	// npc1's cumulative significance should now be 0.0 after being triggered and cleared
	assert.InDelta(t, 0.0, monitor.getCumulativeSignificance("player1", "npc1"), 0.001, "Buffer should be cleared after reaction")

	// Triggering reaction again should clear the buffer
	// We need to manually call checkAndTriggerReaction as HandleActionEvent doesn't guarantee immediate trigger
	monitor.checkAndTriggerReaction("player1", "npc1", mockNPCs["npc1"])
	assert.InDelta(t, 0.0, monitor.getCumulativeSignificance("player1", "npc1"), 0.001, "Buffer should be cleared after reaction")
}