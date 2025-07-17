package dal

import (
	"mud/internal/models"
	"time"
)

// RoomDALInterface defines the methods used by PerceptionFilter on RoomDAL.
type RoomDALInterface interface {
	GetRoomByID(id string) (*models.Room, error)
}

// RaceDALInterface defines the methods used by PerceptionFilter on RaceDAL.
type RaceDALInterface interface {
	GetRaceByID(id string) (*models.Race, error)
}

// ProfessionDALInterface defines the methods used by PerceptionFilter on ProfessionDAL.
type ProfessionDALInterface interface {
	GetProfessionByID(id string) (*models.Profession, error)
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
}

// OwnerDALInterface defines the methods used by ActionSignificanceMonitor on OwnerDAL.
type OwnerDALInterface interface {
	GetAllOwners() ([]*models.Owner, error)
	GetOwnerByID(id string) (*models.Owner, error)
	UpdateOwner(owner *models.Owner) error
}

// QuestmakerDALInterface defines the methods used by ActionSignificanceMonitor on QuestmakerDAL.
type QuestmakerDALInterface interface {
	GetAllQuestmakers() ([]*models.Questmaker, error)
	GetQuestmakerByID(id string) (*models.Questmaker, error)
}