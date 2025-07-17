package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mud/internal/models"
)

// RaceDAL handles database operations for Race entities.
type RaceDAL struct {
	db    *sql.DB
	Cache CacheInterface
}

// NewRaceDAL creates a new RaceDAL.
func NewRaceDAL(db *sql.DB, cache CacheInterface) *RaceDAL {
	return &RaceDAL{db: db, Cache: cache}
}

// CreateRace inserts a new race into the database.
func (d *RaceDAL) CreateRace(race *models.Race) error {
	baseStatsJSON, err := json.Marshal(race.BaseStats)
	if err != nil {
		return fmt.Errorf("failed to marshal base stats: %w", err)
	}

	perceptionBiasesJSON, err := json.Marshal(race.PerceptionBiases)
	if err != nil {
		return fmt.Errorf("failed to marshal perception biases: %w", err)
	}

	query := `
	INSERT INTO Races (id, name, description, owner_id, base_stats, perception_biases)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = d.db.Exec(query,
		race.ID,
		race.Name,
		race.Description,
		race.OwnerID,
		string(baseStatsJSON),
		string(perceptionBiasesJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to create race: %w", err)
	}
	d.Cache.Set(race.ID, race, 300) // Cache for 5 minutes
	return nil
}

// GetRaceByID retrieves a race by its ID.
func (d *RaceDAL) GetRaceByID(id string) (*models.Race, error) {
	if cachedRace, found := d.Cache.Get(id); found {
		if race, ok := cachedRace.(*models.Race); ok {
			return race, nil
		}
	}

	query := `SELECT id, name, description, owner_id, base_stats, perception_biases FROM Races WHERE id = ?`
	row := d.db.QueryRow(query, id)

	race := &models.Race{}
	var baseStatsJSON, perceptionBiasesJSON []byte
	err := row.Scan(
		&race.ID,
		&race.Name,
		&race.Description,
		&race.OwnerID,
		&baseStatsJSON,
		&perceptionBiasesJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Race not found
		}
		return nil, fmt.Errorf("failed to get race by ID: %w", err)
	}

	if err := json.Unmarshal(baseStatsJSON, &race.BaseStats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal base stats for race %s: %w", race.ID, err)
	}

	if err := json.Unmarshal(perceptionBiasesJSON, &race.PerceptionBiases); err != nil {
		// Handle cases where the column might be NULL for older data
		if string(perceptionBiasesJSON) != "null" && string(perceptionBiasesJSON) != "" {
			return nil, fmt.Errorf("failed to unmarshal perception biases for race %s: %w", race.ID, err)
		}
		race.PerceptionBiases = make(map[string]float64) // Initialize to empty map
	}

	d.Cache.Set(race.ID, race, 300) // Cache for 5 minutes
	return race, nil
}

// UpdateRace updates an existing race in the database.
func (d *RaceDAL) UpdateRace(race *models.Race) error {
	baseStatsJSON, err := json.Marshal(race.BaseStats)
	if err != nil {
		return fmt.Errorf("failed to marshal base stats: %w", err)
	}

	perceptionBiasesJSON, err := json.Marshal(race.PerceptionBiases)
	if err != nil {
		return fmt.Errorf("failed to marshal perception biases: %w", err)
	}

	query := `
	UPDATE Races
	SET name = ?, description = ?, owner_id = ?, base_stats = ?, perception_biases = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		race.Name,
		race.Description,
		race.OwnerID,
		string(baseStatsJSON),
		string(perceptionBiasesJSON),
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
	d.Cache.Delete(race.ID) // Invalidate cache on update
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
	d.Cache.Delete(id) // Invalidate cache on delete
	return nil
}

// GetAllRaces retrieves all races from the database.
func (d *RaceDAL) GetAllRaces() ([]*models.Race, error) {
	query := `SELECT id, name, description, owner_id, base_stats, perception_biases FROM Races`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all races: %w", err)
	}
	defer rows.Close()

	var races []*models.Race
	for rows.Next() {
		race := &models.Race{}
		var baseStatsJSON, perceptionBiasesJSON []byte
		err := rows.Scan(
			&race.ID,
			&race.Name,
			&race.Description,
			&race.OwnerID,
			&baseStatsJSON,
			&perceptionBiasesJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan race: %w", err)
		}
		if err := json.Unmarshal(baseStatsJSON, &race.BaseStats); err != nil {
			return nil, fmt.Errorf("failed to unmarshal base stats for race %s: %w", race.ID, err)
		}
		if err := json.Unmarshal(perceptionBiasesJSON, &race.PerceptionBiases); err != nil {
			if string(perceptionBiasesJSON) != "null" && string(perceptionBiasesJSON) != "" {
				return nil, fmt.Errorf("failed to unmarshal perception biases for race %s: %w", race.ID, err)
			}
			race.PerceptionBiases = make(map[string]float64)
		}
		races = append(races, race)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through races: %w", err)
	}

	return races, nil
}
