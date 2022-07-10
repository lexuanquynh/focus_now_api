package database

import "time"

type PassworUsers struct {
	ID        string    `json:"id" sql:"id"`
	UserID    string    `json:"user_id" sql:"userid"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	CreatedAt time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt time.Time `json:"updatedat" sql:"updatedat"`
}

// EarnScoreResponse is the data structure for earnscore table
type EarnScoreResponse struct {
	WaterScore int `json:"water_score"`
	LightScore int `json:"light_score"`
	SeedScore  int `json:"seed_score"`
}
