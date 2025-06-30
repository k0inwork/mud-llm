# Owner System Proposal: Sentient World Guardians

This proposal details the design for the Owner system, where sentient LLM entities act as guardians or overseers of specific game world aspects (e.g., a town, a dungeon, a faction). Owners influence their associated areas and entities based on events and player actions within their domain. This system builds upon the LLM integration principles established in Phase 3.

## 1. Core Concept: Sentient World Guardians

Owners are LLM-driven entities responsible for monitoring and influencing specific aspects or areas of the game world. Unlike Questmakers, whose focus is on a single quest line, Owners have a broader, persistent oversight of their domain. They react to significant events and player actions within their sphere of influence by spending an "Influence Budget" to trigger world events, modify NPC behaviors, or manage resources within their domain.

## 2. Data Models

### 2.1. `Owner` Struct

The existing `Owner` struct will be enhanced to include LLM-specific attributes and an influence budget.

```go
type Owner struct {
    ID                  string           // Unique identifier (e.g., "town_council_owner", "goblin_king_owner")
    Name                string           // Display name
    Description         string           // A brief description of the Owner
    MonitoredAspect     string           // Defines what the Owner primarily monitors (e.g., "location", "faction_reputation", "resource_supply")
    AssociatedID        string           // The ID of the entity/area/faction this Owner is associated with (e.g., "town_square", "goblin_faction")
    LLMPromptContext    string           // Defines personality, goals, and core directives for the LLM
    MemoriesAboutPlayers map[string][]string // Private memories about players (from Phase 3)
    CurrentInfluenceBudget float64        // Current points available for actions
    MaxInfluenceBudget  float64        // Maximum capacity of influence points
    BudgetRegenRate     float64        // Points regenerated per game tick/significant event
    // ... other existing Owner fields ...
}
```

### 2.2. `NPC` Struct (Relationship)

NPCs will continue to have an `OwnerIDs` field, linking them to their respective Owners.

```go
type NPC struct {
    // ... existing NPC fields ...
    OwnerIDs            []string         // IDs of Owners this NPC is associated with
    // ... other existing NPC fields ...
}
```

### 2.3. `Room` Struct (Relationship)

Rooms can be associated with Owners, allowing Owners to monitor and influence events within those rooms.

```go
type Room struct {
    // ... existing Room fields ...
    OwnerID             string           // ID of the Owner primarily responsible for this room (optional)
    // ... other existing Room fields ...
}
```

## 3. Owner Decision Logic & LLM Integration

The Owner's decision-making is entirely driven by its LLM, which acts as its "brain."

### 3.1. LLM Prompt Generation

A dedicated `OwnerMonitor` service (potentially integrated with the `Action Significance Monitor` from Phase 4) will periodically (or reactively) compile a comprehensive prompt for the Owner's LLM. This prompt provides all necessary context for the LLM to make informed decisions:

*   **Owner's Core Identity:** `LLMPromptContext` (personality, goals, current disposition), `MonitoredAspect`, `AssociatedID`.
*   **Relevant World State:** Summary of events, player actions, and NPC states within the Owner's `MonitoredAspect` or `AssociatedID` domain. This includes:
    *   Recent player actions (e.g., theft in town, defeating a faction NPC).
    *   Status of key NPCs within the domain.
    *   Resource levels (e.g., town supplies, faction strength).
    *   Environmental conditions.
*   **Owner's Resources:** `CurrentInfluenceBudget` and `MaxInfluenceBudget`.
*   **Owner's Memories:** Relevant `MemoriesAboutPlayers` or other entities.

**Example Prompt Snippet (Conceptual for "The Town Council" Owner):**

```
"You are The Town Council, the collective consciousness of Port Town. Your primary goal is the prosperity and safety of the town and its citizens. You are concerned by disruptions, crime, and threats to order. Your current influence budget is 40/100.

Recent events in your domain (Port Town):
- Player 'player_alice' was observed stealing from the market stall (minor crime).
- A minor market fluctuation occurred (reason unknown, but impacts town prosperity).
- Innkeeper Bob reported a suspicious character loitering near the docks.

Considering these events and your goals, what actions should you take to maintain order, address threats, or influence the town, using your available influence budget? Respond with a list of tool calls."
```

