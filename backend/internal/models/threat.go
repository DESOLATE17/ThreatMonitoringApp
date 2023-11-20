package models

type Threat struct {
	ThreatId    int    `gorm:"primaryKey" json:"threatId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Count       int    `json:"count"`
	Price       int    `json:"price"`
}
