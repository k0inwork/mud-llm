package dal

import (
	"database/sql"
	"fmt"
	"mud/internal/models"
)

// RoomDAL handles database operations for Room entities.
type RoomDAL struct {
	db    *sql.DB
	cache *Cache
}

// NewRoomDAL creates a new RoomDAL.
func NewRoomDAL(db *sql.DB) *RoomDAL {
	return &RoomDAL{db: db, cache: NewCache()}
}

// CreateRoom inserts a new room into the database.
func (d *RoomDAL) CreateRoom(room *models.Room) error {
	query := `
	INSERT INTO Rooms (id, name, description, exits, owner_id, properties)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := d.db.Exec(query,
		room.ID,
		room.Name,
		room.Description,
		room.Exits,
		room.OwnerID,
		room.Properties,
	)
	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}
	d.cache.Set(room.ID, room, 300) // Cache for 5 minutes
	return nil
}

// GetRoomByID retrieves a room by its ID.
func (d *RoomDAL) GetRoomByID(id string) (*models.Room, error) {
	if cachedRoom, found := d.cache.Get(id); found {
		if room, ok := cachedRoom.(*models.Room); ok {
			return room, nil
		}
	}

	query := `SELECT id, name, description, exits, owner_id, properties FROM Rooms WHERE id = ?`
	row := d.db.QueryRow(query, id)

	room := &models.Room{}
	err := row.Scan(
		&room.ID,
		&room.Name,
		&room.Description,
		&room.Exits,
		&room.OwnerID,
		&room.Properties,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Room not found
		}
		return nil, fmt.Errorf("failed to get room by ID: %w", err)
	}

	d.cache.Set(room.ID, room, 300) // Cache for 5 minutes
	return room, nil
}

// GetAllRooms retrieves all rooms from the database.
func (d *RoomDAL) GetAllRooms() ([]*models.Room, error) {
	query := `SELECT id, name, description, exits, owner_id, properties FROM Rooms`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all rooms: %w", err)
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Description,
			&room.Exits,
			&room.OwnerID,
			&room.Properties,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan room: %w", err)
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating through rooms: %w", err)
	}

	return rooms, nil
}

// UpdateRoom updates an existing room in the database.
func (d *RoomDAL) UpdateRoom(room *models.Room) error {
	query := `
	UPDATE Rooms
	SET name = ?, description = ?, exits = ?, owner_id = ?, properties = ?
	WHERE id = ?
	`

	result, err := d.db.Exec(query,
		room.Name,
		room.Description,
		room.Exits,
		room.OwnerID,
		room.Properties,
		room.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("room with ID %s not found for update", room.ID)
	}
	d.cache.Delete(room.ID) // Invalidate cache on update
	return nil
}

// DeleteRoom deletes a room from the database by its ID.
func (d *RoomDAL) DeleteRoom(id string) error {
	query := `DELETE FROM Rooms WHERE id = ?`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete room: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("room with ID %s not found for deletion", id)
	}
	d.cache.Delete(id) // Invalidate cache on delete
	return nil
}
