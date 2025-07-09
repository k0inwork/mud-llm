# Phase 3 Implementation Plan

This document outlines the concrete steps to implement the features described in the Phase 3 design document.

## 1. LLM Integration Components

- **Task 1.1: LLM API Client**
  - Create `internal/llm/client.go`.
  - Implement a `Client` struct to handle HTTP requests to an OpenAI-compatible API.
  - Implement a `SendPrompt` method that takes a prompt, handles authentication via environment variables (`LLM_API_KEY`), and returns the parsed JSON response.
  - The response from the LLM will be expected in JSON format: `{"narrative": "...", "tool_calls": [...]}`.

- **Task 1.2: Prompt Assembler**
  - Create `internal/llm/prompt_assembler.go`.
  - Implement a function `AssemblePrompt` that takes an entity (NPC, Owner, or Questmaker) and a player context.
  - This function will fetch data from the DAL:
    - Global and scoped lore.
    - The entity's personality prompt (`LLMPromptContext`).
    - The entity's memories about the player.
    - The entity's available tools.
  - It will combine this data into a single, coherent prompt string.

- **Task 1.3: Prompt Caching**
  - Create `internal/llm/cache_manager.go`.
  - Implement a `CacheManager` with an in-memory map to store the generated "base prompts" for each entity.
  - The cache key will be the entity's ID.
  - Implement `Get` and `Set` methods for the cache.
  - Implement a mechanism to invalidate the cache for an entity when its underlying data (lore, memories) is updated. This will be triggered by the DAL.

- **Task 1.4: LLM Service**
  - Create `internal/llm/service.go`.
  - Implement an `LLMService` that orchestrates the process:
    1.  Checks the `CacheManager` for a base prompt.
    2.  If not found, uses the `PromptAssembler` to create it and caches it.
    3.  Appends the dynamic player action to the base prompt.
    4.  Uses the `client.go` to send the final prompt to the LLM API.

## 2. Data Model & DAL Updates

- **Task 2.1: Create Questmaker Model**
  - Create `internal/models/questmaker.go` to define the `Questmaker` struct as specified in the design document.

- **Task 2.2: Update Existing Models**
  - Modify `internal/models/quest.go`: Add `QuestmakerID` (string) and `InfluencePointsMap` (map[string]float64).
  - Modify `internal/models/player_quest_state.go`: Add `LastQuestActionTimestamp` (time.Time) and `QuestmakerInfluenceAccumulated` (float64).

- **Task 2.3: Create Questmaker DAL**
  - Create `internal/dal/questmaker_dal.go`.
  - Implement standard CRUD functions (`CreateQuestmaker`, `GetQuestmakerByID`, `UpdateQuestmaker`, `DeleteQuestmaker`) and a `GetAllQuestmakers` function.
  - Ensure cache invalidation is implemented for update and delete operations.

- **Task 2.4: Update Main DAL Struct**
  - Modify `internal/dal/dal.go`:
    - Add `QuestmakerDAL` to the `DAL` struct.
    - Instantiate the new `QuestmakerDAL` in the `NewDAL` function.

## 3. Tool Execution System

- **Task 3.1: Tool Dispatcher**
  - Create a new file, potentially `internal/server/tool_dispatcher.go`.
  - Implement a `Dispatch` function that takes the `tool_calls` JSON from the LLM response.
  - It will parse the `tool_name` and `parameters`.
  - It will use a map or a switch statement to call the appropriate Go function based on the `tool_name`.

- **Task 3.2: Implement Memory Tools**
  - These functions will be called by the `ToolDispatcher`.
  - In `internal/dal/npc_dal.go`, implement `AddMemoryToNPC(npcID, playerID, memoryString)`. The `NPC_memorize` tool will call this.
  - In `internal/dal/owner_dal.go`, implement `AddMemoryToOwner(ownerID, playerID, memoryString)`. The `OWNER_memorize` tool will call this.
  - In `internal/dal/owner_dal.go`, implement `BroadcastMemoryToDependents(ownerID, playerID, memoryString)`. The `OWNER_memorize_dependables` tool will call this. This function will need to get all NPCs associated with the owner and call `AddMemoryToNPC` for each.

## 4. System Integration & Seeding

- **Task 4.1: Update Seed Data**
  - Modify `internal/dal/seed.go`.
  - Add the new test data specified in `design/phase-3-llm-integration-and-memory/README.md`, including:
    - The example `Questmaker` ("the_spice_lord").
    - The updated `Quest` ("the_missing_shipment") with its new fields.
    - The example NPC (`innkeeper_bob`) and Owner (`town_council_owner`) with their `AvailableTools` and `PersonalityPrompt` fields populated.

- **Task 4.2: Main Server Integration**
  - Modify `main.go`.
  - Initialize the `LLMService`.
  - Add the new `QuestmakerDAL` to the cache pre-warming sequence on startup.
