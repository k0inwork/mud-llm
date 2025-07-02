package dal

import (
	"testing"
	"mud/internal/models"
)

func TestItemDAL_AdvancedQueries(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	itemDAL := NewItemDAL(db)

	// Seed with test data
	item1 := &models.Item{ID: "item1", Name: "Sword", Type: "weapon", Properties: `{"location_room_id": "roomA"}`}
	item2 := &models.Item{ID: "item2", Name: "Shield", Type: "armor", Properties: `{"location_room_id": "roomA"}`}
	item3 := &models.Item{ID: "item3", Name: "Potion", Type: "consumable", Properties: `{"location_room_id": "roomB"}`}

	itemDAL.CreateItem(item1)
	itemDAL.CreateItem(item2)
	itemDAL.CreateItem(item3)

	// Test GetItemsInRoom
	roomAItems, err := itemDAL.GetItemsInRoom("roomA")
	if err != nil {
		t.Fatalf("GetItemsInRoom failed: %v", err)
	}
	if len(roomAItems) != 2 {
		t.Errorf("Expected 2 items in roomA, got %d", len(roomAItems))
	}

	roomBItems, err := itemDAL.GetItemsInRoom("roomB")
	if err != nil {
		t.Fatalf("GetItemsInRoom failed: %v", err)
	}
	if len(roomBItems) != 1 {
		t.Errorf("Expected 1 item in roomB, got %d", len(roomBItems))
	}
}
