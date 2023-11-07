package repository

import (
	"context"
	"threat-monitoring/internal/models"
)

func (r *Repository) SignUp(ctx context.Context, newUser models.User) error {
	return r.db.Create(&newUser).Error
}

func (r *Repository) GetByCredentials(ctx context.Context, user models.User) (models.User, error) {
	err := r.db.First(&user, "login = ? AND password = ?", user.Login, user.Password).Error
	return user, err
}
