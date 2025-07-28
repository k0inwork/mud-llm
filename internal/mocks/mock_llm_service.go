package mocks

import (
	"context"
	"fmt"
	"mud/internal/llm"
	"mud/internal/models"
)

// MockLLMService is a mock implementation of LLMServiceInterface for testing.
type MockLLMService struct {
	ProcessActionFunc func(ctx context.Context, entity interface{}, player *models.PlayerCharacter, playerAction string) (*llm.InnerLLMResponse, error)
	AnalyzeResponseFunc func(ctx context.Context, narrative string, query string) (float64, error)
}

// ProcessAction calls the mock function if set, otherwise returns a default response.
func (m *MockLLMService) ProcessAction(ctx context.Context, entity interface{}, player *models.PlayerCharacter, playerAction string) (*llm.InnerLLMResponse, error) {
	if m.ProcessActionFunc != nil {
		return m.ProcessActionFunc(ctx, entity, player, playerAction)
	}
	// Default mock response
	return &llm.InnerLLMResponse{
		Narrative: fmt.Sprintf("Mock LLM response for action '%s' by %s.", playerAction, player.Name),
		ToolCalls: []llm.ToolCall{},
	},
	nil
}

// AnalyzeResponse calls the mock function if set, otherwise returns a default score.
func (m *MockLLMService) AnalyzeResponse(ctx context.Context, narrative string, query string) (float64, error) {
	if m.AnalyzeResponseFunc != nil {
		return m.AnalyzeResponseFunc(ctx, narrative, query)
	}
	// Default mock score
	return 75.0, nil // A default hostility score for testing
}
