# **Proposal: GoMUD World Editor v2.0**

**Author:** Gemini
**Status:** Definitive
**Version:** 3.0
**Related Documents:**
*   [Design: Dynamic Action Significance System (V2.0)](./dynamic_significance_v2.md)

---

## 1. Objective

To create a unified, web-based World Editor that serves as the central hub for all game design and world-building activities. This document provides a specific, field-level blueprint for this tool, ensuring all interconnected systems (`Owners`, `Quests`, `Perception`, etc.) are manageable from a single, intuitive interface.

## 2. UI Philosophy: "The World at a Glance"

The editor will be designed around a clear, hierarchical, and cross-linked interface. A designer should be able to start at a high-level concept like a "Questline" and effortlessly drill down to the specific NPC who gives the quest, the room where it happens, and the territorial laws that govern that room. All relevant data will be presented in context.

---

## 3. Editor Structure: A Multi-Tabbed Workspace

The editor will be organized into three primary workspaces: **World-Building**, **Quest Design**, and **System Rules**.

### **3.1. Workspace Tab: World-Building**

This workspace is for creating and managing the physical, cultural, and social reality of the game world.

#### **3.1.1. View: Territories**
This is the highest-level geographical and cultural editor.

*   **List View:** A table of all defined territories, showing `ID` and `Name`.
*   **Editor Pane:**
    *   `ID`: `string` (e.g., "bree_land") - The unique identifier.
    *   `Name`: `string` (e.g., "Bree-land") - The display name.
    *   `Description`: `textarea` - A description of the territory's overall feel and culture.
    *   **Territory Owner:** `dropdown` - Populated by `Owner` entities where `MonitoredAspect` is "location". This links the territory's rules to an AI persona (e.g., selecting the "Bree Guardian" Owner).
    *   **Default Perception Biases:** `key-value editor` - Defines the baseline cultural norms for this territory. Each entry has a `key` (e.g., "magic", "subterfuge") and a `value` (a float from -1.0 to 1.0). Example: `magic: -0.1`.
    *   **Rooms in Territory:** `read-only list` - A list of all `Room` entities that have this territory set as their `TerritoryID`, with links to edit each room directly.

#### **3.1.2. View: Rooms**
*   **List View:** A searchable, paginated table of all rooms.
*   **Editor Pane:**
    *   `ID`, `Name`, `Description`: Standard text fields.
    *   `Exits`: A dynamic editor to add/remove exits, specifying `direction`, `TargetRoomID` (via a dropdown of all rooms), `IsLocked`, and `KeyID`.
    *   **Territory (`TerritoryID`):** `dropdown` - Populated from the `Territories` table. This assigns the room to a larger cultural zone, inheriting its perception biases.
    *   **Direct NPC Owner (`OwnerID`):** `dropdown` - Populated by all `NPC` entities. This assigns the specific, "real" owner of the room (e.g., selecting "Barliman Butterbur" for "The Prancing Pony"). This NPC will have a high-priority reaction to events in this room.
    *   **Local Perception Overrides (`PerceptionBiases`):** `key-value editor` - Allows this specific room to have perception biases that override the defaults set by its assigned Territory. Example: A "magic-dead" cellar in Bree could have `magic: -0.9`, even if Bree-land is neutral.
    *   `Properties`: `JSON editor` - For any other dynamic room properties.

#### **3.1.3. View: NPCs**
*   **List View:** A searchable table of all NPCs.
*   **Editor Pane:**
    *   `ID`, `Name`, `Description`, `Health`, `MaxHealth`.
    *   `CurrentRoomID`: `dropdown` - To place the NPC in a room.
    *   `Inventory`: `multi-select` - Populated from the `Items` table.
    *   `ReactionThreshold`: `integer` - The significance score needed to trigger an AI reaction.
    *   `PersonalityPrompt`: `textarea` - The core prompt for the NPC's persona.
    *   **Race:** `dropdown` - Assigns the NPC's race (e.g., "Human"). This automatically applies the base racial perception biases, which are visible but not editable here.
    *   **Profession:** `dropdown` - Assigns the NPC's profession (e.g., "Warrior"). This applies professional biases.
    *   **Abstract Owner Memberships (`OwnerIDs`):** `multi-select` - Populated by `Owner` entities. This shows which high-level groups the NPC belongs to (e.g., "Strider" is a member of the "human_elder" and "warrior_guild_master" Owners).
    *   **Quests Given:** `read-only list` - A list of quests that list this NPC as a starting point or objective, with links to the `Quest` editor.

