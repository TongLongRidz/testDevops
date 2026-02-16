package usecase

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type FormStatusService interface {
	GetAllFormStatuses(ctx context.Context) ([]models.FormStatus, error)
}

type formStatusService struct {
	repo repository.FormStatusRepository
}

func NewFormStatusService(repo repository.FormStatusRepository) FormStatusService {
	return &formStatusService{repo: repo}
}

func (s *formStatusService) GetAllFormStatuses(ctx context.Context) ([]models.FormStatus, error) {
	return s.repo.GetAll(ctx)
}
