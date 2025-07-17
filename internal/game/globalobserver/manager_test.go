package globalobserver

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mud/internal/game/events"
	"mud/internal/game/perception"
	"mud/internal/models"
)

// Mock implementations for testing


type MockPerceptionFilter struct {
	FilterFunc func(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error)
}

func (m *MockPerceptionFilter) Filter(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error) {
	return m.FilterFunc(event, observer)
}

type MockOwnerDAL struct {
	GetAllOwnersFunc func() ([]*models.Owner, error)
	GetOwnerByIDFunc func(id string) (*models.Owner, error)
	UpdateOwnerFunc  func(owner *models.Owner) error
}

func (m *MockOwnerDAL) GetAllOwners() ([]*models.Owner, error) {
	if m.GetAllOwnersFunc != nil {
		return m.GetAllOwnersFunc()
	}
	return nil, nil
}

func (m *MockOwnerDAL) GetOwnerByID(id string) (*models.Owner, error) {
	if m.GetOwnerByIDFunc != nil {
		return m.GetOwnerByIDFunc(id)
	}
	return nil, nil
}

func (m *MockOwnerDAL) UpdateOwner(owner *models.Owner) error {
	if m.UpdateOwnerFunc != nil {
		return m.UpdateOwnerFunc(owner)
	}
	return nil
}

type MockRaceDAL struct{}

func (m *MockRaceDAL) GetRaceByID(id string) (*models.Race, error) { return nil, nil }
func (m *MockRaceDAL) GetAllRaces() ([]*models.Race, error) { return nil, nil }
func (m *MockRaceDAL) CreateRace(race *models.Race) error { return nil }
func (m *MockRaceDAL) UpdateRace(race *models.Race) error { return nil }
func (m *MockRaceDAL) DeleteRace(id string) error { return nil }

type MockProfessionDAL struct{}

func (m *MockProfessionDAL) GetProfessionByID(id string) (*models.Profession, error) { return nil, nil }
func (m *MockProfessionDAL) GetAllProfessions() ([]*models.Profession, error) { return nil, nil }
func (m *MockProfessionDAL) CreateProfession(profession *models.Profession) error { return nil }
func (m *MockProfessionDAL) UpdateProfession(profession *models.Profession) error { return nil }
func (m *MockProfessionDAL) DeleteProfession(id string) error { return nil }

