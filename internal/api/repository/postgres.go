package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"threat-monitoring/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(connectionString string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Threat{})
	if err != nil {
		panic("cant migrate db")
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetThreatByID(threatId int) (models.Threat, error) {
	threat := models.Threat{}

	err := r.db.First(&threat, "threat_id = ?", strconv.Itoa(threatId)).Error
	if err != nil {
		return threat, err
	}

	return threat, nil
}

func (r *Repository) GetThreats() ([]models.Threat, error) {
	threats := make([]models.Threat, 0, 4)

	r.db.Where("is_deleted = ?", false).Find(&threats)

	return threats, nil
}

func (r *Repository) CreateProduct(threat models.Threat) error {
	return r.db.Create(threat).Error
}
