package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"threat-monitoring/internal/api/handler"
	"threat-monitoring/internal/api/repository"
	minio "threat-monitoring/internal/pkg/minio"
	"time"
)

func main() {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	vp := viper.New()
	if err := initConfig(vp); err != nil {
		logger.Fatalf("error initializing configs: %s", err.Error())
	}

	repo, err := repository.NewRepository(logger, vp)
	if err != nil {
		logger.Error(err)
	}

	minioConfig := minio.InitConfig(vp)

	minioClient, err := minio.NewMinioClient(context.Background(), minioConfig, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	handler := handler.NewHandler(repo, minioClient, logger)
	r := handler.InitRoutes()
	r.Run()
}

func initConfig(vp *viper.Viper) error {
	vp.AddConfigPath("./config")
	vp.SetConfigName("config")

	return vp.ReadInConfig()
}
