package dal

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"mud/internal/models"
	"mud/internal/testutils"
)

// MockItemDAL for player_character_dal_test
type MockItemDALForCharacter struct {
	GetItemByIDFunc func(id string) (*models.Item, error)
}

func (m *MockItemDALForCharacter) GetItemByID(id string) (*models.Item, error) {
	if m.GetItemByIDFunc != nil {
		return m.GetItemByIDFunc(id)
	}
	return nil, nil
}
func (m *MockItemDALForCharacter) GetAllItems() ([]*models.Item, error) { return nil, nil }
func (m *MockItemDALForCharacter) CreateItem(item *models.Item) error { return nil }
func (m *MockItemDALForCharacter) UpdateItem(item *models.Item) error { return nil }
func (m *MockItemDALForCharacter) DeleteItem(id string) error { return nil }
func (m *MockItemDALForCharacter) Cache() CacheInterface { return &testutils.MockCache{} }

func setupPlayerCharacterTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE player_accounts (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			hashed_password TEXT NOT NULL,
			email TEXT UNIQUE,
			created_at TIMESTAMP NOT NULL,
			last_login_at TIMESTAMP
		);
		CREATE TABLE player_characters (
			id TEXT PRIMARY KEY,
			player_account_id TEXT NOT NULL,
			name TEXT NOT NULL UNIQUE,
			race_id TEXT NOT NULL,
			profession_id TEXT NOT NULL,
			current_room_id TEXT NOT NULL,
			health INTEGER NOT NULL,
			max_health INTEGER NOT NULL,
			inventory TEXT NOT NULL,
			visited_room_ids TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			last_played_at TIMESTAMP,
			FOREIGN KEY (player_account_id) REFERENCES player_accounts(id)
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func TestPlayerCharacterDAL_CRUD(t *testing.T) {
	db, cleanup := setupPlayerCharacterTestDB(t)
	defer cleanup()

	mockCache := testutils.NewMockCache()
	mockItemDAL := &MockItemDALForCharacter{}
	characterDAL := NewPlayerCharacterDAL(db, mockCache, mockItemDAL)
	accountDAL := NewPlayerAccountDAL(db)

	// Create an account first
	account, err := accountDAL.CreateAccount("testuser", "password123", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	inventory := []string{"item1", "item2"}
	inventoryJSON, _ := json.Marshal(inventory)
	visitedRooms := []string{"room1"}
	visitedRoomsJSON, _ := json.Marshal(visitedRooms)

	newCharacter := &models.PlayerCharacter{
		ID:              uuid.New().String(),
		PlayerAccountID: account.ID,
		Name:            "TestCharacter",
		RaceID:          "human",
		ProfessionID:    "mage",
		CurrentRoomID:   "start_room",
		Health:          100,
		MaxHealth:       100,
		Inventory:       string(inventoryJSON),
		VisitedRoomIDs:  string(visitedRoomsJSON),
		CreatedAt:       time.Now(),
	}

	// Test Create
	err = characterDAL.CreateCharacter(newCharacter)
	if err != nil {
		t.Fatalf("CreateCharacter failed: %v", err)
	}

	// Test Read
	retrievedChar, err := characterDAL.GetCharacterByID(newCharacter.ID)
	if err != nil {
		t.Fatalf("GetCharacterByID failed: %v", err)
	}
	if retrievedChar == nil {
		t.Fatal("GetCharacterByID returned nil")
	}
	if retrievedChar.Name != "TestCharacter" {
		t.Errorf("Expected name %s, got %s", "TestCharacter", retrievedChar.Name)
	}

	// Test Get by Account ID
	chars, err := characterDAL.GetCharactersByAccountID(account.ID)
	if err != nil {
		t.Fatalf("GetCharactersByAccountID failed: %v", err)
	}
	if len(chars) != 1 {
		t.Fatalf("Expected 1 character for account, got %d", len(chars))
	}
	if chars[0].Name != "TestCharacter" {
		t.Errorf("Expected character name %s, got %s", "TestCharacter", chars[0].Name)
	}

	// Test Update
	retrievedChar.Health = 90
	retrievedChar.LastPlayedAt = time.Now()
	err = characterDAL.UpdateCharacter(retrievedChar)
	if err != nil {
		t.Fatalf("UpdateCharacter failed: %v", err)
	}

	updatedChar, err := characterDAL.GetCharacterByID(newCharacter.ID)
	if err != nil {
		t.Fatalf("GetCharacterByID after update failed: %v", err)
	}
	if updatedChar.Health != 90 {
		t.Errorf("Expected health 90, got %d", updatedChar.Health)
	}

	// Test Delete
	err = characterDAL.DeleteCharacter(newCharacter.ID)
	if err != nil {
		t.Fatalf("DeleteCharacter failed: %v", err)
	}

	deletedChar, err := characterDAL.GetCharacterByID(newCharacter.ID)
	if err != nil {
		t.Fatalf("GetCharacterByID after delete failed: %v", err)
	}
	if deletedChar != nil {
		t.Fatal("Character was not deleted")
	}
}

func TestPlayerCharacterDAL_GetCharacterInventory(t *testing.T) {
	db, cleanup := setupPlayerCharacterTestDB(t)
	defer cleanup()

	mockCache := testutils.NewMockCache()
	mockItemDAL := &MockItemDALForCharacter{
		GetItemByIDFunc: func(id string) (*models.Item, error) {
			if id == "sword_of_testing" {
				return &models.Item{ID: "sword_of_testing", Name: "Sword of Testing", Type: "weapon", Properties: "{}"}, nil
			}
			return nil, nil
		},
	}
	characterDAL := NewPlayerCharacterDAL(db, mockCache, mockItemDAL)
	accountDAL := NewPlayerAccountDAL(db)

	account, err := accountDAL.CreateAccount("inv_user", "password", "inv@test.com")
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	inventory := []string{"sword_of_testing"}
	inventoryJSON, _ := json.Marshal(inventory)

	newCharacter := &models.PlayerCharacter{
		ID:              uuid.New().String(),
		PlayerAccountID: account.ID,
		Name:            "InvCharacter",
		RaceID:          "elf",
		ProfessionID:    "ranger",
		CurrentRoomID:   "forest",
		Health:          100,
		MaxHealth:       100,
		Inventory:       string(inventoryJSON),
		VisitedRoomIDs:  "[]",
		CreatedAt:       time.Now(),
	}
	err = characterDAL.CreateCharacter(newCharacter)
	if err != nil {
		t.Fatalf("Failed to create character for inventory test: %v", err)
	}

	items, err := characterDAL.GetCharacterInventory(newCharacter.ID)
	if err != nil {
		t.Fatalf("GetCharacterInventory failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("Expected 1 item in inventory, got %d", len(items))
	}
	if items[0].Name != "Sword of Testing" {
		t.Errorf("Expected item 'Sword of Testing', got '%s'", items[0].Name)
	}
}