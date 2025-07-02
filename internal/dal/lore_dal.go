package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// LoreDAL handles database operations for Lore entities.
type LoreDAL struct {
	db    *sql.DB
	cache *Cache
}

// NewLoreDAL creates a new LoreDAL.
func NewLoreDAL(db *sql.DB) *LoreDAL {
	return &LoreDAL{db: db, cache: NewCache()}
}

// CreateLore inserts a new lore entry into the database.
func (d *LoreDAL) CreateLore(lore *models.Lore) error {
	query := `
	INSERT INTO Lore (id, title, content, scope, associated_id)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		lore.ID,
		lore.Title,
		lore.Content,
		lore.Scope,
		lore.AssociatedID,
	)
	if err != nil {
		return fmt.Errorf("failed to create lore: %w", err)
	}
	d.cache.Set(lore.ID, lore, 300) // Cache for 5 minutes
	return nil
}

// GetLoreByID retrieves a lore entry by its ID.
func (d *LoreDAL) GetLoreByID(id string) (*models.Lore, error) {
	if cachedLore, found := d.cache.Get(id); found {
		if lore, ok := cachedLore.(*models.Lore); ok {
			return lore, nil
		}
	}

	query := `SELECT id, title, content, scope, associated_id FROM Lore WHERE id = ?`
	row := d.db.QueryRow(query, id)

	lore := &models.Lore{}
	err := row.Scan(
		&lore.ID,
		&lore.Title,
		&lore.Content,
		&lore.Scope,
		&lore.AssociatedID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Lore not found
		}
		return nil, fmt.Errorf("failed to get lore by ID: %w", err)
	}

	d.cache.Set(lore.ID, lore, 300) // Cache for 5 minutes
	return lore, nil
}

// UpdateLore updates an existing lore entry in the database.
func (d *LoreDAL) UpdateLore(lore *models.Lore) error {
	query := `
	UPDATE Lore
	SET title = ?, content = ?, scope = ?, associated_id = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		lore.Title,
		lore.Content,
		lore.Scope,
		lore.AssociatedID,
		lore.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update lore: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("lore with ID %s not found for update", lore.ID)
	}

	return nil
}

// DeleteLore deletes a lore entry from the database by its ID.
func (d *LoreDAL) DeleteLore(id string) error {
	query := `DELETE FROM Lore WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete lore: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("lore with ID %s not found for deletion", id)
	}

	return nil
}

// GetLoreByTypeAndAssociatedID retrieves lore entries by type and associated ID.
func (d *LoreDAL) GetLoreByTypeAndAssociatedID(loreType string, associatedID string) ([]*models.Lore, error) {
	query := `SELECT id, title, content, scope, associated_id FROM Lore WHERE scope = ? AND associated_id = ?`
	rows, err := d.db.Query(query, loreType, associatedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lore by type and associated ID: %w", err)
	}
	defer rows.Close()

	var lores []*models.Lore
	for rows.Next() {
		lore := &models.Lore{}
		err := rows.Scan(
			&lore.ID,
			&lore.Title,
			&lore.Content,
			&lore.Scope,
			&lore.AssociatedID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lore row: %w", err)
		}
		lores = append(lores, lore)
	}

	return lores, nil
}

// GetAllLore retrieves all lore entries from the database.
func (d *LoreDAL) GetAllLore() ([]*models.Lore, error) {
	query := `SELECT id, title, content, scope, associated_id FROM Lore`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all lore: %w", err)
	}
	defer rows.Close()

	var lores []*models.Lore
	for rows.Next() {
		lore := &models.Lore{}
		err := rows.Scan(
			&lore.ID,
			&lore.Title,
			&lore.Content,
			&lore.Scope,
			&lore.AssociatedID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lore: %w", err)
		}
		lores = append(lores, lore)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through lore: %w", err)
	}

	return lores, nil
}
