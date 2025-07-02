package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// ClassDAL handles database operations for Class entities.
type ClassDAL struct {
	db    *sql.DB
	cache *Cache
}

// NewClassDAL creates a new ClassDAL.
func NewClassDAL(db *sql.DB) *ClassDAL {
	return &ClassDAL{db: db, cache: NewCache()}
}

// CreateClass inserts a new class into the database.
func (d *ClassDAL) CreateClass(class *models.Class) error {
	query := `
	INSERT INTO Classes (id, name, description, total_levels, parent_class_id, associated_entity_type, associated_entity_id, level_up_rewards)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		class.ID,
		class.Name,
		class.Description,
		class.TotalLevels,
		class.ParentClassID,
		class.AssociatedEntityType,
		class.AssociatedEntityID,
		class.LevelUpRewards,
	)
	if err != nil {
		return fmt.Errorf("failed to create class: %w", err)
	}
	return nil
}

// GetClassByID retrieves a class by its ID.
func (d *ClassDAL) GetClassByID(id string) (*models.Class, error) {
	query := `SELECT id, name, description, total_levels, parent_class_id, associated_entity_type, associated_entity_id, level_up_rewards FROM Classes WHERE id = ?`
	row := d.db.QueryRow(query, id)

	class := &models.Class{}
	err := row.Scan(
		&class.ID,
		&class.Name,
		&class.Description,
		&class.TotalLevels,
		&class.ParentClassID,
		&class.AssociatedEntityType,
		&class.AssociatedEntityID,
		&class.LevelUpRewards,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Class not found
		}
		return nil, fmt.Errorf("failed to get class by ID: %w", err)
	}

	return class, nil
}

// UpdateClass updates an existing class in the database.
func (d *ClassDAL) UpdateClass(class *models.Class) error {
	query := `
	UPDATE Classes
	SET name = ?, description = ?, total_levels = ?, parent_class_id = ?, associated_entity_type = ?, associated_entity_id = ?, level_up_rewards = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		class.Name,
		class.Description,
		class.TotalLevels,
		class.ParentClassID,
		class.AssociatedEntityType,
		class.AssociatedEntityID,
		class.LevelUpRewards,
		class.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update class: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("class with ID %s not found for update", class.ID)
	}

	return nil
}

// DeleteClass deletes a class from the database by its ID.
func (d *ClassDAL) DeleteClass(id string) error {
	query := `DELETE FROM Classes WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete class: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("class with ID %s not found for deletion", id)
	}

	return nil
}

// GetAllClasses retrieves all classes from the database.
func (d *ClassDAL) GetAllClasses() ([]*models.Class, error) {
	query := `SELECT id, name, description, total_levels, parent_class_id, associated_entity_type, associated_entity_id, level_up_rewards FROM Classes`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all classes: %w", err)
	}
	defer rows.Close()

	var classes []*models.Class
	for rows.Next() {
		class := &models.Class{}
		err := rows.Scan(
			&class.ID,
			&class.Name,
			&class.Description,
			&class.TotalLevels,
			&class.ParentClassID,
			&class.AssociatedEntityType,
			&class.AssociatedEntityID,
			&class.LevelUpRewards,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan class: %w", err)
		}
		classes = append(classes, class)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through classes: %w", err)
	}

	return classes, nil
}