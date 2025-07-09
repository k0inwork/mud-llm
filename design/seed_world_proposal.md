# Seed World Proposal: Middle-earth Theme (Expanded)

## 1. Introduction

This document proposes an expanded set of seed data for a Middle-earth themed world, designed to populate the MUD with initial areas, characters, and quest-giving entities. This data will serve as a rich foundation for testing game mechanics, LLM interactions, and overall world immersion.

## 2. Core Principles

*   **Iconic Locations:** Focus on recognizable Middle-earth places.
*   **Character Archetypes:** Include key character types (hero, mentor, commoner, antagonist).
*   **Quest Hooks:** Design entities with inherent potential for quest generation.
*   **Owner Associations:** Link owners to specific areas or concepts for LLM monitoring.
*   **Logical Connectivity:** Ensure rooms form a coherent, navigable map.
*   **Quest Ownership & Control:** Clearly define entities responsible for quest initiation, thematic ownership, and execution control.

## 3. Proposed Seed Data

### 3.1. Rooms (Areas)

Representing distinct locations within Middle-earth.

#### 3.1.1. Bag End (The Shire)
```json
{
  "ID": "bag_end",
  "Name": "Bag End, Hobbiton",
  "Description": "A cozy hobbit-hole, warm and inviting, with a round green door. The smell of pipe-weed and fresh baking lingers in the air. A path leads east.",
  "Exits": { "east": { "Direction": "east", "TargetRoomID": "hobbiton_path" } },
  "OwnerID": "shire_spirit"
}
```

#### 3.1.2. Hobbiton Path (The Shire)
```json
{
  "ID": "hobbiton_path",
  "Name": "Hobbiton Path",
  "Description": "A well-worn path winding through green hills and past other hobbit-holes. The Bywater river glitters nearby.",
  "Exits": {
    "west":  { "Direction": "west", "TargetRoomID": "bag_end" },
    "east":  { "Direction": "east", "TargetRoomID": "bree_road" },
    "south": { "Direction": "south", "TargetRoomID": "green_dragon_inn" }
  },
  "OwnerID": "shire_spirit"
}
```

#### 3.1.3. The Green Dragon Inn (The Shire)
```json
{
  "ID": "green_dragon_inn",
  "Name": "The Green Dragon Inn",
  "Description": "A lively hobbit inn, filled with chatter and the clinking of mugs. A roaring fire warms the common room.",
  "Exits": { "north": { "Direction": "north", "TargetRoomID": "hobbiton_path" } },
  "OwnerID": "shire_spirit"
}
```

#### 3.1.4. Farmer Maggot's Field (The Shire)
```json
{
  "ID": "farmer_maggots_field",
  "Name": "Farmer Maggot's Field",
  "Description": "A large, well-tended field. A small, fenced-off patch in the corner seems to be growing some particularly large mushrooms.",
  "Exits": { "west": { "Direction": "west", "TargetRoomID": "hobbiton_path" } },
  "OwnerID": "shire_spirit"
}
```

#### 3.1.5. Bree Road (Outside Bree)
```json
{
  "ID": "bree_road",
  "Name": "Road to Bree",
  "Description": "A dusty road leading towards the walled town of Bree. Farmland stretches on either side.",
  "Exits": {
    "west": { "Direction": "west", "TargetRoomID": "hobbiton_path" },
    "east": { "Direction": "east", "TargetRoomID": "prancing_pony" }
  },
  "OwnerID": "bree_guardian"
}
```

... (rest of the world data remains the same)

### 3.5. Quests

#### 3.5.1. The Great Mushroom Hunt
```json
{
  "ID": "the_great_mushroom_hunt",
  "Name": "The Great Mushroom Hunt",
  "Description": "Farmer Maggot needs help gathering his prized mushrooms. Gather five of them from his field and bring them to him.",
  "QuestOwnerID": "shire_local_governance",
  "QuestmakerID": "mushroom_hunt_questmaker",
  "Objectives": [
    {"Type": "gather_item", "TargetID": "maggots_prize_mushrooms", "Count": 5, "From": "farmer_maggots_field"},
    {"Type": "deliver_item", "TargetID": "farmer_maggot", "ItemID": "maggots_prize_mushrooms"}
  ],
  "Rewards": {
    "experience": 50,
    "items": [{"item_id": "hobbit_pipe_weed_bundle", "quantity": 1}]
  }
}
```

... (other quests remain the same)