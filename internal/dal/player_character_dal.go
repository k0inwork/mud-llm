package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mud/internal/models"
	"time"
)

// PlayerCharacterDAL handles database operations for PlayerCharacter entities.
type PlayerCharacterDAL struct {
	db      *sql.DB
	cache   CacheInterface
	itemDAL ItemDALInterface
}

func (d *PlayerCharacterDAL) Cache() CacheInterface {
	return d.cache
}

// NewPlayerCharacterDAL creates a new PlayerCharacterDAL.
func NewPlayerCharacterDAL(db *sql.DB, cache CacheInterface, itemDAL ItemDALInterface) *PlayerCharacterDAL {
	return &PlayerCharacterDAL{db: db, cache: cache, itemDAL: itemDAL}
}

// CreateCharacter inserts a new player character into the database.
func (d *PlayerCharacterDAL) CreateCharacter(character *models.PlayerCharacter) error {
	query := `
	INSERT INTO player_characters (id, player_account_id, name, race_id, profession_id, current_room_id, health, max_health, inventory, visited_room_ids, created_at, last_played_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		character.ID,
		character.PlayerAccountID,
		character.Name,
		character.RaceID,
		character.ProfessionID,
		character.CurrentRoomID,
		character.Health,
		character.MaxHealth,
		character.Inventory,
		character.VisitedRoomIDs,
		character.CreatedAt,
		character.LastPlayedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create player character: %w", err)
	}
	d.cache.Set(character.ID, character, 300*time.Second)
	return nil
}

// GetCharacterByID retrieves a player character by their ID.
func (d *PlayerCharacterDAL) GetCharacterByID(id string) (*models.PlayerCharacter, error) {
	if cached, found := d.cache.Get(id); found {
		if character, ok := cached.(*models.PlayerCharacter); ok {
			return character, nil
		}
	}

	query := `SELECT id, player_account_id, name, race_id, profession_id, current_room_id, health, max_health, inventory, visited_room_ids, created_at, last_played_at FROM player_characters WHERE id = ?`
	row := d.db.QueryRow(query, id)

	character := &models.PlayerCharacter{}
	var lastPlayed sql.NullTime

	err := row.Scan(
		&character.ID,
		&character.PlayerAccountID,
		&character.Name,
		&character.RaceID,
		&character.ProfessionID,
		&character.CurrentRoomID,
		&character.Health,
		&character.MaxHealth,
		&character.Inventory,
		&character.VisitedRoomIDs,
		&character.CreatedAt,
		&lastPlayed,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Character not found
		}
		return nil, fmt.Errorf("failed to get player character by ID: %w", err)
	}
	if lastPlayed.Valid {
		character.LastPlayedAt = lastPlayed.Time
	}

	d.cache.Set(character.ID, character, 300*time.Second)
	return character, nil
}

// GetCharactersByAccountID retrieves all characters associated with a player account.
func (d *PlayerCharacterDAL) GetCharactersByAccountID(accountID string) ([]*models.PlayerCharacter, error) {
	query := `SELECT id, player_account_id, name, race_id, profession_id, current_room_id, health, max_health, inventory, visited_room_ids, created_at, last_played_at FROM player_characters WHERE player_account_id = ?`
	rows, err := d.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get characters by account ID: %w", err)
	}
	defer rows.Close()

	var characters []*models.PlayerCharacter
	for rows.Next() {
		character := &models.PlayerCharacter{}
		var lastPlayed sql.NullTime
		err := rows.Scan(
			&character.ID,
			&character.PlayerAccountID,
			&character.Name,
			&character.RaceID,
			&character.ProfessionID,
			&character.CurrentRoomID,
			&character.Health,
			&character.MaxHealth,
			&character.Inventory,
			&character.VisitedRoomIDs,
			&character.CreatedAt,
			&lastPlayed,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player character: %w", err)
		}
		if lastPlayed.Valid {
			character.LastPlayedAt = lastPlayed.Time
		}
		characters = append(characters, character)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through player characters: %w", err)
	}

	return characters, nil
}

// UpdateCharacter updates an existing player character in the database.
func (d *PlayerCharacterDAL) UpdateCharacter(character *models.PlayerCharacter) error {
	query := `
	UPDATE player_characters
	SET name = ?, race_id = ?, profession_id = ?, current_room_id = ?, health = ?, max_health = ?, inventory = ?, visited_room_ids = ?, last_played_at = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		character.Name,
		character.RaceID,
		character.ProfessionID,
		character.CurrentRoomID,
		character.Health,
		character.MaxHealth,
		character.Inventory,
		character.VisitedRoomIDs,
		character.LastPlayedAt,
		character.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update player character: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("player character with ID %s not found for update", character.ID)
	}

	d.cache.Set(character.ID, character, 300*time.Second)
	return nil
}

