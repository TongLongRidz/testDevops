package usecase

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type OrganizationService interface {
	GetByUserID(ctx context.Context, userID uint) (*models.Organization, error)
}

type organizationService struct {
	repo repository.OrganizationRepository
}

func NewOrganizationService(repo repository.OrganizationRepository) OrganizationService {
	return &organizationService{repo: repo}
}

func (s *organizationService) GetByUserID(ctx context.Context, userID uint) (*models.Organization, error) {
	return s.repo.GetByUserID(ctx, userID)
}
