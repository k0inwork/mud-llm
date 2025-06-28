# Phase 1 Design: Core Architecture & Data Structures

## 1. Objectives

This initial phase focuses on laying the foundational groundwork for the entire project. The primary goal is to establish the core data structures and the server architecture, particularly the decoupled presentation layer. This will allow us to have a runnable, albeit simple, server by the end of this phase.

## 2. Key Components to be Implemented

### 2.1. Go Data Structures

All data structures outlined in Section 1 of the main `DESIGN.md` and the `proposal.md` will be implemented in Go. This includes:
*   `Player`
*   `NPC`
*   `Owner`
*   `OwnerTool`
*   `NPCTool`
*   `Room`
*   `Exit`
*   `Item`
*   `Lore`
*   The `MemoriesAboutPlayers` map will be added to both the `NPC` and `Owner` structs.

### 2.2. Server-Side Presentation Layer

This is the most critical component of this phase.
*   **Semantic JSON Format:** A clear, versioned definition of the internal semantic JSON format will be established. This will include schemas for all message types (e.g., `narrative`, `room_update`, `player_stats`).
*   **Telnet Renderer Module:** A Go module responsible for:
    *   Receiving semantic JSON objects from the core game logic.
    *   Maintaining a configurable mapping of `semantic_type` (e.g., `magic_keyword`, `neutral_npc`) to specific ANSI escape codes for color and style.
    *   Constructing the final, raw string with embedded ANSI codes to be sent to the Telnet client.
*   **Main Server Loop:** The main server loop will be structured to correctly route game events through the presentation layer (Core Logic -> Semantic JSON -> Telnet Renderer -> Client).

### 2.3. Basic Telnet Server

*   A basic Telnet server will be created that can handle multiple simultaneous connections.
*   It will be able to send the output from the `TelnetRenderer` to the appropriate client.
*   It will accept basic player input, though command handling will be minimal in this phase.

## 3. Acceptance Criteria

By the end of this phase, the following must be demonstrable:

1.  The Go project compiles successfully with all defined data structures.
2.  A player can connect to the MUD using a standard Telnet client.
3.  The server can manually generate a semantic JSON object (e.g., a hardcoded room description) and have it correctly rendered and displayed to the connected Telnet client with the appropriate colors and styles.
4.  The server can handle multiple connections without crashing.
5.  The project structure is clean, with a clear separation between the (currently minimal) core logic and the presentation layer.
