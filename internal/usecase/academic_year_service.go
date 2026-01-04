package usecase

import (
    "backend/internal/models"
    "context"
    "errors"
)

type AcademicYearRepo interface {
    Create(ctx context.Context, ay *models.AcademicYear) error
    Update(ctx context.Context, ay *models.AcademicYear) error
    Delete(ctx context.Context, id uint) error
    GetByID(ctx context.Context, id uint) (*models.AcademicYear, error)
    GetLatest(ctx context.Context) (*models.AcademicYear, error)
	GetList(ctx context.Context) ([]models.AcademicYear, error)
}

type AcademicYearService struct {
    repo AcademicYearRepo
}

func NewAcademicYearService(r AcademicYearRepo) *AcademicYearService {
    return &AcademicYearService{repo: r}
}

// Create เพิ่มปีการศึกษา
func (s *AcademicYearService) Create(ctx context.Context, ay *models.AcademicYear) (*models.AcademicYear, error) {
    if ay == nil {
        return nil, errors.New("invalid input")
    }
    if err := s.repo.Create(ctx, ay); err != nil {
        return nil, err
    }
    return ay, nil
}

// Update แก้ไขปีการศึกษา (ต้องมี academic_year_id ใน ay)
func (s *AcademicYearService) Update(ctx context.Context, ay *models.AcademicYear) (*models.AcademicYear, error) {
    if ay == nil || ay.AcademicYearID == 0 {
        return nil, errors.New("invalid input")
    }
    // ตรวจว่า record มีจริง
    _, err := s.repo.GetByID(ctx, ay.AcademicYearID)
    if err != nil {
        return nil, err
    }
    if err := s.repo.Update(ctx, ay); err != nil {
        return nil, err
    }
    return ay, nil
}

// Delete ลบปีการศึกษา ตาม id
func (s *AcademicYearService) Delete(ctx context.Context, id uint) error {
    if id == 0 {
        return errors.New("invalid id")
    }
    // ตรวจว่า record มีจริง
    _, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    return s.repo.Delete(ctx, id)
}

func (s *AcademicYearService) GetList(ctx context.Context) ([]models.AcademicYear, error) {
    return s.repo.GetList(ctx)
}

func (s *AcademicYearService) GetByID(ctx context.Context, id uint) (*models.AcademicYear, error) {
    if id == 0 {
        return nil, errors.New("invalid id")
    }
    return s.repo.GetByID(ctx, id)
}
// GetLatest ดึงปีการศึกษาและเทอมล่าสุด
func (s *AcademicYearService) GetLatest(ctx context.Context) (*models.AcademicYear, error) {
    return s.repo.GetLatest(ctx)
}