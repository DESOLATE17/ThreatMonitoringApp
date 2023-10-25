package repository

import "threat-monitoring/internal/models"

func (r *Repository) Register(newUser models.User) error {
	return r.db.Create(&newUser).Error
}
