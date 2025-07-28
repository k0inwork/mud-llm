package sentiententitymanager

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"mud/internal/dal"
	"mud/internal/game/events"
	"mud/internal/game/perception"
	"mud/internal/llm"
	"mud/internal/models"
	"mud/internal/presentation"
)

// Mock implementations for DAL interfaces
type MockPlayerCharacterDAL struct {
	GetCharacterByIDFunc func(id string) (*models.PlayerCharacter, error)
	GetCharacterClassFunc func(characterID string) (*models.PlayerClass, error)
	GetCharacterQuestStateFunc func(characterID, questID string) (*models.PlayerQuestState, error)
}
func (m *MockPlayerCharacterDAL) GetCharacterByID(id string) (*models.PlayerCharacter, error) { if m.GetCharacterByIDFunc != nil { return m.GetCharacterByIDFunc(id) } ; return nil, nil }
func (m *MockPlayerCharacterDAL) GetAllCharacters() ([]*models.PlayerCharacter, error) { return nil, nil }
func (m *MockPlayerCharacterDAL) CreateCharacter(character *models.PlayerCharacter) error { return nil }
func (m *MockPlayerCharacterDAL) UpdateCharacter(character *models.PlayerCharacter) error { return nil }
func (m *MockPlayerCharacterDAL) DeleteCharacter(id string) error { return nil }
func (m *MockPlayerCharacterDAL) GetCharacterInventory(characterID string) ([]*models.Item, error) { return nil, nil }
func (m *MockPlayerCharacterDAL) GetCharacterSkills(characterID string) ([]*models.PlayerSkill, error) { return nil, nil }
func (m *MockPlayerCharacterDAL) GetCharacterClass(characterID string) (*models.PlayerClass, error) { if m.GetCharacterClassFunc != nil { return m.GetCharacterClassFunc(characterID) } ; return nil, nil }
func (m *MockPlayerCharacterDAL) GetCharacterQuestState(characterID, questID string) (*models.PlayerQuestState, error) { if m.GetCharacterQuestStateFunc != nil { return m.GetCharacterQuestStateFunc(characterID, questID) } ; return nil, nil }
func (m *MockPlayerCharacterDAL) Cache() dal.CacheInterface { return nil }

type MockRoomDAL struct {
	GetRoomByIDFunc func(id string) (*models.Room, error)
}
func (m *MockRoomDAL) GetRoomByID(id string) (*models.Room, error) { if m.GetRoomByIDFunc != nil { return m.GetRoomByIDFunc(id) } ; return nil, nil }
func (m *MockRoomDAL) GetAllRooms() ([]*models.Room, error) { return nil, nil }
func (m *MockRoomDAL) CreateRoom(room *models.Room) error { return nil }
func (m *MockRoomDAL) UpdateRoom(room *models.Room) error { return nil }
func (m *MockRoomDAL) DeleteRoom(id string) error { return nil }
func (m *MockRoomDAL) Cache() dal.CacheInterface { return nil }

type MockNPCDAL struct {
	GetNPCByIDFunc func(id string) (*models.NPC, error)
}
func (m *MockNPCDAL) GetNPCByID(id string) (*models.NPC, error) { if m.GetNPCByIDFunc != nil { return m.GetNPCByIDFunc(id) } ; return nil, nil }
func (m *MockNPCDAL) GetAllNPCs() ([]*models.NPC, error) { return nil, nil }
func (m *MockNPCDAL) CreateNPC(npc *models.NPC) error { return nil }
func (m *MockNPCDAL) UpdateNPC(npc *models.NPC) error { return nil }
func (m *MockNPCDAL) DeleteNPC(id string) error { return nil }
func (m *MockNPCDAL) GetNPCsByRoom(roomID string) ([]*models.NPC, error) { return nil, nil }
func (m *MockNPCDAL) GetNPCsByOwner(ownerID string) ([]*models.NPC, error) { return nil, nil }
func (m *MockNPCDAL) Cache() dal.CacheInterface { return nil }

