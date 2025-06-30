# Phase 1 Design: Foundation & Content Tools

## 1. Objectives

This foundational phase establishes the core persistence layer, essential data structures, the server-side presentation mechanism, and the initial content creation tools. The primary goal is to have a functional backend capable of storing all game data in a local database and allowing administrators to create initial world content via a web editor.

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
*   **Go Structs:** Implement the corresponding Go structs for all entities, mirroring the database schema. This includes the `MemoriesAboutPlayers` map on both `NPC` and `Owner` structs.

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

### 2.4. Web Server & Content Editor

*   **Web Server:** A lightweight Go web server will be set up to serve the administrative interface and the content editor.
*   **Direct DAL Interaction:** The web editor's backend will interact directly with the Data Access Layer (DAL) to perform CRUD operations on all game entities (Lore, Rooms, Items, NPCs, Owners, Tools). There will be no separate REST API endpoints for lore or other data; the editor will use the DAL directly.
*   **Editor Front-End:** A simple, single-page web application (HTML, CSS, vanilla JavaScript) will be created. This editor will allow administrators to:
    *   Create, read, update, and delete `Rooms`, `Exits`, `Items`, `NPCs`, and `Owners` directly via the DAL.
    *   Create initial `Lore` entries (Global, Race, Profession, Zone, Faction, Creature, Item Lore).
*   **Initial Test Content:** Basic starting rooms, a few NPCs, and some lore entries will be created using this editor to provide a minimal, runnable world for early testing.

## 3. Acceptance Criteria

1.  The Go project compiles successfully with all defined data structures.
2.  A local database is successfully initialized and can store data for all core entities.
3.  The DAL provides functional CRUD operations for all core entities.
4.  A player can connect to the MUD using a standard Telnet client.
5.  The server can generate a semantic JSON object (e.g., a hardcoded room description) and have it correctly rendered by the Telnet Renderer and displayed to the client with appropriate ANSI codes.
6.  The web editor is accessible via a browser.
7.  Administrators can use the web editor to successfully create, modify, and delete all game entities (Lore, Rooms, Items, NPCs, Owners, Tools) through its interface, with changes persisted in the database.
8.  A minimal set of initial game content (rooms, NPCs, lore) is created and loadable by the server.

## 4. Test Data Requirements

To test the functionality implemented in Phase 1, the following data structures (represented here in a JSON-like format for clarity, but stored in the database) should be creatable and manageable via the web editor:

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

#### 4.5.3. Race Lore

```json
{
  "ID": "human_traits",
  "Type": "race",
  "AssociatedID": "human",
  "Content": "Humans are known for their adaptability and ambition, spreading across Aerthos faster than any other race. They are often seen as resourceful but sometimes impulsive."
}
```

#### 4.5.4. Profession Lore

```json
{
  "ID": "adventurer_code",
  "Type": "profession",
  "AssociatedID": "adventurer",
  "Content": "Adventurers are driven by curiosity and the thrill of discovery. They often seek ancient ruins, forgotten treasures, and new challenges. They are expected to uphold a basic code of conduct, respecting ancient sites and aiding those in need."
}
```