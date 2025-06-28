# Phase 3 Design: LLM Integration & Multi-Layered Memory

## 1. Objectives

This phase is the core of the AI implementation. The goal is to establish a robust connection to the LLM, implement the prompt caching strategy for performance, and build the multi-layered memory system that will form the basis of all NPC and Owner behavior.

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
    1.  Fetch the relevant Global and Scoped Lore.
    2.  Fetch the entity's Personality Prompt.
    3.  Fetch the entity's Memories about the target player.
    4.  Fetch and format the entity's Available Tools.
    5.  Combine all of the above into a single "base prompt" string.
*   **Cache Manager:**
    *   An in-memory cache (e.g., a `map[string]string`) will store the generated "base prompt" for each NPC/Owner.
    *   The cache key could be the NPC/Owner's ID.
    *   The cache for a specific entity will be invalidated (deleted) whenever its underlying memories or lore are updated.
*   **Interaction Flow:**
    1.  On first interaction, the Prompt Assembler generates the base prompt, which is sent to the LLM and stored in the cache.
    2.  On subsequent interactions, the cached base prompt is retrieved, and only the new, dynamic player action is sent to the LLM API (leveraging the API's conversation history).

### 2.3. Tool Dispatcher and Memory Tools

*   **Tool Parser:** A function that takes the raw XML response from the LLM and extracts the JSON from the `<tools>` section.
*   **Tool Dispatcher:** A central function that:
    1.  Receives the parsed tool JSON.
    2.  Uses the `tool_name` to look up and call the corresponding Go function.
    3.  Passes the tool's parameters to the Go function.
*   **Memory Tool Implementation:** The first tools to be implemented will be the memory tools:
    *   `NPC_memorize(player_id, memory_string)`: Adds a string to the `MemoriesAboutPlayers` map of the calling NPC.
    *   `OWNER_memorize(player_id, memory_string)`: Adds a string to the `MemoriesAboutPlayers` map of the calling Owner.
    *   `OWNER_memorize_dependables(player_id, memory_string)`: Iterates through all NPCs associated with the calling Owner and adds the memory string to each of their `MemoriesAboutPlayers` maps. This will require a way to look up NPC dependencies from an Owner.

## 3. Acceptance Criteria

1.  The server can successfully make an API call to a configured LLM endpoint and receive a response.
2.  The prompt caching mechanism works as expected: a base prompt is generated and cached on the first interaction, and subsequent interactions are faster.
3.  Cache invalidation works correctly when a memory is added.
4.  The LLM can successfully call the `NPC_memorize`, `OWNER_memorize`, and `OWNER_memorize_dependables` tools.
5.  The game state correctly reflects the changes made by these tools (i.e., memories are saved to the correct entities).
6.  The system can be tested via a hardcoded trigger (e.g., a player command like `/test-memorize`) that initiates an LLM call.
