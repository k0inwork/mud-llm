package server

import (
	"bytes"
	"encoding/json"
	"mud/internal/dal"
	"mud/internal/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
)

// setupTestServer creates a new AdminWebServer with a temporary database for testing.
func setupTestServer(t *testing.T) (*AdminWebServer, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "testdb_*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file for test database: %v", err)
	}

	db, err := dal.InitDB(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	server := NewAdminWebServer("8080", db)

	return server, func() {
		db.Close()
		os.Remove(tmpfile.Name())
	}
}

func TestRoomAPI(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/rooms", server.handleCreateRoom).Methods("POST")
	api.HandleFunc("/rooms/{id}", server.handleGetRoom).Methods("GET")
	api.HandleFunc("/rooms/{id}", server.handleUpdateRoom).Methods("PUT")
	api.HandleFunc("/rooms/{id}", server.handleDeleteRoom).Methods("DELETE")

	// 1. Test CreateRoom
	room := &models.Room{ID: "test_room", Name: "Test Room"}
	body, _ := json.Marshal(room)
	req, _ := http.NewRequest("POST", "/api/v1/rooms", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// 2. Test GetRoom
	req, _ = http.NewRequest("GET", "/api/v1/rooms/test_room", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedRoom models.Room
	json.Unmarshal(rr.Body.Bytes(), &returnedRoom)
	if returnedRoom.Name != "Test Room" {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedRoom.Name, "Test Room")
	}

	// 3. Test UpdateRoom
	updatedRoom := &models.Room{ID: "test_room", Name: "Updated Test Room"}
	body, _ = json.Marshal(updatedRoom)
	req, _ = http.NewRequest("PUT", "/api/v1/rooms/test_room", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 4. Test DeleteRoom
	req, _ = http.NewRequest("DELETE", "/api/v1/rooms/test_room", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestItemAPI(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/items", server.handleCreateItem).Methods("POST")
	api.HandleFunc("/items/{id}", server.handleGetItem).Methods("GET")
	api.HandleFunc("/items/{id}", server.handleUpdateItem).Methods("PUT")
	api.HandleFunc("/items/{id}", server.handleDeleteItem).Methods("DELETE")

	// 1. Test CreateItem
	item := &models.Item{ID: "test_item", Name: "Test Item"}
	body, _ := json.Marshal(item)
	req, _ := http.NewRequest("POST", "/api/v1/items", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// 2. Test GetItem
	req, _ = http.NewRequest("GET", "/api/v1/items/test_item", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedItem models.Item
	json.Unmarshal(rr.Body.Bytes(), &returnedItem)
	if returnedItem.Name != "Test Item" {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedItem.Name, "Test Item")
	}

	// 3. Test UpdateItem
	updatedItem := &models.Item{ID: "test_item", Name: "Updated Test Item"}
	body, _ = json.Marshal(updatedItem)
	req, _ = http.NewRequest("PUT", "/api/v1/items/test_item", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 4. Test DeleteItem
	req, _ = http.NewRequest("DELETE", "/api/v1/items/test_item", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestNPCAPI(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/npcs", server.handleCreateNPC).Methods("POST")
	api.HandleFunc("/npcs/{id}", server.handleGetNPC).Methods("GET")
	api.HandleFunc("/npcs/{id}", server.handleUpdateNPC).Methods("PUT")
	api.HandleFunc("/npcs/{id}", server.handleDeleteNPC).Methods("DELETE")

	// 1. Test CreateNPC
	npc := &models.NPC{ID: "test_npc", Name: "Test NPC"}
	body, _ := json.Marshal(npc)
	req, _ := http.NewRequest("POST", "/api/v1/npcs", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// 2. Test GetNPC
	req, _ = http.NewRequest("GET", "/api/v1/npcs/test_npc", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedNPC models.NPC
	json.Unmarshal(rr.Body.Bytes(), &returnedNPC)
	if returnedNPC.Name != "Test NPC" {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedNPC.Name, "Test NPC")
	}

	// 3. Test UpdateNPC
	updatedNPC := &models.NPC{ID: "test_npc", Name: "Updated Test NPC"}
	body, _ = json.Marshal(updatedNPC)
	req, _ = http.NewRequest("PUT", "/api/v1/npcs/test_npc", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 4. Test DeleteNPC
	req, _ = http.NewRequest("DELETE", "/api/v1/npcs/test_npc", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestOwnerAPI(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/owners", server.handleCreateOwner).Methods("POST")
	api.HandleFunc("/owners/{id}", server.handleGetOwner).Methods("GET")
	api.HandleFunc("/owners/{id}", server.handleUpdateOwner).Methods("PUT")
	api.HandleFunc("/owners/{id}", server.handleDeleteOwner).Methods("DELETE")

	// 1. Test CreateOwner
	owner := &models.Owner{ID: "test_owner", Name: "Test Owner"}
	body, _ := json.Marshal(owner)
	req, _ := http.NewRequest("POST", "/api/v1/owners", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// 2. Test GetOwner
	req, _ = http.NewRequest("GET", "/api/v1/owners/test_owner", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedOwner models.Owner
	json.Unmarshal(rr.Body.Bytes(), &returnedOwner)
	if returnedOwner.Name != "Test Owner" {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedOwner.Name, "Test Owner")
	}

	// 3. Test UpdateOwner
	updatedOwner := &models.Owner{ID: "test_owner", Name: "Updated Test Owner"}
	body, _ = json.Marshal(updatedOwner)
	req, _ = http.NewRequest("PUT", "/api/v1/owners/test_owner", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 4. Test DeleteOwner
	req, _ = http.NewRequest("DELETE", "/api/v1/owners/test_owner", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestLoreAPI(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/lore", server.handleCreateLore).Methods("POST")
	api.HandleFunc("/lore/{id}", server.handleGetLore).Methods("GET")
	api.HandleFunc("/lore/{id}", server.handleUpdateLore).Methods("PUT")
	api.HandleFunc("/lore/{id}", server.handleDeleteLore).Methods("DELETE")

	// 1. Test CreateLore
	lore := &models.Lore{ID: "test_lore", Title: "Test Lore"}
	body, _ := json.Marshal(lore)
	req, _ := http.NewRequest("POST", "/api/v1/lore", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// 2. Test GetLore
	req, _ = http.NewRequest("GET", "/api/v1/lore/test_lore", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedLore models.Lore
	json.Unmarshal(rr.Body.Bytes(), &returnedLore)
	if returnedLore.Title != "Test Lore" {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedLore.Title, "Test Lore")
	}

	// 3. Test UpdateLore
	updatedLore := &models.Lore{ID: "test_lore", Title: "Updated Test Lore"}
	body, _ = json.Marshal(updatedLore)
	req, _ = http.NewRequest("PUT", "/api/v1/lore/test_lore", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 4. Test DeleteLore
	req, _ = http.NewRequest("DELETE", "/api/v1/lore/test_lore", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}