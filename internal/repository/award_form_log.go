package repository

import (
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type AwardFormLogRepository struct {
	db *gorm.DB
}

func NewAwardFormLogRepository(db *gorm.DB) *AwardFormLogRepository {
	return &AwardFormLogRepository{db: db}
}

func (r *AwardFormLogRepository) Create(ctx context.Context, log *models.AwardFormLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AwardFormLogRepository) GetByFormID(ctx context.Context, formID uint) ([]models.AwardFormLog, error) {
	var logs []models.AwardFormLog
	err := r.db.WithContext(ctx).
		Where("form_id = ?", formID).
		Order("created_at desc").
		Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}
