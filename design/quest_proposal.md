# Questmaker System Proposal: Event-Driven with Influence Budget

This proposal outlines a dynamic Questmaker system where sentient LLM entities influence the game world based on player actions and quest progress. This system addresses the "Undefined Quest System" and "LLM Integration" criticisms from `critics.md` by providing a structured approach to quest management and LLM interaction.

## 1. Core Concept: Sentient Questmakers

Each quest will be associated with a unique "Questmaker" – an LLM-driven entity with its own personality, goals, and a quantifiable "Influence Budget." The Questmaker's primary role is to react to player actions (or inactions) related to its quest by spending its budget to trigger specific, LLM-generated world events, modify NPC behaviors, or send messages.

## 2. Data Models

### 2.1. `Questmaker` Struct

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

### 2.2. `Quest` Struct (Additions)

```go
type Quest struct {
    // ... existing quest fields ...
    QuestmakerID        string           // ID of the associated Questmaker
    InfluencePointsMap  map[string]float64 // Map of player actions to influence points granted (e.g., {"recovered_crate": 20})
}
```

### 2.3. `PlayerQuestState` Struct (Additions)

```go
type PlayerQuestState struct {
    // ... existing player quest fields ...
    QuestID                 string
    PlayerID                string
    LastQuestActionTimestamp time.Time // Timestamp of the last relevant player action for this quest
    QuestmakerInfluenceAccumulated float64 // Points player has "given" to the Questmaker
}
```

## 3. Questmaker Decision Logic & LLM Integration

The Questmaker's decision-making is entirely driven by its LLM, which acts as its "brain."

### 3.1. LLM Prompt Generation

A dedicated `QuestmakerMonitor` service will periodically (or reactively) compile a comprehensive prompt for the Questmaker's LLM. This prompt provides all necessary context for the LLM to make informed decisions:

*   **Questmaker's Core Identity:** `LLMPromptContext` (personality, goals, current disposition).
*   **Player Status:** Current location, recent actions (both quest-related and unrelated), time since last relevant quest action, current quest progress, inventory (especially quest items).
*   **Relevant World State:** Summary of the game world relevant to the quest (e.g., market conditions, NPC locations, environmental factors).
*   **Quest-Specific Entities:** Status and location of key NPCs and items related to the quest.
*   **Questmaker's Resources:** `CurrentInfluenceBudget` and `MaxInfluenceBudget`.

**Example Prompt Snippet (Conceptual):**

```
"You are The Spice Lord, a benevolent spirit of trade. Your goal is to ensure the safe delivery of goods. You are pleased by efficiency and direct action, but angered by delays, theft, and players who wander off-task. Your current influence budget is 25/100.

The player 'player_123' is currently in 'Forest Path'. They last made progress on your quest (tracking goblins to the cave) 60 minutes ago. Since then, they have been 'killed_boar' and 'talked_to_farmer', showing a clear deviation from your objective. You have 0/3 spice crates recovered.

Considering this, what actions should you take to guide or pressure the player, or to influence the world, using your available influence budget? Respond with a list of tool calls."
```

### 3.2. LLM Output: Conceptual Tool Calls

The LLM's output will be a structured JSON object representing a list of "conceptual tool calls." These are not direct API calls but rather a structured instruction set that the game engine's `QuestmakerActionProcessor` will interpret and execute. Each proposed action will have an associated `cost` that the LLM must consider against the Questmaker's `CurrentInfluenceBudget`.

