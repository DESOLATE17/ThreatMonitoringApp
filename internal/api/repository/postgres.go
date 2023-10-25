package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"threat-monitoring/internal/models"
)

type Repository struct {
	db     *gorm.DB
	logger *logrus.Entry
}

func NewRepository(logger *logrus.Logger, vp *viper.Viper) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(vp.GetString("db.connection_string")), &gorm.Config{})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	//Migrate the schema
	err = db.AutoMigrate(&models.Threat{})
	err = db.AutoMigrate(&models.User{})
	err = db.AutoMigrate(&models.MonitoringRequest{})
	err = db.AutoMigrate(&models.MonitoringRequestsThreats{})

	if err != nil {
		logger.Fatal("cant migrate db")
	}

	return &Repository{
		db:     db,
		logger: logger.WithField("component", "repository"),
	}, nil
}
