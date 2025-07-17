package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mud/internal/models"
)

type QuestmakerDAL struct {
	db    *sql.DB
	Cache CacheInterface
}

func NewQuestmakerDAL(db *sql.DB, cache CacheInterface) *QuestmakerDAL {
	return &QuestmakerDAL{db: db, Cache: cache}
}

func (d *QuestmakerDAL) CreateQuestmaker(qm *models.Questmaker) error {
	memoriesJSON, err := json.Marshal(qm.MemoriesAboutPlayers)
	if err != nil {
		return fmt.Errorf("failed to marshal memories: %w", err)
	}
	toolsJSON, err := json.Marshal(qm.AvailableTools)
	if err != nil {
		return fmt.Errorf("failed to marshal tools: %w", err)
	}

	query := `
	INSERT INTO Questmakers (id, name, llm_prompt_context, current_influence_budget, max_influence_budget, budget_regen_rate, memories_about_players, available_tools, reaction_threshold)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = d.db.Exec(query, qm.ID, qm.Name, qm.LLMPromptContext, qm.CurrentInfluenceBudget, qm.MaxInfluenceBudget, qm.BudgetRegenRate, string(memoriesJSON), string(toolsJSON), qm.ReactionThreshold)
	if err != nil {
		return fmt.Errorf("failed to create questmaker: %w", err)
	}
	d.Cache.Set(qm.ID, qm, 300)
	return nil
}

func (d *QuestmakerDAL) GetQuestmakerByID(id string) (*models.Questmaker, error) {
	if cached, found := d.Cache.Get(id); found {
		if qm, ok := cached.(*models.Questmaker); ok {
			return qm, nil
		}
	}

	query := `SELECT id, name, llm_prompt_context, current_influence_budget, max_influence_budget, budget_regen_rate, memories_about_players, available_tools, reaction_threshold FROM Questmakers WHERE id = ?`
	row := d.db.QueryRow(query, id)

	qm := &models.Questmaker{}
	var memoriesJSON, toolsJSON []byte
	err := row.Scan(&qm.ID, &qm.Name, &qm.LLMPromptContext, &qm.CurrentInfluenceBudget, &qm.MaxInfluenceBudget, &qm.BudgetRegenRate, &memoriesJSON, &toolsJSON, &qm.ReactionThreshold)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get questmaker by ID: %w", err)
	}

	if err := json.Unmarshal(memoriesJSON, &qm.MemoriesAboutPlayers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal memories: %w", err)
	}
	if err := json.Unmarshal(toolsJSON, &qm.AvailableTools); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tools: %w", err)
	}

	d.Cache.Set(qm.ID, qm, 300)
	return qm, nil
}

func (d *QuestmakerDAL) UpdateQuestmaker(qm *models.Questmaker) error {
	memoriesJSON, err := json.Marshal(qm.MemoriesAboutPlayers)
	if err != nil {
		return fmt.Errorf("failed to marshal memories: %w", err)
	}
	toolsJSON, err := json.Marshal(qm.AvailableTools)
	if err != nil {
		return fmt.Errorf("failed to marshal tools: %w", err)
	}

	query := `
	UPDATE Questmakers
	SET name = ?, llm_prompt_context = ?, current_influence_budget = ?, max_influence_budget = ?, budget_regen_rate = ?, memories_about_players = ?, available_tools = ?, reaction_threshold = ?
	WHERE id = ?
	`
	_, err = d.db.Exec(query, qm.Name, qm.LLMPromptContext, qm.CurrentInfluenceBudget, qm.MaxInfluenceBudget, qm.BudgetRegenRate, string(memoriesJSON), string(toolsJSON), qm.ReactionThreshold, qm.ID)
	if err != nil {
		return fmt.Errorf("failed to update questmaker: %w", err)
	}
	d.Cache.Delete(qm.ID)
	return nil
}

func (d *QuestmakerDAL) DeleteQuestmaker(id string) error {
	query := `DELETE FROM Questmakers WHERE id = ?`
	_, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete questmaker: %w", err)
	}
	d.Cache.Delete(id)
	return nil
}

func (d *QuestmakerDAL) GetAllQuestmakers() ([]*models.Questmaker, error) {
	query := `SELECT id, name, llm_prompt_context, current_influence_budget, max_influence_budget, budget_regen_rate, memories_about_players, available_tools, reaction_threshold FROM Questmakers`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all questmakers: %w", err)
	}
	defer rows.Close()

	var qms []*models.Questmaker
	for rows.Next() {
		qm := &models.Questmaker{}
		var memoriesJSON, toolsJSON []byte
		err := rows.Scan(&qm.ID, &qm.Name, &qm.LLMPromptContext, &qm.CurrentInfluenceBudget, &qm.MaxInfluenceBudget, &qm.BudgetRegenRate, &memoriesJSON, &toolsJSON, &qm.ReactionThreshold)
		if err != nil {
			return nil, fmt.Errorf("failed to scan questmaker: %w", err)
		}

		if err := json.Unmarshal(memoriesJSON, &qm.MemoriesAboutPlayers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal memories for questmaker %s: %w", qm.ID, err)
		}
		if err := json.Unmarshal(toolsJSON, &qm.AvailableTools); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tools for questmaker %s: %w", qm.ID, err)
		}
		qms = append(qms, qm)
	}
	return qms, nil
}