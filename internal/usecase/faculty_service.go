package usecase

import (
	facultyDTO "backend/internal/dto/faculty_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type FacultyService interface {
	CreateFaculty(ctx context.Context, req *facultyDTO.CreateFaculty) (*models.Faculty, error)
	GetFacultyByID(ctx context.Context, id uint) (*models.Faculty, error)
	GetAllFaculties(ctx context.Context) ([]models.Faculty, error)
	UpdateFaculty(ctx context.Context, id uint, req *facultyDTO.UpdateFaculty) (*models.Faculty, error)
	DeleteFaculty(ctx context.Context, id uint) error
}

type facultyService struct {
	repo repository.FacultyRepository
}

func NewFacultyService(repo repository.FacultyRepository) FacultyService {
	return &facultyService{repo: repo}
}

func (s *facultyService) CreateFaculty(ctx context.Context, req *facultyDTO.CreateFaculty) (*models.Faculty, error) {
	faculty := &models.Faculty{
		FacultyName: req.FacultyName,
	}

	if err := s.repo.Create(ctx, faculty); err != nil {
		return nil, err
	}

	return faculty, nil
}

func (s *facultyService) GetFacultyByID(ctx context.Context, id uint) (*models.Faculty, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *facultyService) GetAllFaculties(ctx context.Context) ([]models.Faculty, error) {
	return s.repo.GetAll(ctx)
}

func (s *facultyService) UpdateFaculty(ctx context.Context, id uint, req *facultyDTO.UpdateFaculty) (*models.Faculty, error) {
	faculty, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	faculty.FacultyName = req.FacultyName

	if err := s.repo.Update(ctx, faculty); err != nil {
		return nil, err
	}

	return faculty, nil
}

func (s *facultyService) DeleteFaculty(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
