package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mud/internal/models"
)

// OwnerDAL handles database operations for Owner entities.
type OwnerDAL struct {
	db    *sql.DB
	cache CacheInterface
}

func (d *OwnerDAL) Cache() CacheInterface {
	return d.cache
}

// NewOwnerDAL creates a new OwnerDAL.
func NewOwnerDAL(db *sql.DB, cache CacheInterface) *OwnerDAL {
	return &OwnerDAL{db: db, cache: cache}
}

// CreateOwner inserts a new owner into the database.
func (d *OwnerDAL) CreateOwner(owner *models.Owner) error {
	memoriesJSON, err := json.Marshal(owner.MemoriesAboutPlayers)
	if err != nil {
		return fmt.Errorf("failed to marshal memories about players: %w", err)
	}
	availableToolsJSON, err := json.Marshal(owner.AvailableTools)
	if err != nil {
		return fmt.Errorf("failed to marshal available tools: %w", err)
	}
	initiatedQuestsJSON, err := json.Marshal(owner.InitiatedQuests)
	if err != nil {
		return fmt.Errorf("failed to marshal initiated quests: %w", err)
	}

	query := `
	INSERT INTO Owners (id, name, description, monitored_aspect, associated_id, llm_prompt_context, memories_about_players, current_influence_budget, max_influence_budget, budget_regen_rate, available_tools, initiated_quests, reaction_threshold)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = d.db.Exec(query,
		owner.ID,
		owner.Name,
		owner.Description,
		owner.MonitoredAspect,
		owner.AssociatedID,
		owner.LLMPromptContext,
		string(memoriesJSON),
		owner.CurrentInfluenceBudget,
		owner.MaxInfluenceBudget,
		owner.BudgetRegenRate,
		string(availableToolsJSON),
		string(initiatedQuestsJSON),
		owner.ReactionThreshold,
	)
	if err != nil {
		return fmt.Errorf("failed to create owner: %w", err)
	}
	d.Cache().Set(owner.ID, owner, 300) // Cache for 5 minutes
	return nil
}

// GetOwnerByID retrieves an owner by their ID.
func (d *OwnerDAL) GetOwnerByID(id string) (*models.Owner, error) {
	if cachedOwner, found := d.Cache().Get(id); found {
		if owner, ok := cachedOwner.(*models.Owner); ok {
			return owner, nil
		}
	}

	query := `SELECT id, name, description, monitored_aspect, associated_id, llm_prompt_context, memories_about_players, current_influence_budget, max_influence_budget, budget_regen_rate, available_tools, initiated_quests, reaction_threshold FROM Owners WHERE id = ?`
	row := d.db.QueryRow(query, id)

	owner := &models.Owner{}
	var memoriesJSON, availableToolsJSON, initiatedQuestsJSON []byte
	err := row.Scan(
		&owner.ID,
		&owner.Name,
		&owner.Description,
		&owner.MonitoredAspect,
		&owner.AssociatedID,
		&owner.LLMPromptContext,
		&memoriesJSON,
		&owner.CurrentInfluenceBudget,
		&owner.MaxInfluenceBudget,
		&owner.BudgetRegenRate,
		&availableToolsJSON,
		&initiatedQuestsJSON,
		&owner.ReactionThreshold,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Owner not found
		}
		return nil, fmt.Errorf("failed to get owner by ID: %w", err)
	}

	if err := json.Unmarshal(memoriesJSON, &owner.MemoriesAboutPlayers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal memories about players: %w", err)
	}
	if err := json.Unmarshal(availableToolsJSON, &owner.AvailableTools); err != nil {
		return nil, fmt.Errorf("failed to unmarshal available tools: %w", err)
	}
	if err := json.Unmarshal(initiatedQuestsJSON, &owner.InitiatedQuests); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initiated quests: %w", err)
	}

	d.Cache().Set(owner.ID, owner, 300) // Cache for 5 minutes
	return owner, nil
}

// UpdateOwner updates an existing owner in the database.
func (d *OwnerDAL) UpdateOwner(owner *models.Owner) error {
	memoriesJSON, err := json.Marshal(owner.MemoriesAboutPlayers)
	if err != nil {
		return fmt.Errorf("failed to marshal memories about players: %w", err)
	}
	availableToolsJSON, err := json.Marshal(owner.AvailableTools)
	if err != nil {
		return fmt.Errorf("failed to marshal available tools: %w", err)
	}
	initiatedQuestsJSON, err := json.Marshal(owner.InitiatedQuests)
	if err != nil {
		return fmt.Errorf("failed to marshal initiated quests: %w", err)
	}

	query := `
	UPDATE Owners
	SET name = ?, description = ?, monitored_aspect = ?, associated_id = ?, llm_prompt_context = ?, memories_about_players = ?, current_influence_budget = ?, max_influence_budget = ?, budget_regen_rate = ?, available_tools = ?, initiated_quests = ?, reaction_threshold = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		owner.Name,
		owner.Description,
		owner.MonitoredAspect,
		owner.AssociatedID,
		owner.LLMPromptContext,
		string(memoriesJSON),
		owner.CurrentInfluenceBudget,
		owner.MaxInfluenceBudget,
		owner.BudgetRegenRate,
		string(availableToolsJSON),
		string(initiatedQuestsJSON),
		owner.ReactionThreshold,
		owner.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update owner: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("owner with ID %s not found for update", owner.ID)
	}
	d.Cache().Delete(owner.ID) // Invalidate cache on update
	return nil
}

