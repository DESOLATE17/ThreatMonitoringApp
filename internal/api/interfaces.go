package api

import "threat-monitoring/internal/models"

type Repo interface {
	GetThreats() ([]models.Threat, error)
	GetThreatByID(threatId int) (models.Threat, error)
}