// DeleteCharacter deletes a player character from the database by their ID.
func (d *PlayerCharacterDAL) DeleteCharacter(id string) error {
	query := `DELETE FROM player_characters WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete player character: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("player character with ID %s not found for deletion", id)
	}

	d.cache.Delete(id)
	return nil
}

// GetAllCharacters retrieves all player characters from the database.
func (d *PlayerCharacterDAL) GetAllCharacters() ([]*models.PlayerCharacter, error) {
	query := `SELECT id, player_account_id, name, race_id, profession_id, current_room_id, health, max_health, inventory, visited_room_ids, created_at, last_played_at FROM player_characters`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all player characters: %w", err)
	}
	defer rows.Close()

	var characters []*models.PlayerCharacter
	for rows.Next() {
		character := &models.PlayerCharacter{}
		var lastPlayed sql.NullTime
		err := rows.Scan(
			&character.ID,
			&character.PlayerAccountID,
			&character.Name,
			&character.RaceID,
			&character.ProfessionID,
			&character.CurrentRoomID,
			&character.Health,
			&character.MaxHealth,
			&character.Inventory,
			&character.VisitedRoomIDs,
			&character.CreatedAt,
			&lastPlayed,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player character: %w", err)
		}
		if lastPlayed.Valid {
			character.LastPlayedAt = lastPlayed.Time
		}
		characters = append(characters, character)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through player characters: %w", err)
	}

	return characters, nil
}

// GetCharacterInventory retrieves all items in a character's inventory.
func (d *PlayerCharacterDAL) GetCharacterInventory(characterID string) ([]*models.Item, error) {
	character, err := d.GetCharacterByID(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character for inventory: %w", err)
	}
	if character == nil {
		return nil, fmt.Errorf("character with ID %s not found", characterID)
	}

	var itemIDs []string
	if err := json.Unmarshal([]byte(character.Inventory), &itemIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal character inventory: %w", err)
	}

	var items []*models.Item
	for _, itemID := range itemIDs {
		item, err := d.itemDAL.GetItemByID(itemID)
		if err != nil {
			return nil, fmt.Errorf("failed to get item %s from inventory: %w", itemID, err)
		}
		if item != nil {
			items = append(items, item)
		}
	}

	return items, nil
}

// GetCharacterClass retrieves the class for a given character.
func (d *PlayerCharacterDAL) GetCharacterClass(characterID string) (*models.PlayerClass, error) {
	query := `SELECT player_id, class_id, level, experience FROM PlayerClasses WHERE player_id = ?`
	row := d.db.QueryRow(query, characterID)

	pc := &models.PlayerClass{}
	err := row.Scan(
		&pc.PlayerID,
		&pc.ClassID,
		&pc.Level,
		&pc.Experience,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Player class not found
		}
		return nil, fmt.Errorf("failed to get player class: %w", err)
	}

	return pc, nil
}

// GetCharacterSkills retrieves all skills for a given character.
func (d *PlayerCharacterDAL) GetCharacterSkills(characterID string) ([]*models.PlayerSkill, error) {
	query := `SELECT player_id, skill_id, percentage, granted_by_entity_type, granted_by_entity_id FROM PlayerSkills WHERE player_id = ?`
	rows, err := d.db.Query(query, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player skills: %w", err)
	}
	defer rows.Close()

	var playerSkills []*models.PlayerSkill
	for rows.Next() {
		ps := &models.PlayerSkill{}
		err := rows.Scan(
			&ps.PlayerID,
			&ps.SkillID,
			&ps.Percentage,
			&ps.GrantedByEntityType,
			&ps.GrantedByEntityID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player skill: %w", err)
		}
		playerSkills = append(playerSkills, ps)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through player skills: %w", err)
	}

	return playerSkills, nil
}

// GetCharacterQuestState retrieves a quest state for a given character.
func (d *PlayerCharacterDAL) GetCharacterQuestState(characterID, questID string) (*models.PlayerQuestState, error) {
	query := `SELECT player_id, quest_id, current_progress, last_action_timestamp, questmaker_influence_accumulated, status FROM PlayerQuestStates WHERE player_id = ? AND quest_id = ?`
	row := d.db.QueryRow(query, characterID, questID)

	pqs := &models.PlayerQuestState{}
	err := row.Scan(
		&pqs.PlayerID,
		&pqs.QuestID,
		&pqs.CurrentProgress,
		&pqs.LastActionTimestamp,
		&pqs.QuestmakerInfluenceAccumulated,
		&pqs.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Player quest state not found
		}
		return nil, fmt.Errorf("failed to get player quest state: %w", err)
	}

	return pqs, nil
}