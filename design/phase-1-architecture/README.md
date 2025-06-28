# Phase 1 Design: Foundation & Content Tools

## 1. Objectives

This initial phase is crucial for establishing the project's core infrastructure. The primary goals are to set up the local database for all game data, implement the foundational Go data structures, build the server-side presentation layer, and develop the web-based editor for content creation. This phase aims to deliver a functional administrative tool and a basic server capable of displaying semantically formatted output.

## 2. Key Components to be Implemented

### 2.1. Database Setup and Data Access Layer (DAL)

*   **Local Database:** A local database (e.g., SQLite) will be chosen and configured to store all persistent game data.
*   **Schema Definition:** The database schema will be defined for all core entities:
    *   `Player`
    *   `NPC`
    *   `Owner`
    *   `Room`
    *   `Exit`
    *   `Item`
    *   `Lore`
    *   Tables for `OwnerTool` and `NPCTool` definitions.
*   **Data Access Layer (DAL):** A Go module will be created to abstract all database interactions. It will provide CRUD (Create, Read, Update, Delete) functions for all defined entities. All other server components will interact with the database exclusively through the DAL.

### 2.2. Go Data Structures

All Go structs corresponding to the database entities will be implemented, including:
*   `Player` (with `Race`, `Profession`, `VisitedRoomIDs`, `Inventory`, `Health`, `MaxHealth`)
*   `NPC` (with `OwnerIDs`, `MemoriesAboutPlayers`, `AvailableTools`, `PersonalityPrompt`, `Inventory`)
*   `Owner` (with `MonitoredAspect`, `AssociatedID`, `AvailableTools`, `MemoriesAboutPlayers`)
*   `OwnerTool`
*   `NPCTool`
*   `Room` (with `OwnerID`, `Items`)
*   `Exit` (with `IsLocked`, `KeyID`)
*   `Item` (with `Attributes`)
*   `Lore` (with `ID`, `Type`, `AssociatedID`, `Content`)

### 2.3. Server-Side Presentation Layer

This layer is responsible for translating semantic game events into client-specific output.
*   **Semantic JSON Format:** A clear, versioned definition of the internal semantic JSON format will be established. This will include schemas for all message types (e.g., `narrative`, `room_update`, `player_stats`). The JSON will contain semantic types (e.g., `magic_keyword`, `neutral_npc`) rather than direct styling information.
*   **Telnet Renderer Module:** A Go module responsible for:
    *   Receiving semantic JSON objects from the core game logic.
    *   Maintaining a configurable mapping of `semantic_type` to specific ANSI escape codes for color and style.
    *   Constructing the final, raw string with embedded ANSI codes to be sent to the Telnet client.
*   **Main Server Loop Integration:** The main server loop will be structured to correctly route game events through the presentation layer (Core Logic -> Semantic JSON -> Telnet Renderer -> Client).

### 2.4. Web Server (Admin & Editor)

*   **Basic Web Server:** A lightweight Go web server will be implemented to serve static HTML/CSS/JS files for the editor interface.
*   **Direct DAL Interaction:** The web editor's backend will interact directly with the Data Access Layer (DAL) to perform CRUD operations on all game entities (Lore, Rooms, Items, NPCs, Owners, Tools).
*   **Editor Front-End:** A simple, single-page web application will be created using standard HTML, CSS, and vanilla JavaScript. It will provide:
    *   Forms for creating and editing all game entities.
    *   Lists of existing entities with edit/delete functionalities.
    *   A user-friendly interface for populating the game world.

### 2.5. Initial Test Content

*   Using the newly developed web editor, a minimal set of test content will be created to validate the system:
    *   At least one starting room.
    *   A simple item.
    *   A basic NPC.
    *   A few global and scoped lore entries.

## 3. Acceptance Criteria

1.  The local database is successfully initialized and accessible by the Go application.
2.  All core Go data structures are defined and can be persisted to and loaded from the database via the DAL.
3.  A player can connect to the MUD using a standard Telnet client.
4.  The server can generate a semantic JSON object (e.g., a hardcoded room description) and have it correctly rendered and displayed to the connected Telnet client with the appropriate ANSI colors and styles.
5.  The web editor is fully functional, allowing administrators to create, read, update, and delete all game entities (Lore, Rooms, Items, NPCs, Owners, Tools) through its interface.
6.  The initial test content is successfully created and loaded into the game world via the editor.
7.  The project compiles successfully, and the server can handle multiple Telnet connections without crashing.