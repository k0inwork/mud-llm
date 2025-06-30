# Phase 6 Design: Concurrency & Final Polish

## 1. Objectives

This final phase focuses on stability, performance, and the overall user experience. The goal is to transform the functional prototype into a robust, polished MUD that can handle multiple players smoothly and provides a high-quality, immersive experience.

## 2. Key Components to be Implemented

### 2.1. Concurrency and Performance

*   **Concurrent LLM Calls:**
    *   All calls to the LLM API will be refactored to run in their own Go routines.
    *   This is critical for actions that can trigger multiple AI responses, such as a `pray` command (prompting multiple Owners) or a `say` command (prompting multiple NPCs in a room).
    *   `sync.WaitGroup` or channels will be used to manage these concurrent calls and ensure the system doesn't block while waiting for responses.
*   **State Management Review:**
    *   A thorough review of the entire codebase will be conducted to identify any shared data that could be subject to race conditions.
    *   `sync.Mutex` or `sync.RWMutex` will be added to protect all shared game state (e.g., `Room` inventories, `Player` stats) from concurrent access by multiple player commands or AI tool calls. All database interactions will continue to go through the DAL, which will also ensure its own concurrency safety.

### 2.2. Conditional AI Delivery

*   A queueing system for AI responses will be implemented.
*   When an LLM response is generated for a player, instead of being sent immediately, it will be placed in a queue for that player.
*   The system will check for conditions before delivering the message. For example, an NPC's reaction to a player entering a room should only be delivered if the player is still in that room. This prevents messages from "following" players who have already moved on.

### 2.3. Final Polish and Content Expansion

*   **Telnet Client Enhancements:** The `TelnetRenderer` (from Phase 1) will be refined. A full set of `semantic_type`s will be defined, and the ANSI color/style map will be completed to ensure the text presentation is clear and aesthetically pleasing.
*   **Web Editor Finalization:** The web editor's UI/UX (from Phase 1) will be polished to make it as intuitive and easy to use as possible for game administrators.
*   **Initial World Content:** The initial world content created in Phase 1 will be expanded to provide a more comprehensive and engaging starting experience. This will include:
    *   A more detailed set of interconnected rooms.
    *   A wider variety of NPCs with distinct personalities, lore, and initial memories.
    *   Multiple Owners with defined domains and relationships.
    *   A simple, guided questline that demonstrates the core AI and mechanics.
    *   A rich set of lore entries to fully flesh out the starting area.

## 3. Acceptance Criteria

1.  The server can handle multiple simultaneous LLM API calls without blocking or crashing.
2.  All shared game state is demonstrably protected by mutexes, and no race conditions can be found during stress testing.
3.  The conditional AI delivery system works correctly (e.g., delayed messages are not delivered if the player has moved to a new room).
4.  The Telnet client's output is well-formatted, colorful, and easy to read, utilizing the full range of defined `semantic_type`s.
5.  The web editor is fully functional and user-friendly, allowing for comprehensive content management.
6.  A new player can log in, experience a stable and responsive game, and complete the expanded starting quest, demonstrating the integration of all core features.

## 4. Test Data Requirements

For Phase 6, the focus shifts from defining new data structures to populating the world with a significant amount of interconnected content to stress-test the system's performance, concurrency, and overall polish. The following outlines the *scale* and *interconnectedness* of data needed, building upon previous phases.

### 4.1. Expanded World Map (Rooms & Exits)

*   **Quantity:** At least 50-100 interconnected rooms, forming several distinct zones (e.g., a town, a forest, a dungeon).
*   **Complexity:** Include rooms with various features:
    *   Rooms with multiple exits (4+ directions).
    *   Rooms with locked exits requiring keys (from Phase 5).
    *   Rooms with hidden elements discoverable by passive skills (from Phase 5).
    *   Rooms with unique descriptions that leverage various `semantic_type`s.

### 4.2. Diverse Population (NPCs & Owners)

*   **Quantity:** At least 20-30 NPCs, distributed across various rooms and zones.
*   **Diversity:**
    *   NPCs with varied `PersonalityPrompt`s (e.g., friendly, hostile, neutral, quest-giver, merchant).
    *   NPCs associated with different `Owners`.
    *   NPCs with initial `MemoriesAboutPlayers` (some positive, some negative) to test reputation.
    *   NPCs with different `ReactionThreshold`s for the Action Significance Monitor.
*   **Owners:** At least 5-10 Owners, each monitoring different `MonitoredAspect`s (locations, races, professions, factions).
    *   Owners with pre-existing `MemoriesAboutPlayers`.
    *   Owners configured to use `OWNER_memorize_dependables` based on certain player actions.

### 4.3. Rich Item Economy

*   **Quantity:** At least 50-100 unique items, including:
    *   Standard items (weapons, armor, consumables).
    *   Quest items.
    *   Keys for locked exits.
    *   Magical items with `is_magical` attributes for `Arcane Sight` testing.
    *   Hidden items with `hidden_by_skill` and `hidden_threshold` attributes.
*   **Distribution:** Items placed in rooms, on NPCs, and in player inventories.

### 4.4. Comprehensive Lore Database

*   **Quantity:** At least 50-100 lore entries covering all `Type`s (global, race, profession, zone, faction, creature, item).
*   **Interconnectedness:** Lore entries should reference each other and provide deep background for the expanded world map, diverse population, and item economy.
*   **LLM Context:** Ensure lore is detailed enough to provide rich context for LLM responses, allowing for nuanced and consistent AI behavior.

### 4.5. Player Data for Stress Testing

*   **Quantity:** Create multiple player accounts (e.g., 5-10) with varied races, professions, and starting locations.
*   **Skills:** Players should have a mix of active and passive skills to test their effects on both perception and interaction.

### 4.6. Testing Scenarios with Expanded Data

*   **Scenario 1: Multi-Player Interaction:** Have several players simultaneously interact with the same NPC or in the same room, triggering multiple concurrent LLM calls. Monitor for race conditions and responsiveness.
*   **Scenario 2: Long-Term Reputation:** Have a player perform a series of actions (both positive and negative) over an extended period, interacting with various NPCs and Owners. Verify that the multi-layered memory system correctly tracks and influences reactions.
*   **Scenario 3: Questline Execution:** Guide a player through the simple questline, ensuring all AI interactions, mechanics (locking, mapping), and skill effects function as intended within the larger world context.
*   **Scenario 4: Content Validation:** Systematically traverse the entire expanded world map, interacting with all NPCs and items, to ensure all content loads correctly and behaves as designed.
*   **Scenario 5: Editor Stress Test:** Perform a large number of rapid create/update/delete operations on various entities via the web editor to ensure DAL caching and database integrity hold up under load.
