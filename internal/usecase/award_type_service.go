package usecase

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type AwardTypeService interface {
	GetAllAwardTypes(ctx context.Context) ([]models.AwardType, error)
}

type awardTypeService struct {
	repo repository.AwardTypeRepository
}

func NewAwardTypeService(repo repository.AwardTypeRepository) AwardTypeService {
	return &awardTypeService{repo: repo}
}

func (s *awardTypeService) GetAllAwardTypes(ctx context.Context) ([]models.AwardType, error) {
	return s.repo.GetAll(ctx)
}
