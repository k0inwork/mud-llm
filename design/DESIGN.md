# GoMUD Project Design Document

## 1. Overview

This document outlines the architectural design for the GoMUD project, an advanced, LLM-driven multi-user dungeon. The goal is to create a dynamic and immersive world where sentient AI entities interact with players in a rich, lore-based environment.

The architecture is designed to be modular, scalable, and maintainable, with a clear separation of concerns between the core game logic, the AI interaction layer, and the client presentation layer. All persistent game data (lore, rooms, items, etc.) will be stored in a local database (e.g., SQLite).

## 2. Core Architectural Principles

*   **Single Source of Truth (Database-Centric):** All game data, including lore, entities, and their states, is stored persistently in a local database (e.g., SQLite). This ensures data consistency and simplifies persistence logic.
*   **Decoupled Presentation:** The server's core logic is completely decoupled from client presentation. The game engine produces a semantic JSON representation of events, which is then translated into client-specific formats by a dedicated server-side rendering layer. This allows for supporting various clients (Telnet, Web, etc.) without altering the core game logic.
*   **API-Driven Content Management:** The web-based administrative editor is a fully decoupled front-end application that interacts with the server via a versioned REST API. This enforces a clean separation between the game server and its management tools, allowing for independent development and deployment of the admin interface.
*   **Sentient, Lore-Driven AI:** All AI entities (NPCs and Owners) are driven by a Large Language Model (LLM). Their behavior, knowledge, and decisions are grounded in the comprehensive lore system stored in the database.
*   **Efficient AI Interaction (Action Significance):** An "Action Significance" model is used to filter and batch player actions, ensuring that the LLM is only triggered for meaningful events. This prevents unnecessary API calls, improves performance, and reduces operational costs.
*   **Clear Command Separation:** Player-initiated commands (e.g., `move`, `unlock`, `use skill`) are handled directly by the Core Game Engine based on game rules. LLM-driven tools are exclusively for AI entities to invoke, allowing them to interact with the game world.
*   **Multi-Layered Memory:** The reputation system is built on a three-tiered memory model (NPC personal, Owner private, and Owner-broadcasted), allowing for complex social dynamics and emergent behavior.

## 3. High-Level Architecture

1.  **Database:** A local database (e.g., SQLite) that stores all persistent world data, including lore, rooms, items, NPCs, and Owners.
2.  **Data Access Layer (DAL):** A Go module that provides a clean, typed API for all database operations. All other modules interact with the database exclusively through the DAL.
3.  **Core Game Engine:** Manages the fundamental state of the MUD and enforces game rules. It uses the DAL to access game data and processes player commands (e.g., `move`, `unlock`). Player-initiated actions are handled here directly.
4.  **AI Interaction Module:**
    *   **Action Significance Monitor:** Tracks player actions from the Core Game Engine, scores them for significance, and batches them for relevant NPCs/Owners.
    *   **LLM Integration Client:** Constructs prompts (using cached data for performance) and handles communication with the LLM API.
    *   **Tool Dispatcher:** Executes Go functions based on tool calls *received from the LLM*.
5.  **Server-Side Presentation Layer:** Contains client-specific renderers (e.g., `TelnetRenderer`) that translate semantic JSON from the core engine into the final format for the client.
6.  **Web Server & Admin API:** A lightweight web server running within the main Go application. It serves the static files for the admin front-end and exposes a versioned **REST API** for all content management operations. The API handlers call the DAL to interact with the database.

## 4. Phased Design Documents

The detailed design for each implementation phase is located in its respective subdirectory.

*   **[Phase 1: Foundation & Content Tools](./phase-1-foundation-and-content-tools/README.md)**
    *   Focuses on establishing the database, core data structures, server-side presentation layer, and the **API-driven web editor**.

*   **[Phase 2: Lore & Data Logic](./phase-2-lore-and-data-logic/README.md)**
    *   Focuses on enhancing the Data Access Layer (DAL) with advanced query methods and in-memory caching for all game data.

*   **[Phase 3: LLM Integration & Memory](./phase-3-llm-integration-and-memory/README.md)**
    *   Focuses on integrating the LLM, implementing prompt caching, and building the multi-layered memory system, clarifying the role of LLM-callable tools.

*   **[Phase 4: Sentient Entities & Action Significance](./phase-4-sentient-entities-and-action-significance/README.md)**
    *   Focuses on bringing NPCs and Owners to life by integrating their AI-driven behavior via the Action Significance model.

*   **[Phase 5: Advanced Mechanics & Skills](./phase-5-mechanics-and-skills/README.md)**
    *   Focuses on implementing player commands for mechanics like locking and skills, and specialized NPC tools that can bypass player rules.

*   **[Phase 6: Concurrency & Final Polish](./phase-6-finalization/README.md)**
    *   Focuses on performance, stability, and the final user experience.