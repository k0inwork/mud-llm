# Phase 2 Design: Lore System & Editor

## 1. Objectives

This phase focuses on building the systems required to create, store, and manage the world's lore. The goal is to empower game administrators with the ability to build a rich, detailed world without modifying any Go code. This involves creating the backend storage solution for lore and building the web-based editor to manage it.

## 2. Key Components to be Implemented

### 2.1. Lore Storage

*   **Storage Mechanism:** A file-based storage system will be used for simplicity and ease of editing outside the game. All lore entries will be stored in a structured format (e.g., a single JSON file `world_lore.json` or a directory of individual files). This can be upgraded to a database in the future if needed.
*   **Data Access Layer:** A Go module will be created to handle all interactions with the lore storage. It will provide functions to:
    *   `GetLoreByID(id string) (*Lore, error)`
    *   `GetLoreByType(type string) ([]*Lore, error)`
    *   `CreateLore(lore *Lore) error`
    *   `UpdateLore(lore *Lore) error`
    *   `DeleteLore(id string) error`
*   **Lore Caching:** An in-memory cache will be implemented to hold lore entries, reducing the need to read from disk on every request. The cache will be invalidated and reloaded whenever a change is made via the web editor.

### 2.2. Web Server & Editor

*   **Basic Web Server:** The existing lightweight web server in `main.go` will be expanded. It will serve static HTML/CSS/JS files for the editor interface.
*   **HTTP API Endpoints:** The web server will expose a set of RESTful API endpoints for managing lore. These endpoints will use the Data Access Layer to interact with the lore storage.
    *   `GET /api/lore`: Returns all lore entries.
    *   `GET /api/lore/:id`: Returns a single lore entry.
    *   `POST /api/lore`: Creates a new lore entry.
    *   `PUT /api/lore/:id`: Updates an existing lore entry.
    *   `DELETE /api/lore/:id`: Deletes a lore entry.
*   **Editor Front-End:** A simple, single-page web application will be created.
    *   It will be built using standard HTML, CSS, and vanilla JavaScript to keep it lightweight.
    *   The interface will provide a form for creating and editing lore entries, including fields for `ID`, `Type`, `AssociatedID`, and `Content`.
    *   It will display a list of all existing lore entries, with buttons to edit or delete them.
    *   All interactions with the server will happen via asynchronous JavaScript calls (e.g., using `fetch`) to the API endpoints.

## 3. Acceptance Criteria

1.  The server can successfully read from and write to the `world_lore.json` file (or other chosen file structure).
2.  The in-memory lore cache is successfully populated on startup and invalidated on any change.
3.  A user can navigate to the web editor in a browser.
4.  The web editor correctly lists all lore entries.
5.  A user can successfully create, update, and delete lore entries through the web interface, and the changes are persisted in the storage file.
6.  The MUD server's core logic can successfully query the Lore module to retrieve specific lore entries (though it will not yet use them in prompts).
