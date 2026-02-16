package usecase

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type CampusService interface {
	GetAllCampuses(ctx context.Context) ([]models.Campus, error)
}

type campusService struct {
	repo repository.CampusRepository
}

func NewCampusService(repo repository.CampusRepository) CampusService {
	return &campusService{repo: repo}
}

func (s *campusService) GetAllCampuses(ctx context.Context) ([]models.Campus, error) {
	return s.repo.GetAll(ctx)
}