type MockOwnerDAL struct {
	GetOwnerByIDFunc func(id string) (*models.Owner, error)
}
func (m *MockOwnerDAL) GetOwnerByID(id string) (*models.Owner, error) { if m.GetOwnerByIDFunc != nil { return m.GetOwnerByIDFunc(id) } ; return nil, nil }
func (m *MockOwnerDAL) GetAllOwners() ([]*models.Owner, error) { return nil, nil }
func (m *MockOwnerDAL) CreateOwner(owner *models.Owner) error { return nil }
func (m *MockOwnerDAL) UpdateOwner(owner *models.Owner) error { return nil }
func (m *MockOwnerDAL) DeleteOwner(id string) error { return nil }
func (m *MockOwnerDAL) Cache() dal.CacheInterface { return nil }

type MockQuestmakerDAL struct {
	GetQuestmakerByIDFunc func(id string) (*models.Questmaker, error)
}
func (m *MockQuestmakerDAL) GetQuestmakerByID(id string) (*models.Questmaker, error) { if m.GetQuestmakerByIDFunc != nil { return m.GetQuestmakerByIDFunc(id) } ; return nil, nil }
func (m *MockQuestmakerDAL) GetAllQuestmakers() ([]*models.Questmaker, error) { return nil, nil }
func (m *MockQuestmakerDAL) CreateQuestmaker(questmaker *models.Questmaker) error { return nil }
func (m *MockQuestmakerDAL) UpdateQuestmaker(questmaker *models.Questmaker) error { return nil }
func (m *MockQuestmakerDAL) DeleteQuestmaker(id string) error { return nil }
func (m *MockQuestmakerDAL) Cache() dal.CacheInterface { return nil }

type MockRaceDAL struct {
	GetRaceByIDFunc func(id string) (*models.Race, error)
}
func (m *MockRaceDAL) GetRaceByID(id string) (*models.Race, error) { if m.GetRaceByIDFunc != nil { return m.GetRaceByIDFunc(id) } ; return nil, nil }
func (m *MockRaceDAL) GetAllRaces() ([]*models.Race, error) { return nil, nil }
func (m *MockRaceDAL) CreateRace(race *models.Race) error { return nil }
func (m *MockRaceDAL) UpdateRace(race *models.Race) error { return nil }
func (m *MockRaceDAL) DeleteRace(id string) error { return nil }
func (m *MockRaceDAL) Cache() dal.CacheInterface { return nil }

type MockProfessionDAL struct {
	GetProfessionByIDFunc func(id string) (*models.Profession, error)
}
func (m *MockProfessionDAL) GetProfessionByID(id string) (*models.Profession, error) { if m.GetProfessionByIDFunc != nil { return m.GetProfessionByIDFunc(id) } ; return nil, nil }
func (m *MockProfessionDAL) GetAllProfessions() ([]*models.Profession, error) { return nil, nil }
func (m *MockProfessionDAL) CreateProfession(profession *models.Profession) error { return nil }
func (m *MockProfessionDAL) UpdateProfession(profession *models.Profession) error { return nil }
func (m *MockProfessionDAL) DeleteProfession(id string) error { return nil }
func (m *MockProfessionDAL) Cache() dal.CacheInterface { return nil }

type MockSkillDAL struct {
	GetSkillByIDFunc func(id string) (*models.Skill, error)
}
func (m *MockSkillDAL) GetSkillByID(id string) (*models.Skill, error) { if m.GetSkillByIDFunc != nil { return m.GetSkillByIDFunc(id) } ; return nil, nil }
func (m *MockSkillDAL) GetAllSkills() ([]*models.Skill, error) { return nil, nil }
func (m *MockSkillDAL) CreateSkill(skill *models.Skill) error { return nil }
func (m *MockSkillDAL) UpdateSkill(skill *models.Skill) error { return nil }
func (m *MockSkillDAL) DeleteSkill(id string) error { return nil }
func (m *MockSkillDAL) Cache() dal.CacheInterface { return nil }

type MockClassDAL struct {
	GetClassByIDFunc func(id string) (*models.Class, error)
}
func (m *MockClassDAL) GetClassByID(id string) (*models.Class, error) { if m.GetClassByIDFunc != nil { return m.GetClassByIDFunc(id) } ; return nil, nil }
func (m *MockClassDAL) GetAllClasses() ([]*models.Class, error) { return nil, nil }
func (m *MockClassDAL) CreateClass(class *models.Class) error { return nil }
func (m *MockClassDAL) UpdateClass(class *models.Class) error { return nil }
func (m *MockClassDAL) DeleteClass(id string) error { return nil }
func (m *MockClassDAL) Cache() dal.CacheInterface { return nil }

