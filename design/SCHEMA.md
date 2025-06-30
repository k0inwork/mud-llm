# GoMUD Database Schema Proposal

This document outlines the proposed database schema for the GoMUD project, adhering to the "Single Source of Truth (Database-Centric)" architectural principle. All persistent game data will be stored in a local SQLite database. This schema will serve as the foundation for the Data Access Layer (DAL).

## 1. Core Entities

### 1.1. `Players` Table

Stores persistent player data.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique player identifier (e.g., UUID)          |
| `name`               | TEXT      | NOT NULL, UNIQUE      | Player's chosen name                           |
| `race_id`            | TEXT      | NOT NULL, FOREIGN KEY | Reference to `Races` table                     |
| `profession_id`      | TEXT      | NOT NULL, FOREIGN KEY | Reference to `Professions` table               |
| `current_room_id`    | TEXT      | NOT NULL, FOREIGN KEY | Current room player is in (Reference to `Rooms` table) |
| `health`             | INTEGER   | NOT NULL              | Current health points                          |
| `max_health`         | INTEGER   | NOT NULL              | Maximum health points                          |
| `inventory`          | JSON      | NOT NULL              | JSON array of item IDs and quantities          |
| `visited_room_ids`   | JSON      | NOT NULL              | JSON array of room IDs visited                 |
| `created_at`         | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Record creation timestamp              |
| `last_login_at`      | TIMESTAMP |                       | Last login timestamp                           |
| `last_logout_at`     | TIMESTAMP |                       | Last logout timestamp                          |

### 1.2. `Rooms` Table

Stores information about game rooms/locations.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique room identifier                         |
| `name`               | TEXT      | NOT NULL              | Display name of the room                       |
| `description`        | TEXT      | NOT NULL              | Detailed description of the room               |
| `exits`              | JSON      | NOT NULL              | JSON object mapping directions to room IDs (e.g., `{"north": "room_id_2"}`) |
| `owner_id`           | TEXT      | FOREIGN KEY           | Optional: ID of the Owner controlling this room |
| `properties`         | JSON      |                       | JSON object for dynamic room properties (e.g., `{"locked_exits": ["north"], "weather": "rain"}`) |

### 1.3. `Items` Table

Stores definitions of all item types.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique item identifier                         |
| `name`               | TEXT      | NOT NULL              | Display name of the item                       |
| `description`        | TEXT      | NOT NULL              | Detailed description of the item               |
| `type`               | TEXT      | NOT NULL              | Item type (e.g., "weapon", "armor", "consumable", "quest_item") |
| `properties`         | JSON      |                       | JSON object for item-specific properties (e.g., `{"damage": 10, "healing": 20}`) |

### 1.4. `NPCs` Table

Stores definitions and current states of Non-Player Characters.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique NPC identifier                          |
| `name`               | TEXT      | NOT NULL              | Display name of the NPC                        |
| `description`        | TEXT      | NOT NULL              | Detailed description of the NPC                |
| `current_room_id`    | TEXT      | NOT NULL, FOREIGN KEY | Current room NPC is in (Reference to `Rooms` table) |
| `health`             | INTEGER   | NOT NULL              | Current health points                          |
| `max_health`         | INTEGER   | NOT NULL              | Maximum health points                          |
| `inventory`          | JSON      | NOT NULL              | JSON array of item IDs and quantities          |
| `owner_ids`          | JSON      | NOT NULL              | JSON array of Owner IDs this NPC is associated with |
| `memories_about_players` | JSON      | NOT NULL              | JSON object mapping player IDs to arrays of memory strings |
| `personality_prompt` | TEXT      | NOT NULL              | Base prompt for LLM defining NPC personality   |
| `available_tools`    | JSON      | NOT NULL              | JSON array of conceptual tools LLM can call    |
| `behavior_state`     | JSON      |                       | JSON object for dynamic behavior state (e.g., `{"disposition": "friendly", "target_player": "player_id"}`) |

### 1.5. `Owners` Table

Stores definitions and current states of LLM-driven world guardians.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique Owner identifier                        |
| `name`               | TEXT      | NOT NULL              | Display name of the Owner                      |
| `description`        | TEXT      | NOT NULL              | Detailed description of the Owner              |
| `monitored_aspect`   | TEXT      | NOT NULL              | Defines what the Owner primarily monitors      |
| `associated_id`      | TEXT      | NOT NULL              | ID of the entity/area/faction this Owner is associated with |
| `llm_prompt_context` | TEXT      | NOT NULL              | Base prompt for LLM defining Owner personality/goals |
| `memories_about_players` | JSON      | NOT NULL              | JSON object mapping player IDs to arrays of private memory strings |
| `current_influence_budget` | REAL      | NOT NULL              | Current points available for actions           |
| `max_influence_budget` | REAL      | NOT NULL              | Maximum capacity of influence points           |
| `budget_regen_rate`  | REAL      | NOT NULL              | Points regenerated per game tick/significant event |
| `available_tools`    | JSON      | NOT NULL              | JSON array of conceptual tools LLM can call    |

### 1.6. `Quests` Table

Stores definitions of quests.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique quest identifier                        |
| `name`               | TEXT      | NOT NULL              | Display name of the quest                      |
| `description`        | TEXT      | NOT NULL              | Detailed description of the quest              |
| `questmaker_id`      | TEXT      | NOT NULL, FOREIGN KEY | ID of the associated Questmaker (Reference to `Questmakers` table) |
| `influence_points_map` | JSON      | NOT NULL              | JSON object mapping player actions to influence points granted |
| `objectives`         | JSON      | NOT NULL              | JSON array of quest objectives (e.g., `[{"type": "kill", "target_id": "goblin_chieftain", "count": 1}]`) |
| `rewards`            | JSON      | NOT NULL              | JSON array of rewards (e.g., `[{"type": "item", "id": "spice_crate", "quantity": 3}]`) |

