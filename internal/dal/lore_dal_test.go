package dal

import (
	"testing"
	"mud/internal/models"
)

func TestLoreDAL_AdvancedQueries(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	loreDAL := NewLoreDAL(db)

	// Seed with test data
	globalLore1 := &models.Lore{ID: "global1", Title: "Global 1", Content: "...", Scope: "global"}
	globalLore2 := &models.Lore{ID: "global2", Title: "Global 2", Content: "...", Scope: "global"}
	zoneLore1 := &models.Lore{ID: "zone1", Title: "Zone 1", Content: "...", Scope: "zone", AssociatedID: "zoneA"}
	zoneLore2 := &models.Lore{ID: "zone2", Title: "Zone 2", Content: "...", Scope: "zone", AssociatedID: "zoneB"}
	
	loreDAL.CreateLore(globalLore1)
	loreDAL.CreateLore(globalLore2)
	loreDAL.CreateLore(zoneLore1)
	loreDAL.CreateLore(zoneLore2)

	// Test GetAllLore and filter for global scope
	allLores, err := loreDAL.GetAllLore()
	if err != nil {
		t.Fatalf("GetAllLore failed: %v", err)
	}

	var globalLores []*models.Lore
	for _, lore := range allLores {
		if lore.Scope == "global" {
			globalLores = append(globalLores, lore)
		}
	}

	if len(globalLores) != 2 {
		t.Errorf("Expected 2 global lores, got %d", len(globalLores))
	}

	// Test GetLoreByTypeAndAssociatedID
	zoneALores, err := loreDAL.GetLoreByTypeAndAssociatedID("zone", "zoneA")
	if err != nil {
		t.Fatalf("GetLoreByTypeAndAssociatedID failed: %v", err)
	}
	if len(zoneALores) != 1 {
		t.Errorf("Expected 1 lore for zoneA, got %d", len(zoneALores))
	}
	if zoneALores[0].ID != "zone1" {
		t.Errorf("Expected lore with ID 'zone1', got '%s'", zoneALores[0].ID)
	}
}
