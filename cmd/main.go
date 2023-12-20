package main

import (
	"github.com/sirupsen/logrus"
	"threat-monitoring/internal/api/handler"
	"time"
)

// @title ThreatMonitoringApp
// @version 1.0
// @description App for serving threats monitoring requests

// @host localhost:8080
// @schemes http
// @BasePath /
func main() {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	handler := handler.NewHandler(logger)
	r := handler.InitRoutes()
	r.Run("192.168.160.12:8080")
}
