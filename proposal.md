# Proposal: Enhancing GoMUD with LLM-Driven Sentient Entities and Advanced Mechanics

## Introduction
This document outlines a proposal for significantly enhancing the existing GoMUD project by integrating concepts from the provided development log of another MUD project. The core idea is to introduce LLM-driven sentient entities (NPCs and "Owners"), a robust tool-calling mechanism, and advanced game mechanics like skills, reputation, and a dynamic world.

The primary focus will be on backend implementation in Go. The architecture will feature a decoupled presentation layer, where the server's core logic generates a semantic JSON output. This JSON is then processed by a client-specific presentation layer on the server before being sent to the client.

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
    *   `MemoriesAboutPlayers map[string][]string`: An Owner's private memories about players.
*   **New `OwnerTool` Struct:**
    *   `Name string`: Unique identifier for the tool (e.g., "gift_item", "heal_player", "OWNER_memorize_dependables").
    *   `Description string`: For the LLM to understand the tool's purpose.
    *   `Parameters map[string]interface{}`: JSON schema-like definition of parameters the tool expects.
*   **New `NPCTool` Struct:** (Mirrors `OwnerTool` for NPC-specific actions)
    *   `Name string`: Unique identifier for the tool (e.g., "examine_item", "pickup_item", "NPC_memorize").
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

## 2. Server-Side Presentation Layer

To support various client types (some "dumb," like Telnet), the server will handle the final presentation rendering. The flow is: **Core Logic -> Semantic JSON -> Client-Specific Renderer -> Formatted Output**.

*   **Core Principle:** The server's main game logic produces a "pure" semantic JSON object, describing *what* happened. This JSON is then passed to a specific rendering pipeline on the server based on the connected client's type.

*   **Semantic JSON (Internal Representation):** This is the intermediate, universal format produced by the game engine.
    *   **Example: Narrative Text**
        ```json
        {
          "type": "narrative",
          "payload": {
            "segments": [
              {"text": "The elven guard narrows his eyes. 'State your business,' he says..."},
              {"text": "ancient magic", "semantic_type": "magic_keyword"},
              {"text": " humming in the air."}
            ]
          }
        }
        ```
    *   **Example: Room Description**
        ```json
        {
          "type": "room_update",
          "payload": {
            "name": "The Whispering Glade",
            "description": "Sunlight filters through the dense canopy...",
            "items": [{"name": "a shimmering potion", "semantic_type": "magical_item"}],
            "npcs": [{"name": "An elven guard", "semantic_type": "neutral_npc"}]
          }
        }
        ```

*   **Server-Side Renderers:** For each client type, there will be a corresponding renderer module on the server.
    *   **Telnet Renderer:** This module receives the semantic JSON. It has a ruleset (e.g., a map or config file) that translates `semantic_type` into ANSI escape codes. It would convert the above `narrative` JSON into a raw string like: `The elven guard... \x1b[36m\x1b[3mancient magic\x1b[0m humming...` before sending it over the wire.
    *   **Websocket/API Renderer:** This module would receive the semantic JSON and could either forward it directly to a smart web client or transform it into HTML on the server before sending it. For example, it could turn the JSON into: `<p>The elven guard... <span class="magic-keyword">ancient magic</span> humming...</p>`.

This architecture keeps the core game logic clean and independent of any client's presentation needs, while still supporting dumb clients by handling the rendering for them.

## 3. World-Building and Lore Integration

To ensure a cohesive and immersive world, we will implement a comprehensive and dynamic lore system. This system will serve as the foundational knowledge base for all sentient entities, shaping their understanding of the world, their place in it, and their interactions with players. All lore will be configurable via the web editor.

*   **The "Lore Book":** We will create a central repository for all lore, structured as a collection of `Lore` objects. This allows for modular and easily editable world-building.

*   **Global Lore:** A set of core truths and historical facts about the world that are known to all sentient beings. This includes:
    *   **Cosmology and Creation Myths:** The story of how the world came to be, the major deities, and the fundamental forces that govern reality.
    *   **Major Historical Events:** The timeline of pivotal moments, such as ancient wars, the rise and fall of empires, and cataclysmic events that have shaped the current landscape.
    *   **The "Total Lore":** A concise, high-level summary of the world's state, providing a common foundation for all NPC and Owner knowledge.

*   **Scoped Lore:** To give entities specialized knowledge, lore will be scoped to specific aspects of the game world. When an LLM is prompted, it will receive the Global Lore plus all relevant Scoped Lore.
    *   **Race Lore:** Details the history, culture, traditions, and social norms of each race.
    *   **Profession Lore:** Provides knowledge specific to a character's profession.
    *   **Zone Lore:** Contains detailed information about specific geographical areas, cities, or dungeons.

*   **Creative Lore Categories (for deeper immersion):**
    *   **Faction Lore:** Knowledge related to specific guilds, secret societies, or political movements.
    *   **Bestiary (Creature) Lore:** Information about the world's flora and fauna.
    *   **Item Lore:** Legendary tales and histories associated with specific powerful artifacts.

*   **Lore Integration with LLMs:** The lore system is not just static text. It will be dynamically injected into the context of LLM prompts. When an NPC or Owner is about to generate a response, the system will:
    1.  Fetch the **Global Lore**.
    2.  Fetch the relevant **Scoped Lore** based on the entity's race, profession, location, and faction affiliations.
    3.  Prepend this collection of lore to the prompt as foundational knowledge, ensuring the LLM's response is deeply rooted in the game's established world.

