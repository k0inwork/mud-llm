package llm

import (
	"context"
	"mud/internal/models"
)

// LLMServiceInterface defines the methods that an LLM service should implement.
type LLMServiceInterface interface {
	ProcessAction(ctx context.Context, entity interface{}, player *models.PlayerCharacter, playerAction string) (*InnerLLMResponse, error)
	AnalyzeResponse(ctx context.Context, narrative string, query string) (float64, error)
}
