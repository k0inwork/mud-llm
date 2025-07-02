package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// ProfessionDAL handles database operations for Profession entities.
type ProfessionDAL struct {
	db    *sql.DB
	cache *Cache
}

// NewProfessionDAL creates a new ProfessionDAL.
func NewProfessionDAL(db *sql.DB) *ProfessionDAL {
	return &ProfessionDAL{db: db, cache: NewCache()}
}

// CreateProfession inserts a new profession into the database.
func (d *ProfessionDAL) CreateProfession(prof *models.Profession) error {
	query := `
	INSERT INTO Professions (id, name, description, base_skills)
	VALUES (?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		prof.ID,
		prof.Name,
		prof.Description,
		prof.BaseSkills,
	)
	if err != nil {
		return fmt.Errorf("failed to create profession: %w", err)
	}
	return nil
}

// GetProfessionByID retrieves a profession by its ID.
func (d *ProfessionDAL) GetProfessionByID(id string) (*models.Profession, error) {
	query := `SELECT id, name, description, base_skills FROM Professions WHERE id = ?`
	row := d.db.QueryRow(query, id)

	prof := &models.Profession{}
	err := row.Scan(
		&prof.ID,
		&prof.Name,
		&prof.Description,
		&prof.BaseSkills,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Profession not found
		}
		return nil, fmt.Errorf("failed to get profession by ID: %w", err)
	}

	return prof, nil
}

// UpdateProfession updates an existing profession in the database.
func (d *ProfessionDAL) UpdateProfession(prof *models.Profession) error {
	query := `
	UPDATE Professions
	SET name = ?, description = ?, base_skills = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		prof.Name,
		prof.Description,
		prof.BaseSkills,
		prof.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update profession: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("profession with ID %s not found for update", prof.ID)
	}

	return nil
}

// DeleteProfession deletes a profession from the database by its ID.
func (d *ProfessionDAL) DeleteProfession(id string) error {
	query := `DELETE FROM Professions WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete profession: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("profession with ID %s not found for deletion", id)
	}

	return nil
}

// GetAllProfessions retrieves all professions from the database.
func (d *ProfessionDAL) GetAllProfessions() ([]*models.Profession, error) {
	query := `SELECT id, name, description, base_skills FROM Professions`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all professions: %w", err)
	}
	defer rows.Close()

	var professions []*models.Profession
	for rows.Next() {
		prof := &models.Profession{}
		err := rows.Scan(
			&prof.ID,
			&prof.Name,
			&prof.Description,
			&prof.BaseSkills,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profession: %w", err)
		}
		professions = append(professions, prof)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through professions: %w", err)
	}

	return professions, nil
}