package models

type Threat struct {
	ThreatId    int    `gorm:"primaryKey" json:"threatId"`
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Count       int    `json:"count"`
	Price       int    `json:"price"`
}