### 3.2. LLM Output: Conceptual Tool Calls

The LLM's output will be a structured JSON object representing a list of "conceptual tool calls." These are interpreted and executed by the game engine's `OwnerActionProcessor`. Each proposed action will have an associated `cost` that the LLM must consider against the Owner's `CurrentInfluenceBudget`.

**Example LLM Output (Structured JSON):**

```json
{
  "owner_id": "town_council_owner",
  "proposed_actions": [
    {
      "tool_name": "change_npc_behavior_to_player",
      "parameters": {
        "npc_id": "town_guard_captain",
        "player_id": "player_alice",
        "behavior_type": "suspicious",
        "memory_entry": "Player 'player_alice' was seen stealing from the market."
      },
      "cost": 15,
      "reason": "Player committed a minor crime, increasing guard vigilance."
    },
    {
      "tool_name": "trigger_world_event",
      "parameters": {
        "event_type": "increased_guard_patrols",
        "details": "Town guards increase patrols in market district for a short period."
      },
      "cost": 20,
      "reason": "Response to minor crime and general security concerns."
    },
    {
      "tool_name": "send_message",
      "parameters": {
        "target_player_id": "player_alice",
        "message": "A stern warning from the Town Council: 'Order must be maintained. Your actions have been noted.'",
        "via_npc_id": "town_crier"
      },
      "cost": 10,
      "reason": "Direct warning to the player for disruptive behavior."
    }
  ]
}
```

## 4. Influence Budget Management

### 4.1. Accumulation

The Owner's `CurrentInfluenceBudget` increases based on events within its domain:

*   **Time-Based Regeneration:** A consistent `BudgetRegenRate` per game tick, representing the Owner's inherent capacity to influence its domain.
*   **Significant Events:** Owners might gain `InfluencePoints` from events that align with their goals (e.g., player completing a beneficial quest in their town, a successful trade caravan arriving).

### 4.2. Spending

An `OwnerActionProcessor` component is responsible for executing the LLM's `proposed_actions`:

1.  It receives the LLM's output.
2.  For each `proposed_action`, it checks if the Owner's `CurrentInfluenceBudget` is sufficient for the `action.cost`.
3.  If sufficient, the action is executed by calling the corresponding game engine system.
4.  The `action.cost` is then deducted from `CurrentInfluenceBudget`.
5.  If insufficient, the action is skipped, and the Owner's LLM might be prompted again with the updated (lower) budget, potentially leading to less costly actions or no actions.

## 5. Conceptual Tools for Owners

Owners utilize a set of conceptual tools to influence their domain. These tools are represented as structured JSON outputs from the LLM, interpreted and executed by the `OwnerActionProcessor`. Each tool has an associated conceptual `cost`.

### 5.1. Memory Management Tools (from Phase 3)

*   **`OWNER_memorize`**
    *   **Purpose:** Records a private memory about a player or other entity, known only to this Owner.
    *   **Parameters:** `entity_id`, `memory_string`.
    *   **Conceptual Cost:** Low (e.g., 5).
*   **`OWNER_memorize_dependables`**
    *   **Purpose:** Broadcasts a memory about a player or entity to all NPCs associated with this Owner, influencing their collective perception.
    *   **Parameters:** `entity_id`, `memory_string`.
    *   **Conceptual Cost:** Medium (e.g., 15-25), as it affects multiple entities.

### 5.2. NPC Interaction Tools

*   **`send_message`**
    *   **Purpose:** To deliver a message to a player or specific NPC within the Owner's domain, potentially via another NPC.
    *   **Parameters:** `target_id` (player or NPC), `message`, `via_npc_id` (optional).
    *   **Conceptual Cost:** Low (e.g., 5-15).
*   **`change_npc_behavior`**
    *   **Purpose:** To modify an NPC's general disposition or behavior within the Owner's domain (e.g., making all guards more vigilant).
    *   **Parameters:** `npc_id` (or `npc_type` for broad changes), `behavior_type` (e.g., "vigilant", "relaxed", "hostile").
    *   **Conceptual Cost:** Medium (e.g., 10-30).
*   **`change_npc_stats`**
    *   **Purpose:** To modify an NPC's core combat or non-combat statistics within the Owner's domain.
    *   **Parameters:** `npc_id`, `stat_changes` (map of stat names to values/modifiers).
    *   **Conceptual Cost:** Medium to High (e.g., 15-40).

### 5.3. World Manipulation Tools

*   **`change_room_info`**
    *   **Purpose:** To modify properties of a specific room within the Owner's domain (e.g., locking/unlocking doors, changing descriptions, adding/removing environmental effects).
    *   **Parameters:** `room_id`, `property_changes` (map of properties to new values).
    *   **Conceptual Cost:** Medium (e.g., 10-30).
*   **`trigger_world_event`**
    *   **Purpose:** To initiate a broader world event that affects multiple players or areas within the Owner's domain (e.g., a market boom/bust, a localized natural disaster, a faction gathering).
    *   **Parameters:** `event_type`, `details` (optional).
    *   **Conceptual Cost:** Medium to High (e.g., 20-50).
*   **`spawn_entity`**
    *   **Purpose:** To dynamically spawn an NPC or item into the world within the Owner's domain (e.g., new guards, a lost artifact, a new merchant).
    *   **Parameters:** `entity_type`, `location`, `quantity` (optional).
    *   **Conceptual Cost:** Medium to High (e.g., 25-60).
*   **`modify_resource`**
    *   **Purpose:** To directly influence a quantifiable resource within the Owner's domain (e.g., town food supply, faction gold reserves).
    *   **Parameters:** `resource_id`, `change_amount` (positive or negative value).
    *   **Conceptual Cost:** Medium (e.g., 10-30).

*   **`grant_passive_skill`**
    *   **Purpose:** To directly grant a passive skill to a player, potentially influencing its initial percentage or cap based on the Owner's attitude.
    *   **Parameters:**
        *   `player_id`: The ID of the target player.
        *   `skill_id`: The ID of the passive skill to grant.
        *   `initial_percentage`: (Optional) The initial percentage for the skill (0-100). If omitted, defaults to 0 or a base value.
    *   **Conceptual Cost:** Medium to High (e.g., 20-70), depending on the power of the skill and the Owner's generosity/goals.
    *   **Example LLM Output:**
        ```json
        {
          "tool_name": "grant_passive_skill",
          "parameters": {
            "player_id": "player_123",
            "skill_id": "town_favor",
            "initial_percentage": 15
          },
          "cost": 30,
          "reason": "Player contributed significantly to town defense, earning the council's favor."
        }
        ```

### 5.4. Note on Tool Exclusivity

While some tools are shared between Owners and Questmakers (e.g., `send_message`, `change_npc_behavior`), it is important to note that the `grant_player_reward` tool is **exclusive to Questmakers**, as it directly relates to player progression within a specific quest narrative. Conversely, the `modify_resource` tool is **exclusive to Owners**, as it pertains to the management of domain-level resources.

## 6. Implementation Considerations

*   **`OwnerMonitor`:** A background service responsible for monitoring events within an Owner's domain, compiling prompts, and sending them to the LLM. This will likely leverage the `Action Significance Monitor` for event filtering.
*   **`OwnerActionProcessor`:** Interprets LLM output and interfaces with core game systems (NPC AI, World Event Manager, Communication System, Room Manager, Resource Manager).
*   **Persistence:** Owner states (including `CurrentInfluenceBudget` and `MemoriesAboutPlayers`) must be persisted in the database.

This system provides a robust framework for dynamic, LLM-driven world management, allowing for emergent narratives and consequences within specific game domains.