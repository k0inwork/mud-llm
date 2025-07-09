# Quest Example: The Great Mushroom Hunt

## 1. Quest Overview

This document outlines a simple quest designed to demonstrate the relationships between Owners, Quest Owners, Questmakers, and how quests are initiated and completed within the game world. The quest involves a simple fetch task within the Shire.

## 2. Quest Narrative

Farmer Maggot, a hobbit known for his prized mushrooms, needs help gathering his latest harvest. He's getting on in years and can't manage it all himself. He asks a passing adventurer for help.

## 3. Entities Involved & Their Roles

*   **Initiating Owner:** `shire_spirit` (monitors the Shire, can initiate local quests).
*   **Quest Owner:** `shire_local_governance` (thematic owner for local Shire affairs).
*   **Questmaker:** `mushroom_hunt_questmaker` (controls the execution of this specific quest).
*   **NPC (Quest Giver/Trigger):** `farmer_maggot` (player interacts with him to start the quest).
*   **Item (Quest Objective):** `maggots_prize_mushrooms` (the item to be gathered).
*   **Room (Quest Location):** `farmer_maggots_field` (where the mushrooms are gathered).

## 4. Proposed Seed Data for `quest_example.md`

### 4.1. Owner Update: `shire_spirit`

Add `"the_great_mushroom_hunt"` to its `InitiatedQuests` list.

```json
{
  "ID": "shire_spirit",
  "Name": "The Spirit of the Shire",
  "InitiatedQuests": ["the_great_mushroom_hunt"]
}
```

### 4.2. Quest Owner Update: `shire_local_governance`

Add `"mushroom_hunt_questmaker"` to its `AssociatedQuestmakerIDs` list.

```json
{
  "ID": "shire_local_governance",
  "AssociatedQuestmakerIDs": ["mushroom_hunt_questmaker"]
}
```

### 4.3. New Questmaker: `mushroom_hunt_questmaker`

This Questmaker is specific to this quest.

```json
{
  "ID": "mushroom_hunt_questmaker",
  "Name": "The Great Mushroom Hunt Controller",
  "LLMPromptContext": "I am the spirit of the harvest. My goal is to ensure Farmer Maggot's prized mushrooms are gathered safely."
}
```

### 4.4. New Quest: `the_great_mushroom_hunt`

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

### 4.5. New Item: `maggots_prize_mushrooms`

This is the specific quest item.

```json
{
  "ID": "maggots_prize_mushrooms",
  "Name": "Farmer Maggot's Prize Mushrooms",
  "Description": "A basket of unusually large and delicious-looking mushrooms.",
  "Type": "quest_item"
}
```

## 5. Quest Triggering (Editor's Perspective)

To trigger this quest, the editor would configure the `shire_spirit` Owner to initiate it via `farmer_maggot`.

*   **Owner:** `shire_spirit`
*   **Initiated Quest Entry:**
    ```json
    {
      "quest_id": "the_great_mushroom_hunt",
      "trigger_type": "on_talk_to_npc",
      "trigger_id": "farmer_maggot"
    }
    ```
