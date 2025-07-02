package dal

import (
	"testing"
	"mud/internal/models"
)

func TestOwnerDAL_AdvancedQueries(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ownerDAL := NewOwnerDAL(db)

	// Seed with test data
	owner1 := &models.Owner{ID: "owner1", Name: "Owner 1", MonitoredAspect: "location", AssociatedID: "town_square"}
	owner2 := &models.Owner{ID: "owner2", Name: "Owner 2", MonitoredAspect: "faction", AssociatedID: "guild_of_mages"}
	owner3 := &models.Owner{ID: "owner3", Name: "Owner 3", MonitoredAspect: "location", AssociatedID: "forest_path"}

	ownerDAL.CreateOwner(owner1)
	ownerDAL.CreateOwner(owner2)
	ownerDAL.CreateOwner(owner3)

	// Test GetOwnersByMonitoredAspect
	locationOwners, err := ownerDAL.GetOwnersByMonitoredAspect("location", "town_square")
	if err != nil {
		t.Fatalf("GetOwnersByMonitoredAspect failed: %v", err)
	}
	if len(locationOwners) != 1 {
		t.Errorf("Expected 1 owner for location 'town_square', got %d", len(locationOwners))
	}
	if locationOwners[0].ID != "owner1" {
		t.Errorf("Expected owner with ID 'owner1', got '%s'", locationOwners[0].ID)
	}

	factionOwners, err := ownerDAL.GetOwnersByMonitoredAspect("faction", "guild_of_mages")
	if err != nil {
		t.Fatalf("GetOwnersByMonitoredAspect failed: %v", err)
	}
	if len(factionOwners) != 1 {
		t.Errorf("Expected 1 owner for faction 'guild_of_mages', got %d", len(factionOwners))
	}
	if factionOwners[0].ID != "owner2" {
		t.Errorf("Expected owner with ID 'owner2', got '%s'", factionOwners[0].ID)
	}
}
