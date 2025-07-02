package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// QuestmakerDAL handles database operations for Questmaker entities.
type QuestmakerDAL struct {
	db    *sql.DB
	cache *Cache
}

// NewQuestmakerDAL creates a new QuestmakerDAL.
func NewQuestmakerDAL(db *sql.DB) *QuestmakerDAL {
	return &QuestmakerDAL{db: db, cache: NewCache()}
}

// CreateQuestmaker inserts a new questmaker into the database.
func (d *QuestmakerDAL) CreateQuestmaker(qm *models.Questmaker) error {
	query := `
	INSERT INTO Questmakers (id, name, llm_prompt_context, current_influence_budget, max_influence_budget, budget_regen_rate, memories_about_players, available_tools)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		qm.ID,
		qm.Name,
		qm.LLMPromptContext,
		qm.CurrentInfluenceBudget,
		qm.MaxInfluenceBudget,
		qm.BudgetRegenRate,
		qm.MemoriesAboutPlayers,
		qm.AvailableTools,
	)
	if err != nil {
		return fmt.Errorf("failed to create questmaker: %w", err)
	}
	return nil
}

// GetQuestmakerByID retrieves a questmaker by their ID.
func (d *QuestmakerDAL) GetQuestmakerByID(id string) (*models.Questmaker, error) {
	query := `SELECT id, name, llm_prompt_context, current_influence_budget, max_influence_budget, budget_regen_rate, memories_about_players, available_tools FROM Questmakers WHERE id = ?`
	row := d.db.QueryRow(query, id)

	qm := &models.Questmaker{}
	err := row.Scan(
		&qm.ID,
		&qm.Name,
		&qm.LLMPromptContext,
		&qm.CurrentInfluenceBudget,
		&qm.MaxInfluenceBudget,
		&qm.BudgetRegenRate,
		&qm.MemoriesAboutPlayers,
		&qm.AvailableTools,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Questmaker not found
		}
		return nil, fmt.Errorf("failed to get questmaker by ID: %w", err)
	}

	return qm, nil
}

// UpdateQuestmaker updates an existing questmaker in the database.
func (d *QuestmakerDAL) UpdateQuestmaker(qm *models.Questmaker) error {
	query := `
	UPDATE Questmakers
	SET name = ?, llm_prompt_context = ?, current_influence_budget = ?, max_influence_budget = ?, budget_regen_rate = ?, memories_about_players = ?, available_tools = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		qm.Name,
		qm.LLMPromptContext,
		qm.CurrentInfluenceBudget,
		qm.MaxInfluenceBudget,
		qm.BudgetRegenRate,
		qm.MemoriesAboutPlayers,
		qm.AvailableTools,
		qm.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update questmaker: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("questmaker with ID %s not found for update", qm.ID)
	}

	return nil
}

// DeleteQuestmaker deletes a questmaker from the database by their ID.
func (d *QuestmakerDAL) DeleteQuestmaker(id string) error {
	query := `DELETE FROM Questmakers WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete questmaker: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("questmaker with ID %s not found for deletion", id)
	}

	return nil
}

// GetAllQuestmakers retrieves all questmakers from the database.
func (d *QuestmakerDAL) GetAllQuestmakers() ([]*models.Questmaker, error) {
	query := `SELECT id, name, llm_prompt_context, current_influence_budget, max_influence_budget, budget_regen_rate, memories_about_players, available_tools FROM Questmakers`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all questmakers: %w", err)
	}
	defer rows.Close()

	var questmakers []*models.Questmaker
	for rows.Next() {
		qm := &models.Questmaker{}
		err := rows.Scan(
			&qm.ID,
			&qm.Name,
			&qm.LLMPromptContext,
			&qm.CurrentInfluenceBudget,
			&qm.MaxInfluenceBudget,
			&qm.BudgetRegenRate,
			&qm.MemoriesAboutPlayers,
			&qm.AvailableTools,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan questmaker: %w", err)
		}
		questmakers = append(questmakers, qm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through questmakers: %w", err)
	}

	return questmakers, nil
}