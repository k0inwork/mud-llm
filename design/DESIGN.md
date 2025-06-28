# GoMUD Project Design Document

## 1. Overview

This document outlines the architectural design for the GoMUD project, an advanced, LLM-driven multi-user dungeon. The goal is to create a dynamic and immersive world where sentient AI entities interact with players in a rich, lore-based environment.

The architecture is designed to be modular, scalable, and maintainable, with a clear separation of concerns between the core game logic, the AI interaction layer, and the client presentation layer.

## 2. Core Architectural Principles

*   **Decoupled Presentation:** The server's core logic is completely decoupled from client presentation. The game engine produces a semantic JSON representation of events, which is then translated into client-specific formats by a dedicated server-side rendering layer. This allows for supporting various clients (Telnet, Web, etc.) without altering the core game logic.
*   **Sentient, Lore-Driven AI:** All AI entities (NPCs and Owners) are driven by a Large Language Model (LLM). Their behavior, knowledge, and decisions are grounded in a comprehensive, editable lore system, creating a consistent and believable world.
*   **Multi-Layered Memory:** The reputation system is built on a three-tiered memory model (NPC personal, Owner private, and Owner-broadcasted), allowing for complex social dynamics and emergent behavior.
*   **Modular, Phased Implementation:** The project is broken down into distinct implementation phases, each with a clear set of objectives. This allows for iterative development and testing.

## 3. High-Level Architecture

The system is composed of the following key modules:

1.  **Core Game Engine:** Manages the fundamental state of the MUD—players, rooms, items, and the enforcement of game rules (e.g., movement, combat).
2.  **Lore & Data Module:** Provides access to all game data, including the world's lore, which is stored in a structured format.
3.  **LLM Integration Module:** Handles all communication with the LLM API. This includes constructing prompts, managing the prompt cache for performance, and parsing the LLM's responses.
4.  **Sentient Entity Manager:** Orchestrates the behavior of NPCs and Owners, triggering LLM prompts based on player actions and game events.
5.  **Server-Side Presentation Layer:** Contains client-specific renderers (e.g., `TelnetRenderer`) that translate the semantic JSON from the core engine into the final format for the client.
6.  **Web Server (Admin & Editor):** A lightweight web server providing an administrative interface and a content editor for managing lore, rooms, items, and other game entities.

## 4. Phased Design Documents

The detailed design for each implementation phase is located in its respective subdirectory.

*   **[Phase 1: Core Architecture & Data Structures](./phase-1-architecture/README.md)**
    *   Focuses on establishing the foundational data models and the server-side presentation layer.

*   **[Phase 2: Lore System & Editor](./phase-2-lore/README.md)**
    *   Focuses on building the systems to store, manage, and edit the world's lore.

*   **[Phase 3: LLM Integration & Caching](./phase-3-llm-integration/README.md)**
    *   Focuses on integrating with the LLM, implementing the prompt caching strategy, and the multi-layered memory tools.

*   **[Phase 4: Sentient NPCs & Owners](./phase-4-sentient-entities/README.md)**
    *   Focuses on bringing the NPCs and Owners to life by integrating their AI-driven behavior.

*   **[Phase 5: Advanced Mechanics & Skills](./phase-5-mechanics-and-skills/README.md)**
    *   Focuses on implementing key gameplay features like locking, mapping, and the skills system.

*   **[Phase 6: Concurrency & Final Polish](./phase-6-finalization/README.md)**
    *   Focuses on performance, stability, and the final user experience.
