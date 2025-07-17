# **Implementation Plan: Dynamic Action Significance (Phase 4)**

**Author:** Gemini
**Status:** Proposed
**Version:** 1.0
**Related Documents:**
*   [Design Proposal: Dynamic Action Significance System (V2.0)](../dynamic_significance_v2.md)
*   [Original Phase 4 Readme](./README.md)

---

## 1. Objective

This document outlines the engineering plan to implement the **Dynamic Action Significance System (V2.0)**. This will replace the current, simpler significance model with the new event-driven, perception-based architecture. The goal is to deliver a robust, testable, and extensible system for emergent AI behavior.

This work constitutes the entirety of **Phase 4: Sentient Entities and Action Significance**.

## 2. Implementation Strategy: A Phased, Bottom-Up Approach

We will implement the system in a bottom-up fashion, starting with the foundational architectural components and progressively integrating them into the existing game logic. This minimizes disruption and allows for testing at each stage.

### **Phase 4.1: Architectural Foundation**

*Goal: Create the core, non-functional scaffolding of the new system.*

1.  **Event Bus Implementation:**
    *   Create a new package: `internal/game/events`.
    *   Create `internal/game/events/event_bus.go`.
    *   Implement a simple, thread-safe, in-memory event bus. It should support `Subscribe` and `Publish` methods for `ActionEvent`. This will be a singleton instance initialized in `main.go`.

2.  **Core Data Structures:**
    *   Create `internal/game/events/action_event.go` to define the `ActionEvent` struct.
    *   Create a new package: `internal/game/perception`.
    *   Create `internal/game/perception/perceived_action.go` to define the `PerceivedAction` struct.

3.  **Model & DAL Updates for Perception:**
    *   Modify the `Race`, `Profession`, and `Room` (Territory) models in `internal/models/` to include data fields for perception biases.
        *   Example for `Race`: `PerceptionBiases map[string]float64 // e.g., {"magic": -0.3, "subterfuge": 0.1}`
    *   Update the corresponding DALs in `internal/dal/` to load this new data.
    *   Update `internal/dal/seed.go` with initial bias data for the existing races and territories.

### **Phase 4.2: The Perception Filter**

*Goal: Implement the logic that translates objective events into subjective perceptions.*

1.  **Perception Service:**
    *   Create `internal/game/perception/filter.go`.
    *   This file will contain the `PerceptionFilter` service. It will have a primary method: `func (f *PerceptionFilter) Apply(event *events.ActionEvent, observer interface{}) *PerceivedAction`.

2.  **Implement Filter Layers:**
    *   **Layer 0 (Sensory):** Implement logic within `Apply` to check for blindness, distance, etc.
    *   **Layer 1 (Bias):** The filter will use the DALs to fetch racial and territorial biases and calculate the base `Clarity`.
    *   **Layer 2 (Knowledge):** Implement logic to compare the observer's profession/class with the action's nature to refine `Clarity`.
    *   **Layer 3 (Skills):** Implement a mechanism to check the observer for passive skills that modify perception and apply their effects to `Clarity`.

### **Phase 4.3: Integration with the Game Loop**

*Goal: Wire the new system into the live game server.*

1.  **Refactor `ActionSignificanceMonitor`:**
    *   Modify `internal/game/actionsignificance/monitor.go`.
    *   The `LogPlayerAction` function will be removed.
    *   A new method, `HandleActionEvent(event *events.ActionEvent)`, will be created. This method will be a subscriber to the event bus.
    *   Inside `HandleActionEvent`, it will iterate through local observers, use the `PerceptionFilter` to get a `PerceivedAction` for each, calculate the significance score using the new formula, and trigger the `SentientEntityManager` as needed.

2.  **Update `TelnetServer`:**
    *   Modify `internal/server/telnet_server.go`.
    *   The `handleCommand` function will no longer call `LogPlayerAction` directly.
    *   Instead, after a command is successfully executed, it will construct the appropriate `ActionEvent` struct with all the ground-truth data.
    *   It will then publish this `ActionEvent` to the global event bus.

3.  **Update `main.go`:**
    *   Initialize the `EventBus`, `PerceptionFilter`, and other new services.
    *   Inject dependencies as required (e.g., the `PerceptionFilter` will need the DALs).
    *   Subscribe the `ActionSignificanceMonitor` to the `EventBus`.

### **Phase 4.4: Asynchronous Propagation & Global Observers**

*Goal: Implement the "ripple effect" for non-local entities.*

1.  **Faction Awareness Service:**
    *   Create a new service, e.g., `internal/game/awareness/faction_monitor.go`.
    *   This service will subscribe to the `EventBus`.
    *   When it receives an event, it will process it in the background (as a goroutine).
    *   It will determine if the action is relevant to any major factions (e.g., an attack on a faction member).
    *   It will then update a new data store (e.g., a `faction_reputations` table) to reflect the change in the player's standing. This avoids direct, immediate reaction and models the slow spread of information.

## 3. Testing Strategy

*   **Unit Tests:** Each new service (`EventBus`, `PerceptionFilter`) will have extensive unit tests. The `PerceptionFilter` tests will be data-driven, with various `ActionEvent` and observer combinations to assert the correct `PerceivedAction` is generated.
*   **Integration Tests:** An integration test will be created to simulate a player action in the `TelnetServer`, asserting that the correct `ActionEvent` is published and that the `ActionSignificanceMonitor` correctly processes it.
*   **End-to-End Testing:** Manual testing will be performed using the scenarios outlined in the original design proposal to validate that NPCs react believably to different actions.

---
This implementation plan provides a clear, step-by-step path to delivering the new dynamic significance system.
