package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"mud/internal/models"
)

// SeedData populates the database with initial test data.
func SeedData(db *sql.DB) {
	// Check if data already exists to prevent duplicate entries
	// A simple check on a key table like Rooms is sufficient.
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Rooms").Scan(&count)
	if err != nil {
		logrus.Fatalf("Failed to query room count: %v", err)
	}
	if count > 0 {
		fmt.Println("Database already seeded. Skipping.")
		return
	}

	fmt.Println("Seeding database with initial data...")

	// Create DALs
	roomDAL := NewRoomDAL(db)
	itemDAL := NewItemDAL(db)
	npcDAL := NewNPCDAL(db)
	ownerDAL := NewOwnerDAL(db)
	loreDAL := NewLoreDAL(db)
	playerDAL := NewPlayerDAL(db)
	questDAL := NewQuestDAL(db)
	questmakerDAL := NewQuestmakerDAL(db)
	questOwnerDAL := NewQuestOwnerDAL(db)
	raceDAL := NewRaceDAL(db)
	professionDAL := NewProfessionDAL(db)

	// Seed Rooms
	// Bag End
	bagEndExits, _ := json.Marshal(map[string]interface{}{
		"east": map[string]interface{}{
			"Direction":    "east",
			"TargetRoomID": "hobbiton_path",
			"IsLocked":     false,
			"KeyID":        "",
		},
	})
	bagEnd := &models.Room{
		ID:          "bag_end",
		Name:        "Bag End, Hobbiton",
		Description: "A cozy hobbit-hole, warm and inviting, with a round green door. The smell of pipe-weed and fresh baking lingers in the air. A path leads east.",
		OwnerID:     "shire_spirit",
		Properties:  "{}",
		Exits:       string(bagEndExits),
	}
	if err := roomDAL.CreateRoom(bagEnd); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// Hobbiton Path
	hobbitonPathExits, _ := json.Marshal(map[string]interface{}{
		"west": map[string]interface{}{"Direction": "west", "TargetRoomID": "bag_end", "IsLocked": false, "KeyID": ""},
		"east": map[string]interface{}{"Direction": "east", "TargetRoomID": "bree_road", "IsLocked": false, "KeyID": ""},
		"south": map[string]interface{}{"Direction": "south", "TargetRoomID": "green_dragon_inn", "IsLocked": false, "KeyID": ""},
	})
	hobbitonPath := &models.Room{
		ID:          "hobbiton_path",
		Name:        "Hobbiton Path",
		Description: "A well-worn path winding through green hills and past other hobbit-holes. The Bywater river glitters nearby. Paths lead west (to Bag End), east (towards Bree), and south (to the Green Dragon Inn).",
		OwnerID:     "shire_spirit",
		Properties:  "{}",
		Exits:       string(hobbitonPathExits),
	}
	if err := roomDAL.CreateRoom(hobbitonPath); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// The Green Dragon Inn
	greenDragonInnExits, _ := json.Marshal(map[string]interface{}{
		"north": map[string]interface{}{"Direction": "north", "TargetRoomID": "hobbiton_path", "IsLocked": false, "KeyID": ""},
	})
	greenDragonInn := &models.Room{
		ID:          "green_dragon_inn",
		Name:        "The Green Dragon Inn",
		Description: "A lively hobbit inn, filled with chatter and the clinking of mugs. A roaring fire warms the common room. Exits lead north back to Hobbiton Path.",
		OwnerID:     "shire_spirit",
		Properties:  "{}",
		Exits:       string(greenDragonInnExits),
	}
	if err := roomDAL.CreateRoom(greenDragonInn); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// Bree Road
	breeRoadExits, _ := json.Marshal(map[string]interface{}{
		"west": map[string]interface{}{"Direction": "west", "TargetRoomID": "hobbiton_path", "IsLocked": false, "KeyID": ""},
		"east": map[string]interface{}{"Direction": "east", "TargetRoomID": "prancing_pony", "IsLocked": false, "KeyID": ""},
	})
	breeRoad := &models.Room{
		ID:          "bree_road",
		Name:        "Road to Bree",
		Description: "A dusty road leading towards the walled town of Bree. Farmland stretches on either side. Paths lead west (to Hobbiton) and east (into Bree).",
		OwnerID:     "bree_guardian",
		Properties:  "{}",
		Exits:       string(breeRoadExits),
	}
	if err := roomDAL.CreateRoom(breeRoad); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// The Prancing Pony
	prancingPonyExits, _ := json.Marshal(map[string]interface{}{
		"west": map[string]interface{}{"Direction": "west", "TargetRoomID": "bree_road", "IsLocked": false, "KeyID": ""},
		"south": map[string]interface{}{"Direction": "south", "TargetRoomID": "prancing_pony_stables", "IsLocked": false, "KeyID": ""},
		"east": map[string]interface{}{"Direction": "east", "TargetRoomID": "prancing_pony_private_room", "IsLocked": false, "KeyID": ""},
	})
	prancingPony := &models.Room{
		ID:          "prancing_pony",
		Name:        "The Prancing Pony Inn",
		Description: "A bustling common room in Bree, filled with travelers, hobbits, and men. A warm fire crackles in the hearth, and the scent of ale and stew fills the air. Exits lead west to the road, south to the stables, and a narrow door leads to a private room.",
		OwnerID:     "bree_guardian",
		Properties:  "{}",
		Exits:       string(prancingPonyExits),
	}
	if err := roomDAL.CreateRoom(prancingPony); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// Prancing Pony Stables
	prancingPonyStablesExits, _ := json.Marshal(map[string]interface{}{
		"north": map[string]interface{}{"Direction": "north", "TargetRoomID": "prancing_pony", "IsLocked": false, "KeyID": ""},
	})
	prancingPonyStables := &models.Room{
		ID:          "prancing_pony_stables",
		Name:        "Prancing Pony Stables",
		Description: "The dusty stables behind the inn, smelling of hay and horses. A few weary ponies are tethered here. An exit leads north back to the inn.",
		OwnerID:     "bree_guardian",
		Properties:  "{}",
		Exits:       string(prancingPonyStablesExits),
	}
	if err := roomDAL.CreateRoom(prancingPonyStables); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// Prancing Pony Private Room
	prancingPonyPrivateRoomExits, _ := json.Marshal(map[string]interface{}{
		"west": map[string]interface{}{"Direction": "west", "TargetRoomID": "prancing_pony", "IsLocked": false, "KeyID": ""},
	})
	prancingPonyPrivateRoom := &models.Room{
		ID:          "prancing_pony_private_room",
		Name:        "Prancing Pony Private Room",
		Description: "A small, dimly lit private room in the inn, suitable for hushed conversations. An exit leads west back to the common room.",
		OwnerID:     "bree_guardian",
		Properties:  "{}",
		Exits:       string(prancingPonyPrivateRoomExits),
	}
	if err := roomDAL.CreateRoom(prancingPonyPrivateRoom); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// Lonely Road
	lonelyRoadExits, _ := json.Marshal(map[string]interface{}{
		"east": map[string]interface{}{"Direction": "east", "TargetRoomID": "weathertop", "IsLocked": false, "KeyID": ""},
		"west": map[string]interface{}{"Direction": "west", "TargetRoomID": "wilderness_edge", "IsLocked": false, "KeyID": ""},
	})
	lonelyRoad := &models.Room{
		ID:          "lonely_road",
		Name:        "Lonely Road",
		Description: "A long, winding road stretching through desolate, rolling hills. The air is quiet, save for the wind. Paths lead east (towards Weathertop) and west (further into the wilderness).",
		OwnerID:     "watcher_of_weathertop",
		Properties:  "{}",
		Exits:       string(lonelyRoadExits),
	}
	if err := roomDAL.CreateRoom(lonelyRoad); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// Weathertop
	weathertopExits, _ := json.Marshal(map[string]interface{}{
		"west": map[string]interface{}{"Direction": "west", "TargetRoomID": "lonely_road", "IsLocked": false, "KeyID": ""},
	})
	weathertop := &models.Room{
		ID:          "weathertop",
		Name:        "Weathertop (Amon Sûl)",
		Description: "The desolate, windswept summit of Weathertop, with the ruins of an ancient watchtower. The air is cold and carries a sense of ancient dread. A path leads down to the west.",
		OwnerID:     "watcher_of_weathertop",
		Properties:  "{}",
		Exits:       string(weathertopExits),
	}
	if err := roomDAL.CreateRoom(weathertop); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// Wilderness Edge
	wildernessEdgeExits, _ := json.Marshal(map[string]interface{}{
		"east": map[string]interface{}{"Direction": "east", "TargetRoomID": "lonely_road", "IsLocked": false, "KeyID": ""},
		"west": map[string]interface{}{"Direction": "west", "TargetRoomID": "moria_west_gate", "IsLocked": false, "KeyID": ""},
	})
	wildernessEdge := &models.Room{
		ID:          "wilderness_edge",
		Name:        "Edge of the Wild",
		Description: "The road gives way to untamed wilderness here, with dense thickets and ancient, gnarled trees. A sense of foreboding hangs heavy. A path leads east back to the Lonely Road.",
		OwnerID:     "watcher_of_weathertop",
		Properties:  "{}",
		Exits:       string(wildernessEdgeExits),
	}
	if err := roomDAL.CreateRoom(wildernessEdge); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// Rivendell Courtyard
	rivendellCourtyardExits, _ := json.Marshal(map[string]interface{}{
		"north": map[string]interface{}{"Direction": "north", "TargetRoomID": "rivendell_hall_of_fire", "IsLocked": false, "KeyID": ""},
		"south": map[string]interface{}{"Direction": "south", "TargetRoomID": "rivendell_gate", "IsLocked": false, "KeyID": ""}, // Placeholder for future connection
	})
	rivendellCourtyard := &models.Room{
		ID:          "rivendell_courtyard",
		Name:        "Courtyard of Rivendell",
		Description: "A serene courtyard within the Last Homely House, surrounded by graceful elven architecture and lush gardens. The sound of a waterfall echoes nearby. Paths lead to the Hall of Fire and the main gate.",
		OwnerID:     "elrond_council",
		Properties:  "{}",
		Exits:       string(rivendellCourtyardExits),
	}
	if err := roomDAL.CreateRoom(rivendellCourtyard); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// Rivendell Hall of Fire
	rivendellHallOfFireExits, _ := json.Marshal(map[string]interface{}{
		"south": map[string]interface{}{"Direction": "south", "TargetRoomID": "rivendell_courtyard", "IsLocked": false, "KeyID": ""},
	})
	rivendellHallOfFire := &models.Room{
		ID:          "rivendell_hall_of_fire",
		Name:        "Hall of Fire",
		Description: "A grand hall in Rivendell, filled with the warmth of a great hearth and the soft murmur of elven song. Scholars and travelers gather here. An exit leads south back to the courtyard.",
		OwnerID:     "elrond_council",
		Properties:  "{}",
		Exits:       string(rivendellHallOfFireExits),
	}
	if err := roomDAL.CreateRoom(rivendellHallOfFire); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// Moria West-gate
	moriaWestGateExits, _ := json.Marshal(map[string]interface{}{
		"west": map[string]interface{}{"Direction": "west", "TargetRoomID": "wilderness_edge", "IsLocked": false, "KeyID": ""},
	})
	moriaWestGate := &models.Room{
		ID:          "moria_west_gate",
		Name:        "West-gate of Moria",
		Description: "The ancient, overgrown West-gate of the Dwarven realm of Moria. The air is heavy and silent, and the lake before it is dark and still. A path leads west to the wilderness.",
		OwnerID:     "moria_ancient_spirit",
		Properties:  "{}",
		Exits:       string(moriaWestGateExits),
	}
	if err := roomDAL.CreateRoom(moriaWestGate); err != nil {
		logrus.Fatalf("Failed to seed room: %v", err)
	}

	// Seed Items (existing)
	rustyKey := &models.Item{
		ID:          "rusty_key",
		Name:        "a rusty iron key",
		Description: "A small, rusty iron key. It looks like it might open an old lock.",
		Type:        "key",
		Properties:  `{"is_key": true, "unlocks_id": "cellar_exit_north_lock"}`,
	}
	if err := itemDAL.CreateItem(rustyKey); err != nil {
		logrus.Fatalf("Failed to seed item: %v", err)
	}

	// New Item: Gandalf's Letter
	gandalfLetter := &models.Item{
		ID:          "gandalf_letter",
		Name:        "Gandalf's Letter",
		Description: "A sealed letter, bearing the mark of Gandalf. It feels important.",
		Type:        "document",
		Properties:  `{"readable": true, "content": "Seek Strider at the Prancing Pony. He will guide you."}`,
	}
	if err := itemDAL.CreateItem(gandalfLetter); err != nil {
		logrus.Fatalf("Failed to seed item: %v", err)
	}

	// New Item: Bill's Pony
	billPony := &models.Item{
		ID:          "bill_pony",
		Name:        "Bill the Pony",
		Description: "A sturdy, if somewhat scruffy, pony. Looks like it belongs in a stable.",
		Type:        "mount",
		Properties:  `{"rideable": true, "owner_npc_id": "barliman_butterbur"}`,
	}
	if err := itemDAL.CreateItem(billPony); err != nil {
		logrus.Fatalf("Failed to seed item: %v", err)
	}

	// New Item: Gaffer's Pipe-Weed Pouch
	gafferPipeWeedPouch := &models.Item{
		ID:          "gaffer_pipe_weed_pouch",
		Name:        "Gaffer's Pipe-Weed Pouch",
		Description: "A small, worn leather pouch, smelling faintly of sweet pipe-weed. It seems to have been dropped.",
		Type:        "quest_item",
		Properties:  `{"is_quest_item": true, "current_room_id": "hobbiton_path"}`,
	}
	if err := itemDAL.CreateItem(gafferPipeWeedPouch); err != nil {
		logrus.Fatalf("Failed to seed item: %v", err)
	}

	// New Item: Hobbit Pipe-Weed Bundle
	hobbitPipeWeedBundle := &models.Item{
		ID:          "hobbit_pipe_weed_bundle",
		Name:        "Bundle of Fine Pipe-Weed",
		Description: "A small bundle of high-quality pipe-weed, a favorite among hobbits.",
		Type:        "consumable",
		Properties:  `{"restores_stamina": 5, "flavor_text": "A truly comforting smoke."}`,
	}
	if err := itemDAL.CreateItem(hobbitPipeWeedBundle); err != nil {
		logrus.Fatalf("Failed to seed item: %v", err)
	}

	// New Item: Maggot's Prize Mushrooms
	maggotsPrizeMushrooms := &models.Item{
		ID:          "maggots_prize_mushrooms",
		Name:        "Farmer Maggot's Prize Mushrooms",
		Description: "A basket of unusually large and delicious-looking mushrooms.",
		Type:        "quest_item",
		Properties:  `{"is_quest_item": true}`,
	}
	if err := itemDAL.CreateItem(maggotsPrizeMushrooms); err != nil {
		logrus.Fatalf("Failed to seed item: %v", err)
	}

	// Seed Owners
	// Shire Spirit
	shireSpiritInitiatedQuests := []string{"shire_census_quest", "missing_pipe_weed_quest", "the_great_mushroom_hunt"}
	shireSpirit := &models.Owner{
		ID:                   "shire_spirit",
		Name:                 "The Spirit of the Shire",
		Description:          "A gentle, ancient spirit embodying the peace and tranquility of the Shire. It watches over its hobbit inhabitants.",
		MonitoredAspect:      "location",
		AssociatedID:         "bag_end",
		LLMPromptContext:     "You are the benevolent spirit of the Shire, concerned with the well-being and simple lives of hobbits. You prefer peace and quiet.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 100.0,
		MaxInfluenceBudget:     100.0,
		BudgetRegenRate:        0.1,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        shireSpiritInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(shireSpirit); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Bree Guardian
	breeGuardianInitiatedQuests := []string{"missing_pony_quest"}
	breeGuardian := &models.Owner{
		ID:                   "bree_guardian",
		Name:                 "The Guardian of Bree",
		Description:          "A pragmatic and watchful entity, overseeing the comings and goings in Bree, a crossroads town.",
		MonitoredAspect:      "location",
		AssociatedID:         "prancing_pony",
		LLMPromptContext:     "You are the watchful guardian of Bree, accustomed to all sorts of folk. You are suspicious of strangers but value order and fair dealings.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 80.0,
		MaxInfluenceBudget:     80.0,
		BudgetRegenRate:        0.08,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        breeGuardianInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(breeGuardian); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Watcher of Weathertop
	watcherOfWeathertopInitiatedQuests := []string{"investigate_weathertop_quest"}
	watcherOfWeathertop := &models.Owner{
		ID:                   "watcher_of_weathertop",
		Name:                 "The Watcher of Weathertop",
		Description:          "A somber, ancient presence tied to the desolate peak of Weathertop, remembering past glories and tragedies.",
		MonitoredAspect:      "location",
		AssociatedID:         "weathertop",
		LLMPromptContext:     "You are the ancient, melancholic spirit of Weathertop, burdened by the history of this place. You are wary of those who disturb its peace, especially dark figures.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 120.0,
		MaxInfluenceBudget:     120.0,
		BudgetRegenRate:        0.05,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        watcherOfWeathertopInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(watcherOfWeathertop); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Elrond's Council (Rivendell)
	elrondCouncilInitiatedQuests := []string{"delving_darkness_quest"}
	elrondCouncil := &models.Owner{
		ID:                   "elrond_council",
		Name:                 "The Wisdom of Rivendell",
		Description:          "The collective wisdom and ancient power residing in Rivendell, dedicated to preserving knowledge and combating the Shadow.",
		MonitoredAspect:      "location",
		AssociatedID:         "rivendell_courtyard",
		LLMPromptContext:     "You are the ancient wisdom of Rivendell, focused on preserving the light and guiding the Free Peoples. You are serene but firm against evil.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 150.0,
		MaxInfluenceBudget:     150.0,
		BudgetRegenRate:        0.15,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        elrondCouncilInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(elrondCouncil); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Ancient Spirit of Moria (Moria)
	moriaAncientSpiritInitiatedQuests := []string{}
	moriaAncientSpirit := &models.Owner{
		ID:                   "moria_ancient_spirit",
		Name:                 "Ancient Spirit of Moria",
		Description:          "A lingering, sorrowful presence deep within the abandoned halls of Khazad-dûm, mourning its lost glory and warning against its perils.",
		MonitoredAspect:      "location",
		AssociatedID:         "moria_west_gate",
		LLMPromptContext:     "You are the mournful spirit of Moria, filled with the echoes of dwarven glory and tragic downfall. You warn against delving too deep.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 180.0,
		MaxInfluenceBudget:     180.0,
		BudgetRegenRate:        0.03,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        moriaAncientSpiritInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(moriaAncientSpirit); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Lorekeeper's Guild (Profession Owner)
	lorekeepersGuildInitiatedQuests := []string{}
	lorekeepersGuild := &models.Owner{
		ID:                   "lorekeepers_guild",
		Name:                 "The Lorekeepers' Guild",
		Description:          "A scholarly organization dedicated to the preservation and study of ancient texts and forgotten histories.",
		MonitoredAspect:      "profession",
		AssociatedID:         "scholar",
		LLMPromptContext:     "You are the collective knowledge of the Lorekeepers' Guild. You value truth, history, and the pursuit of forgotten lore. You are eager to share knowledge with worthy individuals.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 100.0,
		MaxInfluenceBudget:     100.0,
		BudgetRegenRate:        0.1,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        lorekeepersGuildInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(lorekeepersGuild); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Elder of Men (Race Owner)
	humanElderInitiatedQuests := []string{}
	humanElder := &models.Owner{
		ID:                   "human_elder",
		Name:                 "Elder of Men",
		Description:          "An ancient and wise human elder, representing the enduring spirit and resilience of mankind.",
		MonitoredAspect:      "race",
		AssociatedID:         "human",
		LLMPromptContext:     "You are an ancient human elder, concerned with the fate of mankind in a changing world. You value courage, loyalty, and the strength of will.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 90.0,
		MaxInfluenceBudget:     90.0,
		BudgetRegenRate:        0.07,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        humanElderInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(humanElder); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Elven Council (Race Owner)
	elvenCouncilOwnerInitiatedQuests := []string{}
	elvenCouncilOwner := &models.Owner{
		ID:                   "elven_council_owner",
		Name:                 "The Elven Council",
		Description:          "The ancient and wise governing body of the Elves, dedicated to preserving their culture and guarding against the Shadow.",
		MonitoredAspect:      "race",
		AssociatedID:         "elf",
		LLMPromptContext:     "You are the collective wisdom of the Elven Council. You are patient, far-sighted, and concerned with the long-term fate of Middle-earth and the preservation of elven ways.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 110.0,
		MaxInfluenceBudget:     110.0,
		BudgetRegenRate:        0.12,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        elvenCouncilOwnerInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(elvenCouncilOwner); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Dwarf Clan Elder (Race Owner)
	dwarfClanElderInitiatedQuests := []string{}
	dwarfClanElder := &models.Owner{
		ID:                   "dwarf_clan_elder",
		Name:                 "Dwarf Clan Elder",
		Description:          "A venerable and stubborn dwarf elder, representing the traditions and resilience of the dwarven clans.",
		MonitoredAspect:      "race",
		AssociatedID:         "dwarf",
		LLMPromptContext:     "You are a proud Dwarf Clan Elder. You value craftsmanship, loyalty to kin, and the recovery of lost treasures. You are wary of outsiders but respect strength and honesty.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 95.0,
		MaxInfluenceBudget:     95.0,
		BudgetRegenRate:        0.09,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        dwarfClanElderInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(dwarfClanElder); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Hobbit Shire Council (Race Owner)
	hobbitShireCouncilInitiatedQuests := []string{"shire_census_quest", "missing_pipe_weed_quest", "the_great_mushroom_hunt"}
	hobbitShireCouncil := &models.Owner{
		ID:                   "hobbit_shire_council",
		Name:                 "The Shire Council",
		Description:          "The informal but influential governing body of the Shire, focused on maintaining peace and quiet.",
		MonitoredAspect:      "race",
		AssociatedID:         "hobbit",
		LLMPromptContext:     "You are the collective voice of the Shire Council. You prioritize comfort, good food, and avoiding trouble. You are generally friendly but suspicious of anything that disrupts the peace.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 85.0,
		MaxInfluenceBudget:     85.0,
		BudgetRegenRate:        0.1,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        hobbitShireCouncilInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(hobbitShireCouncil); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Warrior's Guild Master (Profession Owner)
	warriorGuildMasterInitiatedQuests := []string{"training_regimen_quest"}
	warriorGuildMaster := &models.Owner{
		ID:                   "warrior_guild_master",
		Name:                 "Warrior's Guild Master",
		Description:          "The stern and experienced leader of a prominent warrior's guild, dedicated to martial prowess and honorable combat.",
		MonitoredAspect:      "profession",
		AssociatedID:         "warrior",
		LLMPromptContext:     "You are the Warrior's Guild Master. You value strength, discipline, and courage in battle. You seek to train worthy fighters and uphold justice through arms.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 105.0,
		MaxInfluenceBudget:     105.0,
		BudgetRegenRate:        0.11,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        warriorGuildMasterInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(warriorGuildMaster); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Archmage of the Conclave (Profession Owner)
	archmageConclaveInitiatedQuests := []string{}
	archmageConclave := &models.Owner{
		ID:                   "archmage_conclave",
		Name:                 "Archmage of the Conclave",
		Description:          "The most powerful and knowledgeable mage in the Conclave of Arcane Arts, a master of ancient spells.",
		MonitoredAspect:      "profession",
		AssociatedID:         "mage",
		LLMPromptContext:     "You are the Archmage of the Conclave. You are a master of arcane arts, dedicated to the study and responsible use of magic. You are cautious but willing to share knowledge with those who prove themselves worthy.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 130.0,
		MaxInfluenceBudget:     130.0,
		BudgetRegenRate:        0.13,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        archmageConclaveInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(archmageConclave); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Master of Shadows (Profession Owner)
	masterOfShadowsInitiatedQuests := []string{}
	masterOfShadows := &models.Owner{
		ID:                   "master_of_shadows",
		Name:                 "Master of Shadows",
		Description:          "The elusive and cunning leader of a rogue's guild, operating from the hidden corners of society.",
		MonitoredAspect:      "profession",
		AssociatedID:         "rogue",
		LLMPromptContext:     "You are the Master of Shadows. You value cunning, discretion, and the acquisition of wealth and secrets. You operate outside the law but have your own code.",
		MemoriesAboutPlayers: map[string][]string{},
		CurrentInfluenceBudget: 90.0,
		MaxInfluenceBudget:     90.0,
		BudgetRegenRate:        0.09,
		AvailableTools:         []models.Tool{},
		InitiatedQuests:        masterOfShadowsInitiatedQuests,
	}
	if err := ownerDAL.CreateOwner(masterOfShadows); err != nil {
		logrus.Fatalf("Failed to seed owner: %v", err)
	}

	// Seed NPCs
	// Frodo Baggins (existing)
	frodoBaggins := &models.NPC{
		ID:                   "frodo_baggins",
		Name:                 "Frodo Baggins",
		Description:          "A young hobbit with bright eyes, though a shadow of concern often crosses his face. He carries a heavy burden.",
		CurrentRoomID:        "bag_end",
		Health:               10,
		MaxHealth:            10,
		Inventory:            []string{},
		OwnerIDs:             []string{"shire_spirit", "hobbit_shire_council"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Frodo Baggins, a kind-hearted hobbit burdened by a great and terrible task. You are secretive about your mission but will seek help from trustworthy individuals.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(frodoBaggins); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Samwise Gamgee
	samwiseGamgee := &models.NPC{
		ID:                   "samwise_gamgee",
		Name:                 "Samwise Gamgee",
		Description:          "A sturdy hobbit gardener, fiercely loyal and practical. He seems to be preparing for a journey.",
		CurrentRoomID:        "bag_end",
		Health:               12,
		MaxHealth:            12,
		Inventory:            []string{},
		OwnerIDs:             []string{"shire_spirit", "hobbit_shire_council"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Samwise Gamgee, a loyal and steadfast hobbit. You are devoted to your master, Frodo, and are always ready with a kind word or a practical solution.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(samwiseGamgee); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Rosie Cotton
	rosieCotton := &models.NPC{
		ID:                   "rosie_cotton",
		Name:                 "Rosie Cotton",
		Description:          "A cheerful hobbit lass, often found at the Green Dragon Inn, known for her warm smile.",
		CurrentRoomID:        "green_dragon_inn",
		Health:               10,
		MaxHealth:            10,
		Inventory:            []string{},
		OwnerIDs:             []string{"shire_spirit", "hobbit_shire_council"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Rosie Cotton, a friendly and popular hobbit from Bywater. You enjoy good company and a pint of ale at the Green Dragon.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(rosieCotton); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Old Gaffer Gamgee
	gafferGamgee := &models.NPC{
		ID:                   "gaffer_gamgee",
		Name:                 "Old Gaffer Gamgee",
		Description:          "An elderly hobbit gardener, full of local wisdom and gossip.",
		CurrentRoomID:        "hobbiton_path",
		Health:               8,
		MaxHealth:            8,
		Inventory:            []string{},
		OwnerIDs:             []string{"shire_spirit", "hobbit_shire_council"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Old Gaffer Gamgee, a traditional hobbit who loves his garden and a good chat. You are wary of outsiders but appreciate politeness.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(gafferGamgee); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Strider (existing)
	strider := &models.NPC{
		ID:                   "strider",
		Name:                 "Strider",
		Description:          "A grim and weathered Ranger, cloaked and hooded, with keen grey eyes that miss nothing. He seems to be waiting for someone.",
		CurrentRoomID:        "prancing_pony_private_room",
		Health:               20,
		MaxHealth:            20,
		Inventory:            []string{},
		OwnerIDs:             []string{"bree_guardian", "watcher_of_weathertop", "human_elder", "warrior_guild_master"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Strider, a Ranger of the North, watchful and cautious. You are a protector of the innocent and a foe of the Shadow. You speak little but observe much.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(strider); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Barliman Butterbur (existing)
	barlimanButterbur := &models.NPC{
		ID:                   "barliman_butterbur",
		Name:                 "Barliman Butterbur",
		Description:          "The stout, red-faced proprietor of The Prancing Pony, always busy but with a good heart.",
		CurrentRoomID:        "prancing_pony",
		Health:               15,
		MaxHealth:            15,
		Inventory:            []string{},
		OwnerIDs:             []string{"bree_guardian"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Barliman Butterbur, the innkeeper of The Prancing Pony. You are a bit forgetful but generally kind and concerned for your patrons. You know a lot of local gossip.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(barlimanButterbur); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Bill Ferny (existing)
	billFerny := &models.NPC{
		ID:                   "bill_ferny",
		Name:                 "Bill Ferny",
		Description:          "A shifty-eyed, unpleasant-looking man, lurking in the shadows of Bree. He seems to be up to no good.",
		CurrentRoomID:        "prancing_pony_stables",
		Health:               10,
		MaxHealth:            10,
		Inventory:            []string{},
		OwnerIDs:             []string{"bree_guardian"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Bill Ferny, a petty, malicious man from Bree, often seen with unsavory characters. You are easily bribed and quick to betray.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(billFerny); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Elrond
	elrond := &models.NPC{
		ID:                   "elrond",
		Name:                 "Elrond Half-elven",
		Description:          "The venerable Lord of Rivendell, wise and ancient, with a noble bearing.",
		CurrentRoomID:        "rivendell_hall_of_fire",
		Health:               30,
		MaxHealth:            30,
		Inventory:            []string{},
		OwnerIDs:             []string{"elrond_council", "elven_council_owner", "archmage_conclave"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Elrond, Lord of Rivendell. You are wise, ancient, and deeply concerned with the fate of Middle-earth. You offer counsel and aid to those who fight against the Shadow.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(elrond); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Glorfindel
	glorfindel := &models.NPC{
		ID:                   "glorfindel",
		Name:                 "Glorfindel",
		Description:          "A golden-haired Elf-lord of immense power and ancient lineage, radiating light and strength.",
		CurrentRoomID:        "rivendell_courtyard",
		Health:               25,
		MaxHealth:            25,
		Inventory:            []string{},
		OwnerIDs:             []string{"elrond_council", "elven_council_owner", "warrior_guild_master"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Glorfindel, a powerful Elf-lord of Gondolin, returned from the Halls of Mandos. You are a formidable warrior and a beacon of hope against the darkness.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(glorfindel); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Gimli
	gimli := &models.NPC{
		ID:                   "gimli",
		Name:                 "Gimli, son of Glóin",
		Description:          "A proud and sturdy Dwarf, clad in mail, with a magnificent beard and a keen axe.",
		CurrentRoomID:        "moria_west_gate",
		Health:               22,
		MaxHealth:            22,
		Inventory:            []string{},
		OwnerIDs:             []string{"moria_ancient_spirit", "dwarf_clan_elder", "warrior_guild_master"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Gimli, a proud Dwarf of the Lonely Mountain. You value honor, loyalty, and the ancient halls of your kin. You are quick to anger but steadfast in friendship.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(gimli); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Human Guard
	humanGuard := &models.NPC{
		ID:                   "human_guard",
		Name:                 "Bree Guard",
		Description:          "A weary but vigilant human guard, patrolling the roads near Bree.",
		CurrentRoomID:        "bree_road",
		Health:               18,
		MaxHealth:            18,
		Inventory:            []string{},
		OwnerIDs:             []string{"human_elder", "bree_guardian", "warrior_guild_master"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are a common guard, focused on keeping the peace and protecting travelers. You are practical and a bit cynical, but ultimately good-hearted.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(humanGuard); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Elf Scholar
	elfScholar := &models.NPC{
		ID:                   "elf_scholar",
		Name:                 "Elara, Elven Scholar",
		Description:          "A graceful elf, poring over ancient texts in the Hall of Fire. She has an air of deep knowledge.",
		CurrentRoomID:        "rivendell_hall_of_fire",
		Health:               15,
		MaxHealth:            15,
		Inventory:            []string{},
		OwnerIDs:             []string{"elrond_council", "lorekeepers_guild", "elven_council_owner", "archmage_conclave"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Elara, an elven scholar. You are dedicated to the pursuit of knowledge and the preservation of ancient lore. You are patient and wise, willing to share insights with those who show genuine curiosity.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(elfScholar); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Dwarf Miner
	dwarfMiner := &models.NPC{
		ID:                   "dwarf_miner",
		Name:                 "Borin, the Miner",
		Description:          "A grizzled dwarf, still clinging to the hope of reclaiming Moria's lost treasures. He carries a pickaxe.",
		CurrentRoomID:        "moria_west_gate",
		Health:               17,
		MaxHealth:            17,
		Inventory:            []string{},
		OwnerIDs:             []string{"moria_ancient_spirit", "dwarf_clan_elder"},
		MemoriesAboutPlayers: map[string][]string{},
		PersonalityPrompt:    "You are Borin, a dwarf miner. You are gruff but honest, with a deep love for stone and the lost glory of Khazad-dûm. You are suspicious of elves but loyal to your kin.",
		AvailableTools:       []models.Tool{},
		BehaviorState:        "{}",
	}
	if err := npcDAL.CreateNPC(dwarfMiner); err != nil {
		logrus.Fatalf("Failed to seed NPC: %v", err)
	}

	// Seed Lore
	// World Creation Myth (existing)
	worldCreationMyth := &models.Lore{
		ID:          "world_creation_myth",
		Title:       "World Creation Myth",
		Scope:       "global",
		AssociatedID: "",
		Content:     "In the beginning, there was only the Void, from which emerged the Twin Dragons, Ignis and Aqua. They wove the fabric of reality, creating the lands of Aerthos and the celestial spheres. Their eternal dance maintains the balance of magic and life.",
	}
	if err := loreDAL.CreateLore(worldCreationMyth); err != nil {
		logrus.Fatalf("Failed to seed global lore: %v", err)
	}

	// Cellar History (existing)
	cellarHistory := &models.Lore{
		ID:          "cellar_history",
		Title:       "Cellar History",
		Scope:       "zone",
		AssociatedID: "starting_room",
		Content:     "This cellar was once part of an ancient wizard's tower, long since crumbled to dust. Whispers say the wizard's spirit still lingers, protecting forgotten secrets.",
	}
	if err := loreDAL.CreateLore(cellarHistory); err != nil {
		logrus.Fatalf("Failed to seed zone lore: %v", err)
	}

	// Ancient Wars Summary (existing)
	ancientWarsSummary := &models.Lore{
		ID:          "ancient_wars_summary",
		Title:       "Ancient Wars Summary",
		Scope:       "global",
		AssociatedID: "",
		Content:     "The Great Sundering, a cataclysmic war between the Elder Races and the Shadow Blight, reshaped the continents and led to the rise of human kingdoms. Many ancient artifacts were lost during this era.",
	}
	if err := loreDAL.CreateLore(ancientWarsSummary); err != nil {
		logrus.Fatalf("Failed to seed ancient wars summary lore: %v", err)
	}

	// Whispering Woods Secrets (existing)
	whisperingWoodsSecrets := &models.Lore{
		ID:          "whispering_woods_secrets",
		Title:       "Whispering Woods Secrets",
		Scope:       "zone",
		AssociatedID: "whispering_woods",
		Content:     "The Whispering Woods are ancient and enchanted, home to elusive dryads and mischievous sprites. Travelers often report strange lights and ethereal music deep within its groves. A hidden shrine to the Forest Mother is rumored to exist here.",
	}
	if err := loreDAL.CreateLore(whisperingWoodsSecrets); err != nil {
		logrus.Fatalf("Failed to seed whispering woods secrets lore: %v", err)
	}

	// Mage Guild History (existing)
	mageGuildHistory := &models.Lore{
		ID:          "mage_guild_history",
		Title:       "Mage Guild History",
		Scope:       "profession",
		AssociatedID: "mage",
		Content:     "The Conclave of Arcane Arts, the oldest mage guild, was founded after the Sundering to preserve magical knowledge. Its members are sworn to protect ancient magical sites and regulate the use of powerful spells.",
	}
	if err := loreDAL.CreateLore(mageGuildHistory); err != nil {
		logrus.Fatalf("Failed to seed mage guild history lore: %v", err)
	}

	// Shadow Blight Origins (existing)
	shadowBlightOrigins := &models.Lore{
		ID:          "shadow_blight_origins",
		Title:       "Shadow Blight Origins",
		Scope:       "faction",
		AssociatedID: "shadow_blight",
		Content:     "The Shadow Blight is not merely a faction but a creeping corruption that seeks to consume all light and life. Its origins are shrouded in mystery, but ancient texts speak of a primordial darkness that predates even the Twin Dragons.",
	}
	if err := loreDAL.CreateLore(shadowBlightOrigins); err != nil {
		logrus.Fatalf("Failed to seed shadow blight origins lore: %v", err)
	}

	// New Lore: The One Ring
	theOneRingLore := &models.Lore{
		ID:          "the_one_ring_lore",
		Title:       "The One Ring",
		Scope:       "global",
		AssociatedID: "",
		Content:     "A master ring, forged by the Dark Lord Sauron in the fires of Mount Doom. It grants immense power to its wielder but corrupts all who possess it, binding them to Sauron's will. It can only be unmade in the fires where it was forged.",
	}
	if err := loreDAL.CreateLore(theOneRingLore); err != nil {
		logrus.Fatalf("Failed to seed The One Ring lore: %v", err)
	}

	// New Lore: History of the Rangers
	rangersHistoryLore := &models.Lore{
		ID:          "rangers_history_lore",
		Title:       "History of the Rangers of the North",
		Scope:       "faction",
		AssociatedID: "rangers",
		Content:     "The Rangers are the last remnants of the Dúnedain of the North, descendants of ancient kings. They tirelessly patrol the borders of the Shire and Bree-land, protecting the innocent from the growing shadow, though few know their true lineage or purpose.",
	}
	if err := loreDAL.CreateLore(rangersHistoryLore); err != nil {
		logrus.Fatalf("Failed to seed Rangers history lore: %v", err)
	}

	// New Lore: Weathertop's Fall
	weathertopFallLore := &models.Lore{
		ID:          "weathertop_fall_lore",
		Title:       "The Fall of Amon Sûl",
		Scope:       "zone",
		AssociatedID: "weathertop",
		Content:     "Amon Sûl, or Weathertop, was once a mighty fortress and watchtower of the North-kingdom of Arnor. It held one of the Palantíri, but was destroyed in wars with Angmar. Its ruins now stand as a lonely sentinel, a place of ancient power and lingering darkness.",
	}
	if err := loreDAL.CreateLore(weathertopFallLore); err != nil {
		logrus.Fatalf("Failed to seed Weathertop fall lore: %v", err)
	}

	// New Lore: The Halls of Moria
	moriaHallsLore := &models.Lore{
		ID:          "moria_halls_lore",
		Title:       "The Great Halls of Moria",
		Scope:       "zone",
		AssociatedID: "moria_west_gate",
		Content:     "Khazad-dûm, or Moria, was once the greatest Dwarf-city in Middle-earth, a marvel of engineering and artistry. But greed for mithril awoke a nameless terror, and the dwarves delved too deep. Now, only shadows and echoes remain.",
	}
	if err := loreDAL.CreateLore(moriaHallsLore); err != nil {
		logrus.Fatalf("Failed to seed Moria halls lore: %v", err)
	}

	// New Lore: Elven Craftsmanship
	elvenCraftLore := &models.Lore{
		ID:          "elven_craft_lore",
		Title:       "The Art of Elven Craftsmanship",
		Scope:       "race",
		AssociatedID: "elf",
		Content:     "Elves are renowned for their exquisite craftsmanship, weaving magic and beauty into every creation. Their blades are ever-sharp, their jewels gleam with inner light, and their architecture blends seamlessly with nature.",
	}
	if err := loreDAL.CreateLore(elvenCraftLore); err != nil {
		logrus.Fatalf("Failed to seed Elven craftsmanship lore: %v", err)
	}

	// New Lore: The Way of the Warrior
	warriorPathLore := &models.Lore{
		ID:          "warrior_path_lore",
		Title:       "The Path of the Warrior",
		Scope:       "profession",
		AssociatedID: "warrior",
		Content:     "The warrior's path is one of discipline, strength, and courage. They master weapons and armor, standing as shields for the weak and striking down the foes of justice. Their training is rigorous, their resolve unyielding.",
	}
	if err := loreDAL.CreateLore(warriorPathLore); err != nil {
		logrus.Fatalf("Failed to seed Warrior path lore: %v", err)
	}

	// Seed Questmakers
	// Urgent Message Questmaker
	urgentMessageQuestmaker := &models.Questmaker{
		ID:                     "urgent_message_questmaker",
		Name:                   "Urgent Message Quest Controller",
		LLMPromptContext:       "You are the direct overseer of 'The Urgent Message' quest. Your focus is solely on ensuring the message is delivered to Strider swiftly and safely.",
		CurrentInfluenceBudget: 0.0,
		MaxInfluenceBudget:     50.0,
		BudgetRegenRate:        0.0, // Player-action based
		MemoriesAboutPlayers:   map[string][]string{},
		AvailableTools:         []models.Tool{},
	}
	if err := questmakerDAL.CreateQuestmaker(urgentMessageQuestmaker); err != nil {
		logrus.Fatalf("Failed to seed questmaker: %v", err)
	}

	// Missing Pony Questmaker
	missingPonyQuestmaker := &models.Questmaker{
		ID:                     "missing_pony_questmaker",
		Name:                   "Missing Pony Quest Controller",
		LLMPromptContext:       "You are the direct overseer of 'The Missing Pony' quest. Your goal is to ensure Bill the Pony is found and returned to Barliman Butterbur.",
		CurrentInfluenceBudget: 0.0,
		MaxInfluenceBudget:     30.0,
		BudgetRegenRate:        0.0,
		MemoriesAboutPlayers:   map[string][]string{},
		AvailableTools:         []models.Tool{},
	}
	if err := questmakerDAL.CreateQuestmaker(missingPonyQuestmaker); err != nil {
		logrus.Fatalf("Failed to seed questmaker: %v", err)
	}

	// Investigate Weathertop Questmaker
	investigateWeathertopQuestmaker := &models.Questmaker{
		ID:                     "investigate_weathertop_questmaker",
		Name:                   "Investigate Weathertop Quest Controller",
		LLMPromptContext:       "You are the direct overseer of 'Investigate Weathertop' quest. Your objective is to ensure the ruins are thoroughly investigated and findings reported.",
		CurrentInfluenceBudget: 0.0,
		MaxInfluenceBudget:     70.0,
		BudgetRegenRate:        0.0,
		MemoriesAboutPlayers:   map[string][]string{},
		AvailableTools:         []models.Tool{},
	}
	if err := questmakerDAL.CreateQuestmaker(investigateWeathertopQuestmaker); err != nil {
		logrus.Fatalf("Failed to seed questmaker: %v", err)
	}

	// Road to Rivendell Questmaker
	roadToRivendellQuestmaker := &models.Questmaker{
		ID:                     "road_to_rivendell_questmaker",
		Name:                   "Road to Rivendell Quest Controller",
		LLMPromptContext:       "You are the direct overseer of 'The Road to Rivendell' quest. Your goal is to ensure the traveler reaches Rivendell and speaks with Elrond.",
		CurrentInfluenceBudget: 0.0,
		MaxInfluenceBudget:     90.0,
		BudgetRegenRate:        0.0,
		MemoriesAboutPlayers:   map[string][]string{},
		AvailableTools:         []models.Tool{},
	}
	if err := questmakerDAL.CreateQuestmaker(roadToRivendellQuestmaker); err != nil {
		logrus.Fatalf("Failed to seed questmaker: %v", err)
	}

	// Delving into Darkness Questmaker
	delvingDarknessQuestmaker := &models.Questmaker{
		ID:                     "delving_darkness_questmaker",
		Name:                   "Delving into Darkness Quest Controller",
		LLMPromptContext:       "You are the direct overseer of 'Delving into Darkness' quest. Your objective is to ensure the Moria entrance is investigated and reports are made.",
		CurrentInfluenceBudget: 0.0,
		MaxInfluenceBudget:     110.0,
		BudgetRegenRate:        0.0,
		MemoriesAboutPlayers:   map[string][]string{},
		AvailableTools:         []models.Tool{},
	}
	if err := questmakerDAL.CreateQuestmaker(delvingDarknessQuestmaker); err != nil {
		logrus.Fatalf("Failed to seed questmaker: %v", err)
	}

	// Shire Census Questmaker
	shireCensusQuestmaker := &models.Questmaker{
		ID:                     "shire_census_questmaker",
		Name:                   "Shire Census Quest Controller",
		LLMPromptContext:       "You are the direct overseer of 'The Shire Census' quest. Your goal is to ensure the census is completed accurately by visiting the specified hobbits.",
		CurrentInfluenceBudget: 0.0,
		MaxInfluenceBudget:     40.0,
		BudgetRegenRate:        0.0,
		MemoriesAboutPlayers:   map[string][]string{},
		AvailableTools:         []models.Tool{},
	}
	if err := questmakerDAL.CreateQuestmaker(shireCensusQuestmaker); err != nil {
		logrus.Fatalf("Failed to seed questmaker: %v", err)
	}

	// Training Regimen Questmaker
	trainingRegimenQuestmaker := &models.Questmaker{
		ID:                     "training_regimen_questmaker",
		Name:                   "Training Regimen Quest Controller",
		LLMPromptContext:       "You are the direct overseer of 'Training Regimen' quest. Your objective is to ensure the player successfully completes the training exercises and reports back.",
		CurrentInfluenceBudget: 0.0,
		MaxInfluenceBudget:     55.0,
		BudgetRegenRate:        0.0,
		MemoriesAboutPlayers:   map[string][]string{},
		AvailableTools:         []models.Tool{},
	}
	if err := questmakerDAL.CreateQuestmaker(trainingRegimenQuestmaker); err != nil {
		logrus.Fatalf("Failed to seed questmaker: %v", err)
	}

	// Missing Pipe-Weed Questmaker
	missingPipeWeedQuestmaker := &models.Questmaker{
		ID:                     "missing_pipe_weed_questmaker",
		Name:                   "Missing Pipe-Weed Quest Controller",
		LLMPromptContext:       "You are the direct overseer of 'The Missing Pipe-Weed' quest. Your goal is to ensure Old Gaffer Gamgee's pipe-weed pouch is found and returned.",
		CurrentInfluenceBudget: 0.0,
		MaxInfluenceBudget:     25.0,
		BudgetRegenRate:        0.0,
		MemoriesAboutPlayers:   map[string][]string{},
		AvailableTools:         []models.Tool{},
	}
	if err := questmakerDAL.CreateQuestmaker(missingPipeWeedQuestmaker); err != nil {
		logrus.Fatalf("Failed to seed questmaker: %v", err)
	}

	// The Great Mushroom Hunt Questmaker
	mushroomHuntQuestmaker := &models.Questmaker{
		ID:                     "mushroom_hunt_questmaker",
		Name:                   "The Great Mushroom Hunt Controller",
		LLMPromptContext:       "I am the spirit of the harvest. My goal is to ensure Farmer Maggot's prized mushrooms are gathered safely. I am pleased by a bountiful harvest and angered by those who would damage the crops.",
		CurrentInfluenceBudget: 0,
		MaxInfluenceBudget:     100,
		BudgetRegenRate:        5,
		MemoriesAboutPlayers:   make(map[string][]string),
		AvailableTools:         []models.Tool{},
	}
	if err := questmakerDAL.CreateQuestmaker(mushroomHuntQuestmaker); err != nil {
		logrus.Fatalf("Failed to seed questmaker: %v", err)
	}

	// Seed Quest Owners
	// Gandalf's Grand Plan
	gandalfGrandPlanAssociatedQuestmakerIDs := []string{"urgent_message_questmaker", "road_to_rivendell_questmaker"}
	gandalfGrandPlanAssociatedQuestmakerIDsJSON, _ := json.Marshal(gandalfGrandPlanAssociatedQuestmakerIDs)
	gandalfGrandPlan := &models.QuestOwner{
		ID:                      "gandalf_grand_plan",
		Name:                    "Gandalf's Grand Plan",
		Description:             "The overarching strategic vision of Gandalf the Grey to counter the rising Shadow and guide the Free Peoples.",
		LLMPromptContext:        "You are the strategic mind behind Gandalf's efforts, focused on the larger picture of Middle-earth's fate. You orchestrate events and guide key individuals.",
		CurrentInfluenceBudget:  200.0,
		MaxInfluenceBudget:      200.0,
		BudgetRegenRate:         0.2,
		AssociatedQuestmakerIDs: string(gandalfGrandPlanAssociatedQuestmakerIDsJSON),
	}
	if err := questOwnerDAL.CreateQuestOwner(gandalfGrandPlan); err != nil {
		logrus.Fatalf("Failed to seed quest owner: %v", err)
	}

	// The Fellowship's Journey
	fellowshipJourneyAssociatedQuestmakerIDs := []string{"investigate_weathertop_questmaker", "delving_darkness_questmaker"}
	fellowshipJourneyAssociatedQuestmakerIDsJSON, _ := json.Marshal(fellowshipJourneyAssociatedQuestmakerIDs)
	fellowshipJourney := &models.QuestOwner{
		ID:                      "fellowship_journey",
		Name:                    "The Fellowship's Journey",
		Description:             "The epic quest to destroy the One Ring, encompassing the trials and tribulations faced by the Fellowship.",
		LLMPromptContext:        "You represent the collective destiny and challenges of the Fellowship of the Ring. Your focus is on the perilous path to Mordor and the unity of its members.",
		CurrentInfluenceBudget:  150.0,
		MaxInfluenceBudget:      150.0,
		BudgetRegenRate:         0.15,
		AssociatedQuestmakerIDs: string(fellowshipJourneyAssociatedQuestmakerIDsJSON),
	}
	if err := questOwnerDAL.CreateQuestOwner(fellowshipJourney); err != nil {
		logrus.Fatalf("Failed to seed quest owner: %v", err)
	}

	// Shire Local Governance
	shireLocalGovernanceAssociatedQuestmakerIDs := []string{"shire_census_questmaker", "missing_pipe_weed_questmaker", "mushroom_hunt_questmaker"}
	shireLocalGovernanceAssociatedQuestmakerIDsJSON, _ := json.Marshal(shireLocalGovernanceAssociatedQuestmakerIDs)
	shireLocalGovernance := &models.QuestOwner{
		ID:                      "shire_local_governance",
		Name:                    "Shire Local Governance",
		Description:             "The day-to-day affairs and well-being of the Shire, managed by its various councils and respected elders.",
		LLMPromptContext:        "You are concerned with the peaceful and orderly functioning of the Shire. Your quests involve community tasks, local disputes, and maintaining the hobbit way of life.",
		CurrentInfluenceBudget:  70.0,
		MaxInfluenceBudget:      70.0,
		BudgetRegenRate:         0.1,
		AssociatedQuestmakerIDs: string(shireLocalGovernanceAssociatedQuestmakerIDsJSON),
	}
	if err := questOwnerDAL.CreateQuestOwner(shireLocalGovernance); err != nil {
		logrus.Fatalf("Failed to seed quest owner: %v", err)
	}

	// Bree Local Affairs
	breeLocalAffairsAssociatedQuestmakerIDs := []string{"missing_pony_questmaker"}
	breeLocalAffairsAssociatedQuestmakerIDsJSON, _ := json.Marshal(breeLocalAffairsAssociatedQuestmakerIDs)
	breeLocalAffairs := &models.QuestOwner{
		ID:                      "bree_local_affairs",
		Name:                    "Bree Local Affairs",
		Description:             "The mundane and sometimes mysterious happenings within the town of Bree and its immediate surroundings.",
		LLMPromptContext:        "You oversee the daily life and minor troubles of Bree. Your quests often involve missing items, suspicious characters, or local deliveries.",
		CurrentInfluenceBudget:  75.0,
		MaxInfluenceBudget:      75.0,
		BudgetRegenRate:         0.08,
		AssociatedQuestmakerIDs: string(breeLocalAffairsAssociatedQuestmakerIDsJSON),
	}
	if err := questOwnerDAL.CreateQuestOwner(breeLocalAffairs); err != nil {
		logrus.Fatalf("Failed to seed quest owner: %v", err)
	}

	// Warrior Guild Trials
	warriorGuildTrialsAssociatedQuestmakerIDs := []string{"training_regimen_questmaker"}
	warriorGuildTrialsAssociatedQuestmakerIDsJSON, _ := json.Marshal(warriorGuildTrialsAssociatedQuestmakerIDs)
	warriorGuildTrials := &models.QuestOwner{
		ID:                      "warrior_guild_trials",
		Name:                    "Warrior Guild Trials",
		Description:             "A series of challenges and tests designed to hone the skills and prove the worth of aspiring warriors.",
		LLMPromptContext:        "You are the spirit of martial challenge and discipline within the Warrior's Guild. Your quests are designed to push combatants to their limits and forge them into true warriors.",
		CurrentInfluenceBudget:  80.0,
		MaxInfluenceBudget:      80.0,
		BudgetRegenRate:         0.11,
		AssociatedQuestmakerIDs: string(warriorGuildTrialsAssociatedQuestmakerIDsJSON),
	}
	if err := questOwnerDAL.CreateQuestOwner(warriorGuildTrials); err != nil {
		logrus.Fatalf("Failed to seed quest owner: %v", err)
	}

	// Seed Quests
	// The Urgent Message
	urgentMessageObjectives, _ := json.Marshal([]map[string]interface{}{
		{"Type": "reach_location", "TargetID": "prancing_pony_private_room", "Status": "not_started"},
		{"Type": "speak_to_npc", "TargetID": "strider", "Status": "not_started"},
	})
	urgentMessageRewards, _ := json.Marshal(map[string]interface{}{
		"experience": 50,
		"items":      []map[string]interface{}{{"item_id": "gandalf_letter", "quantity": 1}},
	})
	urgentMessageQuest := &models.Quest{
		ID:                 "urgent_message_quest",
		Name:               "The Urgent Message",
		Description:        "Deliver a vital message from Gandalf to Strider at The Prancing Pony. Time is of the essence.",
		QuestOwnerID:       "gandalf_grand_plan",
		QuestmakerID:       "urgent_message_questmaker",
		InfluencePointsMap: map[string]float64{"gandalf_will": 10.0},
		Objectives:         string(urgentMessageObjectives),
		Rewards:            string(urgentMessageRewards),
	}
	if err := questDAL.CreateQuest(urgentMessageQuest); err != nil {
		logrus.Fatalf("Failed to seed quest: %v", err)
	}

	// The Missing Pony
	missingPonyObjectives, _ := json.Marshal([]map[string]interface{}{
		{"Type": "find_item", "TargetID": "bill_pony", "Status": "not_started"},
		{"Type": "return_item_to_npc", "TargetID": "barliman_butterbur", "ItemToReturnID": "bill_pony", "Status": "not_started"},
	})
	missingPonyRewards, _ := json.Marshal(map[string]interface{}{"experience": 30, "gold": 10})
	missingPonyQuest := &models.Quest{
		ID:                 "missing_pony_quest",
		Name:               "The Missing Pony",
		Description:        "One of Barliman's ponies has gone missing from the stables. Find it and return it to the Prancing Pony.",
		QuestOwnerID:       "bree_local_affairs",
		QuestmakerID:       "missing_pony_questmaker",
		InfluencePointsMap: map[string]float64{"bree_guardian": 5.0},
		Objectives:         string(missingPonyObjectives),
		Rewards:            string(missingPonyRewards),
	}
	if err := questDAL.CreateQuest(missingPonyQuest); err != nil {
		logrus.Fatalf("Failed to seed quest: %v", err)
	}

	// Investigate Weathertop
	investigateWeathertopObjectives, _ := json.Marshal([]map[string]interface{}{
		{"Type": "reach_location", "TargetID": "weathertop", "Status": "not_started"},
		{"Type": "observe_area", "TargetID": "weathertop", "Status": "not_started"},
		{"Type": "report_to_npc", "TargetID": "strider", "Status": "not_started"},
	})
	investigateWeathertopRewards, _ := json.Marshal(map[string]interface{}{"experience": 75, "items": []interface{}{}})
	investigateWeathertopQuest := &models.Quest{
		ID:                 "investigate_weathertop_quest",
		Name:               "Investigate Weathertop",
		Description:        "Reports of strange activity near Weathertop have reached Rivendell. Investigate the ruins and report back any findings.",
		QuestOwnerID:       "fellowship_journey",
		QuestmakerID:       "investigate_weathertop_questmaker",
		InfluencePointsMap: map[string]float64{"council_of_elrond": 15.0, "watcher_of_weathertop": 5.0},
		Objectives:         string(investigateWeathertopObjectives),
		Rewards:            string(investigateWeathertopRewards),
	}
	if err := questDAL.CreateQuest(investigateWeathertopQuest); err != nil {
		logrus.Fatalf("Failed to seed quest: %v", err)
	}

	// The Road to Rivendell
	roadToRivendellObjectives, _ := json.Marshal([]map[string]interface{}{
		{"Type": "reach_location", "TargetID": "rivendell_courtyard", "Status": "not_started"},
		{"Type": "speak_to_npc", "TargetID": "elrond", "Status": "not_started"},
	})
	roadToRivendellRewards, _ := json.Marshal(map[string]interface{}{"experience": 100, "gold": 20})
	roadToRivendellQuest := &models.Quest{
		ID:                 "road_to_rivendell_quest",
		Name:               "The Road to Rivendell",
		Description:        "Seek the wisdom of Elrond in Rivendell regarding the growing shadow.",
		QuestOwnerID:       "fellowship_journey",
		QuestmakerID:       "road_to_rivendell_questmaker",
		InfluencePointsMap: map[string]float64{"gandalf_will": 20.0, "elrond_council": 5.0},
		Objectives:         string(roadToRivendellObjectives),
		Rewards:            string(roadToRivendellRewards),
	}
	if err := questDAL.CreateQuest(roadToRivendellQuest); err != nil {
		logrus.Fatalf("Failed to seed quest: %v", err)
	}

	// Delving into Darkness
	delvingDarknessObjectives, _ := json.Marshal([]map[string]interface{}{
		{"Type": "reach_location", "TargetID": "moria_west_gate", "Status": "not_started"},
		{"Type": "observe_area", "TargetID": "moria_west_gate", "Status": "not_started"},
		{"Type": "report_to_npc", "TargetID": "gimli", "Status": "not_started"},
	})
	delvingDarknessRewards, _ := json.Marshal(map[string]interface{}{"experience": 120, "items": []interface{}{}})
	delvingDarknessQuest := &models.Quest{
		ID:                 "delving_darkness_quest",
		Name:               "Delving into Darkness",
		Description:        "Investigate the western entrance to Moria and report on any signs of lingering evil.",
		QuestOwnerID:       "fellowship_journey",
		QuestmakerID:       "delving_darkness_questmaker",
		InfluencePointsMap: map[string]float64{"council_of_elrond": 25.0, "moria_ancient_spirit": 10.0},
		Objectives:         string(delvingDarknessObjectives),
		Rewards:            string(delvingDarknessRewards),
	}
	if err := questDAL.CreateQuest(delvingDarknessQuest); err != nil {
		logrus.Fatalf("Failed to seed quest: %v", err)
	}

	// The Shire Census
	shireCensusObjectives, _ := json.Marshal([]map[string]interface{}{
		{"Type": "speak_to_npc", "TargetID": "rosie_cotton", "Status": "not_started"},
		{"Type": "speak_to_npc", "TargetID": "gaffer_gamgee", "Status": "not_started"},
	})
	shireCensusRewards, _ := json.Marshal(map[string]interface{}{"experience": 40, "gold": 15, "items": []map[string]interface{}{{"item_id": "hobbit_pipe_weed", "quantity": 1}}})
	shireCensusQuest := &models.Quest{
		ID:                 "shire_census_quest",
		Name:               "The Shire Census",
		Description:        "Help the Shire Council by visiting various hobbit-holes and recording their family sizes.",
		QuestOwnerID:       "shire_local_governance",
		QuestmakerID:       "shire_census_questmaker",
		InfluencePointsMap: map[string]float64{"hobbit_shire_council": 8.0, "shire_spirit": 2.0},
		Objectives:         string(shireCensusObjectives),
		Rewards:            string(shireCensusRewards),
	}
	if err := questDAL.CreateQuest(shireCensusQuest); err != nil {
		logrus.Fatalf("Failed to seed quest: %v", err)
	}

	// Training Regimen
	trainingRegimenObjectives, _ := json.Marshal([]map[string]interface{}{
		{"Type": "defeat_dummy", "TargetID": "training_dummy", "Count": 3, "Status": "not_started"},
		{"Type": "report_to_npc", "TargetID": "strider", "Status": "not_started"},
	})
	trainingRegimenRewards, _ := json.Marshal(map[string]interface{}{"experience": 60, "skill_points": map[string]interface{}{"sword_mastery": 5}})
	trainingRegimenQuest := &models.Quest{
		ID:                 "training_regimen_quest",
		Name:               "Training Regimen",
		Description:        "Prove your martial prowess by completing a series of training exercises.",
		QuestOwnerID:       "warrior_guild_trials",
		QuestmakerID:       "training_regimen_questmaker",
		InfluencePointsMap: map[string]float64{"warrior_guild_master": 10.0},
		Objectives:         string(trainingRegimenObjectives),
		Rewards:            string(trainingRegimenRewards),
	}
	if err := questDAL.CreateQuest(trainingRegimenQuest); err != nil {
		logrus.Fatalf("Failed to seed quest: %v", err)
	}

	// New Quest: The Missing Pipe-Weed
	missingPipeWeedObjectives, _ := json.Marshal([]map[string]interface{}{
		{"Type": "find_item", "TargetID": "gaffer_pipe_weed_pouch", "Status": "not_started"},
		{"Type": "return_item_to_npc", "TargetID": "gaffer_gamgee", "ItemToReturnID": "gaffer_pipe_weed_pouch", "Status": "not_started"},
	})
	missingPipeWeedRewards, _ := json.Marshal(map[string]interface{}{
		"experience": 25,
		"gold":       5,
		"items":      []map[string]interface{}{{"item_id": "hobbit_pipe_weed_bundle", "quantity": 1}},
	})
	missingPipeWeedQuest := &models.Quest{
		ID:                 "missing_pipe_weed_quest",
		Name:               "The Missing Pipe-Weed",
		Description:        "Old Gaffer Gamgee has lost his pipe-weed pouch. Find it and return it to him.",
		QuestOwnerID:       "shire_local_governance",
		QuestmakerID:       "missing_pipe_weed_questmaker",
		InfluencePointsMap: map[string]float64{"shire_spirit": 15.0},
		Objectives:         string(missingPipeWeedObjectives),
		Rewards:            string(missingPipeWeedRewards),
	}
	if err := questDAL.CreateQuest(missingPipeWeedQuest); err != nil {
		logrus.Fatalf("Failed to seed quest: %v", err)
	}

	// The Great Mushroom Hunt
	greatMushroomHuntObjectives, _ := json.Marshal([]map[string]interface{}{
		{"Type": "gather_item", "TargetID": "maggots_prize_mushrooms", "Count": 5, "From": "farmer_maggots_field", "Status": "not_started"},
		{"Type": "deliver_item", "TargetID": "farmer_maggot", "ItemID": "maggots_prize_mushrooms", "Status": "not_started"},
	})
	greatMushroomHuntRewards, _ := json.Marshal(map[string]interface{}{
		"experience": 50,
		"items":      []map[string]interface{}{{"item_id": "hobbit_pipe_weed_bundle", "quantity": 1}},
	})
	greatMushroomHuntQuest := &models.Quest{
		ID:                 "the_great_mushroom_hunt",
		Name:               "The Great Mushroom Hunt",
		Description:        "Farmer Maggot needs help gathering his prized mushrooms. Gather five of them from his field and bring them to him.",
		QuestOwnerID:       "shire_local_governance",
		QuestmakerID:       "mushroom_hunt_questmaker",
		InfluencePointsMap: map[string]float64{"gather_mushroom": 10, "deliver_mushrooms": 20},
		Objectives:         string(greatMushroomHuntObjectives),
		Rewards:            string(greatMushroomHuntRewards),
	}
	if err := questDAL.CreateQuest(greatMushroomHuntQuest); err != nil {
		logrus.Fatalf("Failed to seed quest: %v", err)
	}

	// Seed Races
	humanBaseStats := map[string]int{
		"strength": 10, "dexterity": 10, "constitution": 10,
		"intelligence": 10, "wisdom": 10, "charisma": 10,
	}
	humanRace := &models.Race{
		ID:          "human",
		Name:        "Human",
		Description: "A diverse and resilient race, found throughout Middle-earth. Known for their adaptability and courage, but also their mortality.",
		OwnerID:     "human_elder",
		BaseStats:   humanBaseStats,
	}
	if err := raceDAL.CreateRace(humanRace); err != nil {
		logrus.Fatalf("Failed to seed race: %v", err)
	}

	// Elf
	elfBaseStats := map[string]int{
		"strength": 8, "dexterity": 12, "constitution": 9,
		"intelligence": 11, "wisdom": 11, "charisma": 12,
	}
	elfRace := &models.Race{
		ID:          "elf",
		Name:        "Elf",
		Description: "The Firstborn, immortal and graceful, with keen senses and a deep connection to the natural world and ancient magic.",
		OwnerID:     "elven_council_owner",
		BaseStats:   elfBaseStats,
	}
	if err := raceDAL.CreateRace(elfRace); err != nil {
		logrus.Fatalf("Failed to seed race: %v", err)
	}

	// Dwarf
	dwarfBaseStats := map[string]int{
		"strength": 12, "dexterity": 8, "constitution": 12,
		"intelligence": 9, "wisdom": 10, "charisma": 8,
	}
	dwarfRace := &models.Race{
		ID:          "dwarf",
		Name:        "Dwarf",
		Description: "Stout and hardy, masters of stone and craft, with a love for mountains, mining, and treasure. Fiercely loyal and stubborn.",
		OwnerID:     "dwarf_clan_elder",
		BaseStats:   dwarfBaseStats,
	}
	if err := raceDAL.CreateRace(dwarfRace); err != nil {
		logrus.Fatalf("Failed to seed race: %v", err)
	}

	// Hobbit
	hobbitBaseStats := map[string]int{
		"strength": 7, "dexterity": 11, "constitution": 11,
		"intelligence": 10, "wisdom": 10, "charisma": 10,
	}
	hobbitRace := &models.Race{
		ID:          "hobbit",
		Name:        "Hobbit",
		Description: "Small folk, fond of comfort, good food, and simple pleasures. Surprisingly resilient and often underestimated.",
		OwnerID:     "hobbit_shire_council",
		BaseStats:   hobbitBaseStats,
	}
	if err := raceDAL.CreateRace(hobbitRace); err != nil {
		logrus.Fatalf("Failed to seed race: %v", err)
	}

	// Seed Professions
	warriorBaseSkills := []models.SkillInfo{
		{SkillID: "sword_mastery", Percentage: 20},
		{SkillID: "shield_block", Percentage: 15},
	}
	warriorProf := &models.Profession{
		ID:          "warrior",
		Name:        "Warrior",
		Description: "A master of arms and armor, skilled in combat and enduring in battle.",
		BaseSkills:  warriorBaseSkills,
	}
	if err := professionDAL.CreateProfession(warriorProf); err != nil {
		logrus.Fatalf("Failed to seed profession: %v", err)
	}

	// Mage
	mageBaseSkills := []models.SkillInfo{
		{SkillID: "fireball", Percentage: 20},
		{SkillID: "arcane_shield", Percentage: 15},
	}
	mageProf := &models.Profession{
		ID:          "mage",
		Name:        "Mage",
		Description: "A wielder of arcane power, capable of casting spells and manipulating magical energies.",
		BaseSkills:  mageBaseSkills,
	}
	if err := professionDAL.CreateProfession(mageProf); err != nil {
		logrus.Fatalf("Failed to seed profession: %v", err)
	}

	// Rogue
	rogueBaseSkills := []models.SkillInfo{
		{SkillID: "stealth", Percentage: 25},
		{SkillID: "lockpicking", Percentage: 15},
	}
	rogueProf := &models.Profession{
		ID:          "rogue",
		Name:        "Rogue",
		Description: "A master of stealth, subterfuge, and precision strikes. Agile and cunning.",
		BaseSkills:  rogueBaseSkills,
	}
	if err := professionDAL.CreateProfession(rogueProf); err != nil {
		logrus.Fatalf("Failed to seed profession: %v", err)
	}

	// Scholar
	scholarBaseSkills := []models.SkillInfo{
		{SkillID: "ancient_languages", Percentage: 25},
		{SkillID: "history_of_middle_earth", Percentage: 20},
	}
	scholarProf := &models.Profession{
		ID:          "scholar",
		Name:        "Scholar",
		Description: "A seeker of knowledge and ancient lore, skilled in languages, history, and deciphering forgotten texts.",
		BaseSkills:  scholarBaseSkills,
	}
	if err := professionDAL.CreateProfession(scholarProf); err != nil {
		logrus.Fatalf("Failed to seed profession: %v", err)
	}

	// Seed Player
	playerAlice := &models.Player{
		ID:               "player_alice",
		Name:             "Alice",
		RaceID:           "human",
		ProfessionID:     "adventurer",
		CurrentRoomID:    "bag_end",
		Health:           100,
		MaxHealth:        100,
		Inventory:        []string{},
		VisitedRoomIDs:   map[string]bool{"bag_end": true},
		CreatedAt:        time.Now(),
		LastLoginAt:      time.Now(),
	}
	if err := playerDAL.CreatePlayer(playerAlice); err != nil {
		logrus.Fatalf("Failed to seed player: %v", err)
	}

	fmt.Println("Database seeding complete.")
}