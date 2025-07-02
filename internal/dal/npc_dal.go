package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// NPCDAL handles database operations for NPC entities.
type NPCDAL struct {
	db    *sql.DB
	cache *Cache
}

// NewNPCDAL creates a new NPCDAL.
func NewNPCDAL(db *sql.DB) *NPCDAL {
	return &NPCDAL{db: db, cache: NewCache()}
}

// CreateNPC inserts a new NPC into the database.
func (d *NPCDAL) CreateNPC(npc *models.NPC) error {
	query := `
	INSERT INTO NPCs (id, name, description, current_room_id, health, max_health, inventory, owner_ids, memories_about_players, personality_prompt, available_tools, behavior_state)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		npc.ID,
		npc.Name,
		npc.Description,
		npc.CurrentRoomID,
		npc.Health,
		npc.MaxHealth,
		npc.Inventory,
		npc.OwnerIDs,
		npc.MemoriesAboutPlayers,
		npc.PersonalityPrompt,
		npc.AvailableTools,
		npc.BehaviorState,
	)
	if err != nil {
		return fmt.Errorf("failed to create NPC: %w", err)
	}
	d.cache.Set(npc.ID, npc, 300) // Cache for 5 minutes
	return nil
}

// GetNPCByID retrieves an NPC by their ID.
func (d *NPCDAL) GetNPCByID(id string) (*models.NPC, error) {
	if cachedNPC, found := d.cache.Get(id); found {
		if npc, ok := cachedNPC.(*models.NPC); ok {
			return npc, nil
		}
	}

	query := `SELECT id, name, description, current_room_id, health, max_health, inventory, owner_ids, memories_about_players, personality_prompt, available_tools, behavior_state FROM NPCs WHERE id = ?`
	row := d.db.QueryRow(query, id)

	npc := &models.NPC{}
	err := row.Scan(
		&npc.ID,
		&npc.Name,
		&npc.Description,
		&npc.CurrentRoomID,
		&npc.Health,
		&npc.MaxHealth,
		&npc.Inventory,
		&npc.OwnerIDs,
		&npc.MemoriesAboutPlayers,
		&npc.PersonalityPrompt,
		&npc.AvailableTools,
		&npc.BehaviorState,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // NPC not found
		}
		return nil, fmt.Errorf("failed to get NPC by ID: %w", err)
	}

	d.cache.Set(npc.ID, npc, 300) // Cache for 5 minutes
	return npc, nil
}

// UpdateNPC updates an existing NPC in the database.
func (d *NPCDAL) UpdateNPC(npc *models.NPC) error {
	query := `
	UPDATE NPCs
	SET name = ?, description = ?, current_room_id = ?, health = ?, max_health = ?, inventory = ?, owner_ids = ?, memories_about_players = ?, personality_prompt = ?, available_tools = ?, behavior_state = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		npc.Name,
		npc.Description,
		npc.CurrentRoomID,
		npc.Health,
		npc.MaxHealth,
		npc.Inventory,
		npc.OwnerIDs,
		npc.MemoriesAboutPlayers,
		npc.PersonalityPrompt,
		npc.AvailableTools,
		npc.BehaviorState,
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
	d.cache.Delete(npc.ID) // Invalidate cache on update
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
	d.cache.Delete(id) // Invalidate cache on delete
	return nil
}

// GetNPCsByRoom retrieves all NPCs in a given room.
func (d *NPCDAL) GetNPCsByRoom(roomID string) ([]*models.NPC, error) {
	// For list queries, caching is more complex. For now, we won't cache list results.
	query := `SELECT id, name, description, current_room_id, health, max_health, inventory, owner_ids, memories_about_players, personality_prompt, available_tools, behavior_state FROM NPCs WHERE current_room_id = ?`
	rows, err := d.db.Query(query, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get NPCs by room: %w", err)
	}
	defer rows.Close()

	var npcs []*models.NPC
	for rows.Next() {
		npc := &models.NPC{}
		err := rows.Scan(
			&npc.ID,
			&npc.Name,
			&npc.Description,
			&npc.CurrentRoomID,
			&npc.Health,
			&npc.MaxHealth,
			&npc.Inventory,
			&npc.OwnerIDs,
			&npc.MemoriesAboutPlayers,
			&npc.PersonalityPrompt,
			&npc.AvailableTools,
			&npc.BehaviorState,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan NPC row: %w", err)
		}
		npcs = append(npcs, npc)
	}

	return npcs, nil
}

// GetNPCsByOwner retrieves all NPCs associated with a given owner.
func (d *NPCDAL) GetNPCsByOwner(ownerID string) ([]*models.NPC, error) {
	query := `SELECT n.id, n.name, n.description, n.current_room_id, n.health, n.max_health, n.inventory, n.owner_ids, n.memories_about_players, n.personality_prompt, n.available_tools, n.behavior_state FROM NPCs n, json_each(n.owner_ids) WHERE json_each.value = ?`
	rows, err := d.db.Query(query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get NPCs by owner: %w", err)
	}
	defer rows.Close()

	var npcs []*models.NPC
	for rows.Next() {
		npc := &models.NPC{}
		err := rows.Scan(
			&npc.ID,
			&npc.Name,
			&npc.Description,
			&npc.CurrentRoomID,
			&npc.Health,
			&npc.MaxHealth,
			&npc.Inventory,
			&npc.OwnerIDs,
			&npc.MemoriesAboutPlayers,
			&npc.PersonalityPrompt,
			&npc.AvailableTools,
			&npc.BehaviorState,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan NPC row: %w", err)
		}
		npcs = append(npcs, npc)
	}

	return npcs, nil
}

// GetAllNPCs retrieves all NPCs from the database.
func (d *NPCDAL) GetAllNPCs() ([]*models.NPC, error) {
	query := `SELECT id, name, description, current_room_id, health, max_health, inventory, owner_ids, memories_about_players, personality_prompt, available_tools, behavior_state FROM NPCs`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all NPCs: %w", err)
	}
	defer rows.Close()

	var npcs []*models.NPC
	for rows.Next() {
		npc := &models.NPC{}
		err := rows.Scan(
			&npc.ID,
			&npc.Name,
			&npc.Description,
			&npc.CurrentRoomID,
			&npc.Health,
			&npc.MaxHealth,
			&npc.Inventory,
			&npc.OwnerIDs,
			&npc.MemoriesAboutPlayers,
			&npc.PersonalityPrompt,
			&npc.AvailableTools,
			&npc.BehaviorState,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan NPC: %w", err)
		}
		npcs = append(npcs, npc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through NPCs: %w", err)
	}

	return npcs, nil
}
