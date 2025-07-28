package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mud/internal/dal"
	"mud/internal/game/actionsignificance"
	"mud/internal/game/events"
	"mud/internal/game/globalobserver"
	"mud/internal/game/perception"
	"mud/internal/game/sentiententitymanager"
	"mud/internal/llm"
	"mud/internal/mocks"
	"mud/internal/models"
)

const (
	telnetServerAddress = "localhost:4000"
)

// setupTestEnvironment initializes a test database, DALs, mock LLM service, and Telnet server.
// It returns the TelnetServer instance and a cleanup function.
func setupTestEnvironment(t *testing.T) (*TelnetServer, *mocks.TestRenderer, string, func()) {
	// 1. Setup a temporary SQLite database
	dbPath := fmt.Sprintf("./test_mud_%s.db", uuid.New().String()[:8])
	db, err := dal.InitDB(dbPath)
	assert.NoError(t, err, "Failed to initialize test database")

	// 2. Seed the database
	dal.SeedData(db)

	// 3. Create DALs
	dals := dal.NewDAL(db)

	// 4. Create Mock LLM Service
	mockLLMService := &mocks.MockLLMService{
		ProcessActionFunc: func(ctx context.Context, entity interface{}, player *models.PlayerCharacter, playerAction string) (*llm.InnerLLMResponse, error) {
			t.Logf("Mock LLM received action: '%s' for entity: %+v", playerAction, entity)

			// Customize mock responses based on entity and action
			entityID := ""
			switch e := entity.(type) {
			case *models.NPC:
				entityID = e.ID
			case *models.Owner:
				entityID = e.ID
			case *models.Questmaker:
				entityID = e.ID
			}

			if entityID == "farmer_maggot" && strings.Contains(playerAction, "talk") {
				return &llm.InnerLLMResponse{
					Narrative: "Ah, a new face! Have you seen my prize mushrooms? Those rascals keep wandering off. If you could gather 5 of them from my field, I'd be much obliged!",
					ToolCalls: []llm.ToolCall{},
				}, nil
			}
			if entityID == "farmer_maggot" && strings.Contains(playerAction, "give mushrooms") {
				return &llm.InnerLLMResponse{
					Narrative: "Splendid! You've found them all! Here's a little something for your trouble.",
					ToolCalls: []llm.ToolCall{},
				}, nil
			}
			if entityID == "mushroom_hunt_questmaker" && strings.Contains(playerAction, "gather") {
				return &llm.InnerLLMResponse{
					Narrative: "Excellent! You've gathered a mushroom. Keep up the good work!",
					ToolCalls: []llm.ToolCall{},
				}, nil
			}
			if entityID == "shire_council_owner" && strings.Contains(playerAction, "deliver_item") {
				return &llm.InnerLLMResponse{
					Narrative: "Well done, young one! Your efforts have brought great joy to Farmer Maggot and the Shire. Your reputation here grows!",
					ToolCalls: []llm.ToolCall{},
				}, nil
			}

			return &llm.InnerLLMResponse{
				Narrative: fmt.Sprintf("Mock LLM response for action '%s' by %s.", playerAction, player.Name),
				ToolCalls: []llm.ToolCall{},
			}, nil
		},
		AnalyzeResponseFunc: func(ctx context.Context, narrative string, query string) (float64, error) {
			// Always return a default score for analysis in tests
			return 75.0, nil
		},
	}

	// 5. Initialize Event Bus
	eventBus := events.NewEventBus()

	// 6. Initialize Perception Filter
	perceptionFilter := perception.NewPerceptionFilter(dals.RoomDAL, dals.RaceDAL, dals.ProfessionDAL)

	// 7. Initialize Tool Dispatcher
	toolDispatcher := NewToolDispatcher(dals)

	// 8. Initialize Telnet Renderer
	telnetRenderer := mocks.NewTestRenderer()

	// 9. Initialize Sentient Entity Manager
	sentientEntityManager := sentiententitymanager.NewSentientEntityManager(mockLLMService, dals.NpcDAL, dals.OwnerDAL, dals.QuestmakerDAL, toolDispatcher, telnetRenderer, eventBus)

	// 10. Initialize Action Significance Monitor
	actionMonitor := actionsignificance.NewMonitor(eventBus, perceptionFilter, dals.NpcDAL, dals.OwnerDAL, dals.QuestmakerDAL, sentientEntityManager)
	actionMonitorEventChannel := make(chan interface{}, 1000)
	eventBus.Subscribe(events.ActionEventType, actionMonitorEventChannel)
	go func() {
		for event := range actionMonitorEventChannel {
			if ae, ok := event.(*events.ActionEvent); ok {
				go actionMonitor.HandleActionEvent(ae)
			}
		}
	}()

	// 11. Initialize Global Observer Manager
	globalObserverManager := globalobserver.NewGlobalObserverManager(eventBus, perceptionFilter, dals.OwnerDAL, dals.RaceDAL, dals.ProfessionDAL)
	globalObserverEventChannel := make(chan interface{}, 1000)
	eventBus.Subscribe(events.ActionEventType, globalObserverEventChannel)
	go func() {
		for event := range globalObserverEventChannel {
			if ae, ok := event.(*events.ActionEvent); ok {
				go globalObserverManager.HandleActionEvent(ae)
			}
		}
	}()

	// 12. Create Telnet Server
	listener, err := net.Listen("tcp", ":0") // Listen on a random available port
	assert.NoError(t, err, "Failed to listen on a random port")
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	telnetServer := NewTelnetServer(listener, telnetRenderer, eventBus, dals, mockLLMService)

	// Start Telnet server in a goroutine
	go telnetServer.Start()

	// Wait for the server to be ready
	<-telnetServer.Ready

	cleanup := func() {
		// Close channels to terminate goroutines
		close(actionMonitorEventChannel)
		close(globalObserverEventChannel)
		// Close the listener to stop the server
		listener.Close()
		// Close the database
		db.Close()
		// Remove the temporary database file
		os.Remove(dbPath)
	}

	return telnetServer, telnetRenderer, port, cleanup
}

