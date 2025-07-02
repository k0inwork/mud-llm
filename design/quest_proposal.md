# Quest System Proposal: Event-Driven with Layered Influence

This proposal outlines a dynamic quest system with a clear separation of responsibilities for quest initiation, thematic ownership, and execution control. It addresses the "Undefined Quest System" and "LLM Integration" criticisms from `critics.md` by providing a structured approach to quest management and LLM interaction, with a refined influence budget model.

## 1. Core Concepts: Owners, Quest Owners, and Questmakers

The system introduces three distinct types of entities involved in the quest lifecycle:

*   **Owners:** (Existing concept) These entities monitor specific aspects of the world (locations, races, professions) and can **initiate** quests by making them available to players. Their influence budget is primarily for their general world-monitoring and initiation activities.
*   **Quest Owners:** (NEW Concept) These are high-level, sentient LLM entities representing thematic ownership or overarching narrative arcs for quests. They **supervise** a series or group of related quests, providing narrative cohesion and driving global world changes. Their influence budget regenerates over time, allowing them to enact broad, strategic impacts.
*   **Questmakers:** (Refined Concept) These are sentient LLM entities responsible for controlling the **execution and progression** of a *single, specific* quest. There is a 1-to-1 relationship between a quest and its questmaker. Their influence budget primarily accumulates from player actions directly relevant to their quest, enabling them to react dynamically to player progress (or lack thereof).

## 2. Data Models

### 2.1. `Owner` Struct (Additions)

```go
type Owner struct {
    // ... existing owner fields ...
    InitiatedQuests      []string         // List of quest IDs this owner can initiate/offer
}
```

### 2.2. `QuestOwner` Struct (NEW)

```go
type QuestOwner struct {
    ID                   string           // Unique identifier (e.g., "gandalf_grand_plan")
    Name                 string           // Display name
    Description          string           // A brief description of the Quest Owner's role/theme
    LLMPromptContext     string           // Defines custom personality, goals, and core directives for the LLM (this is appended to a static system prompt)
    CurrentInfluenceBudget float64        // Current points available for global actions (time-based regeneration)
    MaxInfluenceBudget   float64        // Maximum capacity of influence points
    BudgetRegenRate      float64        // Points regenerated per game tick/significant time interval
    AssociatedQuestmakerIDs []string      // JSON array of Questmaker IDs under this Quest Owner's thematic umbrella
}
```

### 2.3. `Questmaker` Struct (Refined)

```go
type Questmaker struct {
    ID                  string           // Unique identifier (e.g., "urgent_message_questmaker")
    Name                string           // Display name
    LLMPromptContext    string           // Defines custom personality, goals, and core directives for the LLM (this is appended to a static system prompt)
    CurrentInfluenceBudget float64        // Current points available for quest-specific actions (player-action-based regeneration)
    MaxInfluenceBudget  float64        // Maximum capacity of influence points
    BudgetRegenRate     float64        // Points regenerated per game tick/significant player action (will be 0 or very low)
    MemoriesAboutPlayers map[string]string // Specific memories about players related to this quest
    AvailableTools      []string         // List of conceptual tools this Questmaker can use
}
```

### 2.4. `Quest` Struct (Additions)

```go
type Quest struct {
    // ... existing quest fields ...
    QuestOwnerID        string           // ID of the thematic Quest Owner
    QuestmakerID        string           // ID of the associated Questmaker (1-to-1)
    InfluencePointsMap  map[string]float64 // Map of player actions to influence points granted to the Questmaker
}
```

### 2.5. `PlayerQuestState` Struct (Additions)

```go
type PlayerQuestState struct {
    // ... existing player quest fields ...
    QuestID                 string
    PlayerID                string
    LastQuestActionTimestamp time.Time // Timestamp of the last relevant player action for this quest
    QuestmakerInfluenceAccumulated float64 // Points player has "given" to the Questmaker
}
```

## 3. Decision Logic & LLM Integration

