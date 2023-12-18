package api

import (
	"context"
	"threat-monitoring/internal/models"
	"time"
)

type Repo interface {
	GetThreatsList(query, lowPrice, highPrice string) ([]models.Threat, error)
	GetThreatByID(threatId int) (models.Threat, error)
	DeleteThreatByID(threatId int) error
	AddThreat(newThreat models.Threat) error
	UpdateThreat(updateThreat models.Threat) error
	AddThreatToRequest(request models.MonitoringRequestCreateMessage) error
	DeleteThreatImage(threatId int) string

	GetMonitoringRequests(status string, startDate, endDate time.Time, userId int, isAdmin bool) ([]models.MonitoringRequest, error)
	GetMonitoringRequestById(requestId int, userId int, isAdmin bool) (models.MonitoringRequest, []models.Threat, error)
	DeleteMonitoringRequest(id int) error
	UpdateMonitoringRequestClient(id int, status string) error
	UpdateMonitoringRequestAdmin(adminId int, requestId int, status string) error
	GetMonitoringRequestDraft(userId int) (int, error)
	SavePayment(monitoringRequest models.MonitoringRequest) error

	DeleteThreatFromRequest(userId, threatId int) (models.MonitoringRequest, []models.Threat, error)

	SignUp(ctx context.Context, newUser models.User) error
	GetByCredentials(ctx context.Context, user models.User) (models.User, error)
}
