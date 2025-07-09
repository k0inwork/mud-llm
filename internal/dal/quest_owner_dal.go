package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// QuestOwnerDAL handles database operations for QuestOwner entities.
type QuestOwnerDAL struct {
	db    *sql.DB
	Cache *Cache
}

// NewQuestOwnerDAL creates a new QuestOwnerDAL.
func NewQuestOwnerDAL(db *sql.DB) *QuestOwnerDAL {
	return &QuestOwnerDAL{db: db, Cache: NewCache()}
}

// CreateQuestOwner inserts a new quest owner into the database.
func (d *QuestOwnerDAL) CreateQuestOwner(qo *models.QuestOwner) error {
	query := `
	INSERT INTO QuestOwners (id, name, description, llm_prompt_context, current_influence_budget, max_influence_budget, budget_regen_rate, associated_questmaker_ids)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		qo.ID,
		qo.Name,
		qo.Description,
		qo.LLMPromptContext,
		qo.CurrentInfluenceBudget,
		qo.MaxInfluenceBudget,
		qo.BudgetRegenRate,
		qo.AssociatedQuestmakerIDs,
	)
	if err != nil {
		return fmt.Errorf("failed to create quest owner: %w", err)
	}
	d.Cache.Set(qo.ID, qo, 300) // Cache for 5 minutes
	return nil
}

// GetQuestOwnerByID retrieves a quest owner by its ID.
func (d *QuestOwnerDAL) GetQuestOwnerByID(id string) (*models.QuestOwner, error) {
	if cachedQuestOwner, found := d.Cache.Get(id); found {
		if qo, ok := cachedQuestOwner.(*models.QuestOwner); ok {
			return qo, nil
		}
	}

	query := `SELECT id, name, description, llm_prompt_context, current_influence_budget, max_influence_budget, budget_regen_rate, associated_questmaker_ids FROM QuestOwners WHERE id = ?`
	row := d.db.QueryRow(query, id)

	qo := &models.QuestOwner{}
	err := row.Scan(
		&qo.ID,
		&qo.Name,
		&qo.Description,
		&qo.LLMPromptContext,
		&qo.CurrentInfluenceBudget,
		&qo.MaxInfluenceBudget,
		&qo.BudgetRegenRate,
		&qo.AssociatedQuestmakerIDs,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Quest Owner not found
		}
		return nil, fmt.Errorf("failed to get quest owner by ID: %w", err)
	}

	d.Cache.Set(qo.ID, qo, 300) // Cache for 5 minutes
	return qo, nil
}

// UpdateQuestOwner updates an existing quest owner in the database.
func (d *QuestOwnerDAL) UpdateQuestOwner(qo *models.QuestOwner) error {
	query := `
	UPDATE QuestOwners
	SET name = ?, description = ?, llm_prompt_context = ?, current_influence_budget = ?, max_influence_budget = ?, budget_regen_rate = ?, associated_questmaker_ids = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		qo.Name,
		qo.Description,
		qo.LLMPromptContext,
		qo.CurrentInfluenceBudget,
		qo.MaxInfluenceBudget,
		qo.BudgetRegenRate,
		qo.AssociatedQuestmakerIDs,
		qo.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update quest owner: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("quest owner with ID %s not found for update", qo.ID)
	}
	d.Cache.Delete(qo.ID) // Invalidate cache on update
	return nil
}

// DeleteQuestOwner deletes a quest owner from the database by its ID.
func (d *QuestOwnerDAL) DeleteQuestOwner(id string) error {
	query := `DELETE FROM QuestOwners WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete quest owner: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("quest owner with ID %s not found for deletion", id)
	}
	d.Cache.Delete(id) // Invalidate cache on delete
	return nil
}

// GetAllQuestOwners retrieves all quest owners from the database.
func (d *QuestOwnerDAL) GetAllQuestOwners() ([]*models.QuestOwner, error) {
	query := `SELECT id, name, description, llm_prompt_context, current_influence_budget, max_influence_budget, budget_regen_rate, associated_questmaker_ids FROM QuestOwners`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all quest owners: %w", err)
	}
	defer rows.Close()

	var questOwners []*models.QuestOwner
	for rows.Next() {
		qo := &models.QuestOwner{}
		err := rows.Scan(
			&qo.ID,
			&qo.Name,
			&qo.Description,
			&qo.LLMPromptContext,
			&qo.CurrentInfluenceBudget,
			&qo.MaxInfluenceBudget,
			&qo.BudgetRegenRate,
			&qo.AssociatedQuestmakerIDs,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quest owner: %w", err)
		}
		questOwners = append(questOwners, qo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through quest owners: %w", err)
	}

	return questOwners, nil
}