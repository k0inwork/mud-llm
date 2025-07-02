# Quest Example: The Missing Pipe-Weed

## 1. Quest Overview

This document outlines a simple quest designed to demonstrate the relationships between Owners, Quest Owners, Questmakers, and how quests are initiated and completed within the game world. The quest involves a simple fetch task within the Shire.

## 2. Quest Narrative

Old Gaffer Gamgee has misplaced his favorite pipe-weed pouch somewhere near Bag End. He's too old and forgetful to search thoroughly himself and asks a passing adventurer for help.

## 3. Entities Involved & Their Roles

*   **Initiating Owner:** `shire_spirit` (monitors the Shire, can initiate local quests).
*   **Quest Owner:** `shire_local_governance` (thematic owner for local Shire affairs).
*   **Questmaker:** `missing_pipe_weed_questmaker` (controls the execution of this specific quest).
*   **NPC (Quest Giver/Trigger):** `gaffer_gamgee` (player interacts with him to start the quest).
*   **Item (Quest Objective):** `gaffer_pipe_weed_pouch` (the item to be found).
*   **Room (Quest Location):** `hobbiton_path` (where the item is found).

## 4. Proposed Seed Data for `quest_example.md`

### 4.1. Owner Update: `shire_spirit`

Add `"missing_pipe_weed_quest"` to its `InitiatedQuests` list.

```json
{
  "ID": "shire_spirit",
  "Name": "The Spirit of the Shire",
  "Description": "A gentle, ancient spirit embodying the peace and tranquility of the Shire. It watches over its hobbit inhabitants.",
  "MonitoredAspect": "location",
  "AssociatedID": "bag_end",
  "LLMPromptContext": "You are the benevolent spirit of the Shire, concerned with the well-being and simple lives of hobbits. You prefer peace and quiet.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 100.0,
  "MaxInfluenceBudget": 100.0,
  "BudgetRegenRate": 0.1,
  "AvailableTools": [],
  "InitiatedQuests": ["shire_census_quest", "missing_pipe_weed_quest"]
}
```

### 4.2. Quest Owner Update: `shire_local_governance`

Add `"missing_pipe_weed_questmaker"` to its `AssociatedQuestmakerIDs` list.

```json
{
  "ID": "shire_local_governance",
  "Name": "Shire Local Governance",
  "Description": "The day-to-day affairs and well-being of the Shire, managed by its various councils and respected elders.",
  "LLMPromptContext": "You are concerned with the peaceful and orderly functioning of the Shire. Your quests involve community tasks, local disputes, and maintaining the hobbit way of life.",
  "CurrentInfluenceBudget": 70.0,
  "MaxInfluenceBudget": 70.0,
  "BudgetRegenRate": 0.1,
  "AssociatedQuestmakerIDs": ["shire_census_questmaker", "missing_pipe_weed_questmaker"]
}
```

### 4.3. New Questmaker: `missing_pipe_weed_questmaker`

This Questmaker is specific to this quest.

```json
{
  "ID": "missing_pipe_weed_questmaker",
  "Name": "Missing Pipe-Weed Quest Controller",
  "LLMPromptContext": "You are the direct overseer of 'The Missing Pipe-Weed' quest. Your goal is to ensure Old Gaffer Gamgee's pipe-weed pouch is found and returned.",
  "CurrentInfluenceBudget": 0.0,
  "MaxInfluenceBudget": 25.0,
  "BudgetRegenRate": 0.0, // Player-action based
  "MemoriesAboutPlayers": {},
  "AvailableTools": []
}
```

### 4.4. New Quest: `missing_pipe_weed_quest`

