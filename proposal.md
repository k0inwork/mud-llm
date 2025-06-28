# Proposal: Enhancing GoMUD with LLM-Driven Sentient Entities and Advanced Mechanics

## Introduction
This document outlines a proposal for significantly enhancing the existing GoMUD project by integrating concepts from the provided development log of another MUD project. The core idea is to introduce LLM-driven sentient entities (NPCs and "Owners"), a robust tool-calling mechanism, and advanced game mechanics like skills, reputation, and a dynamic world.

The primary focus will be on backend implementation in Go, with the primary client being a standard Telnet client. The server will utilize ANSI escape codes to provide rich text formatting, including colors, bold, and underline, to enhance the player experience. A lightweight web server will also be maintained for administrative tasks, including an in-game editor.

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

## 2. World-Building and Lore Integration

To ensure a cohesive and immersive world, we will implement a comprehensive and dynamic lore system. This system will serve as the foundational knowledge base for all sentient entities, shaping their understanding of the world, their place in it, and their interactions with players. All lore will be configurable via the web editor.

*   **The "Lore Book":** We will create a central repository for all lore, structured as a collection of `Lore` objects. This allows for modular and easily editable world-building.

*   **Global Lore:** A set of core truths and historical facts about the world that are known to all sentient beings. This includes:
    *   **Cosmology and Creation Myths:** The story of how the world came to be, the major deities, and the fundamental forces that govern reality.
    *   **Major Historical Events:** The timeline of pivotal moments, such as ancient wars, the rise and fall of empires, and cataclysmic events that have shaped the current landscape.
    *   **The "Total Lore":** A concise, high-level summary of the world's state, providing a common foundation for all NPC and Owner knowledge.

*   **Scoped Lore:** To give entities specialized knowledge, lore will be scoped to specific aspects of the game world. When an LLM is prompted, it will receive the Global Lore plus all relevant Scoped Lore.
    *   **Race Lore:** Details the history, culture, traditions, and social norms of each race. An Elven NPC will know about their ancient lineage, while a Dwarven NPC will know the history of their mountain holds. This can also include inherent biases or alliances between races.
    *   **Profession Lore:** Provides knowledge specific to a character's profession. A Mage will have access to lore about magical theory and arcane history, a Warrior will know of legendary battles and famous warriors, and a Thief will be aware of the criminal underworld and its key players.
    *   **Zone Lore:** Contains detailed information about specific geographical areas, cities, or dungeons. This includes local history, significant landmarks, hidden secrets, and the types of creatures or factions that inhabit the area. An NPC in a particular city will "know" its layout, its rulers, and its local customs.

*   **Creative Lore Categories (for deeper immersion):**
    *   **Faction Lore:** Knowledge related to specific guilds, secret societies, political movements, or religious cults. An NPC belonging to a faction will be aware of its goals, hierarchy, allies, and enemies, creating opportunities for intrigue and quests.
    *   **Bestiary (Creature) Lore:** Information about the world's flora and fauna. An experienced Ranger might possess knowledge about a monster's weaknesses, habitat, and behavior, which they could share with players.
    *   **Item Lore:** Legendary tales and histories associated with specific powerful artifacts. An NPC might recognize a player's sword as a long-lost blade from a forgotten age.

*   **Lore Integration with LLMs:** The lore system is not just static text. It will be dynamically injected into the context of LLM prompts. When an NPC or Owner is about to generate a response, the system will:
    1.  Fetch the **Global Lore**.
    2.  Fetch the relevant **Scoped Lore** based on the entity's race, profession, location, and faction affiliations.
    3.  Prepend this collection of lore to the prompt as foundational knowledge, ensuring the LLM's response is deeply rooted in the game's established world.

## 3. LLM Integration Strategy (Go-Gemini API & Tool Calling)

This is the most critical and complex part. We will need to:

*   **Go LLM API Client:** Implement a Go client to interact directly with the LLM API. This client will be configurable to support both the Gemini API and OpenAI-compatible endpoints. This will involve:
    *   Making HTTP POST requests to the chosen LLM API endpoint.
    *   Handling API keys securely (e.g., environment variables).
    *   Structuring requests with prompts, system instructions, and tool definitions.
    *   Parsing the LLM's JSON responses.