## 4. LLM Integration Strategy (Go-Gemini API & Tool Calling)

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

## 5. Sentient NPCs

NPCs will become dynamic, AI-driven entities whose behavior is shaped by their knowledge and experiences:

*   **AI-Driven Behavior:**
    *   When a player interacts with an NPC, a prompt will be constructed for the LLM.
    *   The prompt will include the relevant **Lore**, the NPC's `PersonalityPrompt`, its personal `MemoriesAboutPlayers`, and `AvailableTools`.
    *   The LLM's response will drive the NPC's narrative and actions.
*   **Personal Memory:**
    *   NPCs are primarily responsible for their own memories. They can use an `NPC_memorize` tool to record their direct interactions with a player. This forms their personal opinion.

## 6. Sentient Owners

Owners are higher-level, LLM-controlled entities that act as guardians or influencers over their domains:

*   **Role and Scope:** Owners monitor broad aspects like locations, races, or factions.
*   **Tool Invocation:** Owners have powerful tools to influence the world (e.g., `TOOL_GIFT_ITEM`, `TOOL_SMITE_PLAYER`).
*   **Prayer Mechanism:** The `pray` command prompts relevant Owners, who can respond based on their lore-informed perspective.

## 7. Reputation and Multi-Layered Memory System

To create a complex social dynamic, memory and reputation will operate on three distinct levels:

*   **Level 1: NPC Personal Memory:**
    *   An NPC's direct experiences with a player.
    *   Managed via an `NPC_memorize` tool available only to that NPC.
    *   Forms the basis of the NPC's private, firsthand opinion of a player.

*   **Level 2: Owner Private Memory:**
    *   An Owner's private thoughts and long-term judgments about a player.
    *   Managed via an `OWNER_memorize` tool available only to that Owner. This tool writes to the Owner's own `MemoriesAboutPlayers` map.
    *   This allows an Owner to maintain a secret opinion of a player, which might differ from its public actions or the information it shares with its subordinates.

*   **Level 3: Owner-Broadcasted Memory (Reputation):**
    *   An Owner's public declarations or official stance on a player.
    *   Managed via a special tool: `OWNER_memorize_dependables(player_id, memory_string)`.
    *   When used, this tool iterates through all NPCs associated with that Owner and imprints the `memory_string` into their `MemoriesAboutPlayers` map.
    *   This is how factions announce heroes or villains. An NPC will receive this as a directive, which will influence its behavior, potentially overriding its personal experiences. For example, an NPC who personally likes a player might turn cold after their Owner broadcasts a message that the player is now an enemy of the faction.

## 8. Skills as Tools

*   **Active Skills:**
    *   These will be defined as `OwnerTool` or `NPCTool` types.
    *   When a player uses an active skill (e.g., "use minor heal"), this command will be routed to the appropriate handler.
    *   For LLM-driven skills, the LLM will be prompted to "use" the skill, and its tool call will trigger the actual effect.
*   **Passive Skills:**
    *   These will not be explicit tool calls but rather modifiers to player stats or contextual information provided to the LLM.
    *   For example, a "City Lore" passive skill might mean the LLM receives additional context about city events when the player is in a city, leading to more informed NPC/Owner responses.

## 9. Game Mechanics Enhancements

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

## 10. Concurrency Considerations

*   **Go Routines for AI Calls:** When multiple AI interactions are triggered, use Go routines (`go func()`) and `sync.WaitGroup` or channels to dispatch these LLM API calls concurrently.
*   **State Management:** Ensure all shared game state is protected by mutexes (`sync.Mutex` or `sync.RWMutex`) to prevent race conditions.

## 11. Web Server (Admin & Editor)

*   The existing lightweight web server will be maintained.
*   It will primarily serve administrative functions, including:
    *   Displaying server statistics.
    *   Providing a web-based editor for game content (rooms, items, NPCs, Owners, and **Lore**).

## 12. High-Level Implementation Phases

1.  **Phase 1: Core Architecture & Data Structures:**
    *   Implement all new and expanded data structures, including the `MemoriesAboutPlayers` map on the `Owner` struct.
    *   Implement the **Server-Side Presentation Layer** and the initial **Telnet Renderer**.
2.  **Phase 2: Lore System & Editor:**
    *   Build the backend logic to store and retrieve lore entries.
    *   Update the web editor to allow for creating and editing lore.
3.  **Phase 3: Basic LLM Integration & Multi-Layered Memory:**
    *   Implement the Go LLM API client and the tool dispatcher.
    *   Implement the `NPC_memorize`, `OWNER_memorize`, and `OWNER_memorize_dependables` tools.
    *   Test the full memory pipeline.
4.  **Phase 4: Sentient NPCs & Owners:**
    *   Integrate NPC AI responses with `talk` and `say`.
    *   Implement the full `Owner` logic and prayer mechanism.
5.  **Phase 5: Advanced Mechanics & Skills:**
    *   Implement locking, mapping, and conditional AI delivery.
    *   Integrate active and passive skills.
6.  **Phase 6: Concurrency & Final Polish:**
    *   Refactor LLM calls for concurrency and review mutex usage.
    *   Build out any additional client renderers.
    *   Finalize the web editor and polish the overall experience.

This proposal outlines a comprehensive path to evolve the GoMUD into a more dynamic and intelligent multi-user dungeon experience.