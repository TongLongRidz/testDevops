package usecase

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type RoleService interface {
	GetAllRoles(ctx context.Context) ([]models.Role, error)
}

type roleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	return s.repo.GetAll(ctx)
}