// DeleteOwner deletes an owner from the database by their ID.
func (d *OwnerDAL) DeleteOwner(id string) error {
	query := `DELETE FROM Owners WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete owner: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("owner with ID %s not found for deletion", id)
	}
	d.Cache().Delete(id) // Invalidate cache on delete
	return nil
}

// GetOwnersByMonitoredAspect retrieves owners by their monitored aspect and associated ID.
func (d *OwnerDAL) GetOwnersByMonitoredAspect(aspectType string, associatedID string) ([]*models.Owner, error) {
	query := `SELECT id, name, description, monitored_aspect, associated_id, llm_prompt_context, memories_about_players, current_influence_budget, max_influence_budget, budget_regen_rate, available_tools, initiated_quests, reaction_threshold FROM Owners WHERE monitored_aspect = ? AND associated_id = ?`
	rows, err := d.db.Query(query, aspectType, associatedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owners by monitored aspect: %w", err)
	}
	defer rows.Close()

	var owners []*models.Owner
	for rows.Next() {
		owner := &models.Owner{}
		var memoriesJSON, availableToolsJSON, initiatedQuestsJSON []byte
		err := rows.Scan(
			&owner.ID,
			&owner.Name,
			&owner.Description,
			&owner.MonitoredAspect,
			&owner.AssociatedID,
			&owner.LLMPromptContext,
			&memoriesJSON,
			&owner.CurrentInfluenceBudget,
			&owner.MaxInfluenceBudget,
			&owner.BudgetRegenRate,
			&availableToolsJSON,
			&initiatedQuestsJSON,
			&owner.ReactionThreshold,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan owner row: %w", err)
		}

		if err := json.Unmarshal(memoriesJSON, &owner.MemoriesAboutPlayers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal memories about players for Owner %s: %w", owner.ID, err)
		}
		if err := json.Unmarshal(availableToolsJSON, &owner.AvailableTools); err != nil {
			return nil, fmt.Errorf("failed to unmarshal available tools for Owner %s: %w", owner.ID, err)
		}
		if err := json.Unmarshal(initiatedQuestsJSON, &owner.InitiatedQuests); err != nil {
			return nil, fmt.Errorf("failed to unmarshal initiated quests for Owner %s: %w", owner.ID, err)
		}

		owners = append(owners, owner)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through owners: %w", err)
	}

	return owners, nil
}

// GetAllOwners retrieves all owners from the database.
func (d *OwnerDAL) GetAllOwners() ([]*models.Owner, error) {
	query := `SELECT id, name, description, monitored_aspect, associated_id, llm_prompt_context, memories_about_players, current_influence_budget, max_influence_budget, budget_regen_rate, available_tools, initiated_quests, reaction_threshold FROM Owners`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all owners: %w", err)
	}
	defer rows.Close()

	var owners []*models.Owner
	for rows.Next() {
		owner := &models.Owner{}
		var memoriesJSON, availableToolsJSON, initiatedQuestsJSON []byte
		err := rows.Scan(
			&owner.ID,
			&owner.Name,
			&owner.Description,
			&owner.MonitoredAspect,
			&owner.AssociatedID,
			&owner.LLMPromptContext,
			&memoriesJSON,
			&owner.CurrentInfluenceBudget,
			&owner.MaxInfluenceBudget,
			&owner.BudgetRegenRate,
			&availableToolsJSON,
			&initiatedQuestsJSON,
			&owner.ReactionThreshold,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan owner: %w", err)
		}

		if err := json.Unmarshal(memoriesJSON, &owner.MemoriesAboutPlayers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal memories about players for Owner %s: %w", owner.ID, err)
		}
		if err := json.Unmarshal(availableToolsJSON, &owner.AvailableTools); err != nil {
			return nil, fmt.Errorf("failed to unmarshal available tools for Owner %s: %w", owner.ID, err)
		}
		if err := json.Unmarshal(initiatedQuestsJSON, &owner.InitiatedQuests); err != nil {
			return nil, fmt.Errorf("failed to unmarshal initiated quests for Owner %s: %w", owner.ID, err)
		}

		owners = append(owners, owner)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through owners: %w", err)
	}

	return owners, nil
}