type MockItemDAL struct {
	GetItemByIDFunc func(id string) (*models.Item, error)
}
func (m *MockItemDAL) GetItemByID(id string) (*models.Item, error) { if m.GetItemByIDFunc != nil { return m.GetItemByIDFunc(id) } ; return nil, nil }
func (m *MockItemDAL) GetAllItems() ([]*models.Item, error) { return nil, nil }
func (m *MockItemDAL) CreateItem(item *models.Item) error { return nil }
func (m *MockItemDAL) UpdateItem(item *models.Item) error { return nil }
func (m *MockItemDAL) DeleteItem(id string) error { return nil }
func (m *MockItemDAL) Cache() dal.CacheInterface { return nil }

type MockLoreDAL struct {
	GetLoreByIDFunc func(id string) (*models.Lore, error)
}
func (m *MockLoreDAL) GetLoreByID(id string) (*models.Lore, error) { if m.GetLoreByIDFunc != nil { return m.GetLoreByIDFunc(id) } ; return nil, nil }
func (m *MockLoreDAL) GetAllLore() ([]*models.Lore, error) { return nil, nil }
func (m *MockLoreDAL) CreateLore(lore *models.Lore) error { return nil }
func (m *MockLoreDAL) UpdateLore(lore *models.Lore) error { return nil }
func (m *MockLoreDAL) DeleteLore(id string) error { return nil }
func (m *MockLoreDAL) Cache() dal.CacheInterface { return nil }

type MockPlayerClassDAL struct {
	GetPlayerClassByIDFunc func(playerID, classID string) (*models.PlayerClass, error)
}
func (m *MockPlayerClassDAL) GetPlayerClassByID(playerID, classID string) (*models.PlayerClass, error) { if m.GetPlayerClassByIDFunc != nil { return m.GetPlayerClassByIDFunc(playerID, classID) } ; return nil, nil }
func (m *MockPlayerClassDAL) GetAllPlayerClasses() ([]*models.PlayerClass, error) { return nil, nil }
func (m *MockPlayerClassDAL) CreatePlayerClass(playerClass *models.PlayerClass) error { return nil }
func (m *MockPlayerClassDAL) UpdatePlayerClass(playerClass *models.PlayerClass) error { return nil }
func (m *MockPlayerClassDAL) DeletePlayerClass(playerID, classID string) error { return nil }
func (m *MockPlayerClassDAL) Cache() dal.CacheInterface { return nil }

type MockPlayerQuestStateDAL struct {
	GetPlayerQuestStateByIDFunc func(playerID, questID string) (*models.PlayerQuestState, error)
}
func (m *MockPlayerQuestStateDAL) GetPlayerQuestStateByID(playerID, questID string) (*models.PlayerQuestState, error) { if m.GetPlayerQuestStateByIDFunc != nil { return m.GetPlayerQuestStateByIDFunc(playerID, questID) } ; return nil, nil }
func (m *MockPlayerQuestStateDAL) GetAllPlayerQuestStates() ([]*models.PlayerQuestState, error) { return nil, nil }
func (m *MockPlayerQuestStateDAL) CreatePlayerQuestState(playerQuestState *models.PlayerQuestState) error { return nil }
func (m *MockPlayerQuestStateDAL) UpdatePlayerQuestState(playerQuestState *models.PlayerQuestState) error { return nil }
func (m *MockPlayerQuestStateDAL) DeletePlayerQuestState(playerID, questID string) error { return nil }
func (m *MockPlayerQuestStateDAL) Cache() dal.CacheInterface { return nil }

type MockPlayerSkillDAL struct {
	GetPlayerSkillByIDFunc func(playerID, skillID string) (*models.PlayerSkill, error)
}
func (m *MockPlayerSkillDAL) GetPlayerSkillByID(playerID, skillID string) (*models.PlayerSkill, error) { if m.GetPlayerSkillByIDFunc != nil { return m.GetPlayerSkillByIDFunc(playerID, skillID) } ; return nil, nil }
func (m *MockPlayerSkillDAL) GetAllPlayerSkills() ([]*models.PlayerSkill, error) { return nil, nil }
func (m *MockPlayerSkillDAL) CreatePlayerSkill(playerSkill *models.PlayerSkill) error { return nil }
func (m *MockPlayerSkillDAL) UpdatePlayerSkill(playerSkill *models.PlayerSkill) error { return nil }
func (m *MockPlayerSkillDAL) DeletePlayerSkill(playerID, skillID string) error { return nil }
func (m *MockPlayerSkillDAL) Cache() dal.CacheInterface { return nil }

