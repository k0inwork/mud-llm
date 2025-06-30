# Skills and Classes System Proposal

This document outlines a new design for the Skills and Classes system, addressing the "Weak Design: The Skills data model" and "Incompleteness: The design doesn't state how players acquire or level up skills" criticisms from `critics.md`. This proposal introduces a class-tree progression model for skill acquisition and a class-level-based cap for skill percentages.

## 1. Core Concepts

*   **Classes as Skill Trees:** Classes are not just labels but serve as progression paths that unlock and enhance skills. Each class defines a set of skills available to its members.
*   **Skill Acquisition via Class Leveling:** Players acquire new skills or improve existing ones by leveling up their chosen classes. At certain class levels, players may be presented with choices for which skills to acquire or specialize in.
*   **Skill Percentage and Class Level Cap:** Every skill has a percentage (0-100%) representing its effectiveness or mastery. This percentage is capped by the player's current level in the associated class. For example, if a class has a total of 5 levels, each level might add 20% to the maximum achievable skill percentage for skills belonging to that class.

## 2. Data Models

### 2.1. `Class` Table

Defines the available player classes and their associated skill trees. Classes can now also represent progression tracks tied to specific LLM entities (Owners/Questmakers).

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `id` | TEXT | PRIMARY KEY, NOT NULL | Unique class identifier (e.g., "warrior", "mage", "spice_lord_patron") |
| `name` | TEXT | NOT NULL, UNIQUE | Display name of the class |
| `description` | TEXT | NOT NULL | Detailed description of the class |
| `total_levels` | INTEGER | NOT NULL, DEFAULT 5 | Total number of levels for this class (used for skill cap calculation) |
| `parent_class_id` | TEXT | FOREIGN KEY | Optional: ID of the parent class in a class tree (e.g., "fighter" for "warrior") |
| `associated_entity_type` | TEXT | | Optional: Type of LLM entity this class is tied to (e.g., "Questmaker", "Owner") |
| `associated_entity_id` | TEXT | FOREIGN KEY | Optional: ID of the specific LLM entity (e.g., "the_spice_lord", "town_council_owner") |
| `level_up_rewards` | JSON | NOT NULL | JSON object mapping level to skill choices/unlocks (e.g., `{"1": {"unlock_skill": "basic_attack"}, "3": {"choose_skill": ["cleave", "shield_bash"]}}`) |

### 2.2. `Skill` Table

Defines all available skills, their types, and structured effects.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `id` | TEXT | PRIMARY KEY, NOT NULL | Unique skill identifier (e.g., "sword_mastery", "fireball") |
| `name` | TEXT | NOT NULL, UNIQUE | Display name of the skill |
| `description` | TEXT | NOT NULL | Detailed description of the skill |
| `type` | TEXT | NOT NULL | Skill type (e.g., "active", "passive", "utility") |
| `associated_class_id` | TEXT | FOREIGN KEY | Optional: The class this skill primarily belongs to |
| `granted_by_entity_type` | TEXT | | Optional: Type of LLM entity that can grant this skill (e.g., "Questmaker", "Owner") |
| `granted_by_entity_id` | TEXT | FOREIGN KEY | Optional: ID of the specific LLM entity that can grant this skill |
| `effects` | JSON | NOT NULL | JSON array of structured effect objects (see 2.3. `Effect` below) |
| `cost` | INTEGER | | Cost to use skill (e.g., mana, stamina) |
| `cooldown` | INTEGER | | Cooldown in seconds |
| `min_class_level` | INTEGER | NOT NULL, DEFAULT 1 | Minimum class level required to acquire this skill |

### 2.3. `Effect` Structure (within `Skill.effects` JSON)

This structured model replaces the generic map, providing clarity and scalability for skill effects.

```json
[
  {
    "type": "DAMAGE",
    "target": "ENEMY",
    "value_formula": "base_damage + (skill_percentage * 0.5)",
    "damage_type": "physical"
  },
  {
    "type": "HEAL",
    "target": "SELF",
    "value_formula": "base_heal + (skill_percentage * 0.2)"
  },
  {
    "type": "MODIFY_ATTRIBUTE",
    "target": "SELF",
    "attribute": "strength",
    "value_formula": "skill_percentage * 0.1",
    "duration_seconds": 30
  },
  {
    "type": "STATUS_EFFECT",
    "target": "ENEMY",
    "status_id": "stunned",
    "duration_seconds": "skill_percentage * 0.05"
  },
  {
    "type": "UNLOCK_ACTION",
    "action_id": "lockpick_door"
  }
]
```
*   `type`: (Enum) e.g., "DAMAGE", "HEAL", "MODIFY_ATTRIBUTE", "STATUS_EFFECT", "UNLOCK_ACTION".
*   `target`: (Enum) e.g., "SELF", "ENEMY", "ALLY", "AREA".
*   `value_formula`: (String) A formula string that the game engine evaluates, allowing skill percentage to dynamically influence the effect's magnitude.
*   Other parameters specific to the `type` (e.g., `damage_type`, `attribute`, `status_id`, `duration_seconds`).

