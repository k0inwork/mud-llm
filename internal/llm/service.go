package llm

import (
	"context"
	"fmt"
	"mud/internal/dal"
	"mud/internal/models"
	"time"
)

type LLMService struct {
	client      *Client
	cache       *CacheManager
	dal         *dal.DAL
}

func NewLLMService(dal *dal.DAL) *LLMService {
	return &LLMService{
		client:      NewClient(),
		cache:       NewCacheManager(),
		dal:         dal,
	}
}

func (s *LLMService) ProcessAction(ctx context.Context, entity interface{}, player *models.Player, playerAction string) (*InnerLLMResponse, error) {
	entityID, err := getEntityID(entity)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("base_prompt:%s", entityID)
	
	// 1. Check cache for base prompt
	cachedPrompt, found := s.cache.Get(cacheKey)
	var basePrompt string
	if found {
		basePrompt = cachedPrompt.(string)
	} else {
		// 2. If not found, assemble and cache it
		promptData := &PromptData{
			Entity:      entity,
			Player:      player,
			DAL:         s.dal,
		}
		assembledPrompt, err := AssemblePrompt(promptData)
		if err != nil {
			return nil, fmt.Errorf("failed to assemble prompt: %w", err)
		}
		basePrompt = assembledPrompt
		s.cache.Set(cacheKey, basePrompt, 5*time.Minute) // Cache for 5 minutes
	}

	// 3. Append dynamic player action
	finalPrompt := fmt.Sprintf("%s\nPlayer action: %s", basePrompt, playerAction)

	// 4. Send to LLM
	return s.client.SendPrompt(ctx, finalPrompt)
}

func getEntityID(entity interface{}) (string, error) {
	switch v := entity.(type) {
	case *models.NPC:
		return v.ID, nil
	case *models.Owner:
		return v.ID, nil
	case *models.Questmaker:
		return v.ID, nil
	default:
		return "", fmt.Errorf("unknown entity type for getting ID")
	}
}