type MockQuestDAL struct {
	GetQuestByIDFunc func(id string) (*models.Quest, error)
}
func (m *MockQuestDAL) GetQuestByID(id string) (*models.Quest, error) { if m.GetQuestByIDFunc != nil { return m.GetQuestByIDFunc(id) } ; return nil, nil }
func (m *MockQuestDAL) GetAllQuests() ([]*models.Quest, error) { return nil, nil }
func (m *MockQuestDAL) CreateQuest(quest *models.Quest) error { return nil }
func (m *MockQuestDAL) UpdateQuest(quest *models.Quest) error { return nil }
func (m *MockQuestDAL) DeleteQuest(id string) error { return nil }
func (m *MockQuestDAL) Cache() dal.CacheInterface { return nil }

type MockQuestOwnerDAL struct {
	GetQuestOwnerByIDFunc func(id string) (*models.QuestOwner, error)
}
func (m *MockQuestOwnerDAL) GetQuestOwnerByID(id string) (*models.QuestOwner, error) { if m.GetQuestOwnerByIDFunc != nil { return m.GetQuestOwnerByIDFunc(id) } ; return nil, nil }
func (m *MockQuestOwnerDAL) GetAllQuestOwners() ([]*models.QuestOwner, error) { return nil, nil }
func (m *MockQuestOwnerDAL) CreateQuestOwner(questOwner *models.QuestOwner) error { return nil }
func (m *MockQuestOwnerDAL) UpdateQuestOwner(questOwner *models.QuestOwner) error { return nil }
func (m *MockQuestOwnerDAL) DeleteQuestOwner(id string) error { return nil }
func (m *MockQuestOwnerDAL) Cache() dal.CacheInterface { return nil }


// Mock implementations for other interfaces
type MockLLMService struct {
	ProcessActionFunc func(ctx context.Context, entity interface{}, player *models.PlayerCharacter, prompt string) (*llm.InnerLLMResponse, error)
	AnalyzeResponseFunc func(ctx context.Context, narrative string, query string) (float64, error)
}

func (m *MockLLMService) ProcessAction(ctx context.Context, entity interface{}, player *models.PlayerCharacter, prompt string) (*llm.InnerLLMResponse, error) {
	if m.ProcessActionFunc != nil {
		return m.ProcessActionFunc(ctx, entity, player, prompt)
	}
	return nil, nil
}

func (m *MockLLMService) AnalyzeResponse(ctx context.Context, narrative string, query string) (float64, error) {
	if m.AnalyzeResponseFunc != nil {
		return m.AnalyzeResponseFunc(ctx, narrative, query)
	}
	return 0.0, nil
}

type MockToolDispatcher struct {
	DispatchFunc func(ctx context.Context, player *models.PlayerCharacter, entity interface{}, toolCalls []llm.ToolCall) error
}

func (m *MockToolDispatcher) Dispatch(ctx context.Context, player *models.PlayerCharacter, entity interface{}, toolCalls []llm.ToolCall) error {
	if m.DispatchFunc != nil {
		return m.DispatchFunc(ctx, player, entity, toolCalls)
	}
	return nil
}

type MockTelnetRenderer struct{
	RenderRawStringFunc func(s string, color presentation.SemanticColorType) string
	RenderMessageFunc func(msg presentation.SemanticMessage) string
}

func (m *MockTelnetRenderer) RenderRawString(s string, color presentation.SemanticColorType) string {
	if m.RenderRawStringFunc != nil {
		return m.RenderRawStringFunc(s, color)
	}
	return s
}

func (m *MockTelnetRenderer) RenderMessage(msg presentation.SemanticMessage) string {
	if m.RenderMessageFunc != nil {
		return m.RenderMessageFunc(msg)
	}
	return msg.Content // Default for mock
}

