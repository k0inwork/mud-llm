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
        2.  Verify if the player possesses the specified `item` (retrieved via DAL).
        3.  Check if the `item.ID` matches the `Exit.KeyID`.
        4.  If all conditions are met, set `Exit.IsLocked` to `false` in the database via the DAL.
        5.  Send a semantic JSON success/failure message to the player via the Server-Side Presentation Layer.
    *   This action *may* trigger an Action Significance event for relevant NPCs/Owners (e.g., "Player unlocked the treasury door").
*   **NPC `NPC_UNLOCK_EXIT` Tool:**
    *   This is an `NPCTool` callable *only* by the LLM for NPCs. It represents an NPC's ability to bypass player-level restrictions.
    *   The Go function for this tool will:
        1.  Receive the target `exit_id` as a parameter from the LLM.
        2.  Retrieve the `Exit` from the DAL.
        3.  **Bypass the `KeyID` check.** NPCs can unlock doors without needing a physical key, representing their inherent knowledge or abilities.
        4.  Set `Exit.IsLocked` to `false` in the database via the DAL.
        5.  Send a semantic JSON message to relevant players (e.g., "The guard deftly unlocks the heavy iron gate.") via the Server-Side Presentation Layer.

### 2.2. Skills System

*   **Active Skills (Player Commands):**
    *   Player commands for active skills (e.g., `use minor heal`) are handled directly by the **Core Game Engine**.
    *   The engine will:
        1.  Check player prerequisites (e.g., mana, cooldowns, skill level - retrieved via DAL).
        2.  Directly apply game effects (e.g., restore health, deal damage) and update player/entity states in the database via DAL.
        3.  Send a semantic JSON message to the player and relevant observers (e.g., "You feel a surge of warmth as your wounds close.") via the Server-Side Presentation Layer.
    *   These actions *may* trigger an Action Significance event for relevant NPCs/Owners (e.g., "Player used a healing spell").
*   **Passive Skills: A Two-Way Street**
    *   Passive skills are not explicit commands but rather modifiers that affect both how the world perceives the player and how the player perceives the world. Their effects are integrated into the core game logic and prompt construction.
    *   **Effect on NPCs/Owners (via Prompt Assembler):**
        *   The Prompt Assembler (from Phase 3) will be enhanced.
        *   Before constructing a prompt for an NPC/Owner about a player, it will query the player's passive skills from the DAL.
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
2.  The `unlock` command correctly checks for the required `KeyID` item in the player's inventory (retrieved via DAL).
3.  An NPC, when prompted by the LLM, can successfully use the `NPC_UNLOCK_EXIT` tool to unlock a door without needing a key, and this action is reflected in the game world (persisted via DAL).
4.  Passive skills correctly influence how NPCs/Owners perceive the player (e.g., stealth reduces detection, social skills alter reactions in LLM prompts).
5.  Passive skills correctly influence how the player perceives the world (e.g., "Arcane Sight" adds magical details to semantic JSON descriptions).
6.  A player can use an active skill (e.g., "Minor Heal"), and it will correctly apply game effects and send appropriate semantic messages.
7.  The `map` command correctly displays a coherent ASCII map of all rooms the player has visited.
8.  All game state changes related to these mechanics are correctly persisted in the database via the DAL.

## 4. Test Data Requirements

To test the advanced mechanics and skills system in Phase 5, the following data should be configured via the web editor (from Phase 1) and loaded via the DAL (from Phase 2):

### 4.1. Example Player with Skills and Inventory

```json
{
  "ID": "player_elara",
  "Name": "Elara",
  "Race": "elf",
  "Profession": "ranger",
  "CurrentRoomID": "forest_path",
  "VisitedRoomIDs": {
    "forest_path": true,
    "hidden_grove": true,
    "ancient_ruins_entrance": true
  },
  "Inventory": [
    {
      "ID": "forest_key",
      "Name": "a moss-covered key",
      "Description": "A small, moss-covered key that smells faintly of pine.",
      "Attributes": {
        "is_key": true,
        "unlocks_id": "ancient_ruins_entrance_east_lock"
      }
    }
  ],
  "Health": 80,
  "MaxHealth": 100,
  "Skills": {
    "stealth": {"level": 5},
    "arcane_sight": {"level": 1},
    "minor_heal": {"level": 1},
    "keen_eyes": {"level": 3},
    "noble_bearing": {"level": 2}
  }
}
```

