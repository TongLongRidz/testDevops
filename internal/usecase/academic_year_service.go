package usecase

import (
	academicYearDTO "backend/internal/dto/academic_year_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type AcademicYearService interface {
	CreateAcademicYear(ctx context.Context, req *academicYearDTO.CreateAcademicYear) (*models.AcademicYear, error)
	GetAcademicYearByID(ctx context.Context, id uint) (*models.AcademicYear, error)
	GetAllAcademicYears(ctx context.Context) ([]models.AcademicYear, error)
	UpdateAcademicYear(ctx context.Context, id uint, req *academicYearDTO.UpdateAcademicYear) (*models.AcademicYear, error)
	DeleteAcademicYear(ctx context.Context, id uint) error
	ToggleCurrent(ctx context.Context, id uint) (*models.AcademicYear, error)
	ToggleOpenRegister(ctx context.Context, id uint) (*models.AcademicYear, error)
	GetCurrentSemester(ctx context.Context) (*models.AcademicYear, error)
	GetLatestAbleRegister(ctx context.Context) (*models.AcademicYear, error)
}

type academicYearService struct {
	repo repository.AcademicYearRepository
}

func NewAcademicYearService(repo repository.AcademicYearRepository) AcademicYearService {
	return &academicYearService{repo: repo}
}

func (s *academicYearService) CreateAcademicYear(ctx context.Context, req *academicYearDTO.CreateAcademicYear) (*models.AcademicYear, error) {
	academicYear := &models.AcademicYear{
		Year:      req.Year,
		Semester:  req.Semester,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	if err := s.repo.Create(ctx, academicYear); err != nil {
		return nil, err
	}

	return academicYear, nil
}

func (s *academicYearService) GetAcademicYearByID(ctx context.Context, id uint) (*models.AcademicYear, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *academicYearService) GetAllAcademicYears(ctx context.Context) ([]models.AcademicYear, error) {
	return s.repo.GetAll(ctx)
}

func (s *academicYearService) UpdateAcademicYear(ctx context.Context, id uint, req *academicYearDTO.UpdateAcademicYear) (*models.AcademicYear, error) {
	academicYear, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	academicYear.Year = req.Year
	academicYear.Semester = req.Semester
	academicYear.StartDate = req.StartDate
	academicYear.EndDate = req.EndDate

	if err := s.repo.Update(ctx, academicYear); err != nil {
		return nil, err
	}

	return academicYear, nil
}

func (s *academicYearService) DeleteAcademicYear(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *academicYearService) ToggleCurrent(ctx context.Context, id uint) (*models.AcademicYear, error) {
	if err := s.repo.ToggleCurrent(ctx, id); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *academicYearService) ToggleOpenRegister(ctx context.Context, id uint) (*models.AcademicYear, error) {
	if err := s.repo.ToggleOpenRegister(ctx, id); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *academicYearService) GetCurrentSemester(ctx context.Context) (*models.AcademicYear, error) {
	return s.repo.GetCurrentSemester(ctx)
}

func (s *academicYearService) GetLatestAbleRegister(ctx context.Context) (*models.AcademicYear, error) {
	return s.repo.GetLatestAbleRegister(ctx)
}
