package repository

import (
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type FormStatusRepository interface {
	GetAll(ctx context.Context) ([]models.FormStatus, error)
}

type formStatusRepository struct {
	db *gorm.DB
}

func NewFormStatusRepository(db *gorm.DB) FormStatusRepository {
	return &formStatusRepository{db: db}
}

func (r *formStatusRepository) GetAll(ctx context.Context) ([]models.FormStatus, error) {
	var statuses []models.FormStatus
	if err := r.db.WithContext(ctx).Find(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}
