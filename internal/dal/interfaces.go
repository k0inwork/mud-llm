package dal

import (
	"mud/internal/models"
	"time"
)

// PlayerCharacterDALInterface defines the methods used by TelnetServer on PlayerCharacterDAL.
type PlayerCharacterDALInterface interface {
	GetCharacterByID(id string) (*models.PlayerCharacter, error)
	GetAllCharacters() ([]*models.PlayerCharacter, error)
	CreateCharacter(character *models.PlayerCharacter) error
	UpdateCharacter(character *models.PlayerCharacter) error
	DeleteCharacter(id string) error
	GetCharacterInventory(characterID string) ([]*models.Item, error)
	GetCharacterSkills(characterID string) ([]*models.PlayerSkill, error)
	GetCharacterClass(characterID string) (*models.PlayerClass, error)
	GetCharacterQuestState(characterID, questID string) (*models.PlayerQuestState, error)
	GetCharactersByAccountID(accountID string) ([]*models.PlayerCharacter, error)
	Cache() CacheInterface
}

// RoomDALInterface defines the methods used by PerceptionFilter on RoomDAL.
type RoomDALInterface interface {
	GetRoomByID(id string) (*models.Room, error)
	GetAllRooms() ([]*models.Room, error)
	CreateRoom(room *models.Room) error
	UpdateRoom(room *models.Room) error
	DeleteRoom(id string) error
	Cache() CacheInterface
}

// RaceDALInterface defines the methods used by PerceptionFilter on RaceDAL.
type RaceDALInterface interface {
	GetRaceByID(id string) (*models.Race, error)
	GetAllRaces() ([]*models.Race, error)
	CreateRace(race *models.Race) error
	UpdateRace(race *models.Race) error
	DeleteRace(id string) error
	Cache() CacheInterface
}

// ProfessionDALInterface defines the methods used by PerceptionFilter on ProfessionDAL.
type ProfessionDALInterface interface {
	GetProfessionByID(id string) (*models.Profession, error)
	GetAllProfessions() ([]*models.Profession, error)
	CreateProfession(profession *models.Profession) error
	UpdateProfession(profession *models.Profession) error
	DeleteProfession(id string) error
	Cache() CacheInterface
}

// CacheInterface defines the methods required for a cache.
type CacheInterface interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	SetMany(items map[string]interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
}

// NPCDALInterface defines the methods used by ActionSignificanceMonitor on NPCDAL.
type NPCDALInterface interface {
	GetAllNPCs() ([]*models.NPC, error)
	GetNPCByID(id string) (*models.NPC, error)
	UpdateNPC(npc *models.NPC) error
	GetNPCsByOwner(ownerID string) ([]*models.NPC, error)
	CreateNPC(npc *models.NPC) error
	DeleteNPC(id string) error
	GetNPCsByRoom(roomID string) ([]*models.NPC, error)
	Cache() CacheInterface
}

// OwnerDALInterface defines the methods used by ActionSignificanceMonitor on OwnerDAL.
type OwnerDALInterface interface {
	GetAllOwners() ([]*models.Owner, error)
	GetOwnerByID(id string) (*models.Owner, error)
	UpdateOwner(owner *models.Owner) error
	CreateOwner(owner *models.Owner) error
	DeleteOwner(id string) error
	Cache() CacheInterface
}

// QuestmakerDALInterface defines the methods used by ActionSignificanceMonitor on QuestmakerDAL.
type QuestmakerDALInterface interface {
	GetAllQuestmakers() ([]*models.Questmaker, error)
	GetQuestmakerByID(id string) (*models.Questmaker, error)
	CreateQuestmaker(questmaker *models.Questmaker) error
	UpdateQuestmaker(questmaker *models.Questmaker) error
	DeleteQuestmaker(id string) error
	Cache() CacheInterface
}

// SkillDALInterface defines the methods for SkillDAL.
type SkillDALInterface interface {
	GetSkillByID(id string) (*models.Skill, error)
	GetAllSkills() ([]*models.Skill, error)
	CreateSkill(skill *models.Skill) error
	UpdateSkill(skill *models.Skill) error
	DeleteSkill(id string) error
	Cache() CacheInterface
}

// ClassDALInterface defines the methods for ClassDAL.
type ClassDALInterface interface {
	GetClassByID(id string) (*models.Class, error)
	GetAllClasses() ([]*models.Class, error)
	CreateClass(class *models.Class) error
	UpdateClass(class *models.Class) error
	DeleteClass(id string) error
	Cache() CacheInterface
}

// ItemDALInterface defines the methods for ItemDAL.
type ItemDALInterface interface {
	GetItemByID(id string) (*models.Item, error)
	GetAllItems() ([]*models.Item, error)
	CreateItem(item *models.Item) error
	UpdateItem(item *models.Item) error
	DeleteItem(id string) error
	Cache() CacheInterface
}

// LoreDALInterface defines the methods for LoreDAL.
type LoreDALInterface interface {
	GetLoreByID(id string) (*models.Lore, error)
	GetAllLore() ([]*models.Lore, error)
	CreateLore(lore *models.Lore) error
	UpdateLore(lore *models.Lore) error
	DeleteLore(id string) error
	Cache() CacheInterface
}

// PlayerClassDALInterface defines the methods for PlayerClassDAL.
type PlayerClassDALInterface interface {
	GetPlayerClassByID(playerID, classID string) (*models.PlayerClass, error)
	GetAllPlayerClasses() ([]*models.PlayerClass, error)
	CreatePlayerClass(playerClass *models.PlayerClass) error
	UpdatePlayerClass(playerClass *models.PlayerClass) error
	DeletePlayerClass(playerID, classID string) error
	Cache() CacheInterface
}

// PlayerQuestStateDALInterface defines the methods for PlayerQuestStateDAL.
type PlayerQuestStateDALInterface interface {
	GetPlayerQuestStateByID(playerID, questID string) (*models.PlayerQuestState, error)
	GetAllPlayerQuestStates() ([]*models.PlayerQuestState, error)
	CreatePlayerQuestState(playerQuestState *models.PlayerQuestState) error
	UpdatePlayerQuestState(playerQuestState *models.PlayerQuestState) error
	DeletePlayerQuestState(playerID, questID string) error
	Cache() CacheInterface
}

// PlayerSkillDALInterface defines the methods for PlayerSkillDAL.
type PlayerSkillDALInterface interface {
	GetPlayerSkillByID(playerID, skillID string) (*models.PlayerSkill, error)
	GetAllPlayerSkills() ([]*models.PlayerSkill, error)
	CreatePlayerSkill(playerSkill *models.PlayerSkill) error
	UpdatePlayerSkill(playerSkill *models.PlayerSkill) error
	DeletePlayerSkill(playerID, skillID string) error
	Cache() CacheInterface
}

// QuestDALInterface defines the methods for QuestDAL.
type QuestDALInterface interface {
	GetQuestByID(id string) (*models.Quest, error)
	GetAllQuests() ([]*models.Quest, error)
	CreateQuest(quest *models.Quest) error
	UpdateQuest(quest *models.Quest) error
	DeleteQuest(id string) error
	Cache() CacheInterface
}

// QuestOwnerDALInterface defines the methods for QuestOwnerDAL.
type QuestOwnerDALInterface interface {
	GetQuestOwnerByID(id string) (*models.QuestOwner, error)
	GetAllQuestOwners() ([]*models.QuestOwner, error)
	CreateQuestOwner(questOwner *models.QuestOwner) error
	UpdateQuestOwner(questOwner *models.QuestOwner) error
	DeleteQuestOwner(id string) error
	Cache() CacheInterface
}
