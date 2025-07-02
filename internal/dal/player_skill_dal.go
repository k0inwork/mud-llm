package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// PlayerSkillDAL handles database operations for PlayerSkill entities.
type PlayerSkillDAL struct {
	db *sql.DB
}

// NewPlayerSkillDAL creates a new PlayerSkillDAL.
func NewPlayerSkillDAL(db *sql.DB) *PlayerSkillDAL {
	return &PlayerSkillDAL{db: db}
}

// CreatePlayerSkill inserts a new player skill into the database.
func (d *PlayerSkillDAL) CreatePlayerSkill(ps *models.PlayerSkill) error {
	query := `
	INSERT INTO PlayerSkills (player_id, skill_id, percentage, granted_by_entity_type, granted_by_entity_id)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		ps.PlayerID,
		ps.SkillID,
		ps.Percentage,
		ps.GrantedByEntityType,
		ps.GrantedByEntityID,
	)
	if err != nil {
		return fmt.Errorf("failed to create player skill: %w", err)
	}
	return nil
}

// GetPlayerSkill retrieves a player skill by player and skill ID.
func (d *PlayerSkillDAL) GetPlayerSkill(playerID, skillID string) (*models.PlayerSkill, error) {
	query := `SELECT player_id, skill_id, percentage, granted_by_entity_type, granted_by_entity_id FROM PlayerSkills WHERE player_id = ? AND skill_id = ?`
	row := d.db.QueryRow(query, playerID, skillID)

	ps := &models.PlayerSkill{}
	err := row.Scan(
		&ps.PlayerID,
		&ps.SkillID,
		&ps.Percentage,
		&ps.GrantedByEntityType,
		&ps.GrantedByEntityID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Player skill not found
		}
		return nil, fmt.Errorf("failed to get player skill: %w", err)
	}

	return ps, nil
}

// UpdatePlayerSkill updates an existing player skill in the database.
func (d *PlayerSkillDAL) UpdatePlayerSkill(ps *models.PlayerSkill) error {
	query := `
	UPDATE PlayerSkills
	SET percentage = ?, granted_by_entity_type = ?, granted_by_entity_id = ?
	WHERE player_id = ? AND skill_id = ?
	`

	result, err := d.db.Exec(query,
		ps.Percentage,
		ps.GrantedByEntityType,
		ps.GrantedByEntityID,
		ps.PlayerID,
		ps.SkillID,
	)
	if err != nil {
		return fmt.Errorf("failed to update player skill: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("player skill for player %s and skill %s not found for update", ps.PlayerID, ps.SkillID)
	}

	return nil
}

// DeletePlayerSkill deletes a player skill from the database.
func (d *PlayerSkillDAL) DeletePlayerSkill(playerID, skillID string) error {
	query := `DELETE FROM PlayerSkills WHERE player_id = ? AND skill_id = ?`
	result, err := d.db.Exec(query, playerID, skillID)
	if err != nil {
		return fmt.Errorf("failed to delete player skill: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("player skill for player %s and skill %s not found for deletion", playerID, skillID)
	}

	return nil
}
