package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mud/internal/dal"
	"mud/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

// AdminWebServer represents the web server for the admin interface.
type AdminWebServer struct {
	port string
	db   *sql.DB
}

// NewAdminWebServer creates a new AdminWebServer.
func NewAdminWebServer(port string, db *sql.DB) *AdminWebServer {
	return &AdminWebServer{port: port, db: db}
}

// Start begins listening for incoming HTTP connections for the admin interface.
func (s *AdminWebServer) Start() {
	r := mux.NewRouter()
	r.Use(s.authMiddleware)
	api := r.PathPrefix("/api/v1").Subrouter()

	// Status
	api.HandleFunc("/status", s.handleStatus).Methods("GET")

	// Rooms
	api.HandleFunc("/rooms", s.handleCreateRoom).Methods("POST")
	api.HandleFunc("/rooms/{id}", s.handleGetRoom).Methods("GET")
	api.HandleFunc("/rooms/{id}", s.handleUpdateRoom).Methods("PUT")
	api.HandleFunc("/rooms/{id}", s.handleDeleteRoom).Methods("DELETE")

	// Items
	api.HandleFunc("/items", s.handleCreateItem).Methods("POST")
	api.HandleFunc("/items/{id}", s.handleGetItem).Methods("GET")
	api.HandleFunc("/items/{id}", s.handleUpdateItem).Methods("PUT")
	api.HandleFunc("/items/{id}", s.handleDeleteItem).Methods("DELETE")

	// NPCs
	api.HandleFunc("/npcs", s.handleCreateNPC).Methods("POST")
	api.HandleFunc("/npcs/{id}", s.handleGetNPC).Methods("GET")
	api.HandleFunc("/npcs/{id}", s.handleUpdateNPC).Methods("PUT")
	api.HandleFunc("/npcs/{id}", s.handleDeleteNPC).Methods("DELETE")

	// Owners
	api.HandleFunc("/owners", s.handleCreateOwner).Methods("POST")
	api.HandleFunc("/owners/{id}", s.handleGetOwner).Methods("GET")
	api.HandleFunc("/owners/{id}", s.handleUpdateOwner).Methods("PUT")
	api.HandleFunc("/owners/{id}", s.handleDeleteOwner).Methods("DELETE")

	// Lore
	api.HandleFunc("/lore", s.handleCreateLore).Methods("POST")
	api.HandleFunc("/lore/{id}", s.handleGetLore).Methods("GET")
	api.HandleFunc("/lore/{id}", s.handleUpdateLore).Methods("PUT")
	api.HandleFunc("/lore/{id}", s.handleDeleteLore).Methods("DELETE")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates")))

	fmt.Printf("Admin web server listening on port %s\n", s.port)
	log.Fatal(http.ListenAndServe(":"+s.port, r))
}

