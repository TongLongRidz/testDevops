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
	GetCurrent(ctx context.Context) (*models.AcademicYear, error)
	GetOpenRegister(ctx context.Context) (*models.AcademicYear, error)
	ToggleCurrent(ctx context.Context, id uint) error
	ToggleOpenRegister(ctx context.Context, id uint) error
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

func (r *academicYearRepository) GetCurrent(ctx context.Context) (*models.AcademicYear, error) {
	var academicYear models.AcademicYear
	err := r.db.WithContext(ctx).Where("is_current = ?", true).First(&academicYear).Error
	if err != nil {
		return nil, err
	}
	return &academicYear, nil
}

func (r *academicYearRepository) GetOpenRegister(ctx context.Context) (*models.AcademicYear, error) {
	var academicYear models.AcademicYear
	err := r.db.WithContext(ctx).Where("is_open_register = ?", true).First(&academicYear).Error
	if err != nil {
		return nil, err
	}
	return &academicYear, nil
}

// ToggleCurrent: ปิดอันปัจจุบัน แล้วเปิดอันที่เลือก
func (r *academicYearRepository) ToggleCurrent(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. ปิดทั้งหมดที่ is_current = true
		if err := tx.Model(&models.AcademicYear{}).
			Where("is_current = ?", true).
			Update("is_current", false).Error; err != nil {
			return err
		}

		// 2. เปิดอันที่ส่งมา
		if err := tx.Model(&models.AcademicYear{}).
			Where("academic_year_id = ?", id).
			Update("is_current", true).Error; err != nil {
			return err
		}

		return nil
	})
}

// ToggleOpenRegister: ปิดอันปัจจุบัน แล้วเปิดอันที่เลือก
func (r *academicYearRepository) ToggleOpenRegister(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. ปิดทั้งหมดที่ is_open_register = true
		if err := tx.Model(&models.AcademicYear{}).
			Where("is_open_register = ?", true).
			Update("is_open_register", false).Error; err != nil {
			return err
		}

		// 2. เปิดอันที่ส่งมา
		if err := tx.Model(&models.AcademicYear{}).
			Where("academic_year_id = ?", id).
			Update("is_open_register", true).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetCurrentSemester: ดึงข้อมูล academic year ที่ current = true
func (r *academicYearRepository) GetCurrentSemester(ctx context.Context) (*models.AcademicYear, error) {
	var academicYear models.AcademicYear
	err := r.db.WithContext(ctx).Where("is_current = ?", true).First(&academicYear).Error
	if err != nil {
		return nil, err
	}
	return &academicYear, nil
}

// GetLatestAbleRegister: ดึงข้อมูล academic year ที่ current = true และ isOpenRegister = true
func (r *academicYearRepository) GetLatestAbleRegister(ctx context.Context) (*models.AcademicYear, error) {
	var academicYear models.AcademicYear
	err := r.db.WithContext(ctx).Where("is_current = ? AND is_open_register = ?", true, true).First(&academicYear).Error
	if err != nil {
		return nil, err
	}
	return &academicYear, nil
}
