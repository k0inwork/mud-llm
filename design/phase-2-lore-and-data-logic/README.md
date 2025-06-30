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
    *   The cache will be populated on server startup by loading all relevant data from the database.
    *   Cache invalidation mechanisms will be implemented. These will be triggered by update or delete operations performed through the DAL (e.g., by the web editor from Phase 1), ensuring the cache remains consistent with the database.

### 2.2. Data Loading and Initialization

*   The server startup sequence will be refined to ensure all necessary game data (rooms, items, NPCs, Owners, Lore) is loaded from the database into memory (or the DAL's cache) upon application launch.
*   Robust error handling for database connection and initial data loading will be implemented.

## 3. Acceptance Criteria

1.  All game entities can be efficiently queried from the database using the new advanced query methods in the DAL.
2.  The DAL's in-memory caching demonstrably reduces database load for frequently accessed data during runtime.
3.  The server successfully loads all initial game content from the database on startup, making it immediately available to the game engine.
4.  Changes made through the web editor (from Phase 1) correctly invalidate and update the DAL's cache, and these changes are reflected in subsequent data retrievals.
5.  Unit tests are in place for the DAL's query and caching mechanisms, ensuring data integrity and performance.

## 4. Test Data Requirements

To test the enhanced DAL and lore retrieval in Phase 2, the following additional lore entries (beyond those from Phase 1) should be created via the web editor. These will allow for testing specific queries and cache behavior.

### 4.1. Example Lore Entries for Advanced Queries

#### 4.1.1. Global Lore (Additional)

```json
{
  "ID": "ancient_wars_summary",
  "Type": "global",
  "AssociatedID": "",
  "Content": "The Great Sundering, a cataclysmic war between the Elder Races and the Shadow Blight, reshaped the continents and led to the rise of human kingdoms. Many ancient artifacts were lost during this era."
}
```

#### 4.1.2. Zone Lore (Specific to a new zone)

```json
{
  "ID": "whispering_woods_secrets",
  "Type": "zone",
  "AssociatedID": "whispering_woods",
  "Content": "The Whispering Woods are ancient and enchanted, home to elusive dryads and mischievous sprites. Travelers often report strange lights and ethereal music deep within its groves. A hidden shrine to the Forest Mother is rumored to exist here."
}
```

#### 4.1.3. Profession Lore (Specific to a new profession)

```json
{
  "ID": "mage_guild_history",
  "Type": "profession",
  "AssociatedID": "mage",
  "Content": "The Conclave of Arcane Arts, the oldest mage guild, was founded after the Sundering to preserve magical knowledge. Its members are sworn to protect ancient magical sites and regulate the use of powerful spells."
}
```

#### 4.1.4. Faction Lore

```json
{
  "ID": "shadow_blight_origins",
  "Type": "faction",
  "AssociatedID": "shadow_blight",
  "Content": "The Shadow Blight is not merely a faction but a creeping corruption that seeks to consume all light and life. Its origins are shrouded in mystery, but ancient texts speak of a primordial darkness that predates even the Twin Dragons."
}
```

### 4.2. Testing Scenarios with Data

*   **Test `GetAllGlobalLore()`:** Ensure both `world_creation_myth` (from Phase 1) and `ancient_wars_summary` are retrieved.
*   **Test `GetLoreByTypeAndAssociatedID("zone", "whispering_woods")`:** Verify that `whispering_woods_secrets` is returned.
*   **Test Cache Invalidation:** Update `mage_guild_history` via the web editor and then immediately query it to confirm the cached version is updated.
*   **Test `GetNPCsByRoom()` and `GetItemsInRoom()`:** Create a new room and populate it with NPCs and items, then use these DAL methods to retrieve them, ensuring the relationships are correctly handled by the DAL.