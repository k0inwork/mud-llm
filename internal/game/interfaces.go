package game

import (
	"context"
	"mud/internal/game/events"
	"mud/internal/game/perception"
	"mud/internal/llm"
	"mud/internal/models"
	"mud/internal/presentation"
)

// PerceptionFilterInterface defines the methods used by ActionSignificanceMonitor on PerceptionFilter.
type PerceptionFilterInterface interface {
	Filter(event *events.ActionEvent, observer interface{}) (*perception.PerceivedAction, error)
}

// SentientEntityManagerInterface defines the methods used by ActionSignificanceMonitor on SentientEntityManager.
type SentientEntityManagerInterface interface {
	TriggerReaction(observer interface{}, perceivedActions []perception.PerceivedActionRecord) error
}

// LLMServiceInterface defines the methods used by SentientEntityManager on LLMService.
type LLMServiceInterface interface {
	ProcessAction(ctx context.Context, entity interface{}, player *models.PlayerCharacter, playerAction string) (*llm.InnerLLMResponse, error)
	AnalyzeResponse(ctx context.Context, narrative string, query string) (float64, error)
}

// ToolDispatcherInterface defines the methods used by SentientEntityManager on ToolDispatcher.
type ToolDispatcherInterface interface {
	Dispatch(ctx context.Context, player *models.PlayerCharacter, entity interface{}, toolCalls []llm.ToolCall) error
}

// TelnetRendererInterface defines the methods used by SentientEntityManager on TelnetRenderer.
type TelnetRendererInterface interface {
	RenderMessage(msg presentation.SemanticMessage) string
	RenderRawString(s string, color presentation.SemanticColorType) string
}
