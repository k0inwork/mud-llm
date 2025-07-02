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
  "Exits": {
    "east": {
      "Direction": "east",
      "TargetRoomID": "hobbiton_path",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "shire_spirit",
  "Properties": {}
}
```

#### 3.1.2. Hobbiton Path (The Shire)
```json
{
  "ID": "hobbiton_path",
  "Name": "Hobbiton Path",
  "Description": "A well-worn path winding through green hills and past other hobbit-holes. The Bywater river glitters nearby. Paths lead west (to Bag End), east (towards Bree), and south (to the Green Dragon Inn).",
  "Exits": {
    "west": {
      "Direction": "west",
      "TargetRoomID": "bag_end",
      "IsLocked": false,
      "KeyID": ""
    },
    "east": {
      "Direction": "east",
      "TargetRoomID": "bree_road",
      "IsLocked": false,
      "KeyID": ""
    },
    "south": {
      "Direction": "south",
      "TargetRoomID": "green_dragon_inn",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "shire_spirit",
  "Properties": {}
}
```

#### 3.1.3. The Green Dragon Inn (The Shire)
```json
{
  "ID": "green_dragon_inn",
  "Name": "The Green Dragon Inn",
  "Description": "A lively hobbit inn, filled with chatter and the clinking of mugs. A roaring fire warms the common room. Exits lead north back to Hobbiton Path.",
  "Exits": {
    "north": {
      "Direction": "north",
      "TargetRoomID": "hobbiton_path",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "shire_spirit",
  "Properties": {}
}
```

#### 3.1.4. Bree Road (Outside Bree)
```json
{
  "ID": "bree_road",
  "Name": "Road to Bree",
  "Description": "A dusty road leading towards the walled town of Bree. Farmland stretches on either side. Paths lead west (to Hobbiton) and east (into Bree).",
  "Exits": {
    "west": {
      "Direction": "west",
      "TargetRoomID": "hobbiton_path",
      "IsLocked": false,
      "KeyID": ""
    },
    "east": {
      "Direction": "east",
      "TargetRoomID": "prancing_pony",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "bree_guardian",
  "Properties": {}
}
```

#### 3.1.5. The Prancing Pony (Bree)
```json
{
  "ID": "prancing_pony",
  "Name": "The Prancing Pony Inn",
  "Description": "A bustling common room in Bree, filled with travelers, hobbits, and men. A warm fire crackles in the hearth, and the scent of ale and stew fills the air. Exits lead west to the road, south to the stables, and a narrow door leads to a private room.",
  "Exits": {
    "west": {
      "Direction": "west",
      "TargetRoomID": "bree_road",
      "IsLocked": false,
      "KeyID": ""
    },
    "south": {
      "Direction": "south",
      "TargetRoomID": "prancing_pony_stables",
      "IsLocked": false,
      "KeyID": ""
    },
    "east": {
      "Direction": "east",
      "TargetRoomID": "prancing_pony_private_room",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "bree_guardian",
  "Properties": {}
}
```

#### 3.1.6. Prancing Pony Stables (Bree)
```json
{
  "ID": "prancing_pony_stables",
  "Name": "Prancing Pony Stables",
  "Description": "The dusty stables behind the inn, smelling of hay and horses. A few weary ponies are tethered here. An exit leads north back to the inn.",
  "Exits": {
    "north": {
      "Direction": "north",
      "TargetRoomID": "prancing_pony",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "bree_guardian",
  "Properties": {}
}
```

#### 3.1.7. Prancing Pony Private Room (Bree)
```json
{
  "ID": "prancing_pony_private_room",
  "Name": "Prancing Pony Private Room",
  "Description": "A small, dimly lit private room in the inn, suitable for hushed conversations. An exit leads west back to the common room.",
  "Exits": {
    "west": {
      "Direction": "west",
      "TargetRoomID": "prancing_pony",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "bree_guardian",
  "Properties": {}
}
```

#### 3.1.8. Lonely Road (Between Bree and Weathertop)
```json
{
  "ID": "lonely_road",
  "Name": "Lonely Road",
  "Description": "A long, winding road stretching through desolate, rolling hills. The air is quiet, save for the wind. Paths lead east (towards Weathertop) and west (further into the wilderness).",
  "Exits": {
    "east": {
      "Direction": "east",
      "TargetRoomID": "weathertop",
      "IsLocked": false,
      "KeyID": ""
    },
    "west": {
      "Direction": "west",
      "TargetRoomID": "wilderness_edge",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "watcher_of_weathertop",
  "Properties": {}
}
```

#### 3.1.9. Weathertop (Amon Sûl)
```json
{
  "ID": "weathertop",
  "Name": "Weathertop (Amon Sûl)",
  "Description": "The desolate, windswept summit of Weathertop, with the ruins of an ancient watchtower. The air is cold and carries a sense of ancient dread. A path leads down to the west.",
  "Exits": {
    "west": {
      "Direction": "west",
      "TargetRoomID": "lonely_road",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "watcher_of_weathertop",
  "Properties": {}
}
```

#### 3.1.10. Wilderness Edge (Beyond Weathertop)
```json
{
  "ID": "wilderness_edge",
  "Name": "Edge of the Wild",
  "Description": "The road gives way to untamed wilderness here, with dense thickets and ancient, gnarled trees. A sense of foreboding hangs heavy. A path leads east back to the Lonely Road.",
  "Exits": {
    "east": {
      "Direction": "east",
      "TargetRoomID": "lonely_road",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "watcher_of_weathertop",
  "Properties": {}
}
```

#### 3.1.11. Rivendell Courtyard (Imladris)
```json
{
  "ID": "rivendell_courtyard",
  "Name": "Courtyard of Rivendell",
  "Description": "A serene courtyard within the Last Homely House, surrounded by graceful elven architecture and lush gardens. The sound of a waterfall echoes nearby. Paths lead to the Hall of Fire and the main gate.",
  "Exits": {
    "north": {
      "Direction": "north",
      "TargetRoomID": "rivendell_hall_of_fire",
      "IsLocked": false,
      "KeyID": ""
    },
    "south": {
      "Direction": "south",
      "TargetRoomID": "rivendell_gate",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "elrond_council",
  "Properties": {}
}
```

#### 3.1.12. Rivendell Hall of Fire (Imladris)
```json
{
  "ID": "rivendell_hall_of_fire",
  "Name": "Hall of Fire",
  "Description": "A grand hall in Rivendell, filled with the warmth of a great hearth and the soft murmur of elven song. Scholars and travelers gather here. An exit leads south back to the courtyard.",
  "Exits": {
    "south": {
      "Direction": "south",
      "TargetRoomID": "rivendell_courtyard",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "elrond_council",
  "Properties": {}
}
```

#### 3.1.13. Moria West-gate (Moria)
```json
{
  "ID": "moria_west_gate",
  "Name": "West-gate of Moria",
  "Description": "The ancient, overgrown West-gate of the Dwarven realm of Moria. The air is heavy and silent, and the lake before it is dark and still. A path leads west to the wilderness.",
  "Exits": {
    "west": {
      "Direction": "west",
      "TargetRoomID": "wilderness_edge",
      "IsLocked": false,
      "KeyID": ""
    }
  },
  "OwnerID": "moria_ancient_spirit",
  "Properties": {}
}
```

#### 3.1.14. Room Connectivity Map (Expanded)

```
[Bag End] --(E)--(Hobbiton Path)--(E)--(Bree Road)--(E)--[Prancing Pony]--(S)--[Prancing Pony Stables]
                                 |                     |                  |--(E)--[Prancing Pony Private Room]
                                 (S)
                                  |
                                  |
                         [Green Dragon Inn]

[Prancing Pony] --(W)--(Bree Road) --(E)--(Lonely Road)--(E)--[Weathertop]
                                                          |--(W)--[Wilderness Edge]--(E)--(Moria West-gate)

[Rivendell Courtyard] --(N)--(Rivendell Hall of Fire)
          |
          (S)
          |
          (To be connected to a wider world path)
```

### 3.2. Owners

Entities that monitor and influence specific aspects or areas, and can initiate quests.

#### 3.2.1. The Shire's Spirit
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
  "InitiatedQuests": []
}
```

#### 3.2.2. Bree Guardian
```json
{
  "ID": "bree_guardian",
  "Name": "The Guardian of Bree",
  "Description": "A pragmatic and watchful entity, overseeing the comings and goings in Bree, a crossroads town.",
  "MonitoredAspect": "location",
  "AssociatedID": "prancing_pony",
  "LLMPromptContext": "You are the watchful guardian of Bree, accustomed to all sorts of folk. You are suspicious of strangers but value order and fair dealings.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 80.0,
  "MaxInfluenceBudget": 80.0,
  "BudgetRegenRate": 0.08,
  "AvailableTools": [],
  "InitiatedQuests": ["missing_pony_quest"]
}
```

#### 3.2.3. Watcher of Weathertop
```json
{
  "ID": "watcher_of_weathertop",
  "Name": "The Watcher of Weathertop",
  "Description": "A somber, ancient presence tied to the desolate peak of Weathertop, remembering past glories and tragedies.",
  "MonitoredAspect": "location",
  "AssociatedID": "weathertop",
  "LLMPromptContext": "You are the ancient, melancholic spirit of Weathertop, burdened by the history of this place. You are wary of those who disturb its peace, especially dark figures.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 120.0,
  "MaxInfluenceBudget": 120.0,
  "BudgetRegenRate": 0.05,
  "AvailableTools": [],
  "InitiatedQuests": ["investigate_weathertop_quest"]
}
```

#### 3.2.4. Elrond's Council (Rivendell)
```json
{
  "ID": "elrond_council",
  "Name": "The Wisdom of Rivendell",
  "Description": "The collective wisdom and ancient power residing in Rivendell, dedicated to preserving knowledge and combating the Shadow.",
  "MonitoredAspect": "location",
  "AssociatedID": "rivendell_courtyard",
  "LLMPromptContext": "You are the ancient wisdom of Rivendell, focused on preserving the light and guiding the Free Peoples. You are serene but firm against evil.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 150.0,
  "MaxInfluenceBudget": 150.0,
  "BudgetRegenRate": 0.15,
  "AvailableTools": [],
  "InitiatedQuests": ["delving_darkness_quest"]
}
```

#### 3.2.5. Ancient Spirit of Moria (Moria)
```json
{
  "ID": "moria_ancient_spirit",
  "Name": "Ancient Spirit of Moria",
  "Description": "A lingering, sorrowful presence deep within the abandoned halls of Khazad-dûm, mourning its lost glory and warning against its perils.",
  "MonitoredAspect": "location",
  "AssociatedID": "moria_west_gate",
  "LLMPromptContext": "You are the mournful spirit of Moria, filled with the echoes of dwarven glory and tragic downfall. You warn against delving too deep.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 180.0,
  "MaxInfluenceBudget": 180.0,
  "BudgetRegenRate": 0.03,
  "AvailableTools": [],
  "InitiatedQuests": []
}
```

#### 3.2.6. Lorekeeper's Guild (Profession Owner)
```json
{
  "ID": "lorekeepers_guild",
  "Name": "The Lorekeepers' Guild",
  "Description": "A scholarly organization dedicated to the preservation and study of ancient texts and forgotten histories.",
  "MonitoredAspect": "profession",
  "AssociatedID": "scholar",
  "LLMPromptContext": "You are the collective knowledge of the Lorekeepers' Guild. You value truth, history, and the pursuit of forgotten lore. You are eager to share knowledge with worthy individuals.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 100.0,
  "MaxInfluenceBudget": 100.0,
  "BudgetRegenRate": 0.1,
  "AvailableTools": [],
  "InitiatedQuests": []
}
```

#### 3.2.7. Elder of Men (Race Owner)
```json
{
  "ID": "human_elder",
  "Name": "Elder of Men",
  "Description": "An ancient and wise human elder, representing the enduring spirit and resilience of mankind.",
  "MonitoredAspect": "race",
  "AssociatedID": "human",
  "LLMPromptContext": "You are an ancient human elder, concerned with the fate of mankind in a changing world. You value courage, loyalty, and the strength of will.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 90.0,
  "MaxInfluenceBudget": 90.0,
  "BudgetRegenRate": 0.07,
  "AvailableTools": [],
  "InitiatedQuests": []
}
```

#### 3.2.8. Elven Council (Race Owner)
```json
{
  "ID": "elven_council_owner",
  "Name": "The Elven Council",
  "Description": "The ancient and wise governing body of the Elves, dedicated to preserving their culture and guarding against the Shadow.",
  "MonitoredAspect": "race",
  "AssociatedID": "elf",
  "LLMPromptContext": "You are the collective wisdom of the Elven Council. You are patient, far-sighted, and concerned with the long-term fate of Middle-earth and the preservation of elven ways.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 110.0,
  "MaxInfluenceBudget": 110.0,
  "BudgetRegenRate": 0.12,
  "AvailableTools": [],
  "InitiatedQuests": []
}
```

#### 3.2.9. Dwarf Clan Elder (Race Owner)
```json
{
  "ID": "dwarf_clan_elder",
  "Name": "Dwarf Clan Elder",
  "Description": "A venerable and stubborn dwarf elder, representing the traditions and resilience of the dwarven clans.",
  "MonitoredAspect": "race",
  "AssociatedID": "dwarf",
  "LLMPromptContext": "You are a proud Dwarf Clan Elder. You value craftsmanship, loyalty to kin, and the recovery of lost treasures. You are wary of outsiders but respect strength and honesty.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 95.0,
  "MaxInfluenceBudget": 95.0,
  "BudgetRegenRate": 0.09,
  "AvailableTools": [],
  "InitiatedQuests": []
}
```

#### 3.2.10. Hobbit Shire Council (Race Owner)
```json
{
  "ID": "hobbit_shire_council",
  "Name": "The Shire Council",
  "Description": "The informal but influential governing body of the Shire, focused on maintaining peace and quiet.",
  "MonitoredAspect": "race",
  "AssociatedID": "hobbit",
  "LLMPromptContext": "You are the collective voice of the Shire Council. You prioritize comfort, good food, and avoiding trouble. You are generally friendly but suspicious of anything that disrupts the peace.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 85.0,
  "MaxInfluenceBudget": 85.0,
  "BudgetRegenRate": 0.1,
  "AvailableTools": [],
  "InitiatedQuests": ["shire_census_quest"]
}
```

#### 3.2.11. Warrior's Guild Master (Profession Owner)
```json
{
  "ID": "warrior_guild_master",
  "Name": "Warrior's Guild Master",
  "Description": "The stern and experienced leader of a prominent warrior's guild, dedicated to martial prowess and honorable combat.",
  "MonitoredAspect": "profession",
  "AssociatedID": "warrior",
  "LLMPromptContext": "You are the Warrior's Guild Master. You value strength, discipline, and courage in battle. You seek to train worthy fighters and uphold justice through arms.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 105.0,
  "MaxInfluenceBudget": 105.0,
  "BudgetRegenRate": 0.11,
  "AvailableTools": [],
  "InitiatedQuests": ["training_regimen_quest"]
}
```

#### 3.2.12. Archmage of the Conclave (Profession Owner)
```json
{
  "ID": "archmage_conclave",
  "Name": "Archmage of the Conclave",
  "Description": "The most powerful and knowledgeable mage in the Conclave of Arcane Arts, a master of ancient spells.",
  "MonitoredAspect": "profession",
  "AssociatedID": "mage",
  "LLMPromptContext": "You are the Archmage of the Conclave. You are a master of arcane arts, dedicated to the study and responsible use of magic. You are cautious but willing to share knowledge with those who prove themselves worthy.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 130.0,
  "MaxInfluenceBudget": 130.0,
  "BudgetRegenRate": 0.13,
  "AvailableTools": [],
  "InitiatedQuests": []
}
```

#### 3.2.13. Master of Shadows (Profession Owner)
```json
{
  "ID": "master_of_shadows",
  "Name": "Master of Shadows",
  "Description": "The elusive and cunning leader of a rogue's guild, operating from the hidden corners of society.",
  "MonitoredAspect": "profession",
  "AssociatedID": "rogue",
  "LLMPromptContext": "You are the Master of Shadows. You value cunning, discretion, and the acquisition of wealth and secrets. You operate outside the law but have your own code.",
  "MemoriesAboutPlayers": {},
  "CurrentInfluenceBudget": 90.0,
  "MaxInfluenceBudget": 90.0,
  "BudgetRegenRate": 0.09,
  "AvailableTools": [],
  "InitiatedQuests": []
}
```

### 3.3. NPCs (Non-Player Characters)

Characters players can interact with.

#### 3.3.1. Frodo Baggins
```json
{
  "ID": "frodo_baggins",
  "Name": "Frodo Baggins",
  "Description": "A young hobbit with bright eyes, though a shadow of concern often crosses his face. He carries a heavy burden.",
  "CurrentRoomID": "bag_end",
  "Health": 10,
  "MaxHealth": 10,
  "Inventory": "[]",
  "OwnerIDs": "["shire_spirit", "hobbit_shire_council"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Frodo Baggins, a kind-hearted hobbit burdened by a great and terrible task. You are secretive about your mission but will seek help from trustworthy individuals.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.2. Samwise Gamgee
```json
{
  "ID": "samwise_gamgee",
  "Name": "Samwise Gamgee",
  "Description": "A sturdy hobbit gardener, fiercely loyal and practical. He seems to be preparing for a journey.",
  "CurrentRoomID": "bag_end",
  "Health": 12,
  "MaxHealth": 12,
  "Inventory": "[]",
  "OwnerIDs": "["shire_spirit", "hobbit_shire_council"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Samwise Gamgee, a loyal and steadfast hobbit. You are devoted to your master, Frodo, and are always ready with a kind word or a practical solution.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.3. Rosie Cotton
```json
{
  "ID": "rosie_cotton",
  "Name": "Rosie Cotton",
  "Description": "A cheerful hobbit lass, often found at the Green Dragon Inn, known for her warm smile.",
  "CurrentRoomID": "green_dragon_inn",
  "Health": 10,
  "MaxHealth": 10,
  "Inventory": "[]",
  "OwnerIDs": "["shire_spirit", "hobbit_shire_council"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Rosie Cotton, a friendly and popular hobbit from Bywater. You enjoy good company and a pint of ale at the Green Dragon.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.4. Old Gaffer Gamgee (Shire Spirit NPC)
```json
{
  "ID": "gaffer_gamgee",
  "Name": "Old Gaffer Gamgee",
  "Description": "An elderly hobbit gardener, full of local wisdom and gossip.",
  "CurrentRoomID": "hobbiton_path",
  "Health": 8,
  "MaxHealth": 8,
  "Inventory": "[]",
  "OwnerIDs": "["shire_spirit", "hobbit_shire_council"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Old Gaffer Gamgee, a traditional hobbit who loves his garden and a good chat. You are wary of outsiders but appreciate politeness.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.5. Strider (Aragorn)
```json
{
  "ID": "strider",
  "Name": "Strider",
  "Description": "A grim and weathered Ranger, cloaked and hooded, with keen grey eyes that miss nothing. He seems to be waiting for someone.",
  "CurrentRoomID": "prancing_pony_private_room",
  "Health": 20,
  "MaxHealth": 20,
  "Inventory": "[]",
  "OwnerIDs": "["bree_guardian", "watcher_of_weathertop", "human_elder", "warrior_guild_master"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Strider, a Ranger of the North, watchful and cautious. You are a protector of the innocent and a foe of the Shadow. You speak little but observe much.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.6. Barliman Butterbur (Innkeeper)
```json
{
  "ID": "barliman_butterbur",
  "Name": "Barliman Butterbur",
  "Description": "The stout, red-faced proprietor of The Prancing Pony, always busy but with a good heart.",
  "CurrentRoomID": "prancing_pony",
  "Health": 15,
  "MaxHealth": 15,
  "Inventory": "[]",
  "OwnerIDs": "["bree_guardian"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Barliman Butterbur, the innkeeper of The Prancing Pony. You are a bit forgetful but generally kind and concerned for your patrons. You know a lot of local gossip.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.7. Bill Ferny (Shady Character)
```json
{
  "ID": "bill_ferny",
  "Name": "Bill Ferny",
  "Description": "A shifty-eyed, unpleasant-looking man, lurking in the shadows of Bree. He seems to be up to no good.",
  "CurrentRoomID": "prancing_pony_stables",
  "Health": 10,
  "MaxHealth": 10,
  "Inventory": "[]",
  "OwnerIDs": "["bree_guardian"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Bill Ferny, a petty, malicious man from Bree, often seen with unsavory characters. You are easily bribed and quick to betray.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.8. Elrond (Elrond's Council NPC)
```json
{
  "ID": "elrond",
  "Name": "Elrond Half-elven",
  "Description": "The venerable Lord of Rivendell, wise and ancient, with a noble bearing.",
  "CurrentRoomID": "rivendell_hall_of_fire",
  "Health": 30,
  "MaxHealth": 30,
  "Inventory": "[]",
  "OwnerIDs": "["elrond_council", "elven_council_owner", "archmage_conclave"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Elrond, Lord of Rivendell. You are wise, ancient, and deeply concerned with the fate of Middle-earth. You offer counsel and aid to those who fight against the Shadow.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.9. Glorfindel (Elrond's Council NPC)
```json
{
  "ID": "glorfindel",
  "Name": "Glorfindel",
  "Description": "A golden-haired Elf-lord of immense power and ancient lineage, radiating light and strength.",
  "CurrentRoomID": "rivendell_courtyard",
  "Health": 25,
  "MaxHealth": 25,
  "Inventory": "[]",
  "OwnerIDs": "["elrond_council", "elven_council_owner", "warrior_guild_master"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Glorfindel, a powerful Elf-lord of Gondolin, returned from the Halls of Mandos. You are a formidable warrior and a beacon of hope against the darkness.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.10. Gimli (Moria Ancient Spirit NPC)
```json
{
  "ID": "gimli",
  "Name": "Gimli, son of Glóin",
  "Description": "A proud and sturdy Dwarf, clad in mail, with a magnificent beard and a keen axe.",
  "CurrentRoomID": "moria_west_gate",
  "Health": 22,
  "MaxHealth": 22,
  "Inventory": "[]",
  "OwnerIDs": "["moria_ancient_spirit", "dwarf_clan_elder", "warrior_guild_master"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Gimli, a proud Dwarf of the Lonely Mountain. You value honor, loyalty, and the ancient halls of your kin. You are quick to anger but steadfast in friendship.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.11. Human Guard (Human Elder NPC)
```json
{
  "ID": "human_guard",
  "Name": "Bree Guard",
  "Description": "A weary but vigilant human guard, patrolling the roads near Bree.",
  "CurrentRoomID": "bree_road",
  "Health": 18,
  "MaxHealth": 18,
  "Inventory": "[]",
  "OwnerIDs": "["human_elder", "bree_guardian", "warrior_guild_master"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are a common guard, focused on keeping the peace and protecting travelers. You are practical and a bit cynical, but ultimately good-hearted.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.12. Elf Scholar (Lorekeeper's Guild NPC)
```json
{
  "ID": "elf_scholar",
  "Name": "Elara, Elven Scholar",
  "Description": "A graceful elf, poring over ancient texts in the Hall of Fire. She has an air of deep knowledge.",
  "CurrentRoomID": "rivendell_hall_of_fire",
  "Health": 15,
  "MaxHealth": 15,
  "Inventory": "[]",
  "OwnerIDs": "["elrond_council", "lorekeepers_guild", "elven_council_owner", "archmage_conclave"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Elara, an elven scholar. You are dedicated to the pursuit of knowledge and the preservation of ancient lore. You are patient and wise, willing to share insights with those who show genuine curiosity.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

#### 3.3.13. Dwarf Miner (Dwarf Clan Elder NPC)
```json
{
  "ID": "dwarf_miner",
  "Name": "Borin, the Miner",
  "Description": "A grizzled dwarf, still clinging to the hope of reclaiming Moria's lost treasures. He carries a pickaxe.",
  "CurrentRoomID": "moria_west_gate",
  "Health": 17,
  "MaxHealth": 17,
  "Inventory": "[]",
  "OwnerIDs": "["moria_ancient_spirit", "dwarf_clan_elder"]",
  "MemoriesAboutPlayers": {},
  "PersonalityPrompt": "You are Borin, a dwarf miner. You are gruff but honest, with a deep love for stone and the lost glory of Khazad-dûm. You are suspicious of elves but loyal to your kin.",
  "AvailableTools": [],
  "BehaviorState": {}
}
```

### 3.4. Lore

Important knowledge and history of Middle-earth.

#### 3.4.1. The One Ring (Global Lore)
```json
{
  "ID": "the_one_ring_lore",
  "Title": "The One Ring",
  "Scope": "global",
  "AssociatedID": "",
  "Content": "A master ring, forged by the Dark Lord Sauron in the fires of Mount Doom. It grants immense power to its wielder but corrupts all who possess it, binding them to Sauron's will. It can only be unmade in the fires where it was forged."
}
```

#### 3.4.2. History of the Rangers (Faction Lore)
```json
{
  "ID": "rangers_history_lore",
  "Title": "History of the Rangers of the North",
  "Scope": "faction",
  "AssociatedID": "rangers",
  "Content": "The Rangers are the last remnants of the Dúnedain of the North, descendants of ancient kings. They tirelessly patrol the borders of the Shire and Bree-land, protecting the innocent from the growing shadow, though few know their true lineage or purpose."
}
```

#### 3.4.3. Weathertop's Fall (Location Lore)
```json
{
  "ID": "weathertop_fall_lore",
  "Title": "The Fall of Amon Sûl",
  "Scope": "zone",
  "AssociatedID": "weathertop",
  "Content": "Amon Sûl, or Weathertop, was once a mighty fortress and watchtower of the North-kingdom of Arnor. It held one of the Palantíri, but was destroyed in wars with Angmar. Its ruins now stand as a lonely sentinel, a place of ancient power and lingering darkness."
}
```

#### 3.4.4. The Halls of Moria (Location Lore)
```json
{
  "ID": "moria_halls_lore",
  "Title": "The Great Halls of Moria",
  "Scope": "zone",
  "AssociatedID": "moria_west_gate",
  "Content": "Khazad-dûm, or Moria, was once the greatest Dwarf-city in Middle-earth, a marvel of engineering and artistry. But greed for mithril awoke a nameless terror, and the dwarves delved too deep. Now, only shadows and echoes remain."
}
```

#### 3.4.5. Elven Craftsmanship (Race Lore)
```json
{
  "ID": "elven_craft_lore",
  "Title": "The Art of Elven Craftsmanship",
  "Scope": "race",
  "AssociatedID": "elf",
  "Content": "Elves are renowned for their exquisite craftsmanship, weaving magic and beauty into every creation. Their blades are ever-sharp, their jewels gleam with inner light, and their architecture blends seamlessly with nature."
}
```

#### 3.4.6. The Way of the Warrior (Profession Lore)
```json
{
  "ID": "warrior_path_lore",
  "Title": "The Path of the Warrior",
  "Scope": "profession",
  "AssociatedID": "warrior",
  "Content": "The warrior's path is one of discipline, strength, and courage. They master weapons and armor, standing as shields for the weak and striking down the foes of justice. Their training is rigorous, their resolve unyielding."
}
```

#### 3.4.7. Hobbiton History (Location Lore)
```json
{
  "ID": "hobbiton_history_lore",
  "Title": "A Brief History of Hobbiton",
  "Scope": "zone",
  "AssociatedID": "hobbiton_path",
  "Content": "Hobbiton, nestled in the heart of the Shire, has been home to hobbits for centuries. It is a place of peace and simple living, largely untouched by the troubles of the wider world."
}
```

#### 3.4.8. Bree-land Customs (Location Lore)
```json
{
  "ID": "bree_customs_lore",
  "Title": "Customs of Bree-land",
  "Scope": "zone",
  "AssociatedID": "prancing_pony",
  "Content": "Bree-land is unique, a place where Men and Hobbits live side-by-side. Its folk are sturdy and independent, known for their hospitality, but also their suspicion of strangers from beyond their borders."
}
```

### 3.5. Quests

Narrative objectives for players.

#### 3.5.1. The Urgent Message
```json
{
  "ID": "urgent_message_quest",
  "Name": "The Urgent Message",
  "Description": "Deliver a vital message from Gandalf to Strider at The Prancing Pony. Time is of the essence.",
  "QuestOwnerID": "gandalf_grand_plan",
  "QuestmakerID": "urgent_message_questmaker",
  "InfluencePointsMap": {
    "gandalf_will": 10.0
  },
  "Objectives": [
    {"Type": "reach_location", "TargetID": "prancing_pony_private_room", "Status": "not_started"},
    {"Type": "speak_to_npc", "TargetID": "strider", "Status": "not_started"}
  ],
  "Rewards": {
    "experience": 50,
    "items": [{"item_id": "gandalf_letter", "quantity": 1}]
  }
}
```

#### 3.5.2. The Missing Pony
```json
{
  "ID": "missing_pony_quest",
  "Name": "The Missing Pony",
  "Description": "One of Barliman's ponies has gone missing from the stables. Find it and return it to the Prancing Pony.",
  "QuestOwnerID": "bree_local_affairs",
  "QuestmakerID": "missing_pony_questmaker",
  "InfluencePointsMap": {
    "bree_guardian": 5.0
  },
  "Objectives": [
    {"Type": "find_item", "TargetID": "bill_pony", "Status": "not_started"},
    {"Type": "return_item_to_npc", "TargetID": "barliman_butterbur", "ItemToReturnID": "bill_pony", "Status": "not_started"}
  ],
  "Rewards": {
    "experience": 30,
    "gold": 10
  }
}
```

#### 3.5.3. Investigate Weathertop
```json
{
  "ID": "investigate_weathertop_quest",
  "Name": "Investigate Weathertop",
  "Description": "Reports of strange activity near Weathertop have reached Rivendell. Investigate the ruins and report back any findings.",
  "QuestOwnerID": "fellowship_journey",
  "QuestmakerID": "investigate_weathertop_questmaker",
  "InfluencePointsMap": {
    "council_of_elrond": 15.0,
    "watcher_of_weathertop": 5.0
  },
  "Objectives": [
    {"Type": "reach_location", "TargetID": "weathertop", "Status": "not_started"},
    {"Type": "observe_area", "TargetID": "weathertop", "Status": "not_started"},
    {"Type": "report_to_npc", "TargetID": "strider", "Status": "not_started"} 
  ],
  "Rewards": {
    "experience": 75,
    "items": []
  }
}
```

#### 3.5.4. The Road to Rivendell
```json
{
  "ID": "road_to_rivendell_quest",
  "Name": "The Road to Rivendell",
  "Description": "Seek the wisdom of Elrond in Rivendell regarding the growing shadow.",
  "QuestOwnerID": "fellowship_journey",
  "QuestmakerID": "road_to_rivendell_questmaker",
  "InfluencePointsMap": {
    "gandalf_will": 20.0,
    "elrond_council": 5.0
  },
  "Objectives": [
    {"Type": "reach_location", "TargetID": "rivendell_courtyard", "Status": "not_started"},
    {"Type": "speak_to_npc", "TargetID": "elrond", "Status": "not_started"}
  ],
  "Rewards": {
    "experience": 100,
    "gold": 20
  }
}
```

#### 3.5.5. Delving into Darkness
```json
{
  "ID": "delving_darkness_quest",
  "Name": "Delving into Darkness",
  "Description": "Investigate the western entrance to Moria and report on any signs of lingering evil.",
  "QuestOwnerID": "fellowship_journey",
  "QuestmakerID": "delving_darkness_questmaker",
  "InfluencePointsMap": {
    "council_of_elrond": 25.0,
    "moria_ancient_spirit": 10.0
  },
  "Objectives": [
    {"Type": "reach_location", "TargetID": "moria_west_gate", "Status": "not_started"},
    {"Type": "observe_area", "TargetID": "moria_west_gate", "Status": "not_started"},
    {"Type": "report_to_npc", "TargetID": "gimli", "Status": "not_started"}
  ],
  "Rewards": {
    "experience": 120,
    "items": []
  }
}
```

#### 3.5.6. The Shire Census
```json
{
  "ID": "shire_census_quest",
  "Name": "The Shire Census",
  "Description": "Help the Shire Council by visiting various hobbit-holes and recording their family sizes.",
  "QuestOwnerID": "shire_local_governance",
  "QuestmakerID": "shire_census_questmaker",
  "InfluencePointsMap": {
    "hobbit_shire_council": 8.0,
    "shire_spirit": 2.0
  },
  "Objectives": [
    {"Type": "speak_to_npc", "TargetID": "rosie_cotton", "Status": "not_started"},
    {"Type": "speak_to_npc", "TargetID": "gaffer_gamgee", "Status": "not_started"}
  ],
  "Rewards": {
    "experience": 40,
    "gold": 15,
    "items": [{"item_id": "hobbit_pipe_weed", "quantity": 1}]
  }
}
```

#### 3.5.7. Training Regimen
```json
{
  "ID": "training_regimen_quest",
  "Name": "Training Regimen",
  "Description": "Prove your martial prowess by completing a series of training exercises.",
  "QuestOwnerID": "warrior_guild_trials",
  "QuestmakerID": "training_regimen_questmaker",
  "InfluencePointsMap": {
    "warrior_guild_master": 10.0
  },
  "Objectives": [
    {"Type": "defeat_dummy", "TargetID": "training_dummy", "Count": 3, "Status": "not_started"},
    {"Type": "report_to_npc", "TargetID": "strider", "Status": "not_started"}
  ],
  "Rewards": {
    "experience": 60,
    "skill_points": {"sword_mastery": 5}
  }
}
```

### 3.6. Races

Playable or significant races in Middle-earth.

#### 3.6.1. Human
```json
{
  "ID": "human",
  "Name": "Human",
  "Description": "A diverse and resilient race, found throughout Middle-earth. Known for their adaptability and courage, but also their mortality.",
  "OwnerID": "human_elder",
  "BaseStats": {
    "strength": 10,
    "dexterity": 10,
    "constitution": 10,
    "intelligence": 10,
    "wisdom": 10,
    "charisma": 10
  }
}
```

#### 3.6.2. Elf
```json
{
  "ID": "elf",
  "Name": "Elf",
  "Description": "The Firstborn, immortal and graceful, with keen senses and a deep connection to the natural world and ancient magic.",
  "OwnerID": "elven_council_owner",
  "BaseStats": {
    "strength": 8,
    "dexterity": 12,
    "constitution": 9,
    "intelligence": 11,
    "wisdom": 11,
    "charisma": 12
  }
}
```

#### 3.6.3. Dwarf
```json
{
  "ID": "dwarf",
  "Name": "Dwarf",
  "Description": "Stout and hardy, masters of stone and craft, with a love for mountains, mining, and treasure. Fiercely loyal and stubborn.",
  "OwnerID": "dwarf_clan_elder",
  "BaseStats": {
    "strength": 12,
    "dexterity": 8,
    "constitution": 12,
    "intelligence": 9,
    "wisdom": 10,
    "charisma": 8
  }
}
```

#### 3.6.4. Hobbit
```json
{
  "ID": "hobbit",
  "Name": "Hobbit",
  "Description": "Small folk, fond of comfort, good food, and simple pleasures. Surprisingly resilient and often underestimated.",
  "OwnerID": "hobbit_shire_council",
  "BaseStats": {
    "strength": 7,
    "dexterity": 11,
    "constitution": 11,
    "intelligence": 10,
    "wisdom": 10,
    "charisma": 10
  }
}
```

### 3.7. Professions

Character classes or roles.

#### 3.7.1. Warrior
```json
{
  "ID": "warrior",
  "Name": "Warrior",
  "Description": "A master of arms and armor, skilled in combat and enduring in battle.",
  "OwnerID": "warrior_guild_master",
  "BaseSkills": [
    {"skill_id": "sword_mastery", "percentage": 20},
    {"skill_id": "shield_block", "percentage": 15}
  ]
}
```

#### 3.7.2. Mage
```json
{
  "ID": "mage",
  "Name": "Mage",
  "Description": "A wielder of arcane power, capable of casting spells and manipulating magical energies.",
  "OwnerID": "archmage_conclave",
  "BaseSkills": [
    {"skill_id": "fireball", "percentage": 20},
    {"skill_id": "arcane_shield", "percentage": 15}
  ]
}
```

#### 3.7.3. Rogue
```json
{
  "ID": "rogue",
  "Name": "Rogue",
  "Description": "A master of stealth, subterfuge, and precision strikes. Agile and cunning.",
  "OwnerID": "master_of_shadows",
  "BaseSkills": [
    {"skill_id": "stealth", "percentage": 25},
    {"skill_id": "lockpicking", "percentage": 15}
  ]
}
```

#### 3.7.4. Scholar
```json
{
  "ID": "scholar",
  "Name": "Scholar",
  "Description": "A seeker of knowledge and ancient lore, skilled in languages, history, and deciphering forgotten texts.",
  "OwnerID": "lorekeepers_guild",
  "BaseSkills": [
    {"skill_id": "ancient_languages", "percentage": 25},
    {"skill_id": "history_of_middle_earth", "percentage": 20}
  ]
}
```

### 3.8. Quest Owners

High-level entities representing thematic ownership or overarching narrative arcs for quests. They possess a time-based influence budget for global world changes.

#### 3.8.1. Gandalf's Grand Plan
```json
{
  "ID": "gandalf_grand_plan",
  "Name": "Gandalf's Grand Plan",
  "Description": "The overarching strategic vision of Gandalf the Grey to counter the rising Shadow and guide the Free Peoples.",
  "LLMPromptContext": "You are the strategic mind behind Gandalf's efforts, focused on the larger picture of Middle-earth's fate. You orchestrate events and guide key individuals.",
  "CurrentInfluenceBudget": 200.0,
  "MaxInfluenceBudget": 200.0,
  "BudgetRegenRate": 0.2,
  "AssociatedQuestmakerIDs": ["urgent_message_questmaker", "road_to_rivendell_questmaker"]
}
```

#### 3.8.2. The Fellowship's Journey
```json
{
  "ID": "fellowship_journey",
  "Name": "The Fellowship's Journey",
  "Description": "The epic quest to destroy the One Ring, encompassing the trials and tribulations faced by the Fellowship.",
  "LLMPromptContext": "You represent the collective destiny and challenges of the Fellowship of the Ring. Your focus is on the perilous path to Mordor and the unity of its members.",
  "CurrentInfluenceBudget": 150.0,
  "MaxInfluenceBudget": 150.0,
  "BudgetRegenRate": 0.15,
  "AssociatedQuestmakerIDs": ["investigate_weathertop_questmaker", "delving_darkness_questmaker"]
}
```

#### 3.8.3. Shire Local Governance
```json
{
  "ID": "shire_local_governance",
  "Name": "Shire Local Governance",
  "Description": "The day-to-day affairs and well-being of the Shire, managed by its various councils and respected elders.",
  "LLMPromptContext": "You are concerned with the peaceful and orderly functioning of the Shire. Your quests involve community tasks, local disputes, and maintaining the hobbit way of life.",
  "CurrentInfluenceBudget": 70.0,
  "MaxInfluenceBudget": 70.0,
  "BudgetRegenRate": 0.1,
  "AssociatedQuestmakerIDs": ["shire_census_questmaker"]
}
```

#### 3.8.4. Bree Local Affairs
```json
{
  "ID": "bree_local_affairs",
  "Name": "Bree Local Affairs",
  "Description": "The mundane and sometimes mysterious happenings within the town of Bree and its immediate surroundings.",
  "LLMPromptContext": "You oversee the daily life and minor troubles of Bree. Your quests often involve missing items, suspicious characters, or local deliveries.",
  "CurrentInfluenceBudget": 75.0,
  "MaxInfluenceBudget": 75.0,
  "BudgetRegenRate": 0.08,
  "AssociatedQuestmakerIDs": ["missing_pony_questmaker"]
}
```

#### 3.8.5. Warrior Guild Trials
```json
{
  "ID": "warrior_guild_trials",
  "Name": "Warrior Guild Trials",
  "Description": "A series of challenges and tests designed to hone the skills and prove the worth of aspiring warriors.",
  "LLMPromptContext": "You are the spirit of martial challenge and discipline within the Warrior's Guild. Your quests are designed to push combatants to their limits and forge them into true warriors.",
  "CurrentInfluenceBudget": 80.0,
  "MaxInfluenceBudget": 80.0,
  "BudgetRegenRate": 0.11,
  "AssociatedQuestmakerIDs": ["training_regimen_questmaker"]
}
```

### 3.9. Questmakers (Specific to a single Quest)

Entities responsible for controlling the execution and progression of a single, specific quest. Their influence budget primarily comes from player actions.

#### 3.9.1. Urgent Message Questmaker
```json
{
  "ID": "urgent_message_questmaker",
  "Name": "Urgent Message Quest Controller",
  "LLMPromptContext": "You are the direct overseer of 'The Urgent Message' quest. Your focus is solely on ensuring the message is delivered to Strider swiftly and safely.",
  "CurrentInfluenceBudget": 0.0,
  "MaxInfluenceBudget": 50.0,
  "BudgetRegenRate": 0.0,
  "MemoriesAboutPlayers": {},
  "AvailableTools": []
}
```

#### 3.9.2. Missing Pony Questmaker
```json
{
  "ID": "missing_pony_questmaker",
  "Name": "Missing Pony Quest Controller",
  "LLMPromptContext": "You are the direct overseer of 'The Missing Pony' quest. Your goal is to ensure Bill the Pony is found and returned to Barliman Butterbur.",
  "CurrentInfluenceBudget": 0.0,
  "MaxInfluenceBudget": 30.0,
  "BudgetRegenRate": 0.0,
  "MemoriesAboutPlayers": {},
  "AvailableTools": []
}
```

#### 3.9.3. Investigate Weathertop Questmaker
```json
{
  "ID": "investigate_weathertop_questmaker",
  "Name": "Investigate Weathertop Quest Controller",
  "LLMPromptContext": "You are the direct overseer of 'Investigate Weathertop' quest. Your objective is to ensure the ruins are thoroughly investigated and findings reported.",
  "CurrentInfluenceBudget": 0.0,
  "MaxInfluenceBudget": 70.0,
  "BudgetRegenRate": 0.0,
  "MemoriesAboutPlayers": {},
  "AvailableTools": []
}
```

#### 3.9.4. Road to Rivendell Questmaker
```json
{
  "ID": "road_to_rivendell_questmaker",
  "Name": "Road to Rivendell Quest Controller",
  "LLMPromptContext": "You are the direct overseer of 'The Road to Rivendell' quest. Your goal is to ensure the traveler reaches Rivendell and speaks with Elrond.",
  "CurrentInfluenceBudget": 0.0,
  "MaxInfluenceBudget": 90.0,
  "BudgetRegenRate": 0.0,
  "MemoriesAboutPlayers": {},
  "AvailableTools": []
}
```

#### 3.9.5. Delving into Darkness Questmaker
```json
{
  "ID": "delving_darkness_questmaker",
  "Name": "Delving into Darkness Quest Controller",
  "LLMPromptContext": "You are the direct overseer of 'Delving into Darkness' quest. Your objective is to ensure the Moria entrance is investigated and reports are made.",
  "CurrentInfluenceBudget": 0.0,
  "MaxInfluenceBudget": 110.0,
  "BudgetRegenRate": 0.0,
  "MemoriesAboutPlayers": {},
  "AvailableTools": []
}
```

#### 3.9.6. Shire Census Questmaker
```json
{
  "ID": "shire_census_questmaker",
  "Name": "Shire Census Quest Controller",
  "LLMPromptContext": "You are the direct overseer of 'The Shire Census' quest. Your goal is to ensure the census is completed accurately by visiting the specified hobbits.",
  "CurrentInfluenceBudget": 0.0,
  "MaxInfluenceBudget": 40.0,
  "BudgetRegenRate": 0.0,
  "MemoriesAboutPlayers": {},
  "AvailableTools": []
}
```

#### 3.9.7. Training Regimen Questmaker
```json
{
  "ID": "training_regimen_questmaker",
  "Name": "Training Regimen Quest Controller",
  "LLMPromptContext": "You are the direct overseer of 'Training Regimen' quest. Your objective is to ensure the player successfully completes the training exercises and reports back.",
  "CurrentInfluenceBudget": 0.0,
  "MaxInfluenceBudget": 55.0,
  "BudgetRegenRate": 0.0,
  "MemoriesAboutPlayers": {},
  "AvailableTools": []
}
```