package usecase

import (
	academicYearDTO "backend/internal/dto/academic_year_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"fmt"
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
	if err := s.validateAcademicYearRules(ctx, req.Year, req.Semester, nil); err != nil {
		return nil, err
	}

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

	if err := s.validateAcademicYearRules(ctx, req.Year, req.Semester, &id); err != nil {
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

func (s *academicYearService) validateAcademicYearRules(ctx context.Context, year int, semester int, excludeID *uint) error {
	if semester != 1 && semester != 2 {
		return fmt.Errorf("semester must be 1 or 2")
	}

	academicYears, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
	}

	var hasSem1 bool
	var hasSem2 bool
	var maxYear int
	var maxSemester int
	for _, ay := range academicYears {
		if excludeID != nil && ay.AcademicYearID == *excludeID {
			continue
		}

		if ay.Year == year {
			if ay.Semester == 1 {
				hasSem1 = true
			}
			if ay.Semester == 2 {
				hasSem2 = true
			}
		}

		if ay.Year > maxYear || (ay.Year == maxYear && ay.Semester > maxSemester) {
			maxYear = ay.Year
			maxSemester = ay.Semester
		}
	}

	if semester == 1 && hasSem1 {
		return fmt.Errorf("semester 1 already exists for year %d", year)
	}
	if semester == 2 && hasSem2 {
		return fmt.Errorf("semester 2 already exists for year %d", year)
	}
	if semester == 1 && hasSem2 {
		return fmt.Errorf("cannot create semester 1 because semester 2 already exists for year %d", year)
	}
	if semester == 2 && !hasSem1 {
		return fmt.Errorf("semester 2 requires semester 1 for year %d", year)
	}

	if len(academicYears) == 0 {
		return nil
	}

	if year < maxYear {
		return fmt.Errorf("year must be %d or newer", maxYear)
	}
	if year == maxYear {
		if maxSemester == 1 && semester != 2 {
			return fmt.Errorf("next semester must be 2 for year %d", year)
		}
		if maxSemester == 2 {
			return fmt.Errorf("year %d already has semester 1 and 2; next semester must be year %d semester 1", year, year+1)
		}
	}
	if year == maxYear+1 && semester != 1 {
		return fmt.Errorf("new year must start with semester 1")
	}
	if year > maxYear+1 {
		return fmt.Errorf("year must be %d or %d", maxYear, maxYear+1)
	}

	return nil
}
