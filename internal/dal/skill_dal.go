package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// SkillDAL handles database operations for Skill entities.
type SkillDAL struct {
	db    *sql.DB
	Cache *Cache
}

// NewSkillDAL creates a new SkillDAL.
func NewSkillDAL(db *sql.DB, cache *Cache) *SkillDAL {
	return &SkillDAL{db: db, Cache: cache}
}

// CreateSkill inserts a new skill into the database.
func (d *SkillDAL) CreateSkill(skill *models.Skill) error {
	query := `
	INSERT INTO Skills (id, name, category, description, type, associated_class_id, granted_by_entity_type, granted_by_entity_id, effects, cost, cooldown, min_class_level)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		skill.ID,
		skill.Name,
		skill.Category,
		skill.Description,
		skill.Type,
		skill.AssociatedClassID,
		skill.GrantedByEntityType,
		skill.GrantedByEntityID,
		skill.Effects,
		skill.Cost,
		skill.Cooldown,
		skill.MinClassLevel,
	)
	if err != nil {
		return fmt.Errorf("failed to create skill: %w", err)
	}
	return nil
}

// GetSkillByID retrieves a skill by its ID.
func (d *SkillDAL) GetSkillByID(id string) (*models.Skill, error) {
	query := `SELECT id, name, category, description, type, associated_class_id, granted_by_entity_type, granted_by_entity_id, effects, cost, cooldown, min_class_level FROM Skills WHERE id = ?`
	row := d.db.QueryRow(query, id)

	skill := &models.Skill{}
	err := row.Scan(
		&skill.ID,
		&skill.Name,
		&skill.Category,
		&skill.Description,
		&skill.Type,
		&skill.AssociatedClassID,
		&skill.GrantedByEntityType,
		&skill.GrantedByEntityID,
		&skill.Effects,
		&skill.Cost,
		&skill.Cooldown,
		&skill.MinClassLevel,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Skill not found
		}
		return nil, fmt.Errorf("failed to get skill by ID: %w", err)
	}

	return skill, nil
}

// UpdateSkill updates an existing skill in the database.
func (d *SkillDAL) UpdateSkill(skill *models.Skill) error {
	query := `
	UPDATE Skills
	SET name = ?, category = ?, description = ?, type = ?, associated_class_id = ?, granted_by_entity_type = ?, granted_by_entity_id = ?, effects = ?, cost = ?, cooldown = ?, min_class_level = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		skill.Name,
		skill.Category,
		skill.Description,
		skill.Type,
		skill.AssociatedClassID,
		skill.GrantedByEntityType,
		skill.GrantedByEntityID,
		skill.Effects,
		skill.Cost,
		skill.Cooldown,
		skill.MinClassLevel,
		skill.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update skill: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("skill with ID %s not found for update", skill.ID)
	}

	return nil
}

// DeleteSkill deletes a skill from the database by its ID.
func (d *SkillDAL) DeleteSkill(id string) error {
	query := `DELETE FROM Skills WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("skill with ID %s not found for deletion", id)
	}

	return nil
}

// GetAllSkills retrieves all skills from the database.
func (d *SkillDAL) GetAllSkills() ([]*models.Skill, error) {
	query := `SELECT id, name, category, description, type, associated_class_id, granted_by_entity_type, granted_by_entity_id, effects, cost, cooldown, min_class_level FROM Skills`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all skills: %w", err)
	}
	defer rows.Close()

	var skills []*models.Skill
	for rows.Next() {
		skill := &models.Skill{}
		err := rows.Scan(
			&skill.ID,
			&skill.Name,
			&skill.Category,
			&skill.Description,
			&skill.Type,
			&skill.AssociatedClassID,
			&skill.GrantedByEntityType,
			&skill.GrantedByEntityID,
			&skill.Effects,
			&skill.Cost,
			&skill.Cooldown,
			&skill.MinClassLevel,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}
		skills = append(skills, skill)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through skills: %w", err)
	}

	return skills, nil
}