The system employs two layers of LLM-driven decision-making: Quest Owners for strategic, global influence, and Questmakers for tactical, quest-specific control.

### 3.1. LLM Prompt Generation

Dedicated monitor services will compile comprehensive prompts for both Quest Owners and Questmakers. Each prompt will consist of a static, system-defined prefix followed by dynamic context and the entity's custom `LLMPromptContext` from the database.

*   **Static Quest Owner Prompt Prefix:**
    ```
    "You are a Quest Owner, a high-level entity responsible for supervising overarching narrative arcs and driving global world changes. You have a time-based influence budget to enact broad, strategic impacts. You can use the following tools: `trigger_world_event`, `spawn_entity` (for major entities), `change_room_info` (for significant global changes)."
    ```

*   **`QuestOwnerMonitor`:** Periodically (or reactively to major world events/quest completions) compiles a prompt for Quest Owners. This prompt includes:
    1.  The **Static Quest Owner Prompt Prefix**.
    2.  Global Lore: Core truths, history, and cosmology.
    3.  Relevant World State: Summary of major events, faction standings, overall player progress in related quest lines.
    4.  Quest Owner's custom `LLMPromptContext` (overarching goals, disposition).
    5.  Quest Owner's Resources: `CurrentInfluenceBudget`.
    6.  Status of Associated Quests/Questmakers: High-level summaries of quests under their thematic ownership.

*   **Static Questmaker Prompt Prefix:**
    ```
    "You are a Questmaker, responsible for controlling the execution and progression of a single, specific quest. Your influence budget primarily accumulates from player actions directly relevant to your quest, enabling you to react dynamically to player progress (or lack thereof). You can use the following tools: `send_message`, `change_npc_behavior_to_player`, `change_npc_stats`, `grant_player_reward`, `grant_passive_skill`, `QUESTMAKER_memorize`."
    ```

*   **`QuestmakerMonitor`:** Reactively (triggered by player actions related to its specific quest) compiles a prompt for its associated Questmaker. This prompt includes:
    1.  The **Static Questmaker Prompt Prefix**.
    2.  Quest-Specific Lore: Relevant lore entries for the quest's context.
    3.  Questmaker's custom `LLMPromptContext` (personality, goals for *this specific quest*).
    4.  Player Status: Current location, recent actions (both quest-related and unrelated), time since last relevant quest action, current quest progress, and inventory (especially quest items).
    5.  Quest-Specific Entity Status: The current state and location of key NPCs and items directly related to *this quest*.
    6.  Questmaker's Resources: `CurrentInfluenceBudget`.

**Example Quest Owner Prompt Snippet (Conceptual - Full Prompt):**

```
"You are a Quest Owner, a high-level entity responsible for supervising overarching narrative arcs and driving global world changes. You have a time-based influence budget to enact broad, strategic impacts. You can use the following tools: `trigger_world_event`, `spawn_entity` (for major entities), `change_room_info` (for significant global changes).

[Global Lore and World State here...]

You are the strategic mind behind Gandalf's efforts, focused on the larger picture of Middle-earth's fate. You orchestrate events and guide key individuals. The 'Urgent Message' quest is active, and the player 'player_123' has just delivered the message to Strider. The 'Road to Rivendell' quest is now available. Your current influence budget is 150/200. What strategic world changes or new quest initiations should occur?"
```

**Example Questmaker Prompt Snippet (Conceptual - Full Prompt):**

```
"You are a Questmaker, responsible for controlling the execution and progression of a single, specific quest. Your influence budget primarily accumulates from player actions directly relevant to your quest, enabling you to react dynamically to player progress (or lack thereof). You can use the following tools: `send_message`, `change_npc_behavior_to_player`, `change_npc_stats`, `grant_player_reward`, `grant_passive_skill`, `QUESTMAKER_memorize`.

[Quest-Specific Lore and Entity Status here...]

You are the direct overseer of 'The Urgent Message' quest. Your focus is solely on ensuring the message is delivered to Strider swiftly and safely. Player 'player_123' is currently in 'Prancing Pony Private Room' and has just spoken to Strider, completing your primary objective. Your current influence budget is 40/50. What immediate rewards or follow-up actions are appropriate for this player?"
```