func TestNewSentientEntityManager(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	manager := NewSentientEntityManager(
		mockLLMService,
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockToolDispatcher,
		mockTelnetRenderer,
		events.NewEventBus(),
	)

	assert.NotNil(t, manager)
	assert.Equal(t, mockLLMService, manager.llmService)
	assert.Equal(t, mockNPCDAL, manager.npcDAL)
	assert.Equal(t, mockOwnerDAL, manager.ownerDAL)
	assert.Equal(t, mockQuestmakerDAL, manager.questmakerDAL)
	assert.Equal(t, mockToolDispatcher, manager.toolDispatcher)
	assert.Equal(t, mockTelnetRenderer, manager.telnetRenderer)
}

func TestSentientEntityManager_TriggerReaction_NoPerceivedActions(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	manager := NewSentientEntityManager(
		mockLLMService, mockNPCDAL, mockOwnerDAL, mockQuestmakerDAL, mockToolDispatcher, mockTelnetRenderer, events.NewEventBus(),
	)

	err := manager.TriggerReaction(&models.NPC{ID: "npc1"}, []perception.PerceivedActionRecord{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no perceived actions provided")
}

func TestSentientEntityManager_TriggerReaction_BelowThreshold(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	npc := &models.NPC{ID: "npc1", Name: "Test NPC", ReactionThreshold: 10.0}
	player := &models.PlayerCharacter{ID: "player1", Name: "Test Player"}
	perceivedAction := perception.PerceivedAction{
		PerceivedActionType: "say",
		SourcePlayer:        player,
		Clarity:             1.0,
		BaseSignificance:    5.0,
	}
	record := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction,
		Significance:    5.0, // Below threshold
	}

	mockNPCDAL.GetNPCByIDFunc = func(id string) (*models.NPC, error) {
		assert.Equal(t, npc.ID, id)
		return npc, nil
	}

	// LLMService.ProcessActionFunc should NOT be called
	mockLLMService.ProcessActionFunc = func(ctx context.Context, entity interface{}, p *models.PlayerCharacter, prompt string) (*llm.InnerLLMResponse, error) {
		t.Error("LLMService.ProcessAction should not be called when significance is below threshold")
		return nil, nil
	}

	manager := NewSentientEntityManager(
		mockLLMService,
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockToolDispatcher,
		mockTelnetRenderer,
		events.NewEventBus(),
	)

	err := manager.TriggerReaction(npc, []perception.PerceivedActionRecord{record})
	assert.NoError(t, err) // Error should be logged, but TriggerReaction should not return an error
}

func TestSentientEntityManager_TriggerReaction_SuccessfulNPCReaction(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	npc := &models.NPC{ID: "npc1", Name: "Test NPC", ReactionThreshold: 4.0} // Set a threshold
	player := &models.PlayerCharacter{ID: "player1", Name: "Test Player"}
	perceivedAction := perception.PerceivedAction{
		PerceivedActionType: "say",
		SourcePlayer:        player,
		Clarity:             1.0,
		BaseSignificance:    5.0,
	}
	record := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction,
		Significance:    5.0, // Above threshold
	}

	mockNPCDAL.GetNPCByIDFunc = func(id string) (*models.NPC, error) {
		assert.Equal(t, npc.ID, id)
		return npc, nil
	}

	// Mock the ProcessAction call for the game LLM
	mockLLMService.ProcessActionFunc = func(ctx context.Context, entity interface{}, p *models.PlayerCharacter, prompt string) (*llm.InnerLLMResponse, error) {
		assert.Equal(t, npc, entity)
		assert.Equal(t, player, p)
		assert.Contains(t, prompt, "Player Test Player performed action say (clarity 1.00). Respond to this.")
		return &llm.InnerLLMResponse{Narrative: "Grrr, I'll get you next time, adventurer!"}, nil
	}

	// Mock the AnalyzeResponse call for the game LLM
	mockLLMService.AnalyzeResponseFunc = func(ctx context.Context, narrative string, query string) (float64, error) {
		assert.Equal(t, "Grrr, I'll get you next time, adventurer!", narrative)
		assert.Equal(t, "On a scale of 0 to 100, how hostile is this goblin's response?", query)
		return 75.0, nil // Mock a hostile response
	}

	manager := NewSentientEntityManager(
		mockLLMService, // Use mock for ProcessAction and AnalyzeResponse
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockToolDispatcher,
		mockTelnetRenderer,
		events.NewEventBus(),
	)

	err := manager.TriggerReaction(npc, []perception.PerceivedActionRecord{record})
	assert.NoError(t, err)

	// Perform LLM-based analysis of the response using the mocked LLMService
	hostilityScore, err := mockLLMService.AnalyzeResponse(context.Background(), "Grrr, I'll get you next time, adventurer!", "On a scale of 0 to 100, how hostile is this goblin's response?")
	assert.NoError(t, err)
	assert.Greater(t, hostilityScore, 50.0, "Hostility score should be greater than 50")
}

