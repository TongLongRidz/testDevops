package usecase

import (
	departmentDTO "backend/internal/dto/department_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type DepartmentService interface {
	CreateDepartment(ctx context.Context, req *departmentDTO.CreateDepartment) (*models.Department, error)
	GetDepartmentByID(ctx context.Context, id uint) (*models.Department, error)
	GetAllDepartments(ctx context.Context) ([]models.Department, error)
	UpdateDepartment(ctx context.Context, id uint, req *departmentDTO.UpdateDepartment) (*models.Department, error)
	DeleteDepartment(ctx context.Context, id uint) error
	GetDepartmentsByFacultyID(ctx context.Context, facultyID uint) ([]models.Department, error)
}

type departmentService struct {
	repo repository.DepartmentRepository
}

func NewDepartmentService(repo repository.DepartmentRepository) DepartmentService {
	return &departmentService{repo: repo}
}

func (s *departmentService) CreateDepartment(ctx context.Context, req *departmentDTO.CreateDepartment) (*models.Department, error) {
	department := &models.Department{
		FacultyID:      req.FacultyID,
		DepartmentName: req.DepartmentName,
	}

	if err := s.repo.Create(ctx, department); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *departmentService) GetDepartmentByID(ctx context.Context, id uint) (*models.Department, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *departmentService) GetAllDepartments(ctx context.Context) ([]models.Department, error) {
	return s.repo.GetAll(ctx)
}

func (s *departmentService) UpdateDepartment(ctx context.Context, id uint, req *departmentDTO.UpdateDepartment) (*models.Department, error) {
	department, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	department.FacultyID = req.FacultyID
	department.DepartmentName = req.DepartmentName

	if err := s.repo.Update(ctx, department); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *departmentService) DeleteDepartment(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *departmentService) GetDepartmentsByFacultyID(ctx context.Context, facultyID uint) ([]models.Department, error) {
	return s.repo.GetByFacultyID(ctx, facultyID)
}
