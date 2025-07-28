package presentation

// SemanticMessageType defines the type of semantic message.
type SemanticMessageType string

const (
	// Narrative messages
	NarrativeMessage SemanticMessageType = "narrative"
	// Game state updates
	RoomUpdate       SemanticMessageType = "room_update"
	PlayerStatsUpdate SemanticMessageType = "player_stats_update"
	InventoryUpdate  SemanticMessageType = "inventory_update"
	// System messages
	SystemMessage    SemanticMessageType = "system_message"
	ErrorMessage     SemanticMessageType = "error_message"
	// AI-specific messages
	NPCMessage       SemanticMessageType = "npc_message"
	OwnerMessage     SemanticMessageType = "owner_message"
	QuestMessage     SemanticMessageType = "quest_message"
)

// SemanticColorType defines the type of semantic color/style.
type SemanticColorType string

const (
	ColorDefault   SemanticColorType = "default"
	ColorHighlight SemanticColorType = "highlight"
	ColorSuccess   SemanticColorType = "success"
	ColorError     SemanticColorType = "error"
	ColorWarning   SemanticColorType = "warning"
	ColorNarrative SemanticColorType = "narrative"
	ColorNPC       SemanticColorType = "npc"
	ColorPlayer    SemanticColorType = "player"
	ColorItem      SemanticColorType = "item"
	ColorLore      SemanticColorType = "lore"
	ColorQuest     SemanticColorType = "quest"
	ColorOwner     SemanticColorType = "owner"
)

// SemanticMessage represents a universal, semantic JSON message from the server.
type SemanticMessage struct {
	Type    SemanticMessageType `json:"type"`
	Content string              `json:"content"`
	Color   SemanticColorType   `json:"color,omitempty"`
	// Payload can hold additional structured data depending on the message Type
	Payload map[string]interface{} `json:"payload,omitempty"`
}
