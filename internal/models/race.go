package models

// Race defines available player races.
type Race struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID          string                 `json:"owner_id"`     // ID of the Owner associated with this race
	BaseStats        map[string]int         `json:"base_stats"`   // Map of base stats for the race
	PerceptionBiases map[string]float64     `json:"perception_biases"` // Map of perception biases, e.g., {"magic": -0.3}
}