*   **Prompt Engineering for Tool Calling:**
    *   The LLM needs to be instructed to output not just narrative, but also structured tool calls. The log suggests an XML format (`<response><narrative>...</narrative><tools>TOOL_JSON_HERE</tools></response>`). This is a robust approach.
    *   Prompts sent to the LLM will be dynamically constructed to include:
        *   **World Knowledge:** The relevant **Global and Scoped Lore** to ground the entity in the game's universe.
        *   **Situational Context:** The player's action, room description, and relevant NPC/Owner state.
        *   **Personality:** The NPC's/Owner's `PersonalityPrompt` and `MemoriesAboutPlayers`.
        *   **Capabilities:** A list of `AvailableTools` with their descriptions and parameter schemas, explicitly instructing the LLM on how to format tool calls within the XML.
*   **Tool Dispatcher:**
    *   A central Go function will be responsible for parsing the LLM's XML response.
    *   If a `tools` section is present, it will parse the JSON within it.
    *   It will then map the `tool_name` to a corresponding Go function (e.g., `handleGiftItem(player, itemID)`, `handleHealPlayer(player, amount)`).
    *   These Go functions will directly modify the game state.

## 4. Sentient NPCs

NPCs will become dynamic, AI-driven entities whose behavior is shaped by their knowledge and experiences:

*   **AI-Driven Behavior:**
    *   When a player interacts with an NPC (e.g., `talk`, `attack`, `say`), or when an NPC needs to react to an event (e.g., player entering room), a prompt will be constructed for the LLM.
    *   The prompt will include the relevant **Lore**, the NPC's `PersonalityPrompt`, `MemoriesAboutPlayers`, and `AvailableTools`.
    *   The LLM's response will drive the NPC's narrative and actions (via tool calls), ensuring they are consistent with their background and the world's history.
*   **Memory System (`MemoriesAboutPlayers`):**
    *   When a player performs a significant action (e.g., fails a password, attacks an NPC), relevant Owners/NPCs can use a `TOOL_MEMORIZE` to record this information in the NPC's `MemoriesAboutPlayers`.
    *   These memories will be prepended to future AI prompts for that NPC when interacting with the specific player, influencing the NPC's behavior.
*   **Item Interaction Tools:**
    *   Implement Go functions for `TOOL_NPC_EXAMINE_ITEM_IN_ROOM`, `TOOL_NPC_PICK_UP_ITEM_FROM_ROOM`, `TOOL_NPC_DROP_ITEM_TO_ROOM`, `TOOL_GIFT_ITEM`.
    *   These tools will allow NPCs to dynamically interact with items in their environment and with the player's inventory.

## 5. Sentient Owners

Owners are higher-level, LLM-controlled entities that act as guardians or influencers over their domains:

*   **Role and Scope:**
    *   Owners will be associated with `MonitoredAspect`s (locations, races, professions, factions).
    *   When a player enters a room, changes race/profession, or performs an action relevant to an Owner's domain, the relevant Owner(s) will be prompted.
*   **Tool Invocation:**
    *   Owners will have their own `AvailableTools` (e.g., `TOOL_GIFT_ITEM`, `TOOL_HEAL_PLAYER`, `TOOL_SMITE_PLAYER`, `TOOL_BLESS_PLAYER`).
    *   Their LLM responses can include tool calls to affect the player or the world, acting in accordance with their knowledge and personality.
*   **Prayer Mechanism:**
    *   Implement a `pray [message]` command.
    *   This command will identify all relevant Owners (based on player's current room's location owner, player's race owner, player's profession owner).
    *   Each relevant Owner will receive an AI prompt with the player's prayer and context, allowing them to respond and potentially use tools based on their divine or authoritative perspective, shaped by the lore.

## 6. Skills as Tools

Skills will be integrated directly into the LLM's tool-calling mechanism:

*   **Active Skills:**
    *   These will be defined as `OwnerTool` or `NPCTool` types.
    *   When a player uses an active skill (e.g., "use minor heal"), this command will be routed to the appropriate handler.
    *   For LLM-driven skills, the LLM will be prompted to "use" the skill, and its tool call will trigger the actual effect.
*   **Passive Skills:**
    *   These will not be explicit tool calls but rather modifiers to player stats or contextual information provided to the LLM.
    *   For example, a "City Lore" passive skill might mean the LLM receives additional context about city events when the player is in a city, leading to more informed NPC/Owner responses.

## 7. Reputation System

*   **Owner-to-NPC Communication:**
    *   When an Owner "observes" a player (e.g., player enters their territory, performs a significant action), the Owner's LLM can use `TOOL_MEMORIZE` to update the `MemoriesAboutPlayers` of its associated NPCs.
    *   This would involve iterating through NPCs linked to that Owner and updating their memory state.
