package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mud/internal/models"
)

// ProfessionDAL handles database operations for Profession entities.
type ProfessionDAL struct {
	db    *sql.DB
	Cache *Cache
}

// NewProfessionDAL creates a new ProfessionDAL.
func NewProfessionDAL(db *sql.DB) *ProfessionDAL {
	return &ProfessionDAL{db: db, Cache: NewCache()}
}

// CreateProfession inserts a new profession into the database.
func (d *ProfessionDAL) CreateProfession(prof *models.Profession) error {
	baseSkillsJSON, err := json.Marshal(prof.BaseSkills)
	if err != nil {
		return fmt.Errorf("failed to marshal base skills: %w", err)
	}

	query := `
	INSERT INTO Professions (id, name, description, base_skills)
	VALUES (?, ?, ?, ?)
	`

	_, err = d.db.Exec(query,
		prof.ID,
		prof.Name,
		prof.Description,
		string(baseSkillsJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to create profession: %w", err)
	}
	d.Cache.Set(prof.ID, prof, 300)
	return nil
}

// GetProfessionByID retrieves a profession by its ID.
func (d *ProfessionDAL) GetProfessionByID(id string) (*models.Profession, error) {
	if cachedProf, found := d.Cache.Get(id); found {
		if prof, ok := cachedProf.(*models.Profession); ok {
			return prof, nil
		}
	}

	query := `SELECT id, name, description, base_skills FROM Professions WHERE id = ?`
	row := d.db.QueryRow(query, id)

	prof := &models.Profession{}
	var baseSkillsJSON []byte
	err := row.Scan(
		&prof.ID,
		&prof.Name,
		&prof.Description,
		&baseSkillsJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Profession not found
		}
		return nil, fmt.Errorf("failed to get profession by ID: %w", err)
	}

	if err := json.Unmarshal(baseSkillsJSON, &prof.BaseSkills); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base skills for profession %s: %w", prof.ID, err)
	}

	d.Cache.Set(prof.ID, prof, 300)
	return prof, nil
}

// UpdateProfession updates an existing profession in the database.
func (d *ProfessionDAL) UpdateProfession(prof *models.Profession) error {
	baseSkillsJSON, err := json.Marshal(prof.BaseSkills)
	if err != nil {
		return fmt.Errorf("failed to marshal base skills: %w", err)
	}

	query := `
	UPDATE Professions
	SET name = ?, description = ?, base_skills = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		prof.Name,
		prof.Description,
		string(baseSkillsJSON),
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
	d.Cache.Delete(prof.ID)
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
	d.Cache.Delete(id)
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
		var baseSkillsJSON []byte
		err := rows.Scan(
			&prof.ID,
			&prof.Name,
			&prof.Description,
			&baseSkillsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profession: %w", err)
		}
		if err := json.Unmarshal(baseSkillsJSON, &prof.BaseSkills); err != nil {
			return nil, fmt.Errorf("failed to unmarshal base skills for profession %s: %w", prof.ID, err)
		}
		professions = append(professions, prof)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through professions: %w", err)
	}

	return professions, nil
}