### 3.2. LLM Output: Conceptual Tool Calls

The LLM's output will be a structured JSON object representing a list of "conceptual tool calls." These are interpreted and executed by dedicated processors. Each proposed action has an associated `cost`.

## 4. Influence Budget Management

### 4.1. Accumulation

*   **Quest Owners:** `CurrentInfluenceBudget` increases primarily through **time-based regeneration** (`BudgetRegenRate` per game tick/interval), representing their continuous, overarching influence. They may also gain budget upon major quest line completions or significant world events.
*   **Questmakers:** `CurrentInfluenceBudget` increases primarily based on **player actions** directly contributing to their specific quest (as defined in `Quest.InfluencePointsMap`). This allows them to react directly to player engagement. Their `BudgetRegenRate` will be very low or zero, emphasizing player-driven influence.

### 4.2. Spending

*   **`QuestOwnerActionProcessor`:** Executes LLM-proposed actions for Quest Owners. Checks `CurrentInfluenceBudget` against `action.cost`. If sufficient, executes global impact tools.
*   **`QuestmakerActionProcessor`:** Executes LLM-proposed actions for Questmakers. Checks `CurrentInfluenceBudget` against `action.cost`. If sufficient, executes local, quest-specific impact tools.

## 5. Player Action Influence & Layered Reaction

Player actions feed into both Quest Owners and Questmakers, triggering different layers of reaction.

*   **Player Action -> Questmaker Reaction:** Direct player actions (e.g., completing an objective, failing a task) immediately influence the associated Questmaker's budget and prompt its LLM to react with local, quest-specific changes (messages, NPC behavior, minor rewards/penalties).
*   **Quest Progress -> Quest Owner Reaction:** When a Questmaker reports significant progress or completion of its quest, this information is relayed to its `QuestOwner`. The `QuestOwner`'s LLM is then prompted to consider broader implications and enact global changes (major world events, new quest lines becoming available).

## 6. Conceptual Tools

Tools are now explicitly categorized by the entity type that can wield them, reflecting their scope of influence.

### 6.1. Owner Tools (Initiation)

*   **`initiate_quest`**
    *   **Purpose:** To make a specific quest available to a player or the world.
    *   **Parameters:** `quest_id`, `target_player_id` (optional), `trigger_condition` (e.g., "on_talk_to_npc", "on_enter_room").
    *   **Conceptual Cost:** Low (e.g., 0-10), representing the cost of offering a quest.
    *   **Example LLM Output (from an Owner):**
        ```json
        {
          "tool_name": "initiate_quest",
          "parameters": {
            "quest_id": "missing_pony_quest",
            "trigger_condition": "on_talk_to_npc",
            "target_npc_id": "barliman_butterbur"
          },
          "cost": 5,
          "reason": "Barliman is distressed about his pony, making the quest available."
        }
        ```

### 6.2. Quest Owner Tools (Global Impact)

These tools are used by Quest Owners to enact broad, strategic changes in the world.

*   **`trigger_world_event`**
    *   **Purpose:** To initiate a broader world event that affects multiple players or areas.
    *   **Parameters:** `event_type`, `details` (Optional).
    *   **Conceptual Cost:** Medium to High (e.g., 20-50), reflecting significant impact.
    *   **Example LLM Output:**
        ```json
        {
          "tool_name": "trigger_world_event",
          "parameters": {
            "event_type": "nazgul_patrol_increase",
            "details": "The Nazgûl's presence intensifies along the East Road, making travel more perilous."
          },
          "cost": 40,
          "reason": "The One Ring's journey progresses, increasing Sauron's awareness."
        }
        ```

