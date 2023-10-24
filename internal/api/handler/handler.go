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
	r.GET("/threats", h.GetThreatsList)
	r.GET("/threats/:id", h.GetThreatById)
	r.DELETE("/threats/:id", h.DeleteThreat)
	r.POST("/threats", h.AddThreat)
	r.PUT("/threats/:id", h.UpdateThreat)
	r.POST("/threats/request", h.AddThreatToRequest)

	// заявки - мониторинг угроз
	r.GET("/monitoring-requests", h.GetMonitoringRequestsList)
	r.GET("/monitoring-requests/:id", h.GetMonitoringRequestById)
	r.DELETE("/monitoring-requests", h.DeleteMonitoringRequest)
	r.PUT("/monitoring-requests/client", h.UpdateMonitoringRequestClient)
	r.PUT("/monitoring-requests/admin", h.UpdateMonitoringRequestAdmin)

	// м-м

	r.DELETE("/monitoring-request-threats/threats/:threatId", h.DeleteThreatFromRequest)

	return r
}