*   **NPC Reaction:**
    *   When a player interacts with an NPC, the NPC's LLM prompt will include its `MemoriesAboutPlayers` for that player, influencing its narrative and tool usage (e.g., a "liked" player might receive a `gift_item`, a "disliked" player might be `smite_player`).

## 8. Game Mechanics Enhancements

*   **Locking Mechanism:**
    *   Exits can be `IsLocked` and require a `KeyID`.
    *   `movePlayer` will check `IsLocked` and prevent movement.
    *   A `use <item> on <exit>` command (or an LLM tool call like `unlock_exit`) will check for the correct `KeyID` and set `IsLocked` to `false`.
*   **Map Feature:**
    *   The `Player.VisitedRoomIDs` map will track explored rooms.
    *   The backend will maintain this `VisitedRoomIDs` state.
    *   For Telnet clients, a simple ASCII art map could be rendered based on this data.
*   **Conditional AI Delivery:**
    *   AI responses (narrative and tool effects) can be conditional based on player ID or room ID.
    *   This means the server will hold onto AI responses and only deliver them to the client (or apply effects) if the conditions are met (e.g., player is in the correct room for an NPC's reaction). This will be crucial for managing concurrent AI responses.

## 9. Concurrency Considerations

*   **Go Routines for AI Calls:** When multiple AI interactions are triggered (e.g., a player's prayer prompts multiple Owners, or a `say` command prompts multiple NPCs), use Go routines (`go func()`) and `sync.WaitGroup` or channels to dispatch these LLM API calls concurrently.
*   **State Management:** Ensure all shared game state (players, rooms, NPCs, items) is protected by mutexes (`sync.Mutex` or `sync.RWMutex`) to prevent race conditions during concurrent updates from AI tool calls or player actions.

## 10. Web Server (Admin & Editor)

*   The existing lightweight web server will be maintained.
*   It will primarily serve administrative functions, including:
    *   Displaying server statistics.
    *   Providing a web-based editor for game content (rooms, items, NPCs, Owners, and **Lore**). This editor will allow for easy modification and creation of game entities and their foundational knowledge without direct file manipulation.

## 11. High-Level Implementation Phases

A possible phased approach:

1.  **Phase 1: Core Data Structures & Lore System:**
    *   Implement all new and expanded data structures (`Lore`, `NPCTool`, `OwnerTool`, etc.).
    *   Build the backend logic to store and retrieve lore entries from a database or file.
    *   Update the web editor to allow for creating, editing, and deleting global and scoped lore entries.
2.  **Phase 2: Basic Go LLM Integration & Tool Calling Framework:**
    *   Implement the Go LLM API client (configurable for Gemini/OpenAI).
    *   Create the tool dispatcher.
    *   Implement the logic to dynamically construct prompts, injecting the relevant lore.
    *   Test the full pipeline with a simple "pray" command to a single Owner with one tool (e.g., `TOOL_BLESS_PLAYER`) to establish the LLM interaction and tool parsing.
3.  **Phase 3: Sentient NPCs & Basic Tools:**
    *   Expand `NPC` struct logic.
    *   Implement `TOOL_MEMORIZE` and basic item interaction tools (`pickup_item`, `drop_item`).
    *   Integrate NPC AI responses into `talk` and `say` commands.
4.  **Phase 4: Advanced Owner Logic & Skills:**
    *   Implement full `Owner` struct and logic for identifying relevant Owners.
    *   Add more complex Owner tools (e.g., `heal_player`, `smite_player`, `gift_item`).
    *   Integrate active skills as LLM-callable tools.
5.  **Phase 5: Reputation & Advanced Mechanics:**
    *   Implement the reputation system (Owners updating NPC memories).
    *   Implement the locking mechanism for exits.
    *   Track `visitedRoomIds` for map data.
6.  **Phase 6: Concurrency & Robustness:**
    *   Refactor LLM calls to use Go routines for concurrency.
    *   Thoroughly review mutex usage for all shared state.
    *   Implement conditional AI delivery.
7.  **Phase 7: Telnet Client Enhancements & Final Polish:**
    *   Implement ANSI escape code handling for colors, bold, and underline in Telnet output.
    *   Finalize the web editor interface and polish the overall gameplay experience.

This proposal outlines a comprehensive path to evolve the GoMUD into a more dynamic and intelligent multi-user dungeon experience.