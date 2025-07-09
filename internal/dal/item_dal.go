package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// ItemDAL handles database operations for Item entities.
type ItemDAL struct {
	db    *sql.DB
	Cache *Cache
}

// NewItemDAL creates a new ItemDAL.
func NewItemDAL(db *sql.DB) *ItemDAL {
	return &ItemDAL{db: db, Cache: NewCache()}
}

// CreateItem inserts a new item into the database.
func (d *ItemDAL) CreateItem(item *models.Item) error {
	query := `
	INSERT INTO Items (id, name, description, type, properties)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		item.ID,
		item.Name,
		item.Description,
		item.Type,
		item.Properties,
	)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}
	d.Cache.Set(item.ID, item, 300) // Cache for 5 minutes
	return nil
}

// GetItemByID retrieves an item by its ID.
func (d *ItemDAL) GetItemByID(id string) (*models.Item, error) {
	if cachedItem, found := d.Cache.Get(id); found {
		if item, ok := cachedItem.(*models.Item); ok {
			return item, nil
		}
	}

	query := `SELECT id, name, description, type, properties FROM Items WHERE id = ?`
	row := d.db.QueryRow(query, id)

	item := &models.Item{}
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Type,
		&item.Properties,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Item not found
		}
		return nil, fmt.Errorf("failed to get item by ID: %w", err)
	}

	d.Cache.Set(item.ID, item, 300) // Cache for 5 minutes
	return item, nil
}

// UpdateItem updates an existing item in the database.
func (d *ItemDAL) UpdateItem(item *models.Item) error {
	query := `
	UPDATE Items
	SET name = ?, description = ?, type = ?, properties = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		item.Name,
		item.Description,
		item.Type,
		item.Properties,
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("item with ID %s not found for update", item.ID)
	}
	d.Cache.Delete(item.ID) // Invalidate cache on update
	return nil
}

// DeleteItem deletes an item from the database by its ID.
func (d *ItemDAL) DeleteItem(id string) error {
	query := `DELETE FROM Items WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("item with ID %s not found for deletion", id)
	}
	d.Cache.Delete(id) // Invalidate cache on delete
	return nil
}

// GetItemsInRoom retrieves all items located in a specific room.
func (d *ItemDAL) GetItemsInRoom(roomID string) ([]*models.Item, error) {
	// This query assumes that the 'properties' JSON field in the Items table
	// contains a 'location_room_id' key for items that are in a room.
	// This is a simplified approach and might need refinement based on actual item-location mapping.
	query := `SELECT id, name, description, type, properties FROM Items WHERE json_extract(properties, '$.location_room_id') = ?`
	rows, err := d.db.Query(query, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items in room: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		item := &models.Item{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Type,
			&item.Properties,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item row: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// GetAllItems retrieves all items from the database.
func (d *ItemDAL) GetAllItems() ([]*models.Item, error) {
	query := `SELECT id, name, description, type, properties FROM Items`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all items: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		item := &models.Item{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Type,
			&item.Properties,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through items: %w", err)
	}

	return items, nil
}
