package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// RaceDAL handles database operations for Race entities.
type RaceDAL struct {
	db    *sql.DB
	cache *Cache
}

// NewRaceDAL creates a new RaceDAL.
func NewRaceDAL(db *sql.DB) *RaceDAL {
	return &RaceDAL{db: db, cache: NewCache()}
}

// CreateRace inserts a new race into the database.
func (d *RaceDAL) CreateRace(race *models.Race) error {
	query := `
	INSERT INTO Races (id, name, description, base_stats)
	VALUES (?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		race.ID,
		race.Name,
		race.Description,
		race.BaseStats,
	)
	if err != nil {
		return fmt.Errorf("failed to create race: %w", err)
	}
	return nil
}

// GetRaceByID retrieves a race by its ID.
func (d *RaceDAL) GetRaceByID(id string) (*models.Race, error) {
	query := `SELECT id, name, description, base_stats FROM Races WHERE id = ?`
	row := d.db.QueryRow(query, id)

	race := &models.Race{}
	err := row.Scan(
		&race.ID,
		&race.Name,
		&race.Description,
		&race.BaseStats,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Race not found
		}
		return nil, fmt.Errorf("failed to get race by ID: %w", err)
	}

	return race, nil
}

// UpdateRace updates an existing race in the database.
func (d *RaceDAL) UpdateRace(race *models.Race) error {
	query := `
	UPDATE Races
	SET name = ?, description = ?, base_stats = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		race.Name,
		race.Description,
		race.BaseStats,
		race.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update race: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("race with ID %s not found for update", race.ID)
	}

	return nil
}

// DeleteRace deletes a race from the database by its ID.
func (d *RaceDAL) DeleteRace(id string) error {
	query := `DELETE FROM Races WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete race: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("race with ID %s not found for deletion", id)
	}

	return nil
}

// GetAllRaces retrieves all races from the database.
func (d *RaceDAL) GetAllRaces() ([]*models.Race, error) {
	query := `SELECT id, name, description, base_stats FROM Races`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all races: %w", err)
	}
	defer rows.Close()

	var races []*models.Race
	for rows.Next() {
		race := &models.Race{}
		err := rows.Scan(
			&race.ID,
			&race.Name,
			&race.Description,
			&race.BaseStats,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan race: %w", err)
		}
		races = append(races, race)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through races: %w", err)
	}

	return races, nil
}