### 2.4. `PlayerClass` Table

Tracks a player's progression in each class they have acquired.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `player_id` | TEXT | PRIMARY KEY, NOT NULL, FOREIGN KEY | Reference to `Players` table |
| `class_id` | TEXT | PRIMARY KEY, NOT NULL, FOREIGN KEY | Reference to `Classes` table |
| `level` | INTEGER | NOT NULL, DEFAULT 1 | Current level in this class |
| `experience` | INTEGER | NOT NULL, DEFAULT 0 | Current experience towards next level in this class |

### 2.5. `PlayerSkill` Table

Tracks skills learned by players and their current percentage.

| Column Name | Data Type | Constraints | Description |
|---|---|---|---|
| `player_id` | TEXT | PRIMARY KEY, NOT NULL, FOREIGN KEY | Reference to `Players` table |
| `skill_id` | TEXT | PRIMARY KEY, NOT NULL, FOREIGN KEY | Reference to `Skills` table |
| `percentage` | INTEGER | NOT NULL, DEFAULT 0 | Current mastery percentage (0-100) |
| `granted_by_entity_type` | TEXT | | Optional: Type of LLM entity that granted this skill (e.g., "Questmaker", "Owner") |
| `granted_by_entity_id` | TEXT | FOREIGN KEY | Optional: ID of the specific LLM entity that granted this skill |

## 3. Skill Acquisition and Progression Logic

### 3.1. Class Acquisition

*   Players can acquire new classes through various means:
    *   **Starting Class:** Chosen at character creation.
    *   **Trainers:** NPCs who offer to teach new classes (e.g., a "Knight Trainer" offering the "Knight" class).
    *   **Quests:** Completing specific quests might unlock new classes.
    *   **Discovery:** Fulfilling certain hidden criteria (e.g., using a specific weapon type extensively) could unlock a class.

### 3.2. Class Leveling

*   **Experience Gain:** Players gain experience for their active class by performing actions relevant to that class (e.g., a "Warrior" gains XP for combat, a "Mage" for casting spells, a "Crafter" for crafting).
*   **Level Up:** When `PlayerClass.experience` reaches a threshold, `PlayerClass.level` increases.
*   **Skill Unlocks/Choices:** Upon leveling up a class, the `Class.level_up_rewards` JSON is consulted.
    *   `unlock_skill`: The specified skill is automatically added to `PlayerSkill` table for the player with 0% percentage.
    *   `choose_skill`: The player is presented with a choice from a list of skills. The chosen skill is added to `PlayerSkill`.

### 3.3. Skill Percentage Progression

*   **Initial Percentage:** When a skill is acquired, its `percentage` starts at 0%.
*   **Experience/Usage:** Skill percentage increases through repeated use of the skill or by gaining experience in its associated class. The exact mechanism (e.g., direct usage XP, class XP conversion) needs further definition but should be tied to the `PlayerSkill.percentage` field.
*   **Class Level Cap:** The `PlayerSkill.percentage` for any given skill cannot exceed a cap determined by the player's `PlayerClass.level` in the `Skill.associated_class_id`.
    *   **Calculation:** `MaxSkillPercentage = (PlayerClass.level / Class.total_levels) * 100`.
    *   Example: If `Class.total_levels` is 5, and `PlayerClass.level` is 3, the `MaxSkillPercentage` for skills associated with that class is (3/5) * 100 = 60%.

## 4. Integration with Core Game Engine

*   **Skill Usage:** When a player attempts to use a skill, the Core Game Engine will:
    1.  Check if the player has the skill (`PlayerSkill` table).
    2.  Retrieve the skill's current `percentage` from `PlayerSkill`.
    3.  Evaluate the `Skill.effects` formulas using the current `percentage` to determine the actual outcome (damage, healing, duration, etc.).
    4.  Apply the effects to the game world.
