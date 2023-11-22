package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"threat-monitoring/internal/api"
	"threat-monitoring/internal/pkg/minio"
)

type Handler struct {
	repo   api.Repo
	minio  minio.Client
	logger *logrus.Entry
}

func NewHandler(repo api.Repo, minio minio.Client, logger *logrus.Logger) *Handler {
	return &Handler{repo: repo, minio: minio, logger: logger.WithField("component", "handler")}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/styles", "./styles")
	r.Static("/image", "./resources")
	// услуги - угрозы
	api := r.Group("/api")
	api.GET("/threats", h.GetThreatsList)
	api.GET("/threats/:id", h.GetThreatById)
	api.DELETE("/threats/:id", h.DeleteThreat)
	api.POST("/threats", h.AddThreat)
	api.PUT("/threats/:id", h.UpdateThreat)
	api.POST("/threats/request", h.AddThreatToRequest)

	// заявки - мониторинг угроз
	api.GET("/monitoring-requests", h.GetMonitoringRequestsList)
	api.GET("/monitoring-requests/:id", h.GetMonitoringRequestById)
	api.DELETE("/monitoring-requests", h.DeleteMonitoringRequest)
	api.PUT("/monitoring-requests/client", h.UpdateMonitoringRequestClient)
	api.PUT("/monitoring-requests/:id/admin", h.UpdateMonitoringRequestAdmin)

	// м-м

	api.DELETE("/monitoring-request-threats/threats/:threatId", h.DeleteThreatFromRequest)

	return r
}
