package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mud/internal/models"
)

// NPCDAL handles database operations for NPC entities.
type NPCDAL struct {
	db    *sql.DB
	Cache CacheInterface
}

// NewNPCDAL creates a new NPCDAL.
func NewNPCDAL(db *sql.DB, cache CacheInterface) *NPCDAL {
	return &NPCDAL{db: db, Cache: cache}
}

// CreateNPC inserts a new NPC into the database.
func (d *NPCDAL) CreateNPC(npc *models.NPC) error {
	inventoryJSON, err := json.Marshal(npc.Inventory)
	if err != nil {
		return fmt.Errorf("failed to marshal inventory: %w", err)
	}
	ownerIDsJSON, err := json.Marshal(npc.OwnerIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal owner IDs: %w", err)
	}
	memoriesJSON, err := json.Marshal(npc.MemoriesAboutPlayers)
	if err != nil {
		return fmt.Errorf("failed to marshal memories about players: %w", err)
	}
	availableToolsJSON, err := json.Marshal(npc.AvailableTools)
	if err != nil {
		return fmt.Errorf("failed to marshal available tools: %w", err)
	}

	query := `
	INSERT INTO NPCs (id, name, description, current_room_id, health, max_health, inventory, owner_ids, memories_about_players, personality_prompt, available_tools, behavior_state, reaction_threshold, race_id, profession_id)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = d.db.Exec(query,
		npc.ID,
		npc.Name,
		npc.Description,
		npc.CurrentRoomID,
		npc.Health,
		npc.MaxHealth,
		string(inventoryJSON),
		string(ownerIDsJSON),
		string(memoriesJSON),
		npc.PersonalityPrompt,
		string(availableToolsJSON),
		npc.BehaviorState,
		npc.ReactionThreshold,
		npc.RaceID,
		npc.ProfessionID,
	)
	if err != nil {
		return fmt.Errorf("failed to create NPC: %w", err)
	}
	d.Cache.Set(npc.ID, npc, 300) // Cache for 5 minutes
	return nil
}

// GetNPCByID retrieves an NPC by their ID.
func (d *NPCDAL) GetNPCByID(id string) (*models.NPC, error) {
	if cachedNPC, found := d.Cache.Get(id); found {
		if npc, ok := cachedNPC.(*models.NPC); ok {
			return npc, nil
		}
	}

	query := `SELECT id, name, description, current_room_id, health, max_health, inventory, owner_ids, memories_about_players, personality_prompt, available_tools, behavior_state, reaction_threshold, race_id, profession_id FROM NPCs WHERE id = ?`
	row := d.db.QueryRow(query, id)

	npc := &models.NPC{}
	var inventoryJSON, ownerIDsJSON, memoriesJSON, availableToolsJSON []byte
	err := row.Scan(
		&npc.ID,
		&npc.Name,
		&npc.Description,
		&npc.CurrentRoomID,
		&npc.Health,
		&npc.MaxHealth,
		&inventoryJSON,
		&ownerIDsJSON,
		&memoriesJSON,
		&npc.PersonalityPrompt,
		&availableToolsJSON,
		&npc.BehaviorState,
		&npc.ReactionThreshold,
		&npc.RaceID,
		&npc.ProfessionID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // NPC not found
		}
		return nil, fmt.Errorf("failed to get NPC by ID: %w", err)
	}

	if err := json.Unmarshal(inventoryJSON, &npc.Inventory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inventory: %w", err)
	}
	if err := json.Unmarshal(ownerIDsJSON, &npc.OwnerIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal owner IDs: %w", err)
	}
	if err := json.Unmarshal(memoriesJSON, &npc.MemoriesAboutPlayers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal memories about players: %w", err)
	}
	if err := json.Unmarshal(availableToolsJSON, &npc.AvailableTools); err != nil {
		return nil, fmt.Errorf("failed to unmarshal available tools: %w", err)
	}

	d.Cache.Set(npc.ID, npc, 300) // Cache for 5 minutes
	return npc, nil
}

// UpdateNPC updates an existing NPC in the database.
func (d *NPCDAL) UpdateNPC(npc *models.NPC) error {
	inventoryJSON, err := json.Marshal(npc.Inventory)
	if err != nil {
		return fmt.Errorf("failed to marshal inventory: %w", err)
	}
	ownerIDsJSON, err := json.Marshal(npc.OwnerIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal owner IDs: %w", err)
	}
	memoriesJSON, err := json.Marshal(npc.MemoriesAboutPlayers)
	if err != nil {
		return fmt.Errorf("failed to marshal memories about players: %w", err)
	}
	availableToolsJSON, err := json.Marshal(npc.AvailableTools)
	if err != nil {
		return fmt.Errorf("failed to marshal available tools: %w", err)
	}

	query := `
	UPDATE NPCs
	SET name = ?, description = ?, current_room_id = ?, health = ?, max_health = ?, inventory = ?, owner_ids = ?, memories_about_players = ?, personality_prompt = ?, available_tools = ?, behavior_state = ?, reaction_threshold = ?, race_id = ?, profession_id = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		npc.Name,
		npc.Description,
		npc.CurrentRoomID,
		npc.Health,
		npc.MaxHealth,
		string(inventoryJSON),
		string(ownerIDsJSON),
		string(memoriesJSON),
		npc.PersonalityPrompt,
		string(availableToolsJSON),
		npc.BehaviorState,
		npc.ReactionThreshold,
		npc.RaceID,
		npc.ProfessionID,
		npc.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update NPC: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("NPC with ID %s not found for update", npc.ID)
	}
	d.Cache.Delete(npc.ID) // Invalidate cache on update
	return nil
}

