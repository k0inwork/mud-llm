package dal

import (
	"database/sql"
	"encoding/json"
	"mud/internal/models"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates a temporary database for testing.
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// Create a temporary file for the SQLite database
	tmpfile, err := os.CreateTemp("", "testdb_*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file for test database: %v", err)
	}

	// Initialize the database
	db, err := InitDB(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	// Return the database and a cleanup function
	return db, func() {
		db.Close()
		os.Remove(tmpfile.Name())
	}
}

func TestRoomDAL(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	roomDAL := NewRoomDAL(db)

	// 1. Test CreateRoom
	exits, _ := json.Marshal(map[string]interface{}{"north": "room2"})
	newRoom := &models.Room{
		ID:          "room1",
		Name:        "Test Room",
		Description: "A room for testing.",
		Exits:       string(exits),
	}

	err := roomDAL.CreateRoom(newRoom)
	if err != nil {
		t.Fatalf("CreateRoom failed: %v", err)
	}

	// 2. Test GetRoomByID
	retrievedRoom, err := roomDAL.GetRoomByID("room1")
	if err != nil {
		t.Fatalf("GetRoomByID failed: %v", err)
	}
	if retrievedRoom == nil {
		t.Fatal("GetRoomByID failed: room not found")
	}
	if retrievedRoom.Name != "Test Room" {
		t.Errorf("GetRoomByID returned wrong name: got %v want %v", retrievedRoom.Name, "Test Room")
	}

	// 3. Test UpdateRoom
	retrievedRoom.Name = "Updated Test Room"
	err = roomDAL.UpdateRoom(retrievedRoom)
	if err != nil {
		t.Fatalf("UpdateRoom failed: %v", err)
	}

	updatedRoom, err := roomDAL.GetRoomByID("room1")
	if err != nil {
		t.Fatalf("GetRoomByID after update failed: %v", err)
	}
	if updatedRoom.Name != "Updated Test Room" {
		t.Errorf("UpdateRoom failed to update name: got %v want %v", updatedRoom.Name, "Updated Test Room")
	}

	// 4. Test DeleteRoom
	err = roomDAL.DeleteRoom("room1")
	if err != nil {
		t.Fatalf("DeleteRoom failed: %v", err)
	}

	deletedRoom, err := roomDAL.GetRoomByID("room1")
	if err != nil {
		t.Fatalf("GetRoomByID after delete failed: %v", err)
	}
	if deletedRoom != nil {
		t.Fatal("DeleteRoom failed: room was not deleted")
	}
}