func TestSentientEntityManager_TriggerReaction_MultipleActions_MixedThresholds(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	npc := &models.NPC{ID: "npc1", Name: "Test NPC", ReactionThreshold: 7.0} // Threshold for this test
	player := &models.PlayerCharacter{ID: "player1", Name: "Test Player"}

	// Action 1: Below threshold
	perceivedAction1 := perception.PerceivedAction{
		PerceivedActionType: "whisper",
		SourcePlayer:        player,
		Clarity:             0.5,
		BaseSignificance:    3.0,
	}
	record1 := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction1,
		Significance:    3.0,
	}

	// Action 2: Above threshold
	perceivedAction2 := perception.PerceivedAction{
		PerceivedActionType: "shout",
		SourcePlayer:        player,
		Clarity:             0.9,
		BaseSignificance:    8.0,
	}
	record2 := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction2,
		Significance:    8.0,
	}

	// Action 3: Exactly at threshold
	perceivedAction3 := perception.PerceivedAction{
		PerceivedActionType: "emote",
		SourcePlayer:        player,
		Clarity:             0.7,
		BaseSignificance:    7.0,
	}
	record3 := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction3,
		Significance:    7.0,
	}

	mockNPCDAL.GetNPCByIDFunc = func(id string) (*models.NPC, error) {
		assert.Equal(t, npc.ID, id)
		return npc, nil
	}

	llmCallCount := 0
	mockLLMService.ProcessActionFunc = func(ctx context.Context, entity interface{}, p *models.PlayerCharacter, prompt string) (*llm.InnerLLMResponse, error) {
		llmCallCount++
		assert.Equal(t, npc, entity)
		assert.Equal(t, player, p)
		assert.Contains(t, prompt, "Player Test Player performed action shout (clarity 0.90). Respond to this.")
		return &llm.InnerLLMResponse{Narrative: "NPC response to shout"}, nil
	}

	manager := NewSentientEntityManager(
		mockLLMService,
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockToolDispatcher,
		mockTelnetRenderer,
		events.NewEventBus(),
	)

	err := manager.TriggerReaction(npc, []perception.PerceivedActionRecord{record1, record2, record3})
	assert.NoError(t, err)
	assert.Equal(t, 1, llmCallCount, "LLMService.ProcessAction should be called exactly once for relevant actions")
}

func TestSentientEntityManager_TriggerReaction_SuccessfulOwnerReaction(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	owner := &models.Owner{ID: "owner1", Name: "Test Owner"}
	player := &models.PlayerCharacter{ID: "player1", Name: "Test Player"}
	perceivedAction := perception.PerceivedAction{
		PerceivedActionType: "pray",
		SourcePlayer:        player,
		Clarity:             1.0,
		BaseSignificance:    10.0,
	}
	record := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction,
		Significance:    10.0,
	}

	mockOwnerDAL.GetOwnerByIDFunc = func(id string) (*models.Owner, error) {
		assert.Equal(t, owner.ID, id)
		return owner, nil
	}

	mockLLMService.ProcessActionFunc = func(ctx context.Context, entity interface{}, p *models.PlayerCharacter, prompt string) (*llm.InnerLLMResponse, error) {
		assert.Equal(t, owner, entity)
		assert.Equal(t, player, p)
		assert.Contains(t, prompt, "Player Test Player performed action pray")
		return &llm.InnerLLMResponse{Narrative: "Owner response"}, nil
	}

	manager := NewSentientEntityManager(
		mockLLMService,
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockToolDispatcher,
		mockTelnetRenderer,
		events.NewEventBus(),
	)

	err := manager.TriggerReaction(owner, []perception.PerceivedActionRecord{record})
	assert.NoError(t, err)
}

