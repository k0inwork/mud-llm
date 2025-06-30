# Phase 3 Design: LLM Integration, Memory & Foundational Questmakers

## 1. Objectives

This phase is the core of the AI implementation. The goal is to establish a robust connection to the LLM, implement the prompt caching strategy for performance, build the multi-layered memory system that will form the basis of all NPC and Owner behavior, and establish the foundational elements of the **Questmaker system**. This phase also clarifies the distinction between direct player commands and LLM-driven tool usage.

## 2. Key Components to be Implemented

### 2.1. Go LLM API Client

*   An HTTP client will be implemented in Go, specifically for interacting with an OpenAI-compatible API endpoint.
*   The client will handle:
    *   Structuring the JSON request body according to the API specification.
    *   Adding the API key to the request headers securely (read from an environment variable).
    *   Sending the HTTP POST request.
    *   Parsing the JSON response and extracting the narrative and tool-call content.
*   The client's configuration (API endpoint, model name) will be managed via a simple config file or environment variables.

### 2.2. Prompt Construction and Caching

*   **Prompt Assembler:** A module responsible for constructing the full prompt string. It will:
    1.  Fetch the relevant Global and Scoped Lore from the DAL.
    2.  Fetch the entity's Personality Prompt from the DAL.
    3.  Fetch the entity's Memories about the target player from the DAL.
    4.  Fetch and format the entity's Available Tools from the DAL.
    5.  Combine all of the above into a single "base prompt" string.
*   **Cache Manager:**
    *   An in-memory cache (e.g., a `map[string]string`) will store the generated "base prompt" for each NPC/Owner/Questmaker.
    *   The cache key could be the entity's ID.
    *   The cache for a specific entity will be invalidated (deleted) whenever its underlying memories or lore are updated (triggered by DAL update notifications).