*   **Class Leveling:** The `PlayerClass` table will be updated by the game engine when experience thresholds are met.
*   **DAL Integration:** All interactions with `Class`, `Skill`, `PlayerClass`, and `PlayerSkill` tables will be exclusively through the DAL.

## 5. Acceptance Criteria

1.  New `Class`, `Skill`, `PlayerClass`, and `PlayerSkill` data models are defined and persistable via the DAL.
2.  Skills have a structured `effects` field that can be evaluated by the game engine.
3.  Players can acquire classes and level them up, gaining experience relevant to the class.
4.  Upon class level-up, new skills are correctly unlocked or presented as choices to the player.
5.  Skill percentages can be increased through defined mechanisms.
6.  Skill percentages are correctly capped by the player's current level in the associated class.
7.  The game engine can correctly evaluate and apply skill effects based on the skill's percentage.
8.  Unit tests are in place for skill acquisition, leveling, percentage capping, and effect evaluation.

## 6. Test Data Requirements

To test the Skills and Classes system, the following data should be creatable via the web editor:

### 6.1. Example Classes

```json
[
  {
    "ID": "basic_fighter",
    "Name": "Basic Fighter",
    "Description": "A foundational combat class.",
    "TotalLevels": 5,
    "ParentClassID": "",
    "LevelUpRewards": {
      "1": {"unlock_skill": "basic_strike"},
      "2": {"unlock_skill": "defensive_stance"},
      "3": {"choose_skill": ["cleave", "shield_bash"]},
      "5": {"unlock_skill": "battle_fury"}
    }
  },
  {
    "ID": "sword_master",
    "Name": "Sword Master",
    "Description": "A specialized fighter focused on swords.",
    "TotalLevels": 3,
    "ParentClassID": "basic_fighter",
    "LevelUpRewards": {
      "1": {"unlock_skill": "precise_strike"},
      "2": {"choose_skill": ["parry", "riposte"]}
    }
  }
]
```

### 6.2. Example Skills

```json
[
  {
    "ID": "basic_strike",
    "Name": "Basic Strike",
    "Description": "A fundamental melee attack.",
    "Type": "active",
    "AssociatedClassID": "basic_fighter",
    "Effects": [
      {
        "type": "DAMAGE",
        "target": "ENEMY",
        "value_formula": "10 + (skill_percentage * 0.5)",
        "damage_type": "physical"
      }
    ],
    "Cost": 5,
    "Cooldown": 2,
    "MinClassLevel": 1
  },
  {
    "ID": "cleave",
    "Name": "Cleave",
    "Description": "A sweeping attack hitting multiple foes.",
    "Type": "active",
    "AssociatedClassID": "basic_fighter",
    "Effects": [
      {
        "type": "DAMAGE",
        "target": "AREA",
        "value_formula": "15 + (skill_percentage * 0.7)",
        "damage_type": "physical"
      }
    ],
    "Cost": 10,
    "Cooldown": 5,
    "MinClassLevel": 3
  },
  {
    "ID": "tracking_proficiency",
    "Name": "Tracking Proficiency",
    "Description": "Improves ability to follow tracks.",
    "Type": "passive",
    "AssociatedClassID": "ranger", // Assuming a 'ranger' class exists
    "Effects": [
      {
        "type": "MODIFY_ATTRIBUTE",
        "target": "SELF",
        "attribute": "perception",
        "value_formula": "skill_percentage * 0.1"
      }
    ],
    "MinClassLevel": 1
  }
]
```

### 6.3. Testing Scenarios

*   **Class Acquisition:** Test creating a new player and assigning a starting class.
*   **Class Leveling:** Simulate gaining experience for a class and verify level-ups occur correctly.
*   **Skill Unlocks:** Verify that `basic_strike` is unlocked at `basic_fighter` level 1.
*   **Skill Choices:** At `basic_fighter` level 3, verify the player is prompted to choose between `cleave` and `shield_bash`, and the chosen skill is added.
*   **Skill Percentage Cap:**
    *   Increase `basic_strike` percentage for a player with `basic_fighter` level 1 (max 20%). Verify it cannot exceed 20%.
    *   Level up `basic_fighter` to level 2 (max 40%). Verify `basic_strike` can now be increased up to 40%.
*   **Skill Effect Evaluation:** Use `basic_strike` at various percentages (e.g., 10%, 50%, 100%) and verify the damage output matches the `value_formula` calculation.
