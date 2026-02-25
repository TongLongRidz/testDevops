package repository

import (
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type AcademicYearRepository interface {
	Create(ctx context.Context, academicYear *models.AcademicYear) error
	GetByID(ctx context.Context, id uint) (*models.AcademicYear, error)
	GetAll(ctx context.Context) ([]models.AcademicYear, error)
	Update(ctx context.Context, academicYear *models.AcademicYear) error
	Delete(ctx context.Context, id uint) error
	GetCurrentSemester(ctx context.Context) (*models.AcademicYear, error)
	GetLatestAbleRegister(ctx context.Context) (*models.AcademicYear, error)
}

type academicYearRepository struct {
	db *gorm.DB
}

func NewAcademicYearRepository(db *gorm.DB) AcademicYearRepository {
	return &academicYearRepository{db: db}
}

func (r *academicYearRepository) Create(ctx context.Context, academicYear *models.AcademicYear) error {
	return r.db.WithContext(ctx).Create(academicYear).Error
}

func (r *academicYearRepository) GetByID(ctx context.Context, id uint) (*models.AcademicYear, error) {
	var academicYear models.AcademicYear
	err := r.db.WithContext(ctx).Where("academic_year_id = ?", id).First(&academicYear).Error
	if err != nil {
		return nil, err
	}
	return &academicYear, nil
}

func (r *academicYearRepository) GetAll(ctx context.Context) ([]models.AcademicYear, error) {
	var academicYears []models.AcademicYear
	err := r.db.WithContext(ctx).Find(&academicYears).Error
	if err != nil {
		return nil, err
	}
	return academicYears, nil
}

func (r *academicYearRepository) Update(ctx context.Context, academicYear *models.AcademicYear) error {
	return r.db.WithContext(ctx).Save(academicYear).Error
}

func (r *academicYearRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.AcademicYear{}, id).Error
}

// GetCurrentSemester: ดึงข้อมูล academic year ล่าสุด (เรียงปี/เทอมล่าสุด)
func (r *academicYearRepository) GetCurrentSemester(ctx context.Context) (*models.AcademicYear, error) {
	var academicYear models.AcademicYear
	err := r.db.WithContext(ctx).
		Order("year DESC").
		Order("semester DESC").
		First(&academicYear).Error
	if err != nil {
		return nil, err
	}
	return &academicYear, nil
}

// GetLatestAbleRegister: ดึงข้อมูล academic year ล่าสุดที่ใช้รับสมัคร
func (r *academicYearRepository) GetLatestAbleRegister(ctx context.Context) (*models.AcademicYear, error) {
	var academicYear models.AcademicYear
	err := r.db.WithContext(ctx).
		Order("year DESC").
		Order("semester DESC").
		First(&academicYear).Error
	if err != nil {
		return nil, err
	}
	return &academicYear, nil
}
