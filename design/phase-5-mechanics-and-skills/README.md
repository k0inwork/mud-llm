# Phase 5 Design: Advanced Mechanics & Skills

## 1. Objectives

This phase focuses on enriching the gameplay by implementing key interactive mechanics and the skills system. The goal is to move beyond simple social interaction and introduce more complex, rules-based gameplay elements. A key distinction in this phase is that player-initiated actions for these mechanics are handled directly by the Core Game Engine, while NPCs can use specialized LLM-driven tools to interact with them.

## 2. Key Components to be Implemented

### 2.1. Locking Mechanism

*   **Data Model:** The `Exit` struct's `IsLocked` and `KeyID` fields will be fully utilized and persisted via the DAL.
*   **Player `unlock` Command:**
    *   This is a direct player command handled by the **Core Game Engine**.
    *   When a player attempts `unlock <direction> with <item>`, the engine will:
        1.  Check if the target exit `IsLocked`.
        2.  Verify if the player possesses the specified `item`.
        3.  Check if the `item.ID` matches the `Exit.KeyID`.
        4.  If all conditions are met, set `Exit.IsLocked` to `false` in the database via the DAL.
        5.  Send a semantic JSON success/failure message to the player via the Server-Side Presentation Layer.
    *   This action *may* trigger an Action Significance event for relevant NPCs/Owners (e.g., "Player unlocked the treasury door").
*   **NPC `NPC_UNLOCK_EXIT` Tool:**
    *   This is an `NPCTool` callable *only* by the LLM for NPCs.
    *   The Go function for this tool will:
        1.  Receive the target `exit_id` as a parameter from the LLM.
        2.  **Bypass the `KeyID` check.** NPCs can unlock doors without needing a physical key, representing their inherent knowledge or abilities.
        3.  Set `Exit.IsLocked` to `false` in the database via the DAL.
        4.  Send a semantic JSON message to relevant players (e.g., "The guard deftly unlocks the heavy iron gate.").

### 2.2. Skills System

*   **Active Skills (Player Commands):**
    *   Player commands for active skills (e.g., `use minor heal`) are handled directly by the **Core Game Engine**.
    *   The engine will:
        1.  Check player prerequisites (e.g., mana, cooldowns, skill level).
        2.  Directly apply game effects (e.g., restore health, deal damage).
        3.  Send a semantic JSON message to the player and relevant observers (e.g., "You feel a surge of warmth as your wounds close.").
    *   These actions *may* trigger an Action Significance event for relevant NPCs/Owners (e.g., "Player used a healing spell").
*   **Passive Skills: A Two-Way Street**
    *   Passive skills are not explicit commands but rather modifiers that affect both how the world perceives the player and how the player perceives the world. Their effects are integrated into the core game logic and prompt construction.
    *   **Effect on NPCs/Owners (via Prompt Assembler):**
        *   The Prompt Assembler (from Phase 3) will be enhanced.
        *   When constructing a prompt for an NPC/Owner about a player, it will query the player's passive skills from the DAL.
        *   Skills like "Stealth" might cause information about the player's presence or specific actions to be omitted from the context provided to an NPC's LLM, making them less likely to be noticed or reacted to.
        *   Skills like "Noble Bearing" will cause descriptive context (e.g., "The player carries themselves with a noble air.") to be appended to the prompt, influencing the LLM's perception and subsequent narrative/tool usage.
    *   **Effect on the Player's Perception (via Core Game Engine/Semantic JSON):**
        *   The **Core Game Engine** will be responsible for filtering or adding information to the semantic JSON based on the player's passive skills *before* it is passed to the Server-Side Presentation Layer.
        *   Example: A player with the "Arcane Sight" skill might receive extra `semantic_type` data (e.g., `"semantic_type": "magical_aura"`) on items or NPCs in room descriptions, indicating magical properties that other players wouldn't see. A player with "Keen Eyes" might get a higher chance to notice a hidden lever in a room, with the server adding that detail to the room description JSON just for them.

### 2.3. Mapping

*   **Data Model:** The `Player` struct's `VisitedRoomIDs` map will be used to track explored rooms and will be persisted via the DAL.
*   **Core Logic:** Whenever a player successfully enters a new room, the room's ID will be added to their `VisitedRoomIDs` map (updated in DB via DAL).
*   **`map` Command:**
    *   This is a direct player command handled by the **Core Game Engine**.
    *   It retrieves the player's `VisitedRoomIDs` from the DAL.
    *   It generates a semantic JSON object representing the map data (e.g., `{"type": "map_data", "payload": {"visited_rooms": [...], "connections": [...]}}`).
    *   This semantic JSON is then sent to the Server-Side Presentation Layer, which will render it as an ASCII map for Telnet clients or a graphical map for web clients.

## 3. Acceptance Criteria

1.  Player commands for `unlock` and active skills are handled directly by the game engine, without triggering LLM calls for their execution.
2.  The `unlock` command correctly checks for the required `KeyID` item in the player's inventory.
3.  An NPC, when prompted by the LLM, can successfully use the `NPC_UNLOCK_EXIT` tool to unlock a door without needing a key, and this action is reflected in the game world.
4.  Passive skills correctly influence how NPCs/Owners perceive the player (e.g., stealth reduces detection, social skills alter reactions in LLM prompts).
5.  Passive skills correctly influence how the player perceives the world (e.g., "Arcane Sight" adds magical details to semantic JSON descriptions).
6.  A player can use an active skill (e.g., "Minor Heal"), and it will correctly apply game effects and send appropriate semantic messages.
7.  The `map` command correctly displays a coherent ASCII map of all rooms the player has visited.
8.  All game state changes related to these mechanics are correctly persisted in the database via the DAL.