*   **Interaction Flow:**
    1.  On first interaction, the Prompt Assembler generates the base prompt, which is sent to the LLM and stored in the cache.
    2.  On subsequent interactions, the cached base prompt is retrieved, and only the new, dynamic player action is sent to the LLM API (leveraging the API's conversation history).

### 2.3. Tool Dispatcher and Memory Tools

*   **Tool Parser:** A function that takes the raw XML response from the LLM and extracts the JSON from the `<tools>` section.
*   **Tool Dispatcher:** A central Go function that:
    1.  Receives the parsed tool JSON.
    2.  Uses the `tool_name` to look up and call the corresponding Go function.
    3.  Passes the tool's parameters to the Go function.
*   **Memory Tool Implementation (LLM-Callable):** These tools are designed to be invoked *by the LLM* to modify the game state related to memories.
    *   `NPC_memorize(player_id, memory_string)`: Adds a string to the `MemoriesAboutPlayers` map of the calling NPC. This is for an NPC's personal observations.
    *   `OWNER_memorize(player_id, memory_string)`: Adds a string to the `MemoriesAboutPlayers` map of the calling Owner. This is for an Owner's private observations.
    *   `OWNER_memorize_dependables(player_id, memory_string)`: Iterates through all NPCs associated with the calling Owner and adds the memory string to each of their `MemoriesAboutPlayers` maps. This tool is used by Owners to broadcast their opinions.

### 2.4. Player Commands vs. LLM Tools Clarification

*   **Player Commands:** Direct player commands (e.g., `move`, `unlock`, `use skill`) are handled by the **Core Game Engine**. These commands directly modify the game state based on game rules (e.g., checking for keys to unlock a door). They do *not* involve an LLM call for their execution. The outcome of a player command *may* trigger an Action Significance event for relevant NPCs/Owners/Questmakers, which then involves the LLM and its tools.
*   **LLM Tools:** Tools defined in `OwnerTool`, `NPCTool`, and `QuestmakerTool` structs are exclusively for the **LLM to invoke**. They represent actions that AI entities can take within the game world. These tools are never directly called by player commands.

## 3. Questmaker System Integration

This phase lays the groundwork for the dynamic Questmaker system, enabling LLM-driven entities to influence quests and the game world.

### 3.1. Core Concept

Each quest will be associated with a unique "Questmaker" – an LLM-driven entity with its own personality, goals, and a quantifiable "Influence Budget." The Questmaker's primary role is to react to player actions (or inactions) related to its quest by spending its budget to trigger specific, LLM-generated world events, modify NPC behaviors, or send messages.

### 3.2. Data Models (Foundational)

*   **`Questmaker` Struct:**
    ```go
    type Questmaker struct {
        ID                  string           // Unique identifier (e.g., "the_spice_lord")
        Name                string           // Display name
        LLMPromptContext    string           // Defines personality, goals, and core directives for the LLM
        CurrentInfluenceBudget float64        // Current points available for actions
        MaxInfluenceBudget  float64        // Maximum capacity of influence points
        BudgetRegenRate     float64        // Points regenerated per game tick/significant player action
    }
    ```
*   **`Quest` Struct (Additions):**
    ```go
    type Quest struct {
        // ... existing quest fields ...
        QuestmakerID        string           // ID of the associated Questmaker
        InfluencePointsMap  map[string]float64 // Map of player actions to influence points granted (e.g., {"recovered_crate": 20})
    }
    ```
*   **`PlayerQuestState` Struct (Additions):**
    ```go
    type PlayerQuestState struct {
        // ... existing player quest fields ...
        QuestID                 string
        PlayerID                string
        LastQuestActionTimestamp time.Time // Timestamp of the last relevant player action for this quest
        QuestmakerInfluenceAccumulated float64 // Points player has "given" to the Questmaker
    }
    ```

### 3.3. Questmaker Decision Logic & LLM Interaction

A `QuestmakerMonitor` service will compile a comprehensive prompt for the Questmaker's LLM, providing context such as the Questmaker's identity, player status, relevant world state, quest-specific entities, and the Questmaker's current influence budget.

The LLM's output will be a structured JSON object representing a list of "conceptual tool calls." Each proposed action will have an associated `cost` that the LLM must consider against the Questmaker's `CurrentInfluenceBudget`.

**Example LLM Output (Conceptual Tool Calls):**

```json
{
  "questmaker_id": "the_spice_lord",
  "proposed_actions": [
    {
      "tool_name": "send_message",
      "parameters": {
        "target_player_id": "player_123",
        "message": "Captain Elias grows impatient. The market awaits its spices. Do not dally, adventurer!",
        "via_npc_id": "captain_elias"
      },
      "cost": 10,
      "reason": "Player inactivity and deviation from quest line."
    },
    {
      "tool_name": "change_npc_behavior_to_player",
      "parameters": {
        "npc_id": "goblin_scout_1",
        "player_id": "player_123",
        "behavior_type": "more_alert_and_aggressive",
        "memory_entry": "Player 'player_123' is lingering near the caves without engaging."
      },
      "cost": 15,
      "reason": "Player approaching cave but not engaging, allowing goblins to fortify."
    },
    {
      "tool_name": "trigger_world_event",
      "parameters": {
        "event_type": "minor_market_fluctuation",
        "details": "Prices for common goods in Port Town increase slightly due to perceived scarcity."
      },
      "cost": 20,
      "reason": "Initial warning for prolonged delay."
    }
  ]
}
```

### 3.4. Conceptual Tools for Questmakers (Foundational)

This phase will establish the integration points for the following conceptual tools, which will be fully implemented and refined in later phases:

*   **`grant_player_reward`**: Grant skills, spells, or items to the player.
*   **`send_message`**: Deliver messages to a player, potentially via an NPC.
*   **`change_npc_behavior_to_player`**: Modify an NPC's disposition and add memories about a player.
*   **`change_npc_stats`**: Modify an NPC's core combat or non-combat statistics.
*   **`change_room_info`**: Modify properties of a specific room (e.g., lock/unlock doors, change descriptions).
*   **`trigger_world_event`**: Initiate a broader world event.
*   **`spawn_entity`**: Dynamically spawn an NPC or item into the world.

## 4. Acceptance Criteria

1.  The server can successfully make an API call to a configured LLM endpoint and receive a response.
2.  The prompt caching mechanism works as expected: a base prompt is generated and cached on the first interaction, and subsequent interactions are faster.
3.  Cache invalidation works correctly when a memory is added (via DAL update notifications).
4.  The LLM can successfully call the `NPC_memorize`, `OWNER_memorize`, `OWNER_memorize_dependables`, and foundational Questmaker tools (e.g., `send_message`, `change_npc_behavior_to_player`).
5.  The game state correctly reflects the changes made by these tools (i.e., memories are saved to the correct entities in the database via DAL, and basic Questmaker actions are logged/simulated).
6.  The distinction between player commands and LLM tools is clear in the code structure, with player commands handled directly by the game engine and LLM tools only invoked by the AI.
7.  Basic Questmaker data models (`Questmaker`, additions to `Quest`, `PlayerQuestState`) are defined and persistable via the DAL.
8.  The `QuestmakerMonitor` can generate prompts for Questmakers based on player actions and quest state.

## 5. Test Data Requirements

To test the LLM integration, prompt caching, multi-layered memory system, and foundational Questmaker elements in Phase 3, the following data should be present in the database (created via the Phase 1 web editor):

### 5.1. LLM API Configuration

This configuration would typically be stored in environment variables or a server configuration file, not directly in the game database, but it's crucial for testing.

```
LLM_API_ENDPOINT=https://api.openai.com/v1/chat/completions
LLM_API_KEY=sk-YOUR_API_KEY
LLM_MODEL_NAME=gpt-3.5-turbo
```

### 5.2. Example Player (for interaction)

```json
{
  "ID": "player_alice",
  "Name": "Alice",
  "Race": "human",
  "Profession": "adventurer",
  "CurrentRoomID": "starting_room",
  "VisitedRoomIDs": {},
  "Inventory": [],
  "Health": 100,
  "MaxHealth": 100
}
```

### 5.3. Example NPC with Personality and Tools

This NPC will be used to test `NPC_memorize`.

```json
{
  "ID": "innkeeper_bob",
  "Name": "Innkeeper Bob",
  "Description": "A jovial innkeeper with a booming laugh.",
  "OwnerIDs": ["town_council_owner"],
  "MemoriesAboutPlayers": {
    "player_alice": ["Initial impression: Seems friendly."]
  },
  "AvailableTools": [
    {
      "Name": "NPC_memorize",
      "Description": "Records a personal memory about a player.",
      "Parameters": {
        "player_id": {"type": "string"},
        "memory_string": {"type": "string"}
      }
    }
  ],
  "PersonalityPrompt": "You are Innkeeper Bob. You love gossip and are very protective of your regulars. You are generally friendly but can be stern if provoked.",
  "Inventory": []
}
```

### 5.4. Example Owner with Private Memories and Broadcast Tool

This Owner will be used to test `OWNER_memorize` and `OWNER_memorize_dependables`.

```json
{
  "ID": "town_council_owner",
  "Name": "The Town Council",
  "Description": "The collective consciousness of the town's governing body.",
  "MonitoredAspect": "location",
  "AssociatedID": "town_square",
  "MemoriesAboutPlayers": {
    "player_alice": ["Initial assessment: Newcomer, potential asset."]
  },
  "AvailableTools": [
    {
      "Name": "OWNER_memorize",
      "Description": "Records a private memory about a player.",
      "Parameters": {
        "player_id": {"type": "string"},
        "memory_string": {"type": "string"}
      }
    },
    {
      "Name": "OWNER_memorize_dependables",
      "Description": "Broadcasts a memory about a player to all subordinate NPCs.",
      "Parameters": {
        "player_id": {"type": "string"},
        "memory_string": {"type": "string"}
      }
    }
  ]
}
```

### 5.5. Example Questmaker and Quest Data

```json
{
  "ID": "the_spice_lord",
  "Name": "The Spice Lord",
  "LLMPromptContext": "I am the benevolent spirit of trade, guardian of fair commerce. My goal is to ensure the safe delivery of goods. I am pleased by efficiency and direct action, but angered by delays, theft, and players who wander off-task.",
  "CurrentInfluenceBudget": 0,
  "MaxInfluenceBudget": 100,
  "BudgetRegenRate": 5
}
```

```json
{
  "ID": "the_missing_shipment",
  "Name": "The Missing Shipment",
  "QuestmakerID": "the_spice_lord",
  "InfluencePointsMap": {
    "accepted_quest": 5,
    "tracked_goblins_to_cave": 10,
    "recovered_crate": 20,
    "defeated_chieftain": 15
  }
}
```

### 5.6. Example Lore Entries (from Phase 1/2, used in prompts)

Ensure relevant global, zone (`town_square` lore), and profession lore (`adventurer` lore) exists to provide context for the LLM prompts.

### 5.7. Testing Scenarios with Data

*   **Test `NPC_memorize`:** Simulate an LLM call for `innkeeper_bob` where it uses `NPC_memorize` after `player_alice` pays for a drink. Verify `innkeeper_bob`'s memories are updated in the database.
*   **Test `OWNER_memorize`:** Simulate an LLM call for `town_council_owner` where it uses `OWNER_memorize` after `player_alice` completes a minor quest. Verify `town_council_owner`'s private memories are updated.
*   **Test `OWNER_memorize_dependables`:** Simulate an LLM call for `town_council_owner` where it uses `OWNER_memorize_dependables("player_alice", "Has proven to be a reliable ally.")`. Verify that `innkeeper_bob`'s `MemoriesAboutPlayers` for `player_alice` now includes this broadcasted memory.
*   **Test Prompt Caching:** Trigger multiple interactions with `innkeeper_bob` by `player_alice`. Monitor API calls to ensure the base prompt is only sent once, and subsequent calls are smaller (only dynamic content).
*   **Test Questmaker Data Persistence:** Create a `Questmaker` and a `Quest` linked to it. Verify they are correctly persisted in the database.
*   **Test Questmaker Prompt Generation:** Simulate player actions for "The Missing Shipment" quest (e.g., inactivity, recovering a crate). Verify that the `QuestmakerMonitor` generates appropriate prompts for "The Spice Lord" LLM, including player status, quest progress, and influence budget.
*   **Test Questmaker Conceptual Tool Calls:** Simulate an LLM response from "The Spice Lord" containing conceptual tool calls (e.g., `send_message`, `change_npc_behavior_to_player`). Verify that the `QuestmakerActionProcessor` correctly interprets these and logs/simulates their execution, and that the Questmaker's influence budget is appropriately deducted.
*   **Test Influence Budget Accumulation:** Perform positive quest actions (e.g., `recovered_crate`). Verify that "The Spice Lord"'s `CurrentInfluenceBudget` increases as defined in `Quest.InfluencePointsMap`.
*   **Test Influence Budget Spending Logic:** Simulate an LLM response with actions exceeding the `CurrentInfluenceBudget`. Verify that only affordable actions are "executed" and the budget is correctly reduced.