package models

type Threat struct {
	Id          int `gorm:"primaryKey"`
	Name        string
	Description string
	Image       string
	Count       int
	Price       int
}
