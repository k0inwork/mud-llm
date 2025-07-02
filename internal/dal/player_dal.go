package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mud/internal/models"
)

// PlayerDAL handles database operations for Player entities.
type PlayerDAL struct {
	db    *sql.DB
	cache *Cache
}

// NewPlayerDAL creates a new PlayerDAL.
func NewPlayerDAL(db *sql.DB) *PlayerDAL {
	return &PlayerDAL{db: db, cache: NewCache()}
}

// CreatePlayer inserts a new player into the database.
func (d *PlayerDAL) CreatePlayer(player *models.Player) error {
	query := `
	INSERT INTO Players (id, name, race_id, profession_id, current_room_id, health, max_health, inventory, visited_room_ids, created_at, last_login_at, last_logout_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	inventoryJSON := []byte(player.Inventory)
	visitedRoomsJSON := []byte(player.VisitedRoomIDs)

	_, err := d.db.Exec(query,
		player.ID,
		player.Name,
		player.RaceID,
		player.ProfessionID,
		player.CurrentRoomID,
		player.Health,
		player.MaxHealth,
		inventoryJSON,
		visitedRoomsJSON,
		player.CreatedAt,
		player.LastLoginAt,
		player.LastLogoutAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create player: %w", err)
	}
	return nil
}

// GetPlayerByID retrieves a player by their ID.
func (d *PlayerDAL) GetPlayerByID(id string) (*models.Player, error) {
	query := `SELECT id, name, race_id, profession_id, current_room_id, health, max_health, inventory, visited_room_ids, created_at, last_login_at, last_logout_at FROM Players WHERE id = ?`
	row := d.db.QueryRow(query, id)

	player := &models.Player{}
	var inventoryJSON, visitedRoomsJSON []byte

	err := row.Scan(
		&player.ID,
		&player.Name,
		&player.RaceID,
		&player.ProfessionID,
		&player.CurrentRoomID,
		&player.Health,
		&player.MaxHealth,
		&inventoryJSON,
		&visitedRoomsJSON,
		&player.CreatedAt,
		&player.LastLoginAt,
		&player.LastLogoutAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Player not found
		}
		return nil, fmt.Errorf("failed to get player by ID: %w", err)
	}

	// Unmarshal JSON fields
	// Assuming Inventory and VisitedRoomIDs are []string or similar in the Go struct
	// For now, keep them as string in struct and handle JSON directly here.
	// In a real scenario, you'd define specific types for these JSON fields.
	player.Inventory = string(inventoryJSON)
	player.VisitedRoomIDs = string(visitedRoomsJSON)

	return player, nil
}

// UpdatePlayer updates an existing player in the database.
func (d *PlayerDAL) UpdatePlayer(player *models.Player) error {
	query := `
	UPDATE Players
	SET name = ?, race_id = ?, profession_id = ?, current_room_id = ?, health = ?, max_health = ?, inventory = ?, visited_room_ids = ?, last_login_at = ?, last_logout_at = ?
	WHERE id = ?
	`

	inventoryJSON := []byte(player.Inventory)
	visitedRoomsJSON := []byte(player.VisitedRoomIDs)

	result, err := d.db.Exec(query,
		player.Name,
		player.RaceID,
		player.ProfessionID,
		player.CurrentRoomID,
		player.Health,
		player.MaxHealth,
		inventoryJSON,
		visitedRoomsJSON,
		player.LastLoginAt,
		player.LastLogoutAt,
		player.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("player with ID %s not found for update", player.ID)
	}

	return nil
}

// DeletePlayer deletes a player from the database by their ID.
func (d *PlayerDAL) DeletePlayer(id string) error {
	query := `DELETE FROM Players WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("player with ID %s not found for deletion", id)
	}

	return nil
}

// inventoryItem represents a single item entry in the player's inventory JSON.
type inventoryItem struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

// GetPlayerInventory retrieves all items in a player's inventory.
func (d *PlayerDAL) GetPlayerInventory(playerID string) ([]*models.Item, error) {
	player, err := d.GetPlayerByID(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player for inventory: %w", err)
	}
	if player == nil {
		return nil, fmt.Errorf("player with ID %s not found", playerID)
	}

	var playerInventory []inventoryItem
	if err := json.Unmarshal([]byte(player.Inventory), &playerInventory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player inventory JSON: %w", err)
	}

	var items []*models.Item
	itemDAL := NewItemDAL(d.db) // Create a new ItemDAL instance

	for _, invItem := range playerInventory {
		item, err := itemDAL.GetItemByID(invItem.ItemID)
		if err != nil {
			return nil, fmt.Errorf("failed to get item %s from inventory: %w", invItem.ItemID, err)
		}
		if item != nil {
			// For now, we just return the item definition. Quantity can be handled by game logic.
			items = append(items, item)
		}
	}

	return items, nil
}

// GetAllPlayers retrieves all players from the database.
func (d *PlayerDAL) GetAllPlayers() ([]*models.Player, error) {
	query := `SELECT id, name, race_id, profession_id, current_room_id, health, max_health, inventory, visited_room_ids, created_at, last_login_at, last_logout_at FROM Players`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all players: %w", err)
	}
	defer rows.Close()

	var players []*models.Player
	for rows.Next() {
		player := &models.Player{}
		var inventoryJSON, visitedRoomsJSON []byte
		err := rows.Scan(
			&player.ID,
			&player.Name,
			&player.RaceID,
			&player.ProfessionID,
			&player.CurrentRoomID,
			&player.Health,
			&player.MaxHealth,
			&inventoryJSON,
			&visitedRoomsJSON,
			&player.CreatedAt,
			&player.LastLoginAt,
			&player.LastLogoutAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player: %w", err)
		}
		player.Inventory = string(inventoryJSON)
		player.VisitedRoomIDs = string(visitedRoomsJSON)
		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through players: %w", err)
	}

	return players, nil
}
