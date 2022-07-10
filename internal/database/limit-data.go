package database

import "time"

type LimitData struct {
	ID                  string    `json:"id" sql:"id"`
	UserID              string    `json:"user_id" sql:"userid"`
	NumOfSendMail       int       `json:"num_of_send_mail" sql:"numofsendmail"`
	NumOfChangePassword int       `json:"num_of_change_password" sql:"numofchangepassword"`
	NumOfLogin          int       `json:"num_of_login" sql:"numoflogin"`
	CreatedAt           time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt           time.Time `json:"updatedat" sql:"updatedat"`
}

// MultiRatioData is the data structure for multiratio table
type MultiRatioData struct {
	WaterRatio int `json:"water_ratio" sql:"waterratio"`
	LightRatio int `json:"light_ratio" sql:"lightratio"`
	SeedRatio  int `json:"seed_ratio" sql:"seedratio"`
}

// EarnScore is the data structure for earnscore table
type EarnScore struct {
	UserID     string    `json:"user_id" sql:"userid"`
	WaterScore int       `json:"water_score" sql:"waterscore"`
	LightScore int       `json:"light_score" sql:"lightscore"`
	SeedScore  int       `json:"seed_score" sql:"seedscore"`
	CreatedAt  time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt  time.Time `json:"updatedat" sql:"updatedat"`
}
