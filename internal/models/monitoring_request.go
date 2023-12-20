package models

import "time"

type MonitoringRequest struct {
	RequestId     int       `gorm:"primaryKey" json:"requestId"`
	Status        string    `json:"status"`
	CreationDate  time.Time `json:"creationDate"`
	FormationDate time.Time `json:"formationDate"`
	EndingDate    time.Time `json:"endingDate"`
	AdminId       int       `json:"adminId"`
	CreatorId     int       `json:"userId,omitempty"`
	Creator       string    `json:"creator,omitempty"`
	Receipt       string    `json:"receipt"`
}

type RequestAsyncService struct {
	RequestId int    `gorm:"primaryKey" json:"requestId"`
	Receipt   string `json:"receipt"`
	Token     string `json:"Server-Token"`
}

type MonitoringRequestCreateMessage struct {
	CreatorId int `json:"userId"`
	ThreatId  int `json:"threatId"`
}

type MonitoringRequestsThreats struct {
	Id        int
	RequestId int
	ThreatId  int
}

type NewStatus struct {
	Status string `json:"status"`
}
