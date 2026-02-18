package repository

import (
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type OrganizationRepository interface {
	GetByUserID(ctx context.Context, userID uint) (*models.Organization, error)
	Create(ctx context.Context, org *models.Organization) error
	Update(ctx context.Context, org *models.Organization) error
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) GetByUserID(ctx context.Context, userID uint) (*models.Organization, error) {
	var org models.Organization
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) Create(ctx context.Context, org *models.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *organizationRepository) Update(ctx context.Context, org *models.Organization) error {
	return r.db.WithContext(ctx).Save(org).Error
}