// TestTelnetServer_FullFlow tests the entire user flow from connection to in-game.
func TestTelnetServer_FullFlow(t *testing.T) {
	_, renderer, port, cleanup := setupTestEnvironment(t)
	_ = port // Port is not directly used in this test after initial connection
	defer cleanup()

	conn, err := net.Dial("tcp", "localhost:"+port)
	assert.NoError(t, err, "Failed to connect to Telnet server")
	defer conn.Close()

	// --- Welcome and Account Creation ---
	assertEventuallyContains(t, renderer, `[system_message] 
Welcome to GoMUD!
1. Login
2. Create Account

`)
	write(t, conn, "2") // Create Account

    assertEventuallyContains(t, renderer, "[system_message] Enter desired username: \n")
    testUsername := fmt.Sprintf("testuser_%s", uuid.New().String()[:8])
    write(t, conn, testUsername)

    assertEventuallyContains(t, renderer, "[system_message] Enter password: \n")
    write(t, conn, "password123")

    assertEventuallyContains(t, renderer, "[system_message] Enter email (optional): \n")
    testEmail := fmt.Sprintf("test_%s@example.com", uuid.New().String()[:8])
    write(t, conn, testEmail)

    assertEventuallyContains(t, renderer, "[system_message] Account created successfully!\n")
    assertEventuallyContains(t, renderer, "[system_message] \n--- Character Selection ---\nYou have no characters.\n\nType a number to select a character, or 'new' to create one.\n\n")

    write(t, conn, "new")

    assertEventuallyContains(t, renderer, "[system_message] Enter character name: \n")
    testCharName := "TestChar"
    write(t, conn, testCharName)

    // --- In-Game ---
    assertEventuallyContains(t, renderer, fmt.Sprintf("[system_message] Welcome, %s!\n", testCharName))
    assertEventuallyContains(t, renderer, "[room_update] \n--- Bag End, Hobbiton ---\nA cozy hobbit-hole, warm and inviting, with a round green door. The smell of pipe-weed and fresh baking lingers in the air. A path leads east.\nExits: east ()\nNPCs present: Frodo Baggins, Samwise Gamgee\n\n")

    // Send a command
    command := "look"
    write(t, conn, command)

    // Read echo response
    assertEventuallyContains(t, renderer, "[system_message] You typed: look\n")
}