**Example LLM Output (Structured JSON):**

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
      "tool_name": "change_npc_behavior",
      "parameters": {
        "npc_id": "goblin_scout_1",
        "behavior": "more_alert_and_aggressive"
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

## 4. Influence Budget Management

### 4.1. Accumulation

The Questmaker's `CurrentInfluenceBudget` increases based on player actions:

*   **Positive Quest Actions:** When a player performs an action directly contributing to the quest (e.g., accepting the quest, defeating a quest enemy, recovering a quest item), the `QuestManager` grants `InfluencePoints` to the associated Questmaker as defined in the `Quest.InfluencePointsMap`.
    *   Example: `Quest.InfluencePointsMap = {"recovered_crate": 20, "defeated_chieftain": 15}`.
*   **Time-Based Regeneration:** A small amount of `InfluencePoints` can be regenerated over time (e.g., `BudgetRegenRate` per game tick or per significant player action), representing the Questmaker's inherent drive.

### 4.2. Spending

A `QuestmakerActionProcessor` component is responsible for executing the LLM's `proposed_actions`:

1.  It receives the LLM's output.
2.  For each `proposed_action`, it checks if the Questmaker's `CurrentInfluenceBudget` is greater than or equal to the `action.cost`.
3.  If sufficient, the action is executed by calling the corresponding game engine system (e.g., `CommunicationSystem.SendMessage()`, `NPCManager.ChangeBehavior()`, `WorldEventManager.TriggerEvent()`).
4.  The `action.cost` is then deducted from `CurrentInfluenceBudget`.
5.  If insufficient, the action is skipped, and the Questmaker's LLM might be prompted again with the updated (lower) budget, potentially leading to less costly actions or no actions.

## 5. Player Action Influence & Questmaker Reaction

The core of the system is how player actions directly feed into the Questmaker's decision-making process.

### 5.1. Player Progress & Positive Reinforcement

*   **Mechanism:** The `QuestmakerMonitor` detects positive quest-related actions (e.g., recovering an item, completing a stage).
*   **Influence:** These actions grant `InfluencePoints` to the Questmaker.
*   **LLM Reaction:** The LLM is prompted with the positive progress and increased budget. It might then decide to spend budget on "positive" actions:
    *   **Example:** `send_message` (via an NPC) congratulating the player, `trigger_world_event` (e.g., a minor boon or favorable market condition in a relevant town).

### 5.2. Player Inactivity & Negative Reinforcement

*   **Mechanism:** The `QuestmakerMonitor` tracks `time_since_last_action_minutes` in `PlayerQuestState`. If this exceeds a threshold, or if the player deviates significantly from the quest path, the Questmaker is prompted.
*   **Influence:** While no direct "negative points" are accumulated, the *context* of inactivity/deviation in the prompt, combined with the Questmaker's personality, drives the LLM's decision.
*   **LLM Reaction:** The LLM, seeing the player's lack of progress and its own goals, will decide to spend budget on "negative" or "pressure" actions:
    *   **Example:** `send_message` (via an NPC) expressing impatience or subtle threats, `change_npc_behavior` (e.g., making quest enemies more aggressive or patrols more frequent), `trigger_world_event` (e.g., minor market disruptions, increased danger in relevant areas).

### 5.3. Quest Failure & Major Consequences

*   **Mechanism:** If a quest is explicitly failed (abandoned, critical item destroyed, time limit expired), the `QuestmakerMonitor` provides this critical context to the LLM.
*   **Influence:** This is a high-impact event that will likely trigger the Questmaker to spend a significant portion of its budget on severe consequences.
*   **LLM Reaction:** The LLM, given the failure and its personality (e.g., "angered by delays"), will decide on impactful, high-cost actions:
    *   **Example:** `trigger_world_event` (e.g., major market crash, widespread negative reputation for the player), `spawn_entity` (e.g., a powerful, persistent enemy tied to the failure), `change_npc_behavior` (e.g., quest-givers or related NPCs become hostile or refuse interaction).

## 6. Conceptual Tools for Questmakers

Questmakers utilize a set of conceptual tools to influence the game world. These tools are represented as structured JSON outputs from the LLM, which are then interpreted and executed by the `QuestmakerActionProcessor`. Each tool has an associated conceptual `cost` that the LLM must consider against its `CurrentInfluenceBudget`.

### 6.1. Player-Centric Tools

*   **`grant_player_reward`**
    *   **Purpose:** To grant the player skills, spells, or items as a reward or consequence.
    *   **Parameters:**
        *   `player_id`: The ID of the target player.
        *   `reward_type`: Enum (e.g., "skill", "spell", "item").
        *   `reward_id`: The ID of the specific skill, spell, or item to grant.
        *   `quantity`: (Optional) For items, the quantity to grant.
    *   **Conceptual Cost:** Low to High (e.g., 5-50), depending on the power/rarity of the reward.
    *   **Example LLM Output:**
        ```json
        {
          "tool_name": "grant_player_reward",
          "parameters": {
            "player_id": "player_123",
            "reward_type": "skill",
            "reward_id": "tracking_proficiency"
          },
          "cost": 15,
          "reason": "Player successfully tracked goblins, rewarding their diligence."
        }
        ```

*   **`QUESTMAKER_memorize`**
    *   **Purpose:** Records a private memory about a player's actions or progress related to this specific quest. This memory influences the Questmaker's future decisions and interactions with that player regarding the quest.
    *   **Parameters:**
        *   `player_id`: The ID of the player about whom the memory is being recorded.
        *   `memory_string`: A string describing the specific memory (e.g., "Player 'player_123' abandoned the quest for 2 hours to fish.").
    *   **Conceptual Cost:** Low (e.g., 5).
    *   **Example LLM Output:**
        ```json
        {
          "tool_name": "QUESTMAKER_memorize",
          "parameters": {
            "player_id": "player_123",
            "memory_string": "Player 'player_123' recovered the first spice crate efficiently."
          },
          "cost": 5,
          "reason": "Player demonstrated efficiency in quest progress."
        }
        ```

### 6.2. NPC Interaction Tools

*   **`send_message`**
    *   **Purpose:** To deliver a message to a player, potentially via a specific NPC.
    *   **Parameters:**
        *   `target_player_id`: The ID of the player to receive the message.
        *   `message`: The content of the message.
        *   `via_npc_id`: (Optional) The ID of an NPC through whom the message should be delivered.
    *   **Conceptual Cost:** Low (e.g., 5-15).
    *   **Example LLM Output:** (Already in document, kept for consistency)
        ```json
        {
          "tool_name": "send_message",
          "parameters": {
            "target_player_id": "player_123",
            "message": "Captain Elias grows impatient. The market awaits its spices. Do not dally, adventurer!",
            "via_npc_id": "captain_elias"
          },
          "cost": 10,
          "reason": "Player inactivity and deviation from quest line."
        }
        ```

*   **`change_npc_behavior_to_player`**
    *   **Purpose:** To modify an NPC's disposition or behavior specifically towards a player, including adding a memory about that player.
    *   **Parameters:**
        *   `npc_id`: The ID of the target NPC.
        *   `player_id`: The ID of the player whose relationship with the NPC is being modified.
        *   `behavior_type`: Enum (e.g., "friendly", "hostile", "neutral", "fearful", "helpful").
        *   `memory_entry`: (Optional) A string describing the specific memory the NPC gains about the player (e.g., "Player failed to recover spices, causing market disruption.").
    *   **Conceptual Cost:** Medium (e.g., 10-30), depending on the severity of the behavior change.
    *   **Example LLM Output:**
        ```json
        {
          "tool_name": "change_npc_behavior_to_player",
          "parameters": {
            "npc_id": "merchant_guild_rep",
            "player_id": "player_123",
            "behavior_type": "hostile",
            "memory_entry": "Player 'player_123' failed to recover the spice shipment, causing significant losses for the guild."
          },
          "cost": 25,
          "reason": "Quest failure, leading to negative reputation with merchants."
        }
        ```

*   **`change_npc_stats`**
    *   **Purpose:** To modify an NPC's core combat or non-combat statistics.
    *   **Parameters:**
        *   `npc_id`: The ID of the target NPC.
        *   `stat_changes`: A map of stat names to their new values or modifiers (e.g., `{"strength": "+5", "health": "100", "aggression_level": "high"}`).
    *   **Conceptual Cost:** Medium to High (e.g., 15-40), depending on the impact of the stat changes.
    *   **Example LLM Output:**
        ```json
        {
          "tool_name": "change_npc_stats",
          "parameters": {
            "npc_id": "goblin_chieftain",
            "stat_changes": {
              "health": "+50",
              "damage_modifier": "1.2"
            }
          },
          "cost": 35,
          "reason": "Player's prolonged inactivity allowed goblins to fortify their leader."
        }
        ```

### 6.3. World Manipulation Tools

*   **`change_room_info`**
    *   **Purpose:** To modify properties of a specific room, such as locking/unlocking doors, changing descriptions, or adding/removing environmental effects.
    *   **Parameters:**
        *   `room_id`: The ID of the target room.
        *   `property_changes`: A map of room properties to their new values (e.g., `{"door_exit_north_locked": true, "description_add": "A foul stench now fills the air."}`).
    *   **Conceptual Cost:** Medium (e.g., 10-30).
    *   **Example LLM Output:**
        ```json
        {
          "tool_name": "change_room_info",
          "parameters": {
            "room_id": "whispering_caves_entrance",
            "property_changes": {
              "door_exit_north_locked": true,
              "description_add": "A newly erected, crude wooden barricade blocks the path."
            }
          },
          "cost": 20,
          "reason": "Goblins fortified their position due to player's slow progress."
        }
        ```

*   **`trigger_world_event`**
    *   **Purpose:** To initiate a broader world event that affects multiple players or areas.
    *   **Parameters:**
        *   `event_type`: A predefined event type (e.g., "minor_market_fluctuation", "weather_storm", "bandit_raid").
        *   `details`: (Optional) Specific details or parameters for the event.
    *   **Conceptual Cost:** Medium to High (e.g., 20-50), depending on the event's impact.
    *   **Example LLM Output:** (Already in document, kept for consistency)
        ```json
        {
          "tool_name": "trigger_world_event",
          "parameters": {
            "event_type": "minor_market_fluctuation",
            "details": "Prices for common goods in Port Town increase slightly due to perceived scarcity."
          },
          "cost": 20,
          "reason": "Initial warning for prolonged delay."
        }
        ```

*   **`spawn_entity`**
    *   **Purpose:** To dynamically spawn an NPC or item into the world.
    *   **Parameters:**
        *   `entity_type`: The type of entity to spawn (e.g., "goblin_patrol", "cursed_relic").
        *   `location`: The room ID or coordinates where the entity should spawn.
        *   `quantity`: (Optional) For items or groups of NPCs, the quantity to spawn.
    *   **Conceptual Cost:** Medium to High (e.g., 25-60), depending on the power/significance of the spawned entity.
    *   **Example LLM Output:**
        ```json
        {
          "tool_name": "spawn_entity",
          "parameters": {
            "entity_type": "goblin_patrol",
            "location": "forest_path_01",
            "quantity": 3
          },
          "cost": 30,
          "reason": "Increased goblin activity due to player's failure to address the threat."
        }
        ```

## 7. Implementation Considerations

*   **`QuestmakerMonitor`:** A background service that manages prompting the LLMs based on game state changes and time.
*   **`QuestmakerActionProcessor`:** Interprets LLM output and interfaces with core game systems (NPC AI, World Event Manager, Communication System, Player Inventory/Skills, Room Manager).
*   **Tool Definitions:** Clear internal definitions for each "conceptual tool" (e.g., `send_message`, `change_npc_behavior`, `trigger_world_event`, `spawn_entity`) that the LLM can "call," including their parameters and associated `cost` ranges.
*   **Persistence:** Questmaker states and player quest states must be persisted in the database.

This system provides a robust framework for dynamic, LLM-driven quest experiences, directly addressing the design criticisms and adding a unique layer of interactivity to the MUD.
