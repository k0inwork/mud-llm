package perception

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mud/internal/dal"
	"mud/internal/game/events"
	"mud/internal/models"
	"mud/internal/testutils"
)

// MockRoomDAL implements dal.RoomDALInterface for testing.
type MockRoomDAL struct {
	cache dal.CacheInterface
}

func (m *MockRoomDAL) GetRoomByID(id string) (*models.Room, error) {
	if val, ok := m.cache.Get(id); ok {
		if room, isRoom := val.(*models.Room); isRoom {
			return room, nil
		}
	}
	return nil, nil
}

func (m *MockRoomDAL) GetAllRooms() ([]*models.Room, error) { return nil, nil }
func (m *MockRoomDAL) CreateRoom(room *models.Room) error { return nil }
func (m *MockRoomDAL) UpdateRoom(room *models.Room) error { return nil }
func (m *MockRoomDAL) DeleteRoom(id string) error { return nil }
func (m *MockRoomDAL) Cache() dal.CacheInterface { return m.cache }

// MockRaceDAL implements dal.RaceDALInterface for testing.
type MockRaceDAL struct {
	cache dal.CacheInterface
}

func (m *MockRaceDAL) GetRaceByID(id string) (*models.Race, error) {
	if val, ok := m.cache.Get(id); ok {
		if race, isRace := val.(*models.Race); isRace {
			return race, nil
		}
	}
	return nil, nil
}

func (m *MockRaceDAL) GetAllRaces() ([]*models.Race, error) { return nil, nil }
func (m *MockRaceDAL) CreateRace(race *models.Race) error { return nil }
func (m *MockRaceDAL) UpdateRace(race *models.Race) error { return nil }
func (m *MockRaceDAL) DeleteRace(id string) error { return nil }
func (m *MockRaceDAL) Cache() dal.CacheInterface { return m.cache }

// MockProfessionDAL implements dal.ProfessionDALInterface for testing.
type MockProfessionDAL struct {
	cache dal.CacheInterface
}

func (m *MockProfessionDAL) GetProfessionByID(id string) (*models.Profession, error) {
	if val, ok := m.cache.Get(id); ok {
		if profession, isProfession := val.(*models.Profession); isProfession {
			return profession, nil
		}
	}
	return nil, nil
}

func (m *MockProfessionDAL) GetAllProfessions() ([]*models.Profession, error) { return nil, nil }
func (m *MockProfessionDAL) CreateProfession(profession *models.Profession) error { return nil }
func (m *MockProfessionDAL) UpdateProfession(profession *models.Profession) error { return nil }
func (m *MockProfessionDAL) DeleteProfession(id string) error { return nil }
func (m *MockProfessionDAL) Cache() dal.CacheInterface { return m.cache }

