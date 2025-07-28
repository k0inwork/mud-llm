package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// QuestDAL handles database operations for Quest entities.
type QuestDAL struct {
	db    *sql.DB
	cache CacheInterface
}

func (d *QuestDAL) Cache() CacheInterface {
	return d.cache
}

// NewQuestDAL creates a new QuestDAL.
func NewQuestDAL(db *sql.DB, cache CacheInterface) *QuestDAL {
	return &QuestDAL{db: db, cache: cache}
}

// CreateQuest inserts a new quest into the database.
func (d *QuestDAL) CreateQuest(quest *models.Quest) error {
	query := `
	INSERT INTO Quests (id, name, description, quest_owner_id, questmaker_id, influence_points_map, objectives, rewards)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		quest.ID,
		quest.Name,
		quest.Description,
		quest.QuestOwnerID,
		quest.QuestmakerID,
		quest.InfluencePointsMap,
		quest.Objectives,
		quest.Rewards,
	)
	if err != nil {
		return fmt.Errorf("failed to create quest: %w", err)
	}
	return nil
}

// GetQuestByID retrieves a quest by its ID.
func (d *QuestDAL) GetQuestByID(id string) (*models.Quest, error) {
	query := `SELECT id, name, description, quest_owner_id, questmaker_id, influence_points_map, objectives, rewards FROM Quests WHERE id = ?`
	row := d.db.QueryRow(query, id)

	quest := &models.Quest{}
	err := row.Scan(
		&quest.ID,
		&quest.Name,
		&quest.Description,
		&quest.QuestOwnerID,
		&quest.QuestmakerID,
		&quest.InfluencePointsMap,
		&quest.Objectives,
		&quest.Rewards,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Quest not found
		}
		return nil, fmt.Errorf("failed to get quest by ID: %w", err)
	}

	return quest, nil
}

// UpdateQuest updates an existing quest in the database.
func (d *QuestDAL) UpdateQuest(quest *models.Quest) error {
	query := `
	UPDATE Quests
	SET name = ?, description = ?, quest_owner_id = ?, questmaker_id = ?, influence_points_map = ?, objectives = ?, rewards = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		quest.Name,
		quest.Description,
		quest.QuestOwnerID,
		quest.QuestmakerID,
		quest.InfluencePointsMap,
		quest.Objectives,
		quest.Rewards,
		quest.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update quest: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("quest with ID %s not found for update", quest.ID)
	}

	return nil
}

// DeleteQuest deletes a quest from the database by its ID.
func (d *QuestDAL) DeleteQuest(id string) error {
	query := `DELETE FROM Quests WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete quest: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("quest with ID %s not found for deletion", id)
	}

	return nil
}

// GetAllQuests retrieves all quests from the database.
func (d *QuestDAL) GetAllQuests() ([]*models.Quest, error) {
	query := `SELECT id, name, description, quest_owner_id, questmaker_id, influence_points_map, objectives, rewards FROM Quests`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all quests: %w", err)
	}
	defer rows.Close()

	var quests []*models.Quest
	for rows.Next() {
		quest := &models.Quest{}
		err := rows.Scan(
			&quest.ID,
			&quest.Name,
			&quest.Description,
			&quest.QuestOwnerID,
			&quest.QuestmakerID,
			&quest.InfluencePointsMap,
			&quest.Objectives,
			&quest.Rewards,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quest: %w", err)
		}
		quests = append(quests, quest)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through quests: %w", err)
	}

	return quests, nil
}