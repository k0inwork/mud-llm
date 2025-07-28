package main

import (
	"net"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"mud/internal/dal"

	"mud/internal/game/actionsignificance"
	"mud/internal/game/events"
	"mud/internal/game/globalobserver"
	"mud/internal/game/perception"
	"mud/internal/game/sentiententitymanager"
	"mud/internal/llm"
	"mud/internal/presentation"
	"mud/internal/server"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	// Configure Logrus for structured logging
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// Initialize database
	db, err := dal.InitDB("./mud.db")
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	dals := dal.NewDAL(db)
	dal.SeedData(db)

	logrus.Info("Successfully connected to SQLite database and ensured tables exist.")

	// Pre-warm the room cache
	rooms, err := dals.RoomDAL.GetAllRooms()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm room cache: %v", err)
	}
	roomMap := make(map[string]interface{})
	for _, room := range rooms {
		roomMap[room.ID] = room
	}
	dals.RoomDAL.Cache().SetMany(roomMap, 300) // Cache for 5 minutes

	// Pre-warm the item cache
	items, err := dals.ItemDAL.GetAllItems()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm item cache: %v", err)
	}
	itemMap := make(map[string]interface{})
	for _, item := range items {
		itemMap[item.ID] = item
	}
	dals.ItemDAL.Cache().SetMany(itemMap, 300) // Cache for 5 minutes

	// Pre-warm the NPC cache
	npcs, err := dals.NpcDAL.GetAllNPCs()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm NPC cache: %v", err)
	}
	npcMap := make(map[string]interface{})
	for _, npc := range npcs {
		npcMap[npc.ID] = npc
	}
	dals.NpcDAL.Cache().SetMany(npcMap, 300) // Cache for 5 minutes

	// Pre-warm the owner cache
	owners, err := dals.OwnerDAL.GetAllOwners()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm owner cache: %v", err)
	}
	ownerMap := make(map[string]interface{})
	for _, owner := range owners {
		ownerMap[owner.ID] = owner
	}
	dals.OwnerDAL.Cache().SetMany(ownerMap, 300) // Cache for 5 minutes

	// Pre-warm the lore cache
	lores, err := dals.LoreDAL.GetAllLore()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm lore cache: %v", err)
	}
	loreMap := make(map[string]interface{})
	for _, lore := range lores {
		loreMap[lore.ID] = lore
	}
	dals.LoreDAL.Cache().SetMany(loreMap, 300) // Cache for 5 minutes

	// Pre-warm the player character cache
	characters, err := dals.PlayerCharacterDAL.GetAllCharacters()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm player character cache: %v", err)
	}
	characterMap := make(map[string]interface{})
	for _, character := range characters {
		characterMap[character.ID] = character
	}
	dals.PlayerCharacterDAL.Cache().SetMany(characterMap, 300) // Cache for 5 minutes

	// Pre-warm the player quest state cache
	playerQuestStates, err := dals.PlayerQuestState.GetAllPlayerQuestStates()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm player quest state cache: %v", err)
	}
	playerQuestStateMap := make(map[string]interface{})
	for _, pqs := range playerQuestStates {
		playerQuestStateMap[pqs.PlayerID+"-"+pqs.QuestID] = pqs // Composite key for cache
	}
	dals.PlayerQuestState.Cache().SetMany(playerQuestStateMap, 300) // Cache for 5 minutes

	// Pre-warm the quest cache
	quests, err := dals.QuestDAL.GetAllQuests()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm quest cache: %v", err)
	}
	questMap := make(map[string]interface{})
	for _, quest := range quests {
		questMap[quest.ID] = quest
	}
	dals.QuestDAL.Cache().SetMany(questMap, 300) // Cache for 5 minutes

	// Pre-warm the questmaker cache
	questmakers, err := dals.QuestmakerDAL.GetAllQuestmakers()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm questmaker cache: %v", err)
	}
	questmakerMap := make(map[string]interface{})
	for _, qm := range questmakers {
		questmakerMap[qm.ID] = qm
	}
	dals.QuestmakerDAL.Cache().SetMany(questmakerMap, 300) // Cache for 5 minutes

	// Pre-warm the race cache
	races, err := dals.RaceDAL.GetAllRaces()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm race cache: %v", err)
	}
	raceMap := make(map[string]interface{})
	for _, race := range races {
		raceMap[race.ID] = race
	}
	dals.RaceDAL.Cache().SetMany(raceMap, 300) // Cache for 5 minutes

	// Pre-warm the profession cache
	professions, err := dals.ProfessionDAL.GetAllProfessions()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm profession cache: %v", err)
	}
	professionMap := make(map[string]interface{})
	for _, prof := range professions {
		professionMap[prof.ID] = prof
	}
	dals.ProfessionDAL.Cache().SetMany(professionMap, 300) // Cache for 5 minutes

	// Pre-warm the skill cache
	skills, err := dals.SkillDAL.GetAllSkills()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm skill cache: %v", err)
	}
	skillMap := make(map[string]interface{})
	for _, skill := range skills {
		skillMap[skill.ID] = skill
	}
	dals.SkillDAL.Cache().SetMany(skillMap, 300) // Cache for 5 minutes

	// Pre-warm the class cache
	classes, err := dals.ClassDAL.GetAllClasses()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm class cache: %v", err)
	}
	classMap := make(map[string]interface{})
	for _, class := range classes {
		classMap[class.ID] = class
	}
	dals.ClassDAL.Cache().SetMany(classMap, 300) // Cache for 5 minutes

	// Pre-warm the quest owner cache
	questOwners, err := dals.QuestOwnerDAL.GetAllQuestOwners()
	if err != nil {
		logrus.Fatalf("Failed to pre-warm quest owner cache: %v", err)
	}
	questOwnerMap := make(map[string]interface{})
	for _, qo := range questOwners {
		questOwnerMap[qo.ID] = qo
	}
	dals.QuestOwnerDAL.Cache().SetMany(questOwnerMap, 300) // Cache for 5 minutes

	// Initialize LLM Service
	llmClient := llm.NewClient()
	llmService := llm.NewLLMService(llmClient, dals)

	// Initialize Tool Dispatcher
	toolDispatcher := server.NewToolDispatcher(dals)

	// Initialize Event Bus
	eventBus := events.NewEventBus()

	// Initialize Perception Filter
	perceptionFilter := perception.NewPerceptionFilter(dals.RoomDAL, dals.RaceDAL, dals.ProfessionDAL)

	// Initialize Sentient Entity Manager
	telnetRenderer := presentation.NewTelnetRenderer()
	sentientEntityManager := sentiententitymanager.NewSentientEntityManager(llmService, dals.NpcDAL, dals.OwnerDAL, dals.QuestmakerDAL, toolDispatcher, telnetRenderer, eventBus)

	// Initialize Action Significance Monitor
	actionMonitor := actionsignificance.NewMonitor(eventBus, perceptionFilter, dals.NpcDAL, dals.OwnerDAL, dals.QuestmakerDAL, sentientEntityManager)
	actionMonitorEventChannel := make(chan interface{}, 500)
	eventBus.Subscribe(events.ActionEventType, actionMonitorEventChannel)
	go func() {
		for event := range actionMonitorEventChannel {
			if ae, ok := event.(*events.ActionEvent); ok {
				actionMonitor.HandleActionEvent(ae)
			} else {
				logrus.Errorf("main: received unexpected event type on ActionEventType: %T", event)
			}
		}
	}()

	// Initialize Global Observer Manager
	globalObserverManager := globalobserver.NewGlobalObserverManager(eventBus, perceptionFilter, dals.OwnerDAL, dals.RaceDAL, dals.ProfessionDAL)
	globalObserverEventChannel := make(chan interface{}, 100)
	eventBus.Subscribe(events.ActionEventType, globalObserverEventChannel)
	go func() {
		for event := range globalObserverEventChannel {
			if ae, ok := event.(*events.ActionEvent); ok {
				go globalObserverManager.HandleActionEvent(ae)
			} else {
				logrus.Errorf("main: received unexpected event type on ActionEventType for GlobalObserverManager: %T", event)
			}
		}
	}()

	var wg sync.WaitGroup

	// Start Telnet server in a goroutine
	listener, err := net.Listen("tcp", ":4000") // Listen on a specific port for the main server
	if err != nil {
		logrus.Fatalf("Failed to listen on port 4000: %v", err)
	}
	telnetServer := server.NewTelnetServer(listener, telnetRenderer, eventBus, dals, llmService)
	wg.Add(1)
	go func() {
		defer wg.Done()
		telnetServer.Start()
	}()

	// Start Admin Web server in a goroutine
	adminWebServer := server.NewAdminWebServer("8080", db) // Using port 8080 for admin
	wg.Add(1)
	go func() {
		defer wg.Done()
		adminWebServer.Start()
	}()

	logrus.Info("GoMUD server starting...")

	// Keep the main goroutine alive until all servers are done
	wg.Wait()
}