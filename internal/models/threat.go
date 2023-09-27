package models

type Threat struct {
	ThreatId    int `gorm:"primaryKey"`
	Name        string
	Description string
	Image       string
	Count       int
	Price       int
}
