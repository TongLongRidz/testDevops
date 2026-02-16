package repository

import (
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type CampusRepository interface {
	GetAll(ctx context.Context) ([]models.Campus, error)
}

type campusRepository struct {
	db *gorm.DB
}

func NewCampusRepository(db *gorm.DB) CampusRepository {
	return &campusRepository{db: db}
}

func (r *campusRepository) GetAll(ctx context.Context) ([]models.Campus, error) {
	var campuses []models.Campus
	err := r.db.WithContext(ctx).Find(&campuses).Error
	if err != nil {
		return nil, err
	}
	return campuses, nil
}