*   **`spawn_entity`** (for significant, non-quest-specific entities)
    *   **Purpose:** To dynamically spawn a major NPC or item into the world as a strategic development.
    *   **Parameters:** `entity_type`, `location`, `quantity` (Optional).
    *   **Conceptual Cost:** Medium to High (e.g., 25-60), depending on the power/significance.
    *   **Example LLM Output:**
        ```json
        {
          "tool_name": "spawn_entity",
          "parameters": {
            "entity_type": "wandering_merchant_caravan",
            "location": "bree_road",
            "details": "A merchant caravan appears on the road, offering rare goods, but also attracting unwanted attention."
          },
          "cost": 35,
          "reason": "The Shire's peace has been maintained, allowing trade to flourish."
        }
        ```

*   **`change_room_info`** (for significant, non-quest-specific changes)
    *   **Purpose:** To modify properties of a room as a result of broader world events.
    *   **Parameters:** `room_id`, `property_changes`.
    *   **Conceptual Cost:** Medium (e.g., 10-30).
    *   **Example LLM Output:**
        ```json
        {
          "tool_name": "change_room_info",
          "parameters": {
            "room_id": "moria_west_gate",
            "property_changes": {
              "description_add": "A faint, chilling whisper now emanates from within the gate."
            }
          },
          "cost": 25,
          "reason": "The delving into darkness quest has stirred ancient evils within Moria."
        }
        ```

### 6.3. Questmaker Tools (Local Impact)

These tools are used by Questmakers to control the immediate progression and player experience of their single, associated quest.

*   **`send_message`**
    *   **Purpose:** To deliver a message to a player, potentially via a specific NPC.
    *   **Parameters:** `target_player_id`, `message`, `via_npc_id` (Optional).
    *   **Conceptual Cost:** Low (e.g., 5-15).

*   **`change_npc_behavior_to_player`**
    *   **Purpose:** To modify an NPC's disposition or behavior specifically towards a player, including adding a memory about that player.
    *   **Parameters:** `npc_id`, `player_id`, `behavior_type`, `memory_entry` (Optional).
    *   **Conceptual Cost:** Medium (e.g., 10-30).

*   **`change_npc_stats`**
    *   **Purpose:** To modify an NPC's core combat or non-combat statistics relevant to the quest.
    *   **Parameters:** `npc_id`, `stat_changes`.
    *   **Conceptual Cost:** Medium to High (e.g., 15-40).

*   **`grant_player_reward`**
    *   **Purpose:** To grant the player skills, spells, or items as a reward or consequence for quest progress.
    *   **Parameters:** `player_id`, `reward_type`, `reward_id`, `quantity` (Optional).
    *   **Conceptual Cost:** Low to High (e.g., 5-50).

*   **`grant_passive_skill`**
    *   **Purpose:** To directly grant a passive skill to a player, potentially influencing its initial percentage or cap based on quest performance.
    *   **Parameters:** `player_id`, `skill_id`, `initial_percentage` (Optional).
    *   **Conceptual Cost:** Medium to High (e.g., 20-70).

*   **`QUESTMAKER_memorize`**
    *   **Purpose:** Records a private memory about a player's actions or progress related to this specific quest.
    *   **Parameters:** `player_id`, `memory_string`.
    *   **Conceptual Cost:** Low (e.g., 5).

## 7. Implementation Considerations

*   **`QuestOwnerMonitor`:** A background service that periodically prompts Quest Owners' LLMs based on time and major world/quest line changes.
*   **`QuestOwnerActionProcessor`:** Interprets Quest Owner LLM output and interfaces with core game systems for global world changes.
*   **`QuestmakerMonitor`:** A background service that reactively prompts Questmakers' LLMs based on player actions within their specific quest.
*   **`QuestmakerActionProcessor`:** Interprets Questmaker LLM output and interfaces with core game systems for local, quest-specific changes.
*   **Tool Definitions:** Clear internal definitions for each "conceptual tool" with their parameters and associated `cost` ranges.
*   **Persistence:** All Quest Owner, Questmaker, and player quest states must be persisted in the database.

This revised system provides a robust and layered framework for dynamic, LLM-driven quest experiences, directly addressing the design criticisms and adding a unique layer of interactivity to the MUD.