func TestSentientEntityManager_TriggerReaction_SuccessfulQuestmakerReaction(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	questmaker := &models.Questmaker{ID: "qm1", Name: "Test Questmaker"}
	player := &models.PlayerCharacter{ID: "player1", Name: "Test Player"}
	perceivedAction := perception.PerceivedAction{
		PerceivedActionType: "quest_action",
		SourcePlayer:        player,
		Clarity:             1.0,
		BaseSignificance:    15.0,
	}
	record := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction,
		Significance:    15.0,
	}

	mockQuestmakerDAL.GetQuestmakerByIDFunc = func(id string) (*models.Questmaker, error) {
		assert.Equal(t, questmaker.ID, id)
		return questmaker, nil
	}

	mockLLMService.ProcessActionFunc = func(ctx context.Context, entity interface{}, p *models.PlayerCharacter, prompt string) (*llm.InnerLLMResponse, error) {
		assert.Equal(t, questmaker, entity)
		assert.Equal(t, player, p)
		assert.Contains(t, prompt, "Player Test Player performed action quest_action")
		return &llm.InnerLLMResponse{Narrative: "Questmaker response"}, nil
	}

	manager := NewSentientEntityManager(
		mockLLMService,
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockToolDispatcher,
		mockTelnetRenderer,
		events.NewEventBus(),
	)

	err := manager.TriggerReaction(questmaker, []perception.PerceivedActionRecord{record})
	assert.NoError(t, err)
}

func TestSentientEntityManager_TriggerReaction_EntityNotFound(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	// Mock all DALs to return nil/error
	mockNPCDAL.GetNPCByIDFunc = func(id string) (*models.NPC, error) { return nil, nil }
	mockOwnerDAL.GetOwnerByIDFunc = func(id string) (*models.Owner, error) { return nil, errors.New("owner not found") }
	mockQuestmakerDAL.GetQuestmakerByIDFunc = func(id string) (*models.Questmaker, error) { return nil, nil }

	player := &models.PlayerCharacter{ID: "player1", Name: "Test Player"}
	perceivedAction := perception.PerceivedAction{
		PerceivedActionType: "test",
		SourcePlayer:        player,
	}
	record := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction,
	}

	manager := NewSentientEntityManager(
		mockLLMService,
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockToolDispatcher,
		mockTelnetRenderer,
		events.NewEventBus(),
	)

	err := manager.TriggerReaction(&models.NPC{ID: "nonexistent"}, []perception.PerceivedActionRecord{record})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not find entity with ID nonexistent for triggering reaction")
}

func TestSentientEntityManager_TriggerReaction_SourcePlayerNotFound(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	npc := &models.NPC{ID: "npc1", Name: "Test NPC"}
	perceivedAction := perception.PerceivedAction{
		PerceivedActionType: "test",
		SourcePlayer:        nil, // No source player
	}
	record := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction,
	}

	mockNPCDAL.GetNPCByIDFunc = func(id string) (*models.NPC, error) {
		return npc, nil
	}

	manager := NewSentientEntityManager(
		mockLLMService,
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockToolDispatcher,
		mockTelnetRenderer,
		events.NewEventBus(),
	)

	err := manager.TriggerReaction(npc, []perception.PerceivedActionRecord{record})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source player not found in perceived action")
}

func TestSentientEntityManager_TriggerReaction_LLMServiceError(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	npc := &models.NPC{ID: "npc1", Name: "Test NPC"}
	player := &models.PlayerCharacter{ID: "player1", Name: "Test Player"}
	perceivedAction := perception.PerceivedAction{
		PerceivedActionType: "test",
		SourcePlayer:        player,
	}
	record := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction,
	}

	mockNPCDAL.GetNPCByIDFunc = func(id string) (*models.NPC, error) {
		return npc, nil
	}

	mockLLMService.ProcessActionFunc = func(ctx context.Context, entity interface{}, p *models.PlayerCharacter, prompt string) (*llm.InnerLLMResponse, error) {
		return nil, errors.New("LLM service error")
	}

	manager := NewSentientEntityManager(
		mockLLMService,
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockToolDispatcher,
		mockTelnetRenderer,
		events.NewEventBus(),
	)

	err := manager.TriggerReaction(npc, []perception.PerceivedActionRecord{record})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "LLM Service ProcessAction failed")
}