func TestPerceptionFilter_Filter(t *testing.T) {
	mockRooms := map[string]*models.Room{
		"room_shire": {
			ID: "room_shire",
			PerceptionBiases: map[string]float64{
				"magic_action": -0.1, // Changed from "magic"
				"say":          0.1,
				"unknown_action": -0.1, // Added for low clarity test
			},
		},
		"room_bree": {
			ID: "room_bree",
			PerceptionBiases: map[string]float64{
				"subterfuge_action": 0.2, // Changed from "subterfuge"
			},
		},
	}
	mockRaces := map[string]*models.Race{
		"hobbit": {
			ID: "hobbit",
			PerceptionBiases: map[string]float64{
				"magic_action": -0.2, // Changed from "magic"
				"pray":         0.1,
				"subterfuge":   0.2, // Added for skill category test
			},
		},
		"elf": {
			ID: "elf",
			PerceptionBiases: map[string]float64{
				"magic_action": 0.3, // Changed from "magic"
			},
		},
		"human": { // Added human race to mock data
			ID: "human",
			PerceptionBiases: map[string]float64{},
		},
	}
	mockProfessions := map[string]*models.Profession{
		"mage": {
			ID: "mage",
			PerceptionBiases: map[string]float64{
				"magic_action": 0.4, // Changed from "magic"
			},
		},
		"rogue": {
			ID: "rogue",
			PerceptionBiases: map[string]float64{
				"subterfuge_action": 0.3, // Changed from "subterfuge"
			},
		},
		"commoner": { // Added commoner profession to mock data
			ID: "commoner",
			PerceptionBiases: map[string]float64{},
		},
	}

	// Create mock caches for each DAL
	mockRoomCache := testutils.NewMockCache()
	for k, v := range mockRooms {
		mockRoomCache.Data[k] = v
	}

	mockRaceCache := testutils.NewMockCache()
	for k, v := range mockRaces {
		mockRaceCache.Data[k] = v
	}

	mockProfessionCache := testutils.NewMockCache()
	for k, v := range mockProfessions {
		mockProfessionCache.Data[k] = v
	}

	// Create mock DAL instances using the mock caches
	mockRoomDAL := &MockRoomDAL{cache: mockRoomCache}
	mockRaceDAL := &MockRaceDAL{cache: mockRaceCache}
	mockProfessionDAL := &MockProfessionDAL{cache: mockProfessionCache}

	pf := NewPerceptionFilter(mockRoomDAL, mockRaceDAL, mockProfessionDAL)

	player := &models.PlayerCharacter{
		ID:           "player1",
		Name:         "TestPlayer",
		RaceID:       "human", // Default to human for most tests
		ProfessionID: "adventurer",
	}

	tests := []struct {
		name          string
		actionEvent   *events.ActionEvent
		observer      interface{}
		expectedBaseSig float64
		expectedClarity float64
	expectedPerceivedActionType string
	}{
		{
			name: "NPC observes 'say' in same room (no specific biases)",
			actionEvent: &events.ActionEvent{
				ActionType: "say",
				Player:     player,
				Room:       &models.Room{ID: "room_bree"},
				Timestamp:  time.Now(),
			},
			observer: &models.NPC{
				ID:            "npc1",
				CurrentRoomID: "room_bree",
				RaceID:        "human",
				ProfessionID:  "commoner",
			},
			expectedBaseSig: 10.0, // Base for NPC observing 'say'
			expectedClarity: 0.8,  // 1.0 (initial) - 0.2 (dazzled debuff from npc1 in getSkillsAndBuffs)
			expectedPerceivedActionType: "say_general",
		},
		{
			name: "NPC observes 'magic_action' with racial and room bias",
			actionEvent: &events.ActionEvent{
				ActionType: "magic_action",
				Player:     player,
				Room:       &models.Room{ID: "room_shire"},
				Timestamp:  time.Now(),
			},
			observer: &models.NPC{
				ID:            "hobbit_npc",
				CurrentRoomID: "room_shire",
				RaceID:        "hobbit",
				ProfessionID:  "commoner",
			},
			expectedBaseSig: 7.0, // Base for NPC observing 'magic_action'
			expectedClarity: 0.7, // Capped: 1.0 - 0.2 (hobbit magic_action) - 0.1 (shire magic_action) = 0.7
			expectedPerceivedActionType: "magic_action_general", // Corrected expected type
		},
		{
			name: "Owner (location) observes 'subterfuge_action' with room bias",
			actionEvent: &events.ActionEvent{
				ActionType: "subterfuge_action",
				Player:     player,
				Room:       &models.Room{ID: "room_bree"},
				Timestamp:  time.Now(),
			},
			observer: &models.Owner{
				ID:              "bree_owner",
				MonitoredAspect: "location",
				AssociatedID:    "room_bree",
			},
			expectedBaseSig: 6.0, // Base for Owner observing 'subterfuge_action'
			expectedClarity: 1.0, // Capped: 1.0 + 0.2 (bree subterfuge_action) = 1.2 -> 1.0
			expectedPerceivedActionType: "subterfuge_action",
		},
		{
			name: "Owner (race) observes 'magic_action' with racial bias",
			actionEvent: &events.ActionEvent{
				ActionType: "magic_action",
				Player:     player,
				Room:       &models.Room{ID: "some_room"}, // Room doesn't matter for race-based owner
				Timestamp:  time.Now(),
			},
			observer: &models.Owner{
				ID:              "elf_owner",
				MonitoredAspect: "race",
				AssociatedID:    "elf",
			},
			expectedBaseSig: 7.0, // Base for Owner observing 'magic_action'
			expectedClarity: 1.0, // Capped: 1.0 + 0.3 (elf magic_action) = 1.3 -> 1.0
			expectedPerceivedActionType: "magic_action",
		},
		{
			name: "Owner (profession) observes 'subterfuge_action' with profession bias",
			actionEvent: &events.ActionEvent{
				ActionType: "subterfuge_action",
				Player:     player,
				Room:       &models.Room{ID: "some_room"}, // Room doesn't matter for profession-based owner
				Timestamp:  time.Now(),
			},
			observer: &models.Owner{
				ID:              "rogue_owner",
				MonitoredAspect: "profession",
				AssociatedID:    "rogue",
			},
			expectedBaseSig: 6.0, // Base for Owner observing 'subterfuge_action'
			expectedClarity: 1.0, // Capped: 1.0 + 0.3 (rogue subterfuge_action) = 1.3 -> 1.0
			expectedPerceivedActionType: "subterfuge_action",
		},
		{
			name: "Questmaker observes 'attack' (neutral bias)",
			actionEvent: &events.ActionEvent{
				ActionType: "attack",
				Player:     player,
				Room:       &models.Room{ID: "any_room"},
				Timestamp:  time.Now(),
			},
			observer: &models.Questmaker{
				ID: "questmaker1",
			},
			expectedBaseSig: 10.0, // Base for Questmaker observing 'attack'
			expectedClarity: 1.0,  // Neutral bias
			expectedPerceivedActionType: "attack",
		},
		{
			name: "Low clarity action (unclear_action)",
			actionEvent: &events.ActionEvent{
				ActionType: "unknown_action", // Action not in baseActionSignificance
				Player:     player,
				Room:       &models.Room{ID: "room_shire"},
				Timestamp:  time.Now(),
			},
			observer: &models.NPC{
				ID:            "npc_low_clarity",
				CurrentRoomID: "room_shire",
				RaceID:        "human",
				ProfessionID:  "commoner",
			},
			expectedBaseSig: 0.5, // Default low significance
			expectedClarity: 0.9, // Capped: 1.0 - 0.1 (room_shire unknown_action bias) = 0.9
			expectedPerceivedActionType: "unclear_action", // Corrected expected type
		},
		{
			name: "Action with skill, high clarity",
			actionEvent: &events.ActionEvent{
				ActionType: "use_skill",
				SkillUsed: &models.Skill{
					ID:       "fireball_skill",
					Name:     "Fireball",
					Category: "magic",
				},
				Player:    player,
				Room:      &models.Room{ID: "room_bree"},
				Timestamp: time.Now(),
			},
			observer: &models.NPC{
				ID:            "mage_npc",
				CurrentRoomID: "room_bree",
				RaceID:        "elf",
				ProfessionID:  "mage",
			},
			expectedBaseSig: 5.0, // Base for NPC observing 'use_skill'
			expectedClarity: 1.0, // Capped: 1.0 + 0.3 (elf magic_action) + 0.4 (mage magic_action) = 1.7 -> 1.0
			expectedPerceivedActionType: "Fireball",
		},
		{
			name: "Action with skill, moderate clarity",
			actionEvent: &events.ActionEvent{
				ActionType: "use_skill",
				SkillUsed: &models.Skill{
					ID:       "stealth_skill",
					Name:     "Stealth",
					Category: "subterfuge",
				},
				Player:    player,
				Room:      &models.Room{ID: "room_shire"},
				Timestamp: time.Now(),
			},
			observer: &models.NPC{
				ID:            "hobbit_npc_2",
				CurrentRoomID: "room_shire",
				RaceID:        "hobbit",
				ProfessionID:  "commoner",
			},
			expectedBaseSig: 5.0, // Base for NPC observing 'use_skill'
			expectedClarity: 1.0, // Capped: 1.0 + 0.1 (shire say) + 0.2 (hobbit subterfuge) = 1.3 -> 1.0
			expectedPerceivedActionType: "Stealth", // Corrected expected type
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perceivedAction, err := pf.Filter(tt.actionEvent, tt.observer)
			assert.NoError(t, err)
			assert.NotNil(t, perceivedAction)

			assert.InDelta(t, tt.expectedBaseSig, perceivedAction.BaseSignificance, 0.001, "BaseSignificance mismatch")
			assert.InDelta(t, tt.expectedClarity, perceivedAction.Clarity, 0.001, "Clarity mismatch")
			assert.Equal(t, tt.expectedPerceivedActionType, perceivedAction.PerceivedActionType, "PerceivedActionType mismatch")
		})
	}
}