```json
{
  "ID": "missing_pipe_weed_quest",
  "Name": "The Missing Pipe-Weed",
  "Description": "Old Gaffer Gamgee has lost his pipe-weed pouch. Find it and return it to him.",
  "QuestOwnerID": "shire_local_governance",
  "QuestmakerID": "missing_pipe_weed_questmaker",
  "InfluencePointsMap": {
    "found_pipe_weed": 15.0,
    "returned_pipe_weed": 10.0
  },
  "Objectives": [
    {"Type": "find_item", "TargetID": "gaffer_pipe_weed_pouch", "Status": "not_started"},
    {"Type": "return_item_to_npc", "TargetID": "gaffer_gamgee", "ItemToReturnID": "gaffer_pipe_weed_pouch", "Status": "not_started"}
  ],
  "Rewards": {
    "experience": 25,
    "gold": 5,
    "items": [{"item_id": "hobbit_pipe_weed_bundle", "quantity": 1}]
  }
}
```

### 4.5. New Item: `gaffer_pipe_weed_pouch`

This is the specific quest item.

```json
{
  "ID": "gaffer_pipe_weed_pouch",
  "Name": "Gaffer's Pipe-Weed Pouch",
  "Description": "A small, worn leather pouch, smelling faintly of sweet pipe-weed. It seems to have been dropped.",
  "Type": "quest_item",
  "Properties": {
    "is_quest_item": true,
    "current_room_id": "hobbiton_path" // Initial location for the item
  }
}
```

### 4.6. New Item: `hobbit_pipe_weed_bundle`

This is a reward item.

```json
{
  "ID": "hobbit_pipe_weed_bundle",
  "Name": "Bundle of Fine Pipe-Weed",
  "Description": "A small bundle of high-quality pipe-weed, a favorite among hobbits.",
  "Type": "consumable",
  "Properties": {
    "restores_stamina": 5,
    "flavor_text": "A truly comforting smoke."
  }
}
```

## 5. Quest Triggering (Editor's Perspective)

To trigger this quest, the editor would configure the `shire_spirit` Owner to initiate it via `gaffer_gamgee`.

*   **Owner:** `shire_spirit`
*   **Initiated Quest Entry:**
    ```json
    {
      "quest_id": "missing_pipe_weed_quest",
      "trigger_type": "on_talk_to_npc",
      "trigger_id": "gaffer_gamgee",
      "conditions": {} // No specific conditions for this simple quest
    }
    ```

## 6. Quest End Triggering (System's Perspective)

The quest ends when all objectives are marked as `"completed"` in the `PlayerQuestState`. The `Questmaker` (specifically `missing_pipe_weed_questmaker`) would monitor these objectives. Once the final objective (`return_item_to_npc` to `gaffer_gamgee`) is completed, the `Questmaker` would trigger the rewards and mark the quest as `"completed"`.

## 7. Caveats and Future Considerations

*   **Empty Caveats:**
    *   **`AvailableTools` for Questmakers:** For simplicity, the `AvailableTools` for `missing_pipe_weed_questmaker` are empty. In a real scenario, this Questmaker might have tools like `send_message` (to remind the player), `change_npc_behavior` (if Gaffer gets more distressed), or `spawn_entity` (if the pipe-weed moves).
    *   **`InfluencePointsMap` for Owners:** The `InfluencePointsMap` on the `Quest` currently only grants points to the `Questmaker`. We might consider allowing it to grant points to the `QuestOwner` or even the `Initiating Owner` for broader impact tracking.
    *   **Dynamic Quest Initiation Logic:** The `initiate_quest` tool on Owners is conceptual. The actual game logic for how an `on_talk_to_npc` trigger presents the quest to the player (e.g., dialogue options) is not detailed here.
    *   **Quest Failure:** This simple quest doesn't have explicit failure conditions (e.g., time limit, losing the item permanently). More complex quests would require these.
    *   **`AssociatedQuestmakerIDs` on Quest Owners:** This field is currently manually managed in the proposal. In a full editor, this could be automatically populated when a quest is assigned a `QuestOwnerID`.

*   **Simplicity Target:** This example prioritizes ease of understanding for the editor. Complex branching narratives, multiple ways to complete objectives, or dynamic objective generation are beyond this simple example but are supported by the underlying system design.

This document provides a clear blueprint for adding "The Missing Pipe-Weed" quest to the game, demonstrating the interaction between the various LLM-driven entities.