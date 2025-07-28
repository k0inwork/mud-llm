package dal

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

// DAL is a struct that holds all the data access layers.
type DAL struct {
	RoomDAL               RoomDALInterface
	ItemDAL               ItemDALInterface
	NpcDAL                NPCDALInterface
	OwnerDAL              OwnerDALInterface
	LoreDAL               LoreDALInterface
	PlayerAccountDAL      PlayerAccountDALInterface
	PlayerCharacterDAL    PlayerCharacterDALInterface
	QuestDAL              QuestDALInterface
	QuestmakerDAL         QuestmakerDALInterface
	QuestOwnerDAL         QuestOwnerDALInterface
	PlayerQuestState      PlayerQuestStateDALInterface
	RaceDAL               RaceDALInterface
	ProfessionDAL         ProfessionDALInterface
	SkillDAL              SkillDALInterface
	PlayerSkillDAL        PlayerSkillDALInterface
	ClassDAL              ClassDALInterface
	PlayerClassDAL        PlayerClassDALInterface
}

// NewDAL creates a new DAL instance with all its sub-DALs.
func NewDAL(db *sql.DB) *DAL {
	newCache := NewCache()
	itemDAL := NewItemDAL(db, newCache)
	return &DAL{
		RoomDAL:               NewRoomDAL(db, newCache),
		ItemDAL:               itemDAL,
		NpcDAL:                NewNPCDAL(db, newCache),
		OwnerDAL:              NewOwnerDAL(db, newCache),
		LoreDAL:               NewLoreDAL(db, newCache),
		PlayerAccountDAL:      NewPlayerAccountDAL(db),
		PlayerCharacterDAL:    NewPlayerCharacterDAL(db, newCache, itemDAL),
		QuestDAL:              NewQuestDAL(db, newCache),
		QuestmakerDAL:         NewQuestmakerDAL(db, newCache),
		QuestOwnerDAL:         NewQuestOwnerDAL(db, newCache),
		PlayerQuestState:      NewPlayerQuestStateDAL(db, newCache),
		RaceDAL:               NewRaceDAL(db, newCache),
		ProfessionDAL:         NewProfessionDAL(db, newCache),
		SkillDAL:              NewSkillDAL(db, newCache),
		PlayerSkillDAL:        NewPlayerSkillDAL(db, newCache),
		ClassDAL:              NewClassDAL(db, newCache),
		PlayerClassDAL:        NewPlayerClassDAL(db, newCache),
	}
}