func TestNewGlobalObserverManager(t *testing.T) {
	eventBus := events.NewEventBus()
	mockPerceptionFilter := &MockPerceptionFilter{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockRaceDAL := &MockRaceDAL{}
	mockProfessionDAL := &MockProfessionDAL{}

	manager := NewGlobalObserverManager(
		eventBus,
		mockPerceptionFilter,
		mockOwnerDAL,
		mockRaceDAL,
		mockProfessionDAL,
	)

	assert.NotNil(t, manager)
	assert.Equal(t, eventBus, manager.eventBus)
	assert.Equal(t, mockPerceptionFilter, manager.perceptionFilter)
	assert.Equal(t, mockOwnerDAL, manager.ownerDAL)
	assert.Equal(t, mockRaceDAL, manager.raceDAL)
	assert.Equal(t, mockProfessionDAL, manager.professionDAL)

	// To verify subscription, we can publish an event and see if it's handled
	// This is implicitly tested by HandleActionEvent tests, but for NewManager, we just check initialization.
}

func TestGlobalObserverManager_HandleActionEvent_NoGlobalObservers(t *testing.T) {
	eventBus := events.NewEventBus()
	mockPerceptionFilter := &MockPerceptionFilter{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockRaceDAL := &MockRaceDAL{}
	mockProfessionDAL := &MockProfessionDAL{}

	// Mock GetAllOwners to return no global observers
	mockOwnerDAL.GetAllOwnersFunc = func() ([]*models.Owner, error) {
		return []*models.Owner{
			{ID: "owner1", MonitoredAspect: "location", AssociatedID: "room1"},
		}, nil
	}

	manager := NewGlobalObserverManager(
		eventBus,
		mockPerceptionFilter,
		mockOwnerDAL,
		mockRaceDAL,
		mockProfessionDAL,
	)

	// Ensure UpdateOwner is never called
	updateOwnerCalled := false
	mockOwnerDAL.UpdateOwnerFunc = func(owner *models.Owner) error {
		updateOwnerCalled = true
		return nil
	}

	actionEvent := &events.ActionEvent{
		ActionType: "test_action",
		Player:     &models.Player{ID: "player1", RaceID: "human", ProfessionID: "warrior"},
		Room:       &models.Room{ID: "room1"},
		Timestamp:  time.Now(),
	}

	manager.HandleActionEvent(actionEvent)

	// Give some time for goroutines to potentially run (though none should)
	time.Sleep(10 * time.Millisecond)

	assert.False(t, updateOwnerCalled, "UpdateOwner should not be called for non-global observers")
}

func TestGlobalObserverManager_HandleActionEvent_RaceBasedObserver(t *testing.T) {
	eventBus := events.NewEventBus()
	mockPerceptionFilter := &MockPerceptionFilter{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockRaceDAL := &MockRaceDAL{}
	mockProfessionDAL := &MockProfessionDAL{}

	owner := &models.Owner{
		ID:                     "race_owner",
		MonitoredAspect:        "race",
		AssociatedID:           "human",
		CurrentInfluenceBudget: 10.0,
		MaxInfluenceBudget:     100.0,
	}

	mockOwnerDAL.GetAllOwnersFunc = func() ([]*models.Owner, error) {
		return []*models.Owner{owner}, nil
	}

	mockPerceptionFilter.FilterFunc = func(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error) {
		return &perception.PerceivedAction{
			BaseSignificance: 5.0,
			Clarity:          1.0,
		}, nil
	}

	updateOwnerCalled := false
	mockOwnerDAL.UpdateOwnerFunc = func(updatedOwner *models.Owner) error {
		updateOwnerCalled = true
		assert.Equal(t, owner.ID, updatedOwner.ID)
		assert.InDelta(t, 15.0, updatedOwner.CurrentInfluenceBudget, 0.001) // 10 + (5 * 1)
		return nil
	}

	manager := NewGlobalObserverManager(
		eventBus,
		mockPerceptionFilter,
		mockOwnerDAL,
		mockRaceDAL,
		mockProfessionDAL,
	)

	actionEvent := &events.ActionEvent{
		ActionType: "test_action",
		Player:     &models.Player{ID: "player1", RaceID: "human", ProfessionID: "warrior"},
		Room:       &models.Room{ID: "room1"},
		Timestamp:  time.Now(),
	}

	manager.HandleActionEvent(actionEvent)

	// Give time for the goroutine to execute
	time.Sleep(10 * time.Millisecond)

	assert.True(t, updateOwnerCalled, "UpdateOwner should be called for race-based observer")
}

func TestGlobalObserverManager_HandleActionEvent_ProfessionBasedObserver(t *testing.T) {
	eventBus := events.NewEventBus()
	mockPerceptionFilter := &MockPerceptionFilter{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockRaceDAL := &MockRaceDAL{}
	mockProfessionDAL := &MockProfessionDAL{}

	owner := &models.Owner{
		ID:                     "prof_owner",
		MonitoredAspect:        "profession",
		AssociatedID:           "warrior",
		CurrentInfluenceBudget: 20.0,
		MaxInfluenceBudget:     100.0,
	}

	mockOwnerDAL.GetAllOwnersFunc = func() ([]*models.Owner, error) {
		return []*models.Owner{owner}, nil
	}

	mockPerceptionFilter.FilterFunc = func(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error) {
		return &perception.PerceivedAction{
			BaseSignificance: 8.0,
			Clarity:          0.5,
		}, nil
	}

	updateOwnerCalled := false
	mockOwnerDAL.UpdateOwnerFunc = func(updatedOwner *models.Owner) error {
		updateOwnerCalled = true
		assert.Equal(t, owner.ID, updatedOwner.ID)
		assert.InDelta(t, 24.0, updatedOwner.CurrentInfluenceBudget, 0.001) // 20 + (8 * 0.5)
		return nil
	}

	manager := NewGlobalObserverManager(
		eventBus,
		mockPerceptionFilter,
		mockOwnerDAL,
		mockRaceDAL,
		mockProfessionDAL,
	)

	actionEvent := &events.ActionEvent{
		ActionType: "test_action",
		Player:     &models.Player{ID: "player1", RaceID: "human", ProfessionID: "warrior"},
		Room:       &models.Room{ID: "room1"},
		Timestamp:  time.Now(),
	}

	manager.HandleActionEvent(actionEvent)

	// Give time for the goroutine to execute
	time.Sleep(10 * time.Millisecond)

	assert.True(t, updateOwnerCalled, "UpdateOwner should be called for profession-based observer")
}

func TestGlobalObserverManager_HandleActionEvent_OwnerNotObserving(t *testing.T) {
	eventBus := events.NewEventBus()
	mockPerceptionFilter := &MockPerceptionFilter{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockRaceDAL := &MockRaceDAL{}
	mockProfessionDAL := &MockProfessionDAL{}

	owner := &models.Owner{
		ID:                     "non_observing_owner",
		MonitoredAspect:        "race",
		AssociatedID:           "elf", // Player is human
		CurrentInfluenceBudget: 10.0,
		MaxInfluenceBudget:     100.0,
	}

	mockOwnerDAL.GetAllOwnersFunc = func() ([]*models.Owner, error) {
		return []*models.Owner{owner}, nil
	}

	updateOwnerCalled := false
	mockOwnerDAL.UpdateOwnerFunc = func(updatedOwner *models.Owner) error {
		updateOwnerCalled = true
		return nil
	}

	manager := NewGlobalObserverManager(
		eventBus,
		mockPerceptionFilter,
		mockOwnerDAL,
		mockRaceDAL,
		mockProfessionDAL,
	)

	actionEvent := &events.ActionEvent{
		ActionType: "test_action",
		Player:     &models.Player{ID: "player1", RaceID: "human", ProfessionID: "warrior"},
		Room:       &models.Room{ID: "room1"},
		Timestamp:  time.Now(),
	}

	manager.HandleActionEvent(actionEvent)

	// Give time for goroutines to potentially run (though none should for this owner)
	time.Sleep(10 * time.Millisecond)

	assert.False(t, updateOwnerCalled, "UpdateOwner should not be called for non-observing owner")
}

func TestGlobalObserverManager_HandleActionEvent_GetAllOwnersError(t *testing.T) {
	eventBus := events.NewEventBus()
	mockPerceptionFilter := &MockPerceptionFilter{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockRaceDAL := &MockRaceDAL{}
	mockProfessionDAL := &MockProfessionDAL{}

	mockOwnerDAL.GetAllOwnersFunc = func() ([]*models.Owner, error) {
		return nil, errors.New("failed to get owners")
	}

	manager := NewGlobalObserverManager(
		eventBus,
		mockPerceptionFilter,
		mockOwnerDAL,
		mockRaceDAL,
		mockProfessionDAL,
	)

	// Ensure UpdateOwner is never called
	updateOwnerCalled := false
	mockOwnerDAL.UpdateOwnerFunc = func(owner *models.Owner) error {
		updateOwnerCalled = true
		return nil
	}

	actionEvent := &events.ActionEvent{
		ActionType: "test_action",
		Player:     &models.Player{ID: "player1", RaceID: "human", ProfessionID: "warrior"},
		Room:       &models.Room{ID: "room1"},
		Timestamp:  time.Now(),
	}

	manager.HandleActionEvent(actionEvent)

	// Give some time for goroutines to potentially run (though none should)
	time.Sleep(10 * time.Millisecond)

	assert.False(t, updateOwnerCalled, "UpdateOwner should not be called if GetAllOwners fails")
}

func TestGlobalObserverManager_ProcessGlobalObservation_Success(t *testing.T) {
	mockPerceptionFilter := &MockPerceptionFilter{}
	mockOwnerDAL := &MockOwnerDAL{}

	owner := &models.Owner{
		ID:                     "test_owner",
		CurrentInfluenceBudget: 50.0,
		MaxInfluenceBudget:     100.0,
	}

	mockPerceptionFilter.FilterFunc = func(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error) {
		return &perception.PerceivedAction{
			BaseSignificance: 10.0,
			Clarity:          0.8,
		}, nil
	}

	updateOwnerCalled := false
	mockOwnerDAL.UpdateOwnerFunc = func(updatedOwner *models.Owner) error {
		updateOwnerCalled = true
		assert.Equal(t, owner.ID, updatedOwner.ID)
		assert.InDelta(t, 58.0, updatedOwner.CurrentInfluenceBudget, 0.001) // 50 + (10 * 0.8)
		return nil
	}

	manager := &GlobalObserverManager{
		perceptionFilter: mockPerceptionFilter,
		ownerDAL:         mockOwnerDAL,
	}

	actionEvent := &events.ActionEvent{
		ActionType: "test_action",
		Player:     &models.Player{ID: "player1"},
		Room:       &models.Room{ID: "room1"},
		Timestamp:  time.Now(),
	}

	manager.processGlobalObservation(actionEvent, owner)

	assert.True(t, updateOwnerCalled, "UpdateOwner should be called on successful observation")
}

func TestGlobalObserverManager_ProcessGlobalObservation_BudgetCapping(t *testing.T) {
	mockPerceptionFilter := &MockPerceptionFilter{}
	mockOwnerDAL := &MockOwnerDAL{}

	owner := &models.Owner{
		ID:                     "test_owner_cap",
		CurrentInfluenceBudget: 95.0,
		MaxInfluenceBudget:     100.0,
	}

	mockPerceptionFilter.FilterFunc = func(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error) {
		return &perception.PerceivedAction{
			BaseSignificance: 10.0,
			Clarity:          1.0,
		}, nil
	}

	updateOwnerCalled := false
	mockOwnerDAL.UpdateOwnerFunc = func(updatedOwner *models.Owner) error {
		updateOwnerCalled = true
		assert.Equal(t, owner.ID, updatedOwner.ID)
		assert.InDelta(t, 100.0, updatedOwner.CurrentInfluenceBudget, 0.001) // Should be capped at 100
		return nil
	}

	manager := &GlobalObserverManager{
		perceptionFilter: mockPerceptionFilter,
		ownerDAL:         mockOwnerDAL,
	}

	actionEvent := &events.ActionEvent{
		ActionType: "test_action",
		Player:     &models.Player{ID: "player1"},
		Room:       &models.Room{ID: "room1"},
		Timestamp:  time.Now(),
	}

	manager.processGlobalObservation(actionEvent, owner)

	assert.True(t, updateOwnerCalled, "UpdateOwner should be called for budget capping")
}

func TestGlobalObserverManager_ProcessGlobalObservation_FilterError(t *testing.T) {
	mockPerceptionFilter := &MockPerceptionFilter{}
	mockOwnerDAL := &MockOwnerDAL{}

	owner := &models.Owner{
		ID:                     "test_owner_filter_error",
		CurrentInfluenceBudget: 50.0,
		MaxInfluenceBudget:     100.0,
	}

	mockPerceptionFilter.FilterFunc = func(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error) {
		return nil, errors.New("filter error")
	}

	updateOwnerCalled := false
	mockOwnerDAL.UpdateOwnerFunc = func(updatedOwner *models.Owner) error {
		updateOwnerCalled = true
		return nil
	}

	manager := &GlobalObserverManager{
		perceptionFilter: mockPerceptionFilter,
		ownerDAL:         mockOwnerDAL,
	}

	actionEvent := &events.ActionEvent{
		ActionType: "test_action",
		Player:     &models.Player{ID: "player1"},
		Room:       &models.Room{ID: "room1"},
		Timestamp:  time.Now(),
	}

	manager.processGlobalObservation(actionEvent, owner)

	assert.False(t, updateOwnerCalled, "UpdateOwner should not be called if filter fails")
}

func TestGlobalObserverManager_ProcessGlobalObservation_UpdateOwnerError(t *testing.T) {
	mockPerceptionFilter := &MockPerceptionFilter{}
	mockOwnerDAL := &MockOwnerDAL{}

	owner := &models.Owner{
		ID:                     "test_owner_update_error",
		CurrentInfluenceBudget: 50.0,
		MaxInfluenceBudget:     100.0,
	}

	mockPerceptionFilter.FilterFunc = func(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error) {
		return &perception.PerceivedAction{
			BaseSignificance: 10.0,
			Clarity:          1.0,
		}, nil
	}

	mockOwnerDAL.UpdateOwnerFunc = func(updatedOwner *models.Owner) error {
		return errors.New("update owner error")
	}

	manager := &GlobalObserverManager{
		perceptionFilter: mockPerceptionFilter,
		ownerDAL:         mockOwnerDAL,
	}

	actionEvent := &events.ActionEvent{
		ActionType: "test_action",
		Player:     &models.Player{ID: "player1"},
		Room:       &models.Room{ID: "room1"},
		Timestamp:  time.Now(),
	}

	manager.processGlobalObservation(actionEvent, owner)

	// No direct assertion on error logging, but ensure no panic
}