### 1.7. `Questmakers` Table

Stores definitions and current states of LLM-driven quest entities.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique Questmaker identifier                   |
| `name`               | TEXT      | NOT NULL              | Display name of the Questmaker                 |
| `llm_prompt_context` | TEXT      | NOT NULL              | Base prompt for LLM defining Questmaker personality/goals |
| `current_influence_budget` | REAL      | NOT NULL              | Current points available for actions           |
| `max_influence_budget` | REAL      | NOT NULL              | Maximum capacity of influence points           |
| `budget_regen_rate`  | REAL      | NOT NULL              | Points regenerated per game tick/significant event |
| `memories_about_players` | JSON      | NOT NULL              | JSON object mapping player IDs to arrays of private memory strings |
| `available_tools`    | JSON      | NOT NULL              | JSON array of conceptual tools LLM can call    |

### 1.8. `PlayerQuestStates` Table

Stores the dynamic state of a player's progress on active quests.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `player_id`          | TEXT      | PRIMARY KEY, NOT NULL | Reference to `Players` table                   |
| `quest_id`           | TEXT      | PRIMARY KEY, NOT NULL | Reference to `Quests` table                    |
| `current_progress`   | JSON      | NOT NULL              | JSON object tracking objective progress        |
| `last_action_timestamp` | TIMESTAMP | NOT NULL              | Timestamp of the last relevant player action   |
| `questmaker_influence_accumulated` | REAL      | NOT NULL              | Points player has "given" to the Questmaker  |
| `status`             | TEXT      | NOT NULL              | Quest status (e.g., "active", "completed", "failed", "abandoned") |

## 2. Supporting Tables

### 2.1. `Races` Table

Defines available player races.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique race identifier                         |
| `name`               | TEXT      | NOT NULL, UNIQUE      | Display name of the race                       |
| `description`        | TEXT      | NOT NULL              | Description of the race                        |
| `base_stats`         | JSON      | NOT NULL              | JSON object of base stats for the race         |

### 2.2. `Professions` Table

Defines available player professions/classes.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique profession identifier                   |
| `name`               | TEXT      | NOT NULL, UNIQUE      | Display name of the profession                 |
| `description`        | TEXT      | NOT NULL              | Description of the profession                  |
| `base_skills`        | JSON      | NOT NULL              | JSON array of skill IDs granted by profession  |

### 2.3. `Lore` Table

Stores various pieces of lore, categorized by scope.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique lore identifier                         |
| `title`              | TEXT      | NOT NULL              | Title of the lore entry                        |
| `content`            | TEXT      | NOT NULL              | Full text content of the lore                  |
| `scope`              | TEXT      | NOT NULL              | Scope of the lore (e.g., "global", "zone", "faction", "item") |
| `associated_id`      | TEXT      |                       | ID of entity/zone/faction if scope is not global |

### 2.4. `Skills` Table

Defines available skills (active and passive).

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique skill identifier                        |
| `name`               | TEXT      | NOT NULL, UNIQUE      | Display name of the skill                      |
| `description`        | TEXT      | NOT NULL              | Detailed description of the skill              |
| `type`               | TEXT      | NOT NULL              | Skill type (e.g., "active", "passive")       |
| `effects`            | JSON      | NOT NULL              | JSON array of structured effect objects (see Phase 5 for details) |
| `cost`               | INTEGER   |                       | Cost to use skill (e.g., mana, stamina)        |
| `cooldown`           | INTEGER   |                       | Cooldown in seconds                            |

### 2.5. `PlayerSkills` Table

Tracks skills learned by players and their levels.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `player_id`          | TEXT      | PRIMARY KEY, NOT NULL | Reference to `Players` table                   |
| `skill_id`           | TEXT      | PRIMARY KEY, NOT NULL | Reference to `Skills` table                    |
| `level`              | INTEGER   | NOT NULL, DEFAULT 1   | Current level of the skill                     |
| `experience`         | INTEGER   | NOT NULL, DEFAULT 0   | Current experience towards next level          |

### 2.6. `ActionSignificanceConfig` Table

Configures significance scores for player actions.

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `action_type`        | TEXT      | PRIMARY KEY, NOT NULL | Type of player action (e.g., "move", "attack") |
| `score`              | INTEGER   | NOT NULL              | Significance score for this action type        |

### 2.7. `LLMToolDefinitions` Table

Stores definitions of conceptual tools available to LLM entities (NPCs, Owners, Questmakers).

| Column Name          | Data Type | Constraints           | Description                                    |
|----------------------|-----------|-----------------------|------------------------------------------------|
| `id`                 | TEXT      | PRIMARY KEY, NOT NULL | Unique tool identifier (e.g., "send_message") |
| `name`               | TEXT      | NOT NULL, UNIQUE      | Display name of the tool                       |
| `description`        | TEXT      | NOT NULL              | Description of what the tool does              |
| `parameters_schema`  | JSON      | NOT NULL              | JSON schema defining expected parameters       |
| `base_cost`          | REAL      | NOT NULL              | Base influence budget cost to use this tool    |
| `entity_type`        | TEXT      | NOT NULL              | Type of entity that can use this tool (e.g., "NPC", "Owner", "Questmaker", "All") |