// DeleteNPC deletes an NPC from the database by their ID.
func (d *NPCDAL) DeleteNPC(id string) error {
	query := `DELETE FROM NPCs WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete NPC: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("NPC with ID %s not found for deletion", id)
	}
	d.Cache.Delete(id) // Invalidate cache on delete
	return nil
}

// GetNPCsByRoom retrieves all NPCs in a given room.
func (d *NPCDAL) GetNPCsByRoom(roomID string) ([]*models.NPC, error) {
	// For list queries, caching is more complex. For now, we won't cache list results.
	query := `SELECT id, name, description, current_room_id, health, max_health, inventory, owner_ids, memories_about_players, personality_prompt, available_tools, behavior_state, reaction_threshold, race_id, profession_id FROM NPCs WHERE current_room_id = ?`
	rows, err := d.db.Query(query, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get NPCs by room: %w", err)
	}
	defer rows.Close()

	var npcs []*models.NPC
	for rows.Next() {
		npc := &models.NPC{}
		var inventoryJSON, ownerIDsJSON, memoriesJSON, availableToolsJSON []byte
		err := rows.Scan(
			&npc.ID,
			&npc.Name,
			&npc.Description,
			&npc.CurrentRoomID,
			&npc.Health,
			&npc.MaxHealth,
			&inventoryJSON,
			&ownerIDsJSON,
			&memoriesJSON,
			&npc.PersonalityPrompt,
			&availableToolsJSON,
			&npc.BehaviorState,
			&npc.ReactionThreshold,
			&npc.RaceID,
			&npc.ProfessionID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan NPC row: %w", err)
		}

		if err := json.Unmarshal(inventoryJSON, &npc.Inventory); err != nil {
			return nil, fmt.Errorf("failed to unmarshal inventory for NPC %s: %w", npc.ID, err)
		}
		if err := json.Unmarshal(ownerIDsJSON, &npc.OwnerIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal owner IDs for NPC %s: %w", npc.ID, err)
		}
		if err := json.Unmarshal(memoriesJSON, &npc.MemoriesAboutPlayers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal memories about players for NPC %s: %w", npc.ID, err)
		}
		if err := json.Unmarshal(availableToolsJSON, &npc.AvailableTools); err != nil {
			return nil, fmt.Errorf("failed to unmarshal available tools for NPC %s: %w", npc.ID, err)
		}

		npcs = append(npcs, npc)
	}

	return npcs, nil
}

// GetNPCsByOwner retrieves all NPCs associated with a given owner.
func (d *NPCDAL) GetNPCsByOwner(ownerID string) ([]*models.NPC, error) {
	query := `SELECT id, name, description, current_room_id, health, max_health, inventory, owner_ids, memories_about_players, personality_prompt, available_tools, behavior_state, reaction_threshold, race_id, profession_id FROM NPCs WHERE INSTR(owner_ids, ?)`
	rows, err := d.db.Query(query, `"`+ownerID+`"`)
	if err != nil {
		return nil, fmt.Errorf("failed to get NPCs by owner: %w", err)
	}
	defer rows.Close()

	var npcs []*models.NPC
	for rows.Next() {
		npc := &models.NPC{}
		var inventoryJSON, ownerIDsJSON, memoriesJSON, availableToolsJSON []byte
		err := rows.Scan(
			&npc.ID,
			&npc.Name,
			&npc.Description,
			&npc.CurrentRoomID,
			&npc.Health,
			&npc.MaxHealth,
			&inventoryJSON,
			&ownerIDsJSON,
			&memoriesJSON,
			&npc.PersonalityPrompt,
			&availableToolsJSON,
			&npc.BehaviorState,
			&npc.ReactionThreshold,
			&npc.RaceID,
			&npc.ProfessionID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan NPC row: %w", err)
		}

		if err := json.Unmarshal(inventoryJSON, &npc.Inventory); err != nil {
			return nil, fmt.Errorf("failed to unmarshal inventory for NPC %s: %w", npc.ID, err)
		}
		if err := json.Unmarshal(ownerIDsJSON, &npc.OwnerIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal owner IDs for NPC %s: %w", npc.ID, err)
		}
		if err := json.Unmarshal(memoriesJSON, &npc.MemoriesAboutPlayers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal memories about players for NPC %s: %w", npc.ID, err)
		}
		if err := json.Unmarshal(availableToolsJSON, &npc.AvailableTools); err != nil {
			return nil, fmt.Errorf("failed to unmarshal available tools for NPC %s: %w", npc.ID, err)
		}

		npcs = append(npcs, npc)
	}

	return npcs, nil
}

// GetAllNPCs retrieves all NPCs from the database.
func (d *NPCDAL) GetAllNPCs() ([]*models.NPC, error) {
	query := `SELECT id, name, description, current_room_id, health, max_health, inventory, owner_ids, memories_about_players, personality_prompt, available_tools, behavior_state, reaction_threshold, race_id, profession_id FROM NPCs`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all NPCs: %w", err)
	}
	defer rows.Close()

	var npcs []*models.NPC
	for rows.Next() {
		npc := &models.NPC{}
		var inventoryJSON, ownerIDsJSON, memoriesJSON, availableToolsJSON []byte
		err := rows.Scan(
			&npc.ID,
			&npc.Name,
			&npc.Description,
			&npc.CurrentRoomID,
			&npc.Health,
			&npc.MaxHealth,
			&inventoryJSON,
			&ownerIDsJSON,
			&memoriesJSON,
			&npc.PersonalityPrompt,
			&availableToolsJSON,
			&npc.BehaviorState,
			&npc.ReactionThreshold,
			&npc.RaceID,
			&npc.ProfessionID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan NPC: %w", err)
		}

		if err := json.Unmarshal(inventoryJSON, &npc.Inventory); err != nil {
			return nil, fmt.Errorf("failed to unmarshal inventory for NPC %s: %w", npc.ID, err)
		}
		if err := json.Unmarshal(ownerIDsJSON, &npc.OwnerIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal owner IDs for NPC %s: %w", npc.ID, err)
		}
		if err := json.Unmarshal(memoriesJSON, &npc.MemoriesAboutPlayers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal memories about players for NPC %s: %w", npc.ID, err)
		}
		if err := json.Unmarshal(availableToolsJSON, &npc.AvailableTools); err != nil {
			return nil, fmt.Errorf("failed to unmarshal available tools for NPC %s: %w", npc.ID, err)
		}

		npcs = append(npcs, npc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through NPCs: %w", err)
	}

	return npcs, nil
}