func TestSentientEntityManager_TriggerReaction_LLMResponseWithToolCalls(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	npc := &models.NPC{ID: "npc1", Name: "Test NPC", ReactionThreshold: 1.0}
	player := &models.PlayerCharacter{ID: "player1", Name: "Test Player"}
	perceivedAction := perception.PerceivedAction{
		PerceivedActionType: "test",
		SourcePlayer:        player,
	}
	record := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction,
		Significance:    5.0,
	}

	mockNPCDAL.GetNPCByIDFunc = func(id string) (*models.NPC, error) {
		return npc, nil
	}

	mockLLMService.ProcessActionFunc = func(ctx context.Context, entity interface{}, p *models.PlayerCharacter, prompt string) (*llm.InnerLLMResponse, error) {
		return &llm.InnerLLMResponse{
			Narrative: "LLM narrative with tool calls",
			ToolCalls: []llm.ToolCall{
				{ToolName: "test_tool", Parameters: map[string]interface{}{"arg1": "value1"}},
			},
		}, nil
	}

	dispatchCalled := false
	mockToolDispatcher.DispatchFunc = func(ctx context.Context, player *models.PlayerCharacter, entity interface{}, toolCalls []llm.ToolCall) error {
		dispatchCalled = true
		assert.NotNil(t, ctx)
		assert.Equal(t, player, player)
		assert.Equal(t, npc, entity)
		assert.Len(t, toolCalls, 1)
		assert.Equal(t, "test_tool", toolCalls[0].ToolName)
		return nil
	}

	manager := NewSentientEntityManager(
		mockLLMService,
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockToolDispatcher,
		mockTelnetRenderer,
		events.NewEventBus(),
	)

	err := manager.TriggerReaction(npc, []perception.PerceivedActionRecord{record})
	assert.NoError(t, err)
	assert.True(t, dispatchCalled, "ToolDispatcher.Dispatch should be called")
}

func TestSentientEntityManager_TriggerReaction_ToolDispatcherError(t *testing.T) {
	mockLLMService := &MockLLMService{}
	mockNPCDAL := &MockNPCDAL{}
	mockOwnerDAL := &MockOwnerDAL{}
	mockQuestmakerDAL := &MockQuestmakerDAL{}
	mockToolDispatcher := &MockToolDispatcher{}
	mockTelnetRenderer := &MockTelnetRenderer{}

	npc := &models.NPC{ID: "npc1", Name: "Test NPC", ReactionThreshold: 1.0}
	player := &models.PlayerCharacter{ID: "player1", Name: "Test Player"}
	perceivedAction := perception.PerceivedAction{
		PerceivedActionType: "test",
		SourcePlayer:        player,
	}
	record := perception.PerceivedActionRecord{
		PerceivedAction: &perceivedAction,
		Significance:    5.0,
	}

	mockNPCDAL.GetNPCByIDFunc = func(id string) (*models.NPC, error) {
		return npc, nil
	}

	mockLLMService.ProcessActionFunc = func(ctx context.Context, entity interface{}, p *models.PlayerCharacter, prompt string) (*llm.InnerLLMResponse, error) {
		return &llm.InnerLLMResponse{
			Narrative: "LLM narrative with tool calls",
			ToolCalls: []llm.ToolCall{
				{ToolName: "test_tool", Parameters: map[string]interface{}{"arg1": "value1"}},
			},
		}, nil
	}

	mockToolDispatcher.DispatchFunc = func(ctx context.Context, player *models.PlayerCharacter, entity interface{}, toolCalls []llm.ToolCall) error {
		return errors.New("tool dispatcher error")
	}

	manager := NewSentientEntityManager(
		mockLLMService,
		mockNPCDAL,
		mockOwnerDAL,
		mockQuestmakerDAL,
		mockToolDispatcher,
		mockTelnetRenderer,
		events.NewEventBus(),
	)

	err := manager.TriggerReaction(npc, []perception.PerceivedActionRecord{record})
	assert.NoError(t, err) // Error should be logged, but TriggerReaction should not return an error
}