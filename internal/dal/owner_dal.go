package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// OwnerDAL handles database operations for Owner entities.
type OwnerDAL struct {
	db    *sql.DB
	cache *Cache
}

// NewOwnerDAL creates a new OwnerDAL.
func NewOwnerDAL(db *sql.DB) *OwnerDAL {
	return &OwnerDAL{db: db, cache: NewCache()}
}

// CreateOwner inserts a new owner into the database.
func (d *OwnerDAL) CreateOwner(owner *models.Owner) error {
	query := `
	INSERT INTO Owners (id, name, description, monitored_aspect, associated_id, llm_prompt_context, memories_about_players, current_influence_budget, max_influence_budget, budget_regen_rate, available_tools)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		owner.ID,
		owner.Name,
		owner.Description,
		owner.MonitoredAspect,
		owner.AssociatedID,
		owner.LLMPromptContext,
		owner.MemoriesAboutPlayers,
		owner.CurrentInfluenceBudget,
		owner.MaxInfluenceBudget,
		owner.BudgetRegenRate,
		owner.AvailableTools,
	)
	if err != nil {
		return fmt.Errorf("failed to create owner: %w", err)
	}
	d.cache.Set(owner.ID, owner, 300) // Cache for 5 minutes
	return nil
}

// GetOwnerByID retrieves an owner by their ID.
func (d *OwnerDAL) GetOwnerByID(id string) (*models.Owner, error) {
	if cachedOwner, found := d.cache.Get(id); found {
		if owner, ok := cachedOwner.(*models.Owner); ok {
			return owner, nil
		}
	}

	query := `SELECT id, name, description, monitored_aspect, associated_id, llm_prompt_context, memories_about_players, current_influence_budget, max_influence_budget, budget_regen_rate, available_tools FROM Owners WHERE id = ?`
	row := d.db.QueryRow(query, id)

	owner := &models.Owner{}
	err := row.Scan(
		&owner.ID,
		&owner.Name,
		&owner.Description,
		&owner.MonitoredAspect,
		&owner.AssociatedID,
		&owner.LLMPromptContext,
		&owner.MemoriesAboutPlayers,
		&owner.CurrentInfluenceBudget,
		&owner.MaxInfluenceBudget,
		&owner.BudgetRegenRate,
		&owner.AvailableTools,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Owner not found
		}
		return nil, fmt.Errorf("failed to get owner by ID: %w", err)
	}

	d.cache.Set(owner.ID, owner, 300) // Cache for 5 minutes
	return owner, nil
}

// UpdateOwner updates an existing owner in the database.
func (d *OwnerDAL) UpdateOwner(owner *models.Owner) error {
	query := `
	UPDATE Owners
	SET name = ?, description = ?, monitored_aspect = ?, associated_id = ?, llm_prompt_context = ?, memories_about_players = ?, current_influence_budget = ?, max_influence_budget = ?, budget_regen_rate = ?, available_tools = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		owner.Name,
		owner.Description,
		owner.MonitoredAspect,
		owner.AssociatedID,
		owner.LLMPromptContext,
		owner.MemoriesAboutPlayers,
		owner.CurrentInfluenceBudget,
		owner.MaxInfluenceBudget,
		owner.BudgetRegenRate,
		owner.AvailableTools,
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
	d.cache.Delete(owner.ID) // Invalidate cache on update
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
	d.cache.Delete(id) // Invalidate cache on delete
	return nil
}

// GetOwnersByMonitoredAspect retrieves owners by their monitored aspect and associated ID.
func (d *OwnerDAL) GetOwnersByMonitoredAspect(aspectType string, associatedID string) ([]*models.Owner, error) {
	query := `SELECT id, name, description, monitored_aspect, associated_id, llm_prompt_context, memories_about_players, current_influence_budget, max_influence_budget, budget_regen_rate, available_tools FROM Owners WHERE monitored_aspect = ? AND associated_id = ?`
	rows, err := d.db.Query(query, aspectType, associatedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owners by monitored aspect: %w", err)
	}
	defer rows.Close()

	var owners []*models.Owner
	for rows.Next() {
		owner := &models.Owner{}
		err := rows.Scan(
			&owner.ID,
			&owner.Name,
			&owner.Description,
			&owner.MonitoredAspect,
			&owner.AssociatedID,
			&owner.LLMPromptContext,
			&owner.MemoriesAboutPlayers,
			&owner.CurrentInfluenceBudget,
			&owner.MaxInfluenceBudget,
			&owner.BudgetRegenRate,
			&owner.AvailableTools,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan owner row: %w", err)
		}
		owners = append(owners, owner)
	}

	return owners, nil
}

// GetAllOwners retrieves all owners from the database.
func (d *OwnerDAL) GetAllOwners() ([]*models.Owner, error) {
	query := `SELECT id, name, description, monitored_aspect, associated_id, llm_prompt_context, memories_about_players, current_influence_budget, max_influence_budget, budget_regen_rate, available_tools FROM Owners`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all owners: %w", err)
	}
	defer rows.Close()

	var owners []*models.Owner
	for rows.Next() {
		owner := &models.Owner{}
		err := rows.Scan(
			&owner.ID,
			&owner.Name,
			&owner.Description,
			&owner.MonitoredAspect,
			&owner.AssociatedID,
			&owner.LLMPromptContext,
			&owner.MemoriesAboutPlayers,
			&owner.CurrentInfluenceBudget,
			&owner.MaxInfluenceBudget,
			&owner.BudgetRegenRate,
			&owner.AvailableTools,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan owner: %w", err)
		}
		owners = append(owners, owner)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through owners: %w", err)
	}

	return owners, nil
}
