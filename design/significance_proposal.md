### **Proposal: Dynamic Action Significance Scoring System**

This proposal outlines a system for dynamically calculating the significance of player actions, moving beyond static scores to a context-aware model. The goal is to make entity reactions more believable, nuanced, and emergent, creating a richer player experience.

#### **1. Core Philosophy**

An action's importance is not inherent to the action itself, but is defined by **who is performing it, who (or what) it is being performed on, where it is happening, and what the relationships are between these elements.**

A whisper in a throne room is more significant than a shout in a tavern. A simple healing spell cast on a faction leader is more significant than a powerful fireball cast on a sewer rat. This system aims to model that reality.

#### **2. The Scoring Formula**

The final significance score will be calculated using a base score modified by a series of contextual multipliers.

`Final Score = BaseScore * TargetModifier * FactionModifier * LocationModifier * PlayerModifier`

-   **Base Score:** A default value for the action type (e.g., `say` = 5, `attack` = 20, `use_skill` = 15).
-   **Multipliers:** Values that increase or decrease the score based on context. A multiplier of `1.0` is neutral. A multiplier of `0.0` makes the action completely insignificant to that entity.

#### **3. Implementation: The `ActionContext`**

To calculate the score, the `ActionSignificanceMonitor` needs comprehensive information about the event. We will introduce a new struct, `ActionContext`, to be passed into the `LogPlayerAction` function.

```go
// A new struct to hold all relevant data for scoring
type ActionContext struct {
    Player          *models.Player
    ActionType      string      // e.g., "say", "attack", "use_skill"
    TargetEntity    interface{} // The direct target (NPC, Player, Item)
    SkillUsed       *models.Skill
    Room            *models.Room
    DALs            *dal.DALs // To fetch related info like owners, factions
}

// The modified function signature in the monitor
func (m *ActionSignificanceMonitor) LogPlayerAction(ctx ActionContext) {
    // ... dynamic scoring logic ...
}
```

The `TelnetServer` will be responsible for gathering and populating this `ActionContext` object for each command.

#### **4. Contextual Modifiers: The Dynamic Core**

Here is a breakdown of how different game elements will provide multipliers.

##### **A. Target-Based Modifiers (`TargetModifier`)**

This considers the state and nature of the direct target of the action.

| Condition                               | Example Scenario                               | Multiplier | Rationale                                                              |
| --------------------------------------- | ---------------------------------------------- | ---------- | ---------------------------------------------------------------------- |
| **Target is Sleeping/Unconscious**      | `say` to a sleeping guard                      | `0.1x`     | The target is unlikely to notice.                                      |
| **Target is in Combat**                 | `say` to an NPC fighting a dragon              | `0.2x`     | The target is distracted.                                              |
| **Target is a Faction Leader**          | `attack` the Thieves' Guild Master             | `3.0x`     | Actions against important figures have major consequences.             |
| **Target is an "Owner" or "Questmaker"**| `talk` to a Questmaker                         | `1.5x`     | These entities are inherently more attentive to direct interaction.    |
| **Target is an inanimate object**       | `attack` a training dummy                      | `0.0x`     | The dummy doesn't care.                                                |
| **Target is an object owned by someone**| `get` a sword from a shop display              | `5.0x`     | This is theft! The owner will be highly interested.                    |

##### **B. Faction & Relationship Modifiers (`FactionModifier`)**

This is calculated for *each potential observer*, not just the direct target.

| Player's Faction Status with Observer | Example Scenario                               | Multiplier | Rationale                                                              |
| ------------------------------------- | ---------------------------------------------- | ---------- | ---------------------------------------------------------------------- |
| **Ally**                              | An ally sees you `attack` their sworn enemy.   | `1.2x`     | They approve and pay close attention.                                  |
| **Ally**                              | An ally sees you `attack` another ally.        | `10.0x`    | Betrayal! This is extremely significant.                               |
| **Neutral**                           | A neutral guard sees you `say` something.      | `1.0x`     | Standard interaction.                                                  |
| **Enemy**                             | An enemy sees you `use_skill` "Stealth".       | `2.0x`     | They are suspicious of your every move.                                |
| **Hated**                             | A hated enemy sees you do anything at all.     | `2.5x`     | They are actively hostile and watching for any excuse to act.          |