### 4.2. Example Room with Locked Exit

```json
{
  "ID": "ancient_ruins_entrance",
  "Name": "Ancient Ruins Entrance",
  "Description": "The entrance to ancient, crumbling ruins. A heavy, stone door blocks the path to the east.",
  "OwnerID": "ruins_guardian_owner",
  "Items": [],
  "Exits": [
    {
      "Direction": "west",
      "TargetRoomID": "forest_path",
      "IsLocked": false,
      "KeyID": ""
    },
    {
      "Direction": "east",
      "TargetRoomID": "ancient_ruins_hall",
      "IsLocked": true,
      "KeyID": "ancient_ruins_key",
      "LockID": "ancient_ruins_entrance_east_lock" 
    }
  ]
}
```

### 4.3. Example NPC with `NPC_UNLOCK_EXIT` Tool

```json
{
  "ID": "ruins_guardian_golem",
  "Name": "Stone Golem",
  "Description": "A silent, moss-covered stone golem standing guard.",
  "OwnerIDs": ["ruins_guardian_owner"],
  "MemoriesAboutPlayers": {},
  "AvailableTools": [
    {
      "Name": "NPC_memorize",
      "Description": "Records a personal memory about a player.",
      "Parameters": {
        "player_id": {"type": "string"},
        "memory_string": {"type": "string"}
      }
    },
    {
      "Name": "NPC_UNLOCK_EXIT",
      "Description": "Unlocks a specified exit, bypassing key requirements.",
      "Parameters": {
        "exit_id": {"type": "string"} 
      }
    }
  ],
  "PersonalityPrompt": "You are a silent, ancient guardian of the ruins, programmed to protect its secrets. You only act on direct commands from your Owner or to neutralize threats.",
  "Inventory": []
}
```

### 4.4. Example Room with Hidden Item (for Passive Skill `keen_eyes`)

```json
{
  "ID": "hidden_grove",
  "Name": "Hidden Grove",
  "Description": "A secluded grove, filled with ancient trees. A faint shimmer can be seen near the largest oak.",
  "OwnerID": "",
  "Items": [
    {
      "ID": "shimmering_orb",
      "Name": "a shimmering orb",
      "Description": "A small orb that pulses with faint magical energy.",
      "Attributes": {
        "is_magical": true,
        "hidden_by_skill": "keen_eyes", 
        "hidden_threshold": 2 
      }
    }
  ],
  "Exits": [
    {
      "Direction": "south",
      "TargetRoomID": "forest_path",
      "IsLocked": false,
      "KeyID": ""
    }
  ]
}
```

### 4.5. Example Skill Definitions (Stored in a dedicated `Skills` table)

These define the properties of skills, not instances of them on players.

