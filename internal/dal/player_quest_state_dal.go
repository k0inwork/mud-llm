package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
	
)

// PlayerQuestStateDAL handles database operations for PlayerQuestState entities.
type PlayerQuestStateDAL struct {
	db    *sql.DB
	cache CacheInterface
}

func (d *PlayerQuestStateDAL) Cache() CacheInterface {
	return d.cache
}

// NewPlayerQuestStateDAL creates a new PlayerQuestStateDAL.
func NewPlayerQuestStateDAL(db *sql.DB, cache CacheInterface) *PlayerQuestStateDAL {
	return &PlayerQuestStateDAL{db: db, cache: cache}
}

// CreatePlayerQuestState inserts a new player quest state into the database.
func (d *PlayerQuestStateDAL) CreatePlayerQuestState(pqs *models.PlayerQuestState) error {
	query := `
	INSERT INTO PlayerQuestStates (player_id, quest_id, current_progress, last_action_timestamp, questmaker_influence_accumulated, status)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		pqs.PlayerID,
		pqs.QuestID,
		pqs.CurrentProgress,
		pqs.LastActionTimestamp,
		pqs.QuestmakerInfluenceAccumulated,
		pqs.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create player quest state: %w", err)
	}
	return nil
}

// GetPlayerQuestState retrieves a player quest state by player and quest ID.
func (d *PlayerQuestStateDAL) GetPlayerQuestStateByID(playerID, questID string) (*models.PlayerQuestState, error) {
	query := `SELECT player_id, quest_id, current_progress, last_action_timestamp, questmaker_influence_accumulated, status FROM PlayerQuestStates WHERE player_id = ? AND quest_id = ?`
	row := d.db.QueryRow(query, playerID, questID)

	pqs := &models.PlayerQuestState{}
	err := row.Scan(
		&pqs.PlayerID,
		&pqs.QuestID,
		&pqs.CurrentProgress,
		&pqs.LastActionTimestamp,
		&pqs.QuestmakerInfluenceAccumulated,
		&pqs.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Player quest state not found
		}
		return nil, fmt.Errorf("failed to get player quest state: %w", err)
	}

	return pqs, nil
}

// UpdatePlayerQuestState updates an existing player quest state in the database.
func (d *PlayerQuestStateDAL) UpdatePlayerQuestState(pqs *models.PlayerQuestState) error {
	query := `
	UPDATE PlayerQuestStates
	SET current_progress = ?, last_action_timestamp = ?, questmaker_influence_accumulated = ?, status = ?
	WHERE player_id = ? AND quest_id = ?
	`

	result, err := d.db.Exec(query,
		pqs.CurrentProgress,
		pqs.LastActionTimestamp,
		pqs.QuestmakerInfluenceAccumulated,
		pqs.Status,
		pqs.PlayerID,
		pqs.QuestID,
	)
	if err != nil {
		return fmt.Errorf("failed to update player quest state: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("player quest state for player %s and quest %s not found for update", pqs.PlayerID, pqs.QuestID)
	}

	return nil
}

// DeletePlayerQuestState deletes a player quest state from the database.
func (d *PlayerQuestStateDAL) DeletePlayerQuestState(playerID, questID string) error {
	query := `DELETE FROM PlayerQuestStates WHERE player_id = ? AND quest_id = ?`
	result, err := d.db.Exec(query, playerID, questID)
	if err != nil {
		return fmt.Errorf("failed to delete player quest state: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("player quest state for player %s and quest %s not found for deletion", playerID, questID)
	}

	return nil
}

// GetAllPlayerQuestStates retrieves all player quest states from the database.
func (d *PlayerQuestStateDAL) GetAllPlayerQuestStates() ([]*models.PlayerQuestState, error) {
	query := `SELECT player_id, quest_id, current_progress, last_action_timestamp, questmaker_influence_accumulated, status FROM PlayerQuestStates`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all player quest states: %w", err)
	}
	defer rows.Close()

	var playerQuestStates []*models.PlayerQuestState
	for rows.Next() {
		pqs := &models.PlayerQuestState{}
		err := rows.Scan(
			&pqs.PlayerID,
			&pqs.QuestID,
			&pqs.CurrentProgress,
			&pqs.LastActionTimestamp,
			&pqs.QuestmakerInfluenceAccumulated,
			&pqs.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player quest state: %w", err)
		}
		playerQuestStates = append(playerQuestStates, pqs)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through player quest states: %w", err)
	}

	return playerQuestStates, nil
}