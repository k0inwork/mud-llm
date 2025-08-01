# Phase 1 Design: Foundation & Content Tools

## 1. Objectives

This foundational phase establishes the core persistence layer, essential data structures, the server-side presentation mechanism, and the initial content creation tools. The primary goal is to have a functional backend capable of storing all game data in a local database and allowing administrators to create initial world content via a decoupled, API-driven web editor.

## 2. Key Components to be Implemented

### 2.1. Database Setup and Core Data Structures

*   **Local Database:** Initialize and configure a local database (e.g., SQLite) to store all persistent game data. This will be the single source of truth for the MUD.
*   **Schema Definition:** Define the database schema for all core game entities, including:
    *   `Players`
    *   `NPCs`
    *   `Owners`
    *   `Rooms`
    *   `Exits`
    *   `Items`
    *   `Lore`
    *   Tables for `OwnerTools`, `NPCTools`, and `Memories` (linking to Players, NPCs, and Owners).
*   **Go Structs:** Implement the corresponding Go structs for all entities, mirroring the database schema.
*   **Player State Persistence:** The player's dynamic state (current health, inventory, location, etc.) will be loaded from the database upon login and saved periodically (e.g., every few minutes) or upon logout via the DAL. This ensures continuity of player experience.

### 2.2. Data Access Layer (DAL) - Initial Implementation

*   A basic Data Access Layer (DAL) will be implemented in Go. This module will provide fundamental CRUD (Create, Read, Update, Delete) operations for all core entities directly with the database.
*   All other server components will interact with game data exclusively through the DAL.

### 2.3. Server-Side Presentation Layer

*   **Semantic JSON Format:** Define the universal, semantic JSON format for all server-to-client communication. This format describes *what* is being communicated (e.g., `narrative`, `room_update`, `player_stats`) using semantic types (e.g., `magic_keyword`, `neutral_npc`) rather than specific styling.
*   **Telnet Renderer Module:** A Go module responsible for:
    *   Receiving semantic JSON objects from the core game logic.
    *   Translating `semantic_type`s into specific ANSI escape codes for color and style based on a configurable ruleset.
    *   Constructing the final, raw string with embedded ANSI codes to be sent to the Telnet client.
*   **Main Server Loop Integration:** The main server loop will be structured to pass game events through the presentation layer (Core Logic -> Semantic JSON -> Telnet Renderer -> Client).

### 2.4. Web Server & Admin REST API

*   **Web Server:** A lightweight Go web server will be set up, running in the same process as the MUD server but on a different port. It will serve the static HTML/CSS/JS files for the editor.
*   **Admin REST API:** The web server will expose a versioned RESTful API (e.g., `/api/v1/`) for all content management. This API will be the *only* way the web editor interacts with the server for data operations.
    *   **Authentication:** A simple, token-based authentication middleware will protect all API endpoints to prevent unauthorized access.
    *   **Endpoints:** The API will provide standard CRUD endpoints for all major game entities. Examples include:
        *   `GET /api/v1/lore`, `POST /api/v1/lore`, `GET /api/v1/lore/:id`, `PUT /api/v1/lore/:id`, `DELETE /api/v1/lore/:id`
        *   Similar endpoints for `/rooms`, `/items`, `/npcs`, `/owners`, `/players`, `/tools` (for `OwnerTool` and `NPCTool` definitions).
    *   **API Handlers:** The Go functions that handle these API requests will call the appropriate DAL functions to interact with the database.
*   **Editor Front-End:** A simple, single-page web application (HTML, CSS, vanilla JavaScript) will be created. It will communicate exclusively through the REST API to create, read, update, and delete all game entities.

### 2.5. Initial Test Content

*   Using the newly developed web editor, a minimal set of test content (rooms, items, NPCs, lore) will be created and stored in the database to validate the entire system.

## 3. Acceptance Criteria

1.  The Go project compiles successfully with all defined data structures.
2.  A local database is successfully initialized and can store data for all core entities.
3.  The DAL provides functional CRUD operations for all core entities.
4.  A player can connect to the MUD using a standard Telnet client.
5.  The server can display data from the database to the Telnet client, rendered through the semantic JSON -> Telnet Renderer pipeline.
6.  The web editor is accessible via a browser and can successfully perform all CRUD operations on game entities by making authenticated calls to the REST API.
7.  The initial test content is successfully created and persisted in the database via the web editor.

## 4. Test Data Requirements

To test the functionality implemented in Phase 1, the following data structures (represented here in a JSON-like format for clarity) should be creatable and manageable via the web editor. This initial test data will be stored in a seed file (e.g., `seed.sql` or `seed.json`) that the application can use to populate a fresh database for testing.

### 4.1. Example Room

```json
{
  "ID": "starting_room",
  "Name": "The Dusty Cellar",
  "Description": "A small, dusty cellar with a single flickering torch. The air is damp and smells of old earth. A wooden door leads north.",
  "OwnerID": "cellar_owner",
  "Items": [],
  "Exits": [
    {
      "Direction": "north",
      "TargetRoomID": "cellar_exit_north",
      "IsLocked": false,
      "KeyID": ""
    }
  ]
}
```

### 4.2. Example Item

```json
{
  "ID": "rusty_key",
  "Name": "a rusty iron key",
  "Description": "A small, rusty iron key. It looks like it might open an old lock.",
  "Attributes": {
    "is_key": true,
    "unlocks_id": "cellar_exit_north_lock"
  }
}
```

### 4.3. Example NPC

```json
{
  "ID": "old_man_greg",
  "Name": "Old Man Greg",
  "Description": "An old man with a long, white beard, hunched over a workbench.",
  "OwnerIDs": ["cellar_owner"],
  "MemoriesAboutPlayers": {},
  "AvailableTools": [
    {
      "Name": "NPC_memorize",
      "Description": "Records a memory about a player.",
      "Parameters": {
        "player_id": {"type": "string"},
        "memory_string": {"type": "string"}
      }
    }
  ],
  "PersonalityPrompt": "You are Old Man Greg, a reclusive but kind old man who lives in the cellar. You are wary of strangers but will help those who seem genuine. You are very knowledgeable about local history.",
  "Inventory": []
}
```

### 4.4. Example Owner

```json
{
  "ID": "cellar_owner",
  "Name": "The Spirit of the Cellar",
  "Description": "A benevolent spirit that watches over the cellar and its inhabitants.",
  "MonitoredAspect": "location",
  "AssociatedID": "starting_room",
  "MemoriesAboutPlayers": {},
  "AvailableTools": [
    {
      "Name": "OWNER_memorize",
      "Description": "Records a private memory about a player.",
      "Parameters": {
        "player_id": {"type": "string"},
        "memory_string": {"type": "string"}
      }
    }
  ]
}
```

### 4.5. Example Lore Entries

#### 4.5.1. Global Lore

```json
{
  "ID": "world_creation_myth",
  "Type": "global",
  "AssociatedID": "",
  "Content": "In the beginning, there was only the Void, from which emerged the Twin Dragons, Ignis and Aqua. They wove the fabric of reality, creating the lands of Aerthos and the celestial spheres. Their eternal dance maintains the balance of magic and life."
}
```

#### 4.5.2. Zone Lore

```json
{
  "ID": "cellar_history",
  "Type": "zone",
  "AssociatedID": "starting_room",
  "Content": "This cellar was once part of an ancient wizard's tower, long since crumbled to dust. Whispers say the wizard's spirit still lingers, protecting forgotten secrets."
}
```