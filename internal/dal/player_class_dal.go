package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// PlayerClassDAL handles database operations for PlayerClass entities.
type PlayerClassDAL struct {
	db    *sql.DB
	cache CacheInterface
}

func (d *PlayerClassDAL) Cache() CacheInterface {
	return d.cache
}

// NewPlayerClassDAL creates a new PlayerClassDAL.
func NewPlayerClassDAL(db *sql.DB, cache CacheInterface) *PlayerClassDAL {
	return &PlayerClassDAL{db: db, cache: cache}
}

// CreatePlayerClass inserts a new player class entry into the database.
func (d *PlayerClassDAL) CreatePlayerClass(pc *models.PlayerClass) error {
	query := `
	INSERT INTO PlayerClasses (player_id, class_id, level, experience)
	VALUES (?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		pc.PlayerID,
		pc.ClassID,
		pc.Level,
		pc.Experience,
	)
	if err != nil {
		return fmt.Errorf("failed to create player class: %w", err)
	}
	return nil
}

// GetPlayerClass retrieves a player class entry by player and class ID.
func (d *PlayerClassDAL) GetPlayerClassByID(playerID, classID string) (*models.PlayerClass, error) {
	query := `SELECT player_id, class_id, level, experience FROM PlayerClasses WHERE player_id = ? AND class_id = ?`
	row := d.db.QueryRow(query, playerID, classID)

	pc := &models.PlayerClass{}
	err := row.Scan(
		&pc.PlayerID,
		&pc.ClassID,
		&pc.Level,
		&pc.Experience,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Player class not found
		}
		return nil, fmt.Errorf("failed to get player class: %w", err)
	}

	return pc, nil
}

// UpdatePlayerClass updates an existing player class entry in the database.
func (d *PlayerClassDAL) UpdatePlayerClass(pc *models.PlayerClass) error {
	query := `
	UPDATE PlayerClasses
	SET level = ?, experience = ?
	WHERE player_id = ? AND class_id = ?
	`

	result, err := d.db.Exec(query,
		pc.Level,
		pc.Experience,
		pc.PlayerID,
		pc.ClassID,
	)
	if err != nil {
		return fmt.Errorf("failed to update player class: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("player class for player %s and class %s not found for update", pc.PlayerID, pc.ClassID)
	}

	return nil
}

// DeletePlayerClass deletes a player class entry from the database.
func (d *PlayerClassDAL) DeletePlayerClass(playerID, classID string) error {
	query := `DELETE FROM PlayerClasses WHERE player_id = ? AND class_id = ?`
	result, err := d.db.Exec(query, playerID, classID)
	if err != nil {
		return fmt.Errorf("failed to delete player class: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("player class for player %s and class %s not found for deletion", playerID, classID)
	}

	return nil
}

// GetAllPlayerClasses retrieves all player class entries from the database.
func (d *PlayerClassDAL) GetAllPlayerClasses() ([]*models.PlayerClass, error) {
	query := `SELECT player_id, class_id, level, experience FROM PlayerClasses`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all player classes: %w", err)
	}
	defer rows.Close()

	var playerClasses []*models.PlayerClass
	for rows.Next() {
		pc := &models.PlayerClass{}
		err := rows.Scan(
			&pc.PlayerID,
			&pc.ClassID,
			&pc.Level,
			&pc.Experience,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player class: %w", err)
		}
		playerClasses = append(playerClasses, pc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through player classes: %w", err)
	}

	return playerClasses, nil
}