##### **C. Location-Based Modifiers (`LocationModifier`)**

This applies to the **owner of the territory/room**.

| Location Type                           | Example Scenario                               | Multiplier | Rationale                                                              |
| --------------------------------------- | ---------------------------------------------- | ---------- | ---------------------------------------------------------------------- |
| **Private Property (Home, Shop)**       | Player `uses_skill` "Break Lock" on a chest.   | `8.0x`     | A direct, destructive, and illegal act on the owner's property.        |
| **Public Territory (Town Square)**      | Player `attacks` another player.               | `4.0x`     | The town guard (owner) is very interested in violence in their square. |
| **Sacred Ground (Temple, Shrine)**      | Player `uses_skill` "Corrupt Altar".           | `20.0x`    | An act of extreme desecration. The deity/owner will surely notice.     |
| **Wilderness**                          | Player `attacks` a monster.                    | `0.5x`     | The "owner" (e.g., nature spirit) is less concerned with daily squabbles. |

##### **D. Player-Based Modifiers (`PlayerModifier`)**

This relates to the player's own skills and background, creating personalized attention.

| Condition                               | Example Scenario                               | Observer                               | Multiplier | Rationale                                                              |
| --------------------------------------- | ---------------------------------------------- | -------------------------------------- | ---------- | ---------------------------------------------------------------------- |
| **Skill taught by an Owner**            | Player uses "Ancient Smithing" taught by the Dwarven King. | The Dwarven King                       | `1.5x`     | The mentor is interested in how their student uses their teachings.    |
| **Profession-related Skill**            | A "Lumberjack" uses "Fell Tree" skill.         | The Lumberjack Guildmaster             | `1.3x`     | The head of the profession keeps tabs on their members' activities.    |
| **Class-defining Skill**                | A "Paladin" uses "Lay on Hands".               | High-level Paladins, Clerics, Deities  | `1.2x`     | The character's actions reflect on their entire class.                 |

#### **5. Putting It All Together: A Scenario**

**Scenario:** A player, who is a "Rogue" and an ally of the Thieves' Guild, uses the skill "Disable Trap" on a chest inside "Morgan's Marvelous Metals" (a shop). The shop owner, Morgan, is neutral to the player. The city guard is also neutral.

**Action:** `use_skill` ("Disable Trap")
**Base Score:** 15

**1. Significance for Morgan (Shop Owner):**
-   `TargetModifier`: The chest is his property (`5.0x`).
-   `FactionModifier`: He is neutral (`1.0x`).
-   `LocationModifier`: It's his private shop (`8.0x`).
-   `PlayerModifier`: The skill is a standard Rogue skill (`1.0x`).
-   **Final Score for Morgan:** `15 * 5.0 * 1.0 * 8.0 * 1.0 = 600`. **EXTREMELY SIGNIFICANT.** Morgan will almost certainly react.

**2. Significance for the City Guard (Owner of the street outside):**
-   `TargetModifier`: Not their chest (`1.0x`).
-   `FactionModifier`: They are neutral (`1.0x`).
-   `LocationModifier`: The action is not in their direct territory (it's inside a shop), but they are nearby. We can apply a proximity-based "awareness" modifier (`0.3x`).
-   `PlayerModifier`: (`1.0x`).
-   **Final Score for Guard:** `15 * 1.0 * 1.0 * 0.3 * 1.0 = 4.5`. **Low significance.** They won't notice unless something else happens.

**3. Significance for the Thieves' Guild Master (Faction Leader):**
-   `TargetModifier`: Not their chest (`1.0x`).
-   `FactionModifier`: The player is an ally (`1.2x`).
-   `LocationModifier`: Not their territory (`1.0x`).
-   `PlayerModifier`: The player is using a core skill of their profession (`1.3x`).
-   **Final Score for Guild Master:** `15 * 1.0 * 1.2 * 1.0 * 1.3 = 23.4`. **Moderately significant.** The Guild Master might hear about their agent's activities later.

This system allows a single player action to ripple through the world and mean different things to different entities, creating a web of interconnected consequences.
