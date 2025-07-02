package dal

import (
	"testing"
	"mud/internal/models"
)

func TestPlayerDAL_AdvancedQueries(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	playerDAL := NewPlayerDAL(db)
	itemDAL := NewItemDAL(db)

	// Seed items for inventory
	item1 := &models.Item{ID: "sword_of_testing", Name: "Sword of Testing", Type: "weapon", Properties: "{}"}
	item2 := &models.Item{ID: "shield_of_dev", Name: "Shield of Dev", Type: "armor", Properties: "{}"}
	itemDAL.CreateItem(item1)
	itemDAL.CreateItem(item2)

	// Seed player with inventory
	playerInventory := `[{"item_id": "sword_of_testing", "quantity": 1}, {"item_id": "shield_of_dev", "quantity": 1}]`
	newPlayer := &models.Player{
		ID:            "player1",
		Name:          "TestPlayer",
		RaceID:        "human",
		ProfessionID:  "warrior",
		CurrentRoomID: "starting_room",
		Health:        100,
		MaxHealth:     100,
		Inventory:     playerInventory,
		VisitedRoomIDs: "[]",
	}
	playerDAL.CreatePlayer(newPlayer)

	// Test GetPlayerInventory
	inventoryItems, err := playerDAL.GetPlayerInventory("player1")
	if err != nil {
		t.Fatalf("GetPlayerInventory failed: %v", err)
	}

	if len(inventoryItems) != 2 {
		t.Errorf("Expected 2 items in inventory, got %d", len(inventoryItems))
	}

	foundSword := false
	foundShield := false
	for _, item := range inventoryItems {
		if item.ID == "sword_of_testing" {
			foundSword = true
		}
		if item.ID == "shield_of_dev" {
			foundShield = true
		}
	}

	if !foundSword || !foundShield {
		t.Errorf("Did not find expected items in inventory. Sword found: %t, Shield found: %t", foundSword, foundShield)
	}
}
