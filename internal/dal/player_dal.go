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
	Cache *Cache
}

// NewPlayerDAL creates a new PlayerDAL.
func NewPlayerDAL(db *sql.DB) *PlayerDAL {
	return &PlayerDAL{db: db, Cache: NewCache()}
}

// CreatePlayer inserts a new player into the database.
func (d *PlayerDAL) CreatePlayer(player *models.Player) error {
	inventoryJSON, err := json.Marshal(player.Inventory)
	if err != nil {
		return fmt.Errorf("failed to marshal inventory: %w", err)
	}
	visitedRoomsJSON, err := json.Marshal(player.VisitedRoomIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal visited rooms: %w", err)
	}

	query := `
	INSERT INTO Players (id, name, race_id, profession_id, current_room_id, health, max_health, inventory, visited_room_ids, created_at, last_login_at, last_logout_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = d.db.Exec(query,
		player.ID,
		player.Name,
		player.RaceID,
		player.ProfessionID,
		player.CurrentRoomID,
		player.Health,
		player.MaxHealth,
		string(inventoryJSON),
		string(visitedRoomsJSON),
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

	if err := json.Unmarshal(inventoryJSON, &player.Inventory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player inventory: %w", err)
	}
	if err := json.Unmarshal(visitedRoomsJSON, &player.VisitedRoomIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal visited rooms: %w", err)
	}

	return player, nil
}

// UpdatePlayer updates an existing player in the database.
func (d *PlayerDAL) UpdatePlayer(player *models.Player) error {
	inventoryJSON, err := json.Marshal(player.Inventory)
	if err != nil {
		return fmt.Errorf("failed to marshal inventory: %w", err)
	}
	visitedRoomsJSON, err := json.Marshal(player.VisitedRoomIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal visited rooms: %w", err)
	}

	query := `
	UPDATE Players
	SET name = ?, race_id = ?, profession_id = ?, current_room_id = ?, health = ?, max_health = ?, inventory = ?, visited_room_ids = ?, last_login_at = ?, last_logout_at = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		player.Name,
		player.RaceID,
		player.ProfessionID,
		player.CurrentRoomID,
		player.Health,
		player.MaxHealth,
		string(inventoryJSON),
		string(visitedRoomsJSON),
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

// GetPlayerInventory retrieves all items in a player's inventory.
func (d *PlayerDAL) GetPlayerInventory(playerID string) ([]*models.Item, error) {
	player, err := d.GetPlayerByID(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player for inventory: %w", err)
	}
	if player == nil {
		return nil, fmt.Errorf("player with ID %s not found", playerID)
	}

	var items []*models.Item
	itemDAL := NewItemDAL(d.db) // Create a new ItemDAL instance

	for _, itemID := range player.Inventory {
		item, err := itemDAL.GetItemByID(itemID)
		if err != nil {
			return nil, fmt.Errorf("failed to get item %s from inventory: %w", itemID, err)
		}
		if item != nil {
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
		if err := json.Unmarshal(inventoryJSON, &player.Inventory); err != nil {
			return nil, fmt.Errorf("failed to unmarshal inventory for player %s: %w", player.ID, err)
		}
		if err := json.Unmarshal(visitedRoomsJSON, &player.VisitedRoomIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal visited rooms for player %s: %w", player.ID, err)
		}
		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through players: %w", err)
	}

	return players, nil
}
