package api

import (
	"threat-monitoring/internal/models"
	"time"
)

type Repo interface {
	GetThreatsList(query string) ([]models.Threat, error)
	GetThreatByID(threatId int) (models.Threat, error)
	DeleteThreatByID(threatId int) error
	AddThreat(newThreat models.Threat) error
	UpdateThreat(updateThreat models.Threat) error

	GetMonitoringRequests(status string, startDate, endDate time.Time) ([]models.MonitoringRequest, error)
	GetMonitoringRequestById(requestId int) (models.MonitoringRequest, []models.Threat, error)
	DeleteMonitoringRequest(id int) error
	UpdateMonitoringRequestClient(id int, status string) error
	UpdateMonitoringRequestAdmin(id int, status string) error
	GetMonitoringRequestDraft(userId int) (int, error)

	AddThreatToRequest(request models.MonitoringRequestCreateMessage) error
	DeleteThreatFromRequest(requestId, threatId int) error
}
