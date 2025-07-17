# Current Project Backlog

This document summarizes the completed work for Phase 2 and Phase 4, and outlines the current state of the project.

## Phase 2: Lore & Data Logic (Completed)

All objectives for Phase 2 were successfully met:

*   **Enhanced Data Access Layer (DAL):**
    *   Implemented `GetAll` methods for all core entities.
    *   Implemented advanced query methods like `GetNPCsByRoom`, `GetNPCsByOwner`, `GetOwnersByMonitoredAspect`, `GetItemsInRoom`, `GetPlayerInventory`.
*   **In-Memory Caching within DAL:**
    *   Implemented a generic `Cache` struct with `Set`, `SetMany`, `Get`, `Delete`, and `Clear` methods.
    *   Integrated cache pre-warming for all relevant DAL entities on server startup.
    *   Implemented cache invalidation on `Update` and `Delete` operations within the DALs.
*   **Data Loading and Initialization:**
    *   The server startup sequence (`main.go`) correctly initializes the database and pre-warms the caches.
*   **Refinements during Phase 2 (Design for Phase 3):**
    *   Clarified Quest Roles: Introduced `QuestOwner` entity, refined `Owner` and `Questmaker` roles.
    *   Updated Relationships: Modified `Owner` to include `InitiatedQuests`, and `Quest` to include `QuestOwnerID`.
    *   Budget Allocation: Re-assigned budget types and tool exclusivity.

## Phase 4: Sentient Entities and Action Significance (Completed)

The primary goal of Phase 4 was to make entities "sentient" by implementing a dynamic, perception-based "Action Significance" model to enable intelligent and context-aware reactions to player actions. This phase is now largely complete, with core functionalities implemented and tested.

### Objectives Achieved:

*   **Dynamic Action Significance Model:**
    *   Transitioned from static action scores to a dynamic, perception-based model.
    *   `ActionSignificanceMonitor` no longer holds `baseSignificanceScores`; `PerceptionFilter` now determines `BaseSignificance`.
*   **Event-Driven Architecture:**
    *   Implemented an event-driven, two-speed architecture: synchronous for local reactions and asynchronous for global propagation.
    *   `EventBus` now uses `EventType` for publishing and subscribing.
*   **Perception Filtering:**
    *   Implemented `PerceptionFilter` to process `ActionEvent` into `PerceivedAction` (subjective interpretation).
    *   `PerceivedAction` now includes `BaseSignificance`.
    *   Perception is influenced by layered factors: physical sensory limits, innate racial biases, territorial/cultural biases, professional knowledge/experience, and explicit passive skills/buffs.
    *   Extended perception filters to apply to players as observers.
*   **Sentient Entity Reaction Triggering:**
    *   `ReactionThreshold` in `NPC`, `Owner`, and `Questmaker` models determines when an entity reacts.
    *   `ActionSignificanceMonitor` triggers reactions based on cumulative significance.
*   **Dependency Injection and Testability:**
    *   Introduced interfaces (`RoomDALInterface`, `RaceDALInterface`, `ProfessionDALInterface`, `CacheInterface`, `NPCDALInterface`, `OwnerDALInterface`, `QuestmakerDALInterface`, `PerceptionFilterInterface`, `SentientEntityManagerInterface`, `LLMServiceInterface`, `ToolDispatcherInterface`, `TelnetRendererInterface`) for better testability and dependency injection.
    *   Refactored `ActionSignificanceMonitor`, `PerceptionFilter`, `GlobalObserverManager`, and `SentientEntityManager` to use these interfaces.
*   **Database Schema Updates:**
    *   Added `RaceID` and `ProfessionID` to `NPC` model and `NPCs` table.
    *   Added `Category` to `Skill` model and `Skills` table.
    *   Updated `Room` entries with `TerritoryID` and `PerceptionBiases`.
    *   Updated `Quest` model to store `InfluencePointsMap` as JSON string.
*   **Initial Unit Test Coverage:**
    *   Created and refined unit tests for `ActionSignificanceMonitor`, `PerceptionFilter`, `GlobalObserverManager`, and `SentientEntityManager`.
    *   Addressed various compilation and runtime errors during development and testing.
*   **Automated Testing Script:**
    *   Created `run_tests.sh` to automate starting the server, running all tests (including Telnet server tests), and cleaning up the database and server process.

### Remaining TODOs for Phase 4 (or subsequent phases):

1.  **Complete Unit Tests:**
    *   Create and implement comprehensive unit tests for `SentientEntityManager` (initial tests are present, but more detailed scenarios are needed).
2.  **Player Perception Analysis:**
    *   Analyze how the actions/events modifications by perception filters on NPCs and other sentient entities may be expanded to players (perceiving the NPCs or other players). (Initial implementation for players as observers is done, but further analysis and specific tests are needed).
3.  **End-to-End Flow Testing:**
    *   Conduct thorough manual and automated end-to-end testing of the entire system, focusing on player interactions triggering sentient entity reactions and LLM responses.
4.  **LLM Integration Refinement:**
    *   Implement the `toolDispatcher.Dispatch` calls within `SentientEntityManager` (currently commented out in tests).
    *   Refine prompt assembly in `llm/service.go` to leverage `PerceivedAction` details more effectively.
5.  **Influence Budget Mechanics:**
    *   Implement the regeneration logic for `CurrentInfluenceBudget` for Owners and Questmakers.
    *   Define and implement specific tools for Owners and Questmakers to spend their influence budget.
6.  **Behavior State Management:**
    *   Develop and integrate logic for `BehaviorState` in `NPC` and other sentient entities to influence their reactions and actions over time.
7.  **Advanced Perception Layers:**
    *   Implement Layer 0 (Physical Sensory Check) and Layer 3 (Explicit Modifiers - Passive Skills & Buffs) in `PerceptionFilter`.
    *   Integrate Skill Proficiency into perception calculations.