// InitDB initializes the database and creates tables if they don't exist.
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure connection is established
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create tables
	schema := `
	CREATE TABLE IF NOT EXISTS player_accounts (
		id TEXT PRIMARY KEY NOT NULL,
		username TEXT NOT NULL UNIQUE,
		hashed_password TEXT NOT NULL,
		email TEXT UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_login_at TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS player_characters (
		id TEXT PRIMARY KEY NOT NULL,
		player_account_id TEXT NOT NULL,
		name TEXT NOT NULL UNIQUE,
		race_id TEXT NOT NULL,
		profession_id TEXT NOT NULL,
		current_room_id TEXT NOT NULL,
		health INTEGER NOT NULL,
		max_health INTEGER NOT NULL,
		inventory TEXT NOT NULL,
		visited_room_ids TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_played_at TIMESTAMP,
		FOREIGN KEY (player_account_id) REFERENCES player_accounts(id)
	);

	CREATE TABLE IF NOT EXISTS Rooms (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		exits JSON NOT NULL,
		owner_id TEXT,
		territory_id TEXT,
		properties JSON,
		perception_biases JSON
	);

	CREATE TABLE IF NOT EXISTS Items (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		type TEXT NOT NULL,
		properties JSON
	);

	CREATE TABLE IF NOT EXISTS NPCs (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		current_room_id TEXT NOT NULL,
		health INTEGER NOT NULL,
		max_health INTEGER NOT NULL,
		inventory JSON NOT NULL,
		owner_ids JSON NOT NULL,
		memories_about_players JSON NOT NULL,
		personality_prompt TEXT NOT NULL,
		available_tools JSON NOT NULL,
		behavior_state JSON,
		reaction_threshold INTEGER NOT NULL DEFAULT 0,
		race_id TEXT,
		profession_id TEXT
	);

	CREATE TABLE IF NOT EXISTS Owners (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		monitored_aspect TEXT NOT NULL,
		associated_id TEXT NOT NULL,
		llm_prompt_context TEXT NOT NULL,
		memories_about_players JSON NOT NULL,
		current_influence_budget REAL NOT NULL,
		max_influence_budget REAL NOT NULL,
		budget_regen_rate REAL NOT NULL,
		available_tools JSON NOT NULL,
		initiated_quests JSON NOT NULL,
		reaction_threshold INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS Quests (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		quest_owner_id TEXT NOT NULL,
		questmaker_id TEXT NOT NULL,
		influence_points_map JSON NOT NULL,
		objectives JSON NOT NULL,
		rewards JSON NOT NULL
	);

	CREATE TABLE IF NOT EXISTS Questmakers (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		llm_prompt_context TEXT NOT NULL,
		current_influence_budget REAL NOT NULL,
		max_influence_budget REAL NOT NULL,
		budget_regen_rate REAL NOT NULL,
		memories_about_players JSON NOT NULL,
		available_tools JSON NOT NULL,
		reaction_threshold INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS QuestOwners (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		llm_prompt_context TEXT NOT NULL,
		current_influence_budget REAL NOT NULL,
		max_influence_budget REAL NOT NULL,
		budget_regen_rate REAL NOT NULL,
		associated_questmaker_ids JSON NOT NULL
	);

	CREATE TABLE IF NOT EXISTS PlayerQuestStates (
		player_id TEXT NOT NULL,
		quest_id TEXT NOT NULL,
		current_progress JSON NOT NULL,
		last_action_timestamp TIMESTAMP NOT NULL,
		questmaker_influence_accumulated REAL NOT NULL,
		status TEXT NOT NULL,
		PRIMARY KEY (player_id, quest_id)
	);

	CREATE TABLE IF NOT EXISTS Races (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		owner_id TEXT,
		base_stats JSON NOT NULL,
		perception_biases JSON
	);

	CREATE TABLE IF NOT EXISTS Professions (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		base_skills JSON NOT NULL,
		perception_biases JSON
	);

	CREATE TABLE IF NOT EXISTS Lore (
		id TEXT PRIMARY KEY NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		scope TEXT NOT NULL,
		associated_id TEXT
	);

	CREATE TABLE IF NOT EXISTS Skills (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		type TEXT NOT NULL,
		category TEXT,
		associated_class_id TEXT,
		granted_by_entity_type TEXT,
		granted_by_entity_id TEXT,
		effects JSON NOT NULL,
		cost INTEGER,
		cooldown INTEGER,
		min_class_level INTEGER
	);

	CREATE TABLE IF NOT EXISTS PlayerSkills (
		player_id TEXT NOT NULL,
		skill_id TEXT NOT NULL,
		percentage INTEGER NOT NULL DEFAULT 0,
		granted_by_entity_type TEXT,
		granted_by_entity_id TEXT,
		PRIMARY KEY (player_id, skill_id)
	);

	CREATE TABLE IF NOT EXISTS Classes (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		total_levels INTEGER NOT NULL DEFAULT 5,
		parent_class_id TEXT,
		associated_entity_type TEXT,
		associated_entity_id TEXT,
		level_up_rewards JSON NOT NULL
	);

	CREATE TABLE IF NOT EXISTS PlayerClasses (
		player_id TEXT NOT NULL,
		class_id TEXT NOT NULL,
		level INTEGER NOT NULL DEFAULT 1,
		experience INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (player_id, class_id)
	);

	CREATE TABLE IF NOT EXISTS ActionSignificanceConfig (
		action_type TEXT PRIMARY KEY NOT NULL,
		score INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS LLMToolDefinitions (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		parameters_schema JSON NOT NULL,
		base_cost REAL NOT NULL,
		entity_type TEXT NOT NULL
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		logrus.Fatalf("Error creating tables: %v", err)
	}

	return db, nil
}