func TestPerceptionFilter_determinePerceivedActionType(t *testing.T) {
	// Create a mock DAL with empty data for this test, as it doesn't rely on specific DAL data
	mockRoomDAL := &MockRoomDAL{cache: testutils.NewMockCache()}
	mockRaceDAL := &MockRaceDAL{cache: testutils.NewMockCache()}
	mockProfessionDAL := &MockProfessionDAL{cache: testutils.NewMockCache()}
	pf := NewPerceptionFilter(mockRoomDAL, mockRaceDAL, mockProfessionDAL)

	player := &models.PlayerCharacter{ID: "p1", Name: "Player1"}
	skill := &models.Skill{ID: "s1", Name: "Sneak", Category: "subterfuge"}

	tests := []struct {
		name              string
		actionEvent       *events.ActionEvent
		clarity           float64
		expectedPerceivedActionType string
	}{
		{
			name: "High clarity, no skill",
			actionEvent: &events.ActionEvent{
				ActionType: "say",
				Player:     player,
			},
			clarity: 0.95,
			expectedPerceivedActionType: "say",
		},
		{
			name: "High clarity, with skill",
			actionEvent: &events.ActionEvent{
				ActionType: "use_skill",
				SkillUsed:  skill,
				Player:     player,
			},
			clarity: 0.95,
			expectedPerceivedActionType: "Sneak",
		},
		{
			name: "Moderate clarity, no skill",
			actionEvent: &events.ActionEvent{
				ActionType: "say",
				Player:     player,
			},
			clarity: 0.6,
			expectedPerceivedActionType: "say_general",
		},
		{
			name: "Moderate clarity, with skill",
			actionEvent: &events.ActionEvent{
				ActionType: "use_skill",
				SkillUsed:  skill,
				Player:     player,
			},
			clarity: 0.6,
			expectedPerceivedActionType: "subterfuge_action",
		},
		{
			name: "Low clarity, no skill",
			actionEvent: &events.ActionEvent{
				ActionType: "walk",
				Player:     player,
			},
			clarity: 0.3,
			expectedPerceivedActionType: "unclear_action",
		},
		{
			name: "Low clarity, with skill (category known)",
			actionEvent: &events.ActionEvent{
				ActionType: "use_skill",
				SkillUsed:  skill,
				Player:     player,
			},
			clarity: 0.3,
			expectedPerceivedActionType: "strange_subterfuge",
		},
		{
			name: "Low clarity, with skill (no category)",
			actionEvent: &events.ActionEvent{
				ActionType: "use_skill",
				SkillUsed:  &models.Skill{ID: "s2", Name: "Unknown Skill"},
				Player:     player,
			},
			clarity: 0.3,
			expectedPerceivedActionType: "unclear_action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pf.determinePerceivedActionType(tt.actionEvent, tt.clarity)
			assert.Equal(t, tt.expectedPerceivedActionType, result)
		})
	}
}