#### **3.1.4. View: Items**
*   **List View:** A searchable table of all items.
*   **Editor Pane:** `ID`, `Name`, `Description`, `Type` (dropdown: key, document, weapon, etc.), `Properties` (JSON editor).

### **3.2. Workspace Tab: Quest & Narrative Design**

This workspace visualizes the entire quest hierarchy, from grand narratives to individual tasks.

#### **3.2.1. View: Quest Owners (The "Why")**
Represents the grand narratives or strategic goals in the world.
*   **List View:** All `QuestOwners` (e.g., "Gandalf's Grand Plan").
*   **Editor Pane:**
    *   `ID`, `Name`, `Description`.
    *   `LLMPromptContext`: `textarea` - The core persona for this narrative force.
    *   `InfluenceBudget`, `MaxInfluenceBudget`, `BudgetRegenRate`: `float` fields.
    *   **Associated Questmakers:** `multi-select` - Populated from the `Questmakers` table. This links the grand plan to the specific `Questmaker` AIs that will execute it.

#### **3.2.2. View: Questmakers (The "How")**
The AI "dungeon masters" for specific questlines.
*   **List View:** All `Questmakers` (e.g., "Urgent Message Quest Controller").
*   **Editor Pane:**
    *   `ID`, `Name`, `LLMPromptContext`, `ReactionThreshold`.
    *   **Parent Quest Owner:** `read-only linked field` - Showing which `QuestOwner` this AI serves.
    *   **Managed Quests:** `multi-select` - Populated from the `Quests` table. Links this AI to the actual `Quest` definitions it is responsible for managing.

#### **3.2.3. View: Quests (The "What")**
The concrete quest data.
*   **List View:** All defined `Quests`.
*   **Editor Pane:**
    *   `ID`, `Name`, `Description`.
    *   **Owning Questmaker:** `dropdown` - Assigns a `Questmaker` AI to this quest.
    *   **Objectives:** `dynamic list editor` - A UI to add/remove objectives. Each objective would have a `Type` dropdown (e.g., "reach_location", "speak_to_npc", "gather_item") and a `TargetID` dropdown populated by the relevant table (rooms, npcs, items).
    *   **Rewards:** `dynamic editor` - For XP, gold, and item rewards (with item selection from a dropdown).

#### **3.2.4. View: Lore**
*   **List View:** All `Lore` entries.
*   **Editor Pane:** `ID`, `Title`, `Content` (textarea), `Scope` (dropdown: global, zone, faction, etc.), `AssociatedID`.

### **3.3. Workspace Tab: System Rules & Biases**

This is for managing the global, data-driven rules of the game's subjective reality.

#### **3.3.1. View: Races**
*   **List View:** All playable/NPC races.
*   **Editor Pane:**
    *   `ID`, `Name`, `Description`, `BaseStats` (key-value editor).
    *   **Racial Perception Biases:** `key-value editor` - Defines the innate perceptual biases for the race (e.g., Elves get `magic: 0.3`).

#### **3.3.2. View: Professions**
*   **List View:** All defined professions.
*   **Editor Pane:**
    *   `ID`, `Name`, `Description`, `BaseSkills` (dynamic list of skill dropdowns).
    *   **Professional Perception Biases:** `key-value editor` - Defines the perceptual modifiers from training (e.g., Mages get `magic: 0.5`).

#### **3.3.3. View: Owners (Abstract Personas)**
The master editor for all high-level `Owner` personas.
*   **List View:** All abstract `Owner` entities.
*   **Editor Pane:**
    *   `ID`, `Name`, `Description`, `LLMPromptContext`, `ReactionThreshold`.
    *   **Monitored Aspect:** `dropdown` - `location`, `race`, `profession`. This determines what this Owner is.
    *   **Associated Entities:** `read-only list` - Shows which Territories, Races, or Professions are linked to this Owner, providing a clear overview of its domain.

---

## 4. Technical Implementation

The technical sketch remains the same: a Go backend serving a RESTful API to a lightweight, dynamic HTML/JS frontend using a library like HTMX or Alpine.js. This comprehensive and specific design ensures that all interconnected systems are manageable from a single, intuitive interface.