```json
[
  {
    "ID": "stealth",
    "Name": "Stealth",
    "Type": "passive",
    "Description": "Allows the player to move unseen and unheard.",
    "Effects": {
      "npc_perception_modifier": {"type": "omission", "data": ["player_presence", "player_action_details"]},
      "llm_prompt_modifier": "When player has 'stealth' skill, reduce their visibility in NPC context."
    }
  },
  {
    "ID": "arcane_sight",
    "Name": "Arcane Sight",
    "Type": "passive",
    "Description": "Reveals magical auras and hidden enchantments.",
    "Effects": {
      "player_perception_modifier": {"type": "add_semantic_type", "data": {"item_attribute": "is_magical", "semantic_type": "magical_aura"}}
    }
  },
  {
    "ID": "minor_heal",
    "Name": "Minor Heal",
    "Type": "active",
    "Description": "Restores a small amount of health.",
    "Cost": {"mana": 10},
    "Cooldown": 5,
    "Effect": {"type": "heal", "amount": 20}
  },
  {
    "ID": "keen_eyes",
    "Name": "Keen Eyes",
    "Type": "passive",
    "Description": "Improves observation skills, revealing hidden details.",
    "Effects": {
      "player_perception_modifier": {"type": "reveal_hidden_item", "data": {"attribute": "hidden_by_skill", "threshold_attribute": "hidden_threshold"}}
    }
  },
  {
    "ID": "noble_bearing",
    "Name": "Noble Bearing",
    "Type": "passive",
    "Description": "Your dignified demeanor commands respect.",
    "Effects": {
      "llm_prompt_modifier": "When player has 'noble_bearing' skill, add 'The player carries themselves with a noble and commanding presence.' to the NPC/Owner context."
    }
  }
]
```

### 4.6. Testing Scenarios with Data

*   **Scenario 1: Player Unlock with Key:**
    1.  Ensure `player_elara` is in `ancient_ruins_entrance` and has `forest_key` in inventory.
    2.  `player_elara` attempts `unlock east with forest_key`.
    *Expected:* Exit `ancient_ruins_entrance_east_lock` becomes `IsLocked: false`. Player receives success message.

*   **Scenario 2: Player Unlock without Key (Failure):**
    1.  Ensure `player_elara` is in `ancient_ruins_entrance` but does *not* have `forest_key`.
    2.  `player_elara` attempts `unlock east with non_existent_key`.
    *Expected:* Exit remains locked. Player receives failure message.

*   **Scenario 3: NPC Unlock (Bypass Key):**
    1.  Ensure `ancient_ruins_entrance_east_lock` is `IsLocked: true`.
    2.  Simulate `ruins_guardian_golem`'s LLM calling `NPC_UNLOCK_EXIT("ancient_ruins_entrance_east_lock")`.
    *Expected:* Exit `ancient_ruins_entrance_east_lock` becomes `IsLocked: false`. Relevant players receive a message about the golem unlocking the door.

*   **Scenario 4: Player Active Skill Usage:**
    1.  Ensure `player_elara` has `minor_heal` skill and sufficient mana.
    2.  Set `player_elara`'s health to 50.
    3.  `player_elara` uses `use minor heal`.
    *Expected:* `player_elara`'s health increases by 20 (to 70). Player receives semantic message about healing.

*   **Scenario 5: Passive Skill - Player Perception (`keen_eyes`):**
    1.  Ensure `player_elara` (with `keen_eyes` level 3) is in `hidden_grove`.
    2.  Player `look`s around.
    *Expected:* The semantic JSON for `hidden_grove`'s description includes `shimmering_orb` with its details, as `keen_eyes` level 3 meets/exceeds `hidden_threshold` 2.
    3.  Create a new player `player_bob` with `keen_eyes` level 1.
    4.  `player_bob` enters `hidden_grove` and `look`s.
    *Expected:* The semantic JSON for `hidden_grove`'s description does *not* include `shimmering_orb`.

*   **Scenario 6: Passive Skill - NPC Perception (`stealth`):**
    1.  Place `player_elara` (with `stealth` skill) in a room with an NPC that normally reacts to `move` actions.
    2.  `player_elara` `move`s quietly.
    *Expected:* The Action Significance Monitor, when preparing the context for the NPC's LLM, omits or reduces the significance of `player_elara`'s movement due to `stealth` skill, potentially preventing an LLM trigger.

*   **Scenario 7: Passive Skill - NPC Reaction (`noble_bearing`):**
    1.  Place `player_elara` (with `noble_bearing` skill) in a room with a new NPC.
    2.  `player_elara` `talk`s to the NPC.
    *Expected:* The LLM prompt for the NPC includes the `noble_bearing` descriptive context, and the NPC's response is more deferential or respectful than it would be to a player without this skill.