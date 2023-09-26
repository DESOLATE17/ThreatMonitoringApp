package main

import (
	log "github.com/sirupsen/logrus"
	"threat-monitoring/internal/api/handler"
	"threat-monitoring/internal/api/repository"
	"threat-monitoring/internal/pkg"
)

func main() {
	dsn, err := pkg.GetConnectionString()
	if err != nil {
		log.Error(err)
	}
	log.Info(dsn)

	repo, err := repository.NewRepository(dsn)
	if err != nil {
		log.Error(err)
	}

	handler := handler.NewHandler(repo)
	r := handler.InitRoutes()
	r.Run()
}