// TestTelnetServer_QuestingFlow tests a basic questing scenario.
func TestTelnetServer_QuestingFlow(t *testing.T) {
	_, renderer, port, cleanup := setupTestEnvironment(t)
	defer cleanup()

	conn, err := net.Dial("tcp", "localhost:"+port)
	assert.NoError(t, err, "Failed to connect to Telnet server")
	defer conn.Close()

	// --- Account and Character Creation ---
	assertEventuallyContains(t, renderer, `[system_message] 
Welcome to GoMUD!
1. Login
2. Create Account

`)
	write(t, conn, "2") // Create Account

	assertEventuallyContains(t, renderer, "[system_message] Enter desired username: \n")
	testUsername := fmt.Sprintf("questuser_%s", uuid.New().String()[:8])
	write(t, conn, testUsername)

	assertEventuallyContains(t, renderer, "[system_message] Enter password: \n")
	write(t, conn, "password123")

	assertEventuallyContains(t, renderer, "[system_message] Enter email (optional): \n")
	testEmail := fmt.Sprintf("quest_%s@example.com", uuid.New().String()[:8])
	write(t, conn, testEmail)

	assertEventuallyContains(t, renderer, "[system_message] Account created successfully!\n")
	assertEventuallyContains(t, renderer, "[system_message] \n--- Character Selection ---\nYou have no characters.\n\nType a number to select a character, or 'new' to create one.\n\n")

	write(t, conn, "new")

	assertEventuallyContains(t, renderer, "[system_message] Enter character name: \n")
	testCharName := "QuestPlayer"
	write(t, conn, testCharName)

	assertEventuallyContains(t, renderer, fmt.Sprintf("[system_message] Welcome, %s!\n", testCharName))

	// --- Questing Flow ---
	// Initial room description (Bag End)
	assertEventuallyContains(t, renderer, "[room_update] \n--- Bag End, Hobbiton ---\nA cozy hobbit-hole, warm and inviting, with a round green door. The smell of pipe-weed and fresh baking lingers in the air. A path leads east.\nExits: east ()\nNPCs present: Frodo Baggins, Samwise Gamgee\n\n")

	// Move to Hobbiton Path to find Farmer Maggot
	write(t, conn, "east")
	assertEventuallyContains(t, renderer, "[system_message] You typed: east\n")
	assertEventuallyContains(t, renderer, "[room_update] \n--- Hobbiton Path ---")

	// Talk to Farmer Maggot to initiate quest
	write(t, conn, "talk farmer_maggot")
	assertEventuallyContains(t, renderer, "[system_message] You try to talk to farmer_maggot.\n")
	assertEventuallyContains(t, renderer, "[narrative] Farmer Maggot says: Ah, a new face! Have you seen my prize mushrooms? Those rascals keep wandering off. If you could gather 5 of them from my field, I'd be much obliged!\n")

	// Gather mushrooms (5 times)
	for i := 0; i < 5; i++ {
		write(t, conn, "gather mushrooms")
		assertEventuallyContains(t, renderer, "[system_message] You attempt to gather mushrooms.\n")
		// Expect a message from the questmaker/owner about progress
		assertEventuallyContains(t, renderer, "[narrative] The Great Mushroom Hunt Controller says: Excellent! You've gathered a mushroom. Keep up the good work!\n")
	}

	// Give mushrooms to Farmer Maggot to complete quest
	write(t, conn, "give mushrooms to farmer_maggot")
	assertEventuallyContains(t, renderer, "[system_message] You try to give mushrooms to farmer_maggot.\n")
	assertEventuallyContains(t, renderer, "[narrative] Farmer Maggot says: Splendid! You've found them all! Here's a little something for your trouble.\n")
	assertEventuallyContains(t, renderer, "[narrative] The Shire Council says: Well done, young one! Your efforts have brought great joy to Farmer Maggot and the Shire. Your reputation here grows!\n")

	// Send a final command to ensure server is still responsive
	write(t, conn, "look")
	assertEventuallyContains(t, renderer, `[system_message] You typed: look`)
}

// Helper functions for testing

func assertEventuallyContains(t *testing.T, renderer *mocks.TestRenderer, expected string) {
	assert.Eventually(t, func() bool {
		return renderer.ContainsMessage(expected)
	}, 5*time.Second, 100*time.Millisecond, "Expected to find message:\n%s\n\nActual messages:\n%s", expected, renderer.AllMessages())
}

func write(t *testing.T, conn net.Conn, content string) {
	_, err := fmt.Fprintf(conn, "%s\n", content)
	assert.NoError(t, err, "Failed to write to connection")
	time.Sleep(500 * time.Millisecond) // Give server a moment to process
}