func (s *AdminWebServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Auth check for %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (s *AdminWebServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Admin API Status: OK"))
}

// Generic Handlers
func (s *AdminWebServer) handleCreate(w http.ResponseWriter, r *http.Request, model interface{}, createFunc func(interface{}) error) {
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := createFunc(model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model)
}

func (s *AdminWebServer) handleGet(w http.ResponseWriter, r *http.Request, getFunc func(string) (interface{}, error)) {
	id := mux.Vars(r)["id"]
	model, err := getFunc(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if model == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(model)
}

func (s *AdminWebServer) handleUpdate(w http.ResponseWriter, r *http.Request, model interface{}, updateFunc func(interface{}) error) {
	id := mux.Vars(r)["id"]
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// This is a bit of a hack to set the ID from the URL
	// A better solution would use reflection
	switch v := model.(type) {
	case *models.Room:
		v.ID = id
	case *models.Item:
		v.ID = id
	case *models.NPC:
		v.ID = id
	case *models.Owner:
		v.ID = id
	case *models.Lore:
		v.ID = id
	}

	if err := updateFunc(model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *AdminWebServer) handleDelete(w http.ResponseWriter, r *http.Request, deleteFunc func(string) error) {
	id := mux.Vars(r)["id"]
	if err := deleteFunc(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Room Handlers
func (s *AdminWebServer) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var room models.Room
	s.handleCreate(w, r, &room, func(m interface{}) error { return dal.NewRoomDAL(s.db).CreateRoom(m.(*models.Room)) })
}
func (s *AdminWebServer) handleGetRoom(w http.ResponseWriter, r *http.Request) {
	s.handleGet(w, r, func(id string) (interface{}, error) { return dal.NewRoomDAL(s.db).GetRoomByID(id) })
}
func (s *AdminWebServer) handleUpdateRoom(w http.ResponseWriter, r *http.Request) {
	var room models.Room
	s.handleUpdate(w, r, &room, func(m interface{}) error { return dal.NewRoomDAL(s.db).UpdateRoom(m.(*models.Room)) })
}
func (s *AdminWebServer) handleDeleteRoom(w http.ResponseWriter, r *http.Request) {
	s.handleDelete(w, r, dal.NewRoomDAL(s.db).DeleteRoom)
}

// Item Handlers
func (s *AdminWebServer) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	s.handleCreate(w, r, &item, func(m interface{}) error { return dal.NewItemDAL(s.db).CreateItem(m.(*models.Item)) })
}
func (s *AdminWebServer) handleGetItem(w http.ResponseWriter, r *http.Request) {
	s.handleGet(w, r, func(id string) (interface{}, error) { return dal.NewItemDAL(s.db).GetItemByID(id) })
}
func (s *AdminWebServer) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	s.handleUpdate(w, r, &item, func(m interface{}) error { return dal.NewItemDAL(s.db).UpdateItem(m.(*models.Item)) })
}
func (s *AdminWebServer) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	s.handleDelete(w, r, dal.NewItemDAL(s.db).DeleteItem)
}

// NPC Handlers
func (s *AdminWebServer) handleCreateNPC(w http.ResponseWriter, r *http.Request) {
	var npc models.NPC
	s.handleCreate(w, r, &npc, func(m interface{}) error { return dal.NewNPCDAL(s.db).CreateNPC(m.(*models.NPC)) })
}
func (s *AdminWebServer) handleGetNPC(w http.ResponseWriter, r *http.Request) {
	s.handleGet(w, r, func(id string) (interface{}, error) { return dal.NewNPCDAL(s.db).GetNPCByID(id) })
}
func (s *AdminWebServer) handleUpdateNPC(w http.ResponseWriter, r *http.Request) {
	var npc models.NPC
	s.handleUpdate(w, r, &npc, func(m interface{}) error { return dal.NewNPCDAL(s.db).UpdateNPC(m.(*models.NPC)) })
}
func (s *AdminWebServer) handleDeleteNPC(w http.ResponseWriter, r *http.Request) {
	s.handleDelete(w, r, dal.NewNPCDAL(s.db).DeleteNPC)
}

// Owner Handlers
func (s *AdminWebServer) handleCreateOwner(w http.ResponseWriter, r *http.Request) {
	var owner models.Owner
	s.handleCreate(w, r, &owner, func(m interface{}) error { return dal.NewOwnerDAL(s.db).CreateOwner(m.(*models.Owner)) })
}
func (s *AdminWebServer) handleGetOwner(w http.ResponseWriter, r *http.Request) {
	s.handleGet(w, r, func(id string) (interface{}, error) { return dal.NewOwnerDAL(s.db).GetOwnerByID(id) })
}
func (s *AdminWebServer) handleUpdateOwner(w http.ResponseWriter, r *http.Request) {
	var owner models.Owner
	s.handleUpdate(w, r, &owner, func(m interface{}) error { return dal.NewOwnerDAL(s.db).UpdateOwner(m.(*models.Owner)) })
}
func (s *AdminWebServer) handleDeleteOwner(w http.ResponseWriter, r *http.Request) {
	s.handleDelete(w, r, dal.NewOwnerDAL(s.db).DeleteOwner)
}

// Lore Handlers
func (s *AdminWebServer) handleCreateLore(w http.ResponseWriter, r *http.Request) {
	var lore models.Lore
	s.handleCreate(w, r, &lore, func(m interface{}) error { return dal.NewLoreDAL(s.db).CreateLore(m.(*models.Lore)) })
}
func (s *AdminWebServer) handleGetLore(w http.ResponseWriter, r *http.Request) {
	s.handleGet(w, r, func(id string) (interface{}, error) { return dal.NewLoreDAL(s.db).GetLoreByID(id) })
}
func (s *AdminWebServer) handleUpdateLore(w http.ResponseWriter, r *http.Request) {
	var lore models.Lore
	s.handleUpdate(w, r, &lore, func(m interface{}) error { return dal.NewLoreDAL(s.db).UpdateLore(m.(*models.Lore)) })
}
func (s *AdminWebServer) handleDeleteLore(w http.ResponseWriter, r *http.Request) {
	s.handleDelete(w, r, dal.NewLoreDAL(s.db).DeleteLore)
}