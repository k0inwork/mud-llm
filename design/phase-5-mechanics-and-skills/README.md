# Phase 5 Design: Advanced Mechanics & Skills

## 1. Objectives

This phase focuses on enriching the gameplay by implementing key interactive mechanics and the skills system. The goal is to move beyond simple social interaction and introduce more complex, rules-based gameplay elements that are still deeply integrated with the LLM-driven AI.

## 2. Key Components to be Implemented

### 2.1. Locking Mechanism

*   **Data Model:** The `Exit` struct's `IsLocked` and `KeyID` fields will be utilized.
*   **Core Logic:** The `movePlayer` function in the core game engine will be updated to:
    1.  Check if the target exit `IsLocked`.
    2.  If it is, prevent the player from moving and send a "The way is locked." message (via the presentation layer).
*   **`unlock` Command/Tool:**
    *   A new player command, `unlock <direction> with <item>`, will be implemented.
    *   This will trigger a new tool, `TOOL_UNLOCK_EXIT(player_id, exit_direction, item_id)`.
    *   The Go function for this tool will:
        1.  Verify the player has the specified item.
        2.  Check if the item's `ID` matches the exit's `KeyID`.
        3.  If they match, set the exit's `IsLocked` to `false` and send a success message.
        4.  Otherwise, send a failure message.

### 2.2. Skills System

*   **Passive Skills Implementation:**
    *   **Effect on NPCs/Owners:** The `Prompt Assembler` will be modified. Before generating a prompt, it will check the player's passive skills. Skills like "Noble Bearing" will cause extra text to be appended to the context (e.g., "The player carries themselves with a noble air."). Skills like "Stealth" might cause the entire interaction to be aborted if a check fails.
    *   **Effect on Player:** The core game logic will be modified. Before generating the semantic JSON for a room description, it will check the player's skills. A player with "Arcane Sight" will have extra `semantic_type` data added to magical items in the JSON, which the renderer can then use to display them differently.
*   **Active Skills Implementation:**
    *   Active skills will be defined as `OwnerTool` or `NPCTool` entries in the lore/data files.
    *   A new `use <skill_name>` command will be implemented.
    *   This command will trigger the LLM to use the corresponding tool, allowing the AI to narrate the skill's use while the tool's Go function applies the mechanical effect (e.g., healing, dealing damage).

### 2.3. Mapping

*   **Data Model:** The `Player` struct's `VisitedRoomIDs` map will be used.
*   **Core Logic:** Whenever a player successfully enters a new room, the room's ID will be added to their `VisitedRoomIDs` map.
*   **`map` Command:** A new `map` command will be implemented that generates and displays a simple ASCII map of the player's visited rooms. This logic will reside within the Telnet presentation layer, as it's a purely visual feature.

## 3. Acceptance Criteria

1.  A player cannot move through a locked exit.
2.  Using the correct key item with the `unlock` command successfully unlocks an exit.
3.  A passive skill like "Noble Bearing" can be shown to positively influence an NPC's initial reaction.
4.  A passive skill like "Arcane Sight" causes magical items to be displayed differently to the player.
5.  A player can use an active skill like "Minor Heal" (implemented as a tool), and it will correctly restore health and display a narrative from the LLM.
6.  The `map` command displays a coherent ASCII map of all rooms the player has visited.
