# Phase 2 Design: Lore & Data Logic

## 1. Objectives

This phase focuses on solidifying the Data Access Layer (DAL) and ensuring efficient retrieval of all game data, particularly lore, for use by the LLM and core game logic. While basic CRUD operations were established in Phase 1 for the web editor, Phase 2 deepens the DAL's capabilities to handle complex queries and optimize data access.

## 2. Key Components to be Implemented

### 2.1. Enhanced Data Access Layer (DAL)

*   **Advanced Query Methods:** The DAL will be extended to provide more sophisticated query capabilities beyond basic CRUD, enabling efficient retrieval of specific data sets required by the game engine and LLM integration.
    *   `GetLoreByTypeAndAssociatedID(loreType string, associatedID string) ([]*Lore, error)`: For retrieving specific scoped lore.
    *   `GetAllGlobalLore() ([]*Lore, error)`: For retrieving all global lore.
    *   `GetNPCsByRoom(roomID string) ([]*NPC, error)`
    *   `GetNPCsByOwner(ownerID string) ([]*NPC, error)`
    *   `GetOwnersByMonitoredAspect(aspectType string, associatedID string) ([]*Owner, error)`
    *   `GetItemsInRoom(roomID string) ([]*Item, error)`
    *   `GetPlayerInventory(playerID string) ([]*Item, error)`
*   **In-Memory Caching within DAL:** A caching layer will be implemented directly within the DAL for frequently accessed, relatively static data (e.g., lore entries, tool definitions, static entity properties). This will minimize direct database hits during gameplay.
    *   The cache will be populated on server startup.
    *   Cache invalidation mechanisms will be implemented, triggered by updates or deletions of data via the web editor (Phase 1).

### 2.2. Data Loading and Initialization

*   The server startup sequence will be refined to ensure all necessary game data (rooms, items, NPCs, Owners, Lore) is loaded from the database into memory (or the DAL's cache) upon application launch.
*   Error handling for database connection and initial data loading will be robust.

## 3. Acceptance Criteria

1.  All game entities can be efficiently queried from the database using the new advanced query methods in the DAL.
2.  The DAL's in-memory caching demonstrably reduces database load for frequently accessed data.
3.  The server successfully loads all initial game content from the database on startup, making it immediately available to the game engine.
4.  The web editor (from Phase 1) continues to function correctly with the enhanced DAL, and changes made through the editor correctly invalidate and update the DAL's cache.
5.  Unit tests are in place for the DAL's query and caching mechanisms.