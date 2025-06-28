# Phase 3 Design: LLM Integration & Memory

## 1. Objectives

This phase is the core of the AI implementation. The goal is to establish a robust connection to the LLM, implement the prompt caching strategy for performance, and build the multi-layered memory system that will form the basis of all NPC and Owner behavior. This phase also clarifies the distinction between direct player commands and LLM-driven tool usage.

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
    *   An in-memory cache (e.g., a `map[string]string`) will store the generated "base prompt" for each NPC/Owner.
    *   The cache key could be the NPC/Owner's ID.
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

*   **Player Commands:** Direct player commands (e.g., `move`, `unlock`, `use skill`) are handled by the **Core Game Engine**. These commands directly modify the game state based on game rules (e.g., checking for keys to unlock a door). They do *not* involve an LLM call for their execution.
*   **LLM Tools:** Tools defined in `OwnerTool` and `NPCTool` structs are exclusively for the **LLM to invoke**. They represent actions that AI entities can take within the game world. The outcome of a player command *may* trigger an AI reaction, which then involves the LLM and its tools.

## 3. Acceptance Criteria

1.  The server can successfully make an API call to a configured LLM endpoint and receive a response.
2.  The prompt caching mechanism works as expected: a base prompt is generated and cached on the first interaction, and subsequent interactions are faster.
3.  Cache invalidation works correctly when a memory is added (via DAL update notifications).
4.  The LLM can successfully call the `NPC_memorize`, `OWNER_memorize`, and `OWNER_memorize_dependables` tools.
5.  The game state correctly reflects the changes made by these tools (i.e., memories are saved to the correct entities in the database via DAL).
6.  The distinction between player commands and LLM tools is clear in the code structure, with player commands handled directly by the game engine.