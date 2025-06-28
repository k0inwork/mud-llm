# Phase 4 Design: Sentient NPCs & Owners

## 1. Objectives

With the core AI infrastructure in place, this phase focuses on bringing the NPCs and Owners to life. The goal is to move from testing with hardcoded triggers to having the AI entities react dynamically and autonomously to player actions and game events.

## 2. Key Components to be Implemented

### 2.1. Trigger and Event System

*   A simple event-driven system will be implemented. Key game actions will generate events that the Sentient Entity Manager can listen for.
*   **Required Events:**
    *   `PlayerEntersRoom(player_id, room_id)`
    *   `PlayerSays(player_id, message)`
    *   `PlayerTalksToNPC(player_id, npc_id)`
    *   `PlayerAttacksNPC(player_id, npc_id)`
    *   `PlayerPrays(player_id, message)`

### 2.2. Sentient Entity Manager

*   This manager is the central brain that orchestrates AI responses.
*   It will have handlers that subscribe to the events from the trigger system.
*   **Event Handler Logic:**
    *   **On `PlayerEntersRoom`:** The manager will identify the NPCs in the room and decide if any should react. For now, this can be a simple random chance or based on a personality flag (e.g., "Vigilant"). If an NPC reacts, an LLM prompt is triggered.
    *   **On `PlayerSays`:** The manager will identify all NPCs in the room and prompt each of them with the player's message, allowing them to overhear and potentially react.
    *   **On `PlayerTalksToNPC`:** The manager will trigger a prompt specifically for the targeted NPC.
    *   **On `PlayerPrays`:** The manager will identify the relevant Owners (based on the player's location, race, profession) and trigger a prompt for each of them.

### 2.3. Integrating AI Responses

*   The narrative part of the LLM's response needs to be fed back into the presentation layer.
*   The `HandleLLMResponse` function will be created to:
    1.  Take the `narrative` string from the LLM response.
    2.  Create a semantic JSON object (e.g., `{ "type": "narrative", "payload": { "segments": [...] } }`).
    3.  Send this JSON to the appropriate player's client via the `TelnetRenderer`.
*   The tool calls from the LLM response will be sent to the Tool Dispatcher as implemented in Phase 3.

## 3. Acceptance Criteria

1.  When a player enters a room, at least one NPC in that room can be observed to react by speaking.
2.  When a player uses the `say` command, NPCs in the room can be observed to react to the message.
3.  When a player uses the `talk` command on an NPC, the NPC provides a direct, relevant response.
4.  When a player uses the `pray` command, at least one Owner can be observed to respond with a message.
5.  All AI-generated narrative is correctly formatted and displayed to the player via the server-side presentation layer.
6.  The `OWNER_memorize_dependables` tool works as expected, allowing an Owner's prayer response to influence the memories of its subordinate NPCs.
