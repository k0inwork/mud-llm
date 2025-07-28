package server

import (
	"context"
	"fmt"
	"mud/internal/dal"
	"mud/internal/llm"
	"mud/internal/models"
)

type ToolDispatcher struct {
	dal *dal.DAL
}

func NewToolDispatcher(dal *dal.DAL) *ToolDispatcher {
	return &ToolDispatcher{dal: dal}
}

type ToolCall struct {
	ToolName   string                 `json:"tool_name"`
	Parameters map[string]interface{} `json:"parameters"`
}

func (td *ToolDispatcher) Dispatch(ctx context.Context, player *models.PlayerCharacter, entity interface{}, toolCalls []llm.ToolCall) error {
	for _, call := range toolCalls {
		switch call.ToolName {
		case "NPC_memorize":
			if err := td.handleNPCMemorize(player, entity, call.Parameters); err != nil {
				return err
			}
		case "OWNER_memorize":
			if err := td.handleOwnerMemorize(player, entity, call.Parameters); err != nil {
				return err
			}
		case "OWNER_memorize_dependables":
			if err := td.handleOwnerMemorizeDependables(player, entity, call.Parameters); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown tool: %s", call.ToolName)
		}
	}

	return nil
}

func (td *ToolDispatcher) handleNPCMemorize(player *models.PlayerCharacter, entity interface{}, params map[string]interface{}) error {
	npcID, ok := params["npc_id"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid npc_id")
	}
	// Use player.ID directly if available, otherwise get from params
	var playerID string
	if player != nil {
		playerID = player.ID
	} else {
		pid, ok := params["player_id"].(string)
		if !ok {
			return fmt.Errorf("missing or invalid player_id")
		}
		playerID = pid
	}

	memory, ok := params["memory_string"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid memory_string")
	}

	npc, err := td.dal.NpcDAL.GetNPCByID(npcID)
	if err != nil {
		return err
	}
	if npc == nil {
		return fmt.Errorf("npc not found: %s", npcID)
	}

	if npc.MemoriesAboutPlayers == nil {
		npc.MemoriesAboutPlayers = make(map[string][]string)
	}
	npc.MemoriesAboutPlayers[playerID] = append(npc.MemoriesAboutPlayers[playerID], memory)

	return td.dal.NpcDAL.UpdateNPC(npc)
}

func (td *ToolDispatcher) handleOwnerMemorize(player *models.PlayerCharacter, entity interface{}, params map[string]interface{}) error {
	ownerID, ok := params["owner_id"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid owner_id")
	}
	// Use player.ID directly if available, otherwise get from params
	var playerID string
	if player != nil {
		playerID = player.ID
	} else {
		pid, ok := params["player_id"].(string)
		if !ok {
			return fmt.Errorf("missing or invalid player_id")
		}
		playerID = pid
	}
	memory, ok := params["memory_string"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid memory_string")
	}

	owner, err := td.dal.OwnerDAL.GetOwnerByID(ownerID)
	if err != nil {
		return err
	}
	if owner == nil {
		return fmt.Errorf("owner not found: %s", ownerID)
	}

	if owner.MemoriesAboutPlayers == nil {
		owner.MemoriesAboutPlayers = make(map[string][]string)
	}
	owner.MemoriesAboutPlayers[playerID] = append(owner.MemoriesAboutPlayers[playerID], memory)

	return td.dal.OwnerDAL.UpdateOwner(owner)
}

func (td *ToolDispatcher) handleOwnerMemorizeDependables(player *models.PlayerCharacter, entity interface{}, params map[string]interface{}) error {
	ownerID, ok := params["owner_id"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid owner_id")
	}
	// Use player.ID directly if available, otherwise get from params
	var playerID string
	if player != nil {
		playerID = player.ID
	} else {
		pid, ok := params["player_id"].(string)
		if !ok {
			return fmt.Errorf("missing or invalid player_id")
		}
		playerID = pid
	}
	memory, ok := params["memory_string"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid memory_string")
	}

	npcs, err := td.dal.NpcDAL.GetNPCsByOwner(ownerID)
	if err != nil {
		return err
	}

	for _, npc := range npcs {
		if npc.MemoriesAboutPlayers == nil {
			npc.MemoriesAboutPlayers = make(map[string][]string)
		}
		npc.MemoriesAboutPlayers[playerID] = append(npc.MemoriesAboutPlayers[playerID], memory)
		if err := td.dal.NpcDAL.UpdateNPC(npc); err != nil {
			return err
		}
	}

	return nil
}


