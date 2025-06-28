# Proposal: Enhancing GoMUD with LLM-Driven Sentient Entities and Advanced Mechanics

## Introduction
This document outlines a proposal for significantly enhancing the existing GoMUD project by integrating concepts from the provided development log of another MUD project. The core idea is to introduce LLM-driven sentient entities (NPCs and "Owners"), a robust tool-calling mechanism, and advanced game mechanics like skills, reputation, and a dynamic world.

The primary focus will be on backend implementation in Go. The architecture will feature a decoupled presentation layer, where the server generates a semantic JSON output, which is then interpreted by different clients (e.g., a Telnet client for ANSI text, a web client for HTML).

## 1. Core Data Model Enhancements

To support the new features, the existing Go data structures will need significant expansion:

*   **`Player` Struct:**
    *   `Race string`: Player's race (e.g., "Human", "Elf").
    *   `Profession string`: Player's profession (e.g., "Adventurer", "Warrior").
    *   `VisitedRoomIDs map[int]bool`: A map to track visited rooms for map generation/fog of war.
    *   `Inventory map[string]*Item`: (Already exists, but ensure it's robust for item manipulation).
    *   `Health int`, `MaxHealth int`: For combat/healing.
*   **`NPC` Struct:**
    *   `OwnerIDs []string`: List of IDs of Owners associated with this NPC (location, race, profession).
    *   `MemoriesAboutPlayers map[string][]string`: A map where keys are player IDs and values are lists of strings representing memories/reputation about that player.
    *   `AvailableTools []NPCTool`: A list of tools this NPC's controlling LLM can invoke.
    *   `PersonalityPrompt string`: Base prompt for the LLM to define NPC behavior.
    *   `Inventory map[string]*Item`: NPCs can hold items.
*   **New `Owner` Struct:**
    *   `ID string`
    *   `Name string`
    *   `Description string`: For LLM context.
    *   `MonitoredAspect string`: "location", "race", or "profession".
    *   `AssociatedID string`: The ID of the location, race, or profession this owner monitors.
    *   `AvailableTools []OwnerTool`: Tools this Owner's controlling LLM can invoke.
*   **New `OwnerTool` Struct:**
    *   `Name string`: Unique identifier for the tool (e.g., "gift_item", "heal_player").
    *   `Description string`: For the LLM to understand the tool's purpose.
    *   `Parameters map[string]interface{}`: JSON schema-like definition of parameters the tool expects.
*   **New `NPCTool` Struct:** (Mirrors `OwnerTool` for NPC-specific actions)
    *   `Name string`: Unique identifier for the tool (e.g., "examine_item", "pickup_item").
    *   `Description string`: For the LLM to understand the tool's purpose.
    *   `Parameters map[string]interface{}`: JSON schema-like definition of parameters the tool expects.
*   **`Room` Struct:**
    *   `OwnerID string`: ID of the Owner associated with this room's location.
    *   `Items map[string]*Item`: Items present in the room.
*   **`Exit` Struct:**
    *   `IsLocked bool`: Whether the exit is currently locked.
    *   `KeyID string`: The ID of the item required to unlock this exit (optional).
*   **`Item` Struct:**
    *   `ID string`, `Name string`, `Description string`.
    *   `Attributes map[string]interface{}`: Flexible field for item-specific properties (e.g., `{"healAmount": 10}`, `{"isKeyFor": "treasuryGate"}`).
*   **New `Lore` Struct:**
    *   `ID string`: Unique identifier for the lore entry.
    *   `Type string`: The type of lore, e.g., "global", "race", "profession", "zone", "faction", "creature".
    *   `AssociatedID string`: The specific ID this lore is associated with (e.g., "elf", "warrior", "the_dark_forest"). For "global" lore, this can be empty.
    *   `Content string`: The text of the lore itself.

## 2. Unified Presentation Layer (Semantic JSON)

To ensure maximum flexibility for future clients (e.g., web, graphical), the server's output will be decoupled from its presentation. The server will not send raw text; instead, it will send a structured JSON object that describes the content semantically.

*   **Structure:** All messages to the client will be a JSON object with a `type` and a `payload`.
*   **Example 1: Narrative Text**
    ```json
    {
      "type": "narrative",
      "payload": {
        "segments": [
          {"text": "The elven guard narrows his eyes. 'State your business,' he says, his hand resting on the hilt of his sword. You can feel the "},
          {"text": "ancient magic", "style": ["italic", "color:cyan"]},
          {"text": " humming in the air around him."}
        ]
      }
    }
    ```
*   **Example 2: Room Description**
    ```json
    {
      "type": "room_update",
      "payload": {
        "name": "The Whispering Glade",
        "description": "Sunlight filters through the dense canopy above, illuminating a small, peaceful clearing.",
        "exits": ["north", "west"],
        "items": [{"name": "a shimmering potion", "style": ["bold"]}],
        "npcs": [{"name": "An elven guard", "style": ["color:green"]}]
      }
    }
    ```
*   **Client Responsibility:**
    *   **Telnet Client:** Will parse this JSON and translate it into text with the appropriate ANSI escape codes for color and style.
    *   **Web Client:** Will parse this JSON and render it as styled HTML elements.

## 3. World-Building and Lore Integration

A comprehensive and dynamic lore system will serve as the foundational knowledge base for all sentient entities. All lore will be configurable via the web editor.

*   **The "Lore Book":** A central repository for all `Lore` objects.
*   **Global Lore:** Core truths and history known to all (creation myths, major wars, etc.).
*   **Scoped Lore:** Specialized knowledge for races, professions, zones, factions, creatures, and items.
*   **Lore Integration with LLMs:** When an NPC or Owner is prompted, the system will dynamically fetch and prepend the Global Lore and all relevant Scoped Lore to the prompt, ensuring all responses are deeply rooted in the game's established world.

## 4. LLM Integration Strategy (Go-Gemini API & Tool Calling)

*   **Go LLM API Client:** A configurable client for Gemini and OpenAI-compatible endpoints.
*   **Prompt Engineering:** Prompts will be dynamically constructed with World Knowledge (Lore), Situational Context, Personality, and Capabilities (Tools).
*   **Tool Dispatcher:** A central function to parse LLM XML responses and map tool names to Go functions that modify game state.

## 5. Sentient NPCs

*   **AI-Driven Behavior:** NPC actions and dialogue will be driven by LLM responses that are informed by lore, memory, and personality.
*   **Memory System:** Significant player actions can be recorded in an NPC's `MemoriesAboutPlayers` map using a `TOOL_MEMORIZE`.
*   **Item Interaction Tools:** NPCs will have tools to interact with items in their environment.

## 6. Sentient Owners

*   **Role and Scope:** Higher-level entities overseeing locations, races, professions, or factions.
*   **Tool Invocation:** Owners have powerful tools to influence the world (e.g., `TOOL_GIFT_ITEM`, `TOOL_SMITE_PLAYER`).
*   **Prayer Mechanism:** The `pray` command prompts relevant Owners, who can respond based on their lore-informed perspective.

## 7. Reputation System

*   **Owner-to-NPC Communication:** When an Owner observes a significant player action, it can use `TOOL_MEMORIZE` to update the memories of its associated NPCs.
*   **Directives and Impressions:** Furthermore, Owners can actively send directives about a player to their subordinate NPCs, especially for players who have gained significant favor or disfavor. This can be achieved via a new tool, `TOOL_IMPRINT_MEMORY_ON_NPC(npc_id, player_id, memory_string)`, allowing an Owner's opinion to directly and immediately influence an NPC's behavior.
*   **NPC Reaction:** An NPC's behavior will be influenced by its own memories and the impressions passed down from its Owners.

## 8. Skills as Tools

*   **Active Skills:** Defined as `OwnerTool` or `NPCTool` types, callable by the LLM.
*   **Passive Skills:** Act as modifiers or provide additional context to the LLM.

## 9. Game Mechanics Enhancements

*   **Locking Mechanism:** Exits can be locked and require specific keys.
*   **Map Feature:** `Player.VisitedRoomIDs` will be tracked for map generation.
*   **Conditional AI Delivery:** AI responses can be held and delivered only when specific conditions are met.

## 10. Concurrency Considerations

*   Use Go routines for concurrent LLM API calls.
*   Protect all shared game state with mutexes.

## 11. Web Server (Admin & Editor)

*   The web server will provide an editor for all core game content: **Lore**, rooms, items, NPCs, Owners, and their associations.

## 12. High-Level Implementation Phases

1.  **Phase 1: Core Architecture & Data Structures:**
    *   Implement all new and expanded data structures (`Lore`, `NPCTool`, `OwnerTool`, etc.).
    *   Implement the **Unified Presentation Layer**, ensuring all server output is in the semantic JSON format.
2.  **Phase 2: Lore System & Editor:**
    *   Build the backend logic to store and retrieve lore entries.
    *   Update the web editor to allow for creating and editing lore.
3.  **Phase 3: Basic LLM Integration & Tool Calling:**
    *   Implement the Go LLM API client and the tool dispatcher.
    *   Implement the logic to inject lore into the LLM prompt.
    *   Test with a simple "pray" command.
4.  **Phase 4: Sentient NPCs & Basic Tools:**
    *   Integrate NPC AI responses with `talk` and `say`.
    *   Implement `TOOL_MEMORIZE` and item interaction tools.
5.  **Phase 5: Advanced Owner Logic & Reputation:**
    *   Implement the full `Owner` logic.
    *   Implement the `TOOL_IMPRINT_MEMORY_ON_NPC` and other Owner tools.
    *   Integrate active skills.
6.  **Phase 6: Advanced Mechanics & Concurrency:**
    *   Implement locking, mapping, and conditional AI delivery.
    *   Refactor LLM calls for concurrency and review mutex usage.
7.  **Phase 7: Client Implementation & Final Polish:**
    *   Build the **Telnet client** to interpret the semantic JSON.
    *   Finalize the web editor and polish the overall experience.

This proposal outlines a comprehensive path to evolve the GoMUD into a more dynamic and intelligent multi-user dungeon experience.