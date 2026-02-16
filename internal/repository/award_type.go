package repository

import (
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type AwardTypeRepository interface {
	GetAll(ctx context.Context) ([]models.AwardType, error)
}

type awardTypeRepository struct {
	db *gorm.DB
}

func NewAwardTypeRepository(db *gorm.DB) AwardTypeRepository {
	return &awardTypeRepository{db: db}
}

func (r *awardTypeRepository) GetAll(ctx context.Context) ([]models.AwardType, error) {
	var types []models.AwardType
	if err := r.db.WithContext(ctx).Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}
