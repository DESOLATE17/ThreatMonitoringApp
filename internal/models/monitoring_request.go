package models

import "time"

type MonitoringRequest struct {
	RequestId     int `gorm:"primaryKey"`
	Status        string
	CreationDate  time.Time
	FormationDate time.Time
	EndingDate    time.Time
	AdminId       int
}
