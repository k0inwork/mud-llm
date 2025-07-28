package llm

import (
	"fmt"
	"mud/internal/dal"
	"mud/internal/game/perception"
	"mud/internal/models"
	"strings"
)

type PromptData struct {
	Entity        interface{}
	Player        *models.PlayerCharacter
	Room          *models.Room
	RecentActions []*perception.PerceivedAction
	LoreEntries   []*models.Lore
	DAL           *dal.DAL
	PlayerAction  string
}

func AssemblePrompt(data *PromptData) (string, error) {
	var sb strings.Builder

	// 1. Add entity's personality
	personality, err := getEntityPersonality(data.Entity)
	if err != nil {
		return "", err
	}
	sb.WriteString(fmt.Sprintf("Your personality: %s\n\n", personality))

	// 2. Add available tools
	tools, err := getEntityTools(data.Entity)
	if err != nil {
		return "", err
	}
	if len(tools) > 0 {
		sb.WriteString("You have the following tools available:\n")
		for _, tool := range tools {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", tool.Name, tool.Description))
		}
		sb.WriteString("\n")
	}

	// 3. Add memories about the player
	memories, err := getEntityMemories(data.Entity, data.Player.ID)
	if err != nil {
		return "", err
	}
	if len(memories) > 0 {
		sb.WriteString("Your memories about this player:\n")
		for _, memory := range memories {
			sb.WriteString(fmt.Sprintf("- %s\n", memory))
		}
		sb.WriteString("\n")
	}

	// 4. Add relevant lore
	if len(data.LoreEntries) > 0 {
		sb.WriteString("Relevant lore:\n")
		for _, l := range data.LoreEntries {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", l.Title, l.Content))
		}
		sb.WriteString("\n")
	}

	// 5. Add recent actions
	if len(data.RecentActions) > 0 {
		sb.WriteString("Recent perceived actions by the player:\n")
		for _, action := range data.RecentActions {
			sb.WriteString(fmt.Sprintf("- %s (Significance: %.2f)\n", action.PerceivedActionType, action.BaseSignificance))
		}
		sb.WriteString("\n")
	}

	// 6. Add the player's action
	sb.WriteString(fmt.Sprintf("The player's action: %s\n", data.PlayerAction))

	return sb.String(), nil
}

func getEntityPersonality(entity interface{}) (string, error) {
	switch v := entity.(type) {
	case *models.NPC:
		return v.PersonalityPrompt, nil
	case *models.Owner:
		return v.LLMPromptContext, nil
	case *models.Questmaker:
		return v.LLMPromptContext, nil
	default:
		return "", fmt.Errorf("unknown entity type for personality")
	}
}

func getEntityTools(entity interface{}) ([]models.Tool, error) {
	switch v := entity.(type) {
	case *models.NPC:
		return v.AvailableTools, nil
	case *models.Owner:
		return v.AvailableTools, nil
	case *models.Questmaker:
		return v.AvailableTools, nil
	default:
		return nil, fmt.Errorf("unknown entity type for tools")
	}
}

func getEntityMemories(entity interface{}, playerID string) ([]string, error) {
	switch v := entity.(type) {
	case *models.NPC:
		if memories, ok := v.MemoriesAboutPlayers[playerID]; ok {
			return memories, nil
		}
		return nil, nil
	case *models.Owner:
		if memories, ok := v.MemoriesAboutPlayers[playerID]; ok {
			return memories, nil
		}
		return nil, nil
	case *models.Questmaker:
		if memories, ok := v.MemoriesAboutPlayers[playerID]; ok {
			return memories, nil
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown entity type for memories")
	}
}
