package repository

import (
	"backend/internal/models"
	"context"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(ctx context.Context, department *models.Department) error
	GetByID(ctx context.Context, id uint) (*models.Department, error)
	GetAll(ctx context.Context) ([]models.Department, error)
	Update(ctx context.Context, department *models.Department) error
	Delete(ctx context.Context, id uint) error
	GetByFacultyID(ctx context.Context, facultyID uint) ([]models.Department, error)
	GetByName(ctx context.Context, name string) (*models.Department, error)
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(ctx context.Context, department *models.Department) error {
	return r.db.WithContext(ctx).Create(department).Error
}

func (r *departmentRepository) GetByID(ctx context.Context, id uint) (*models.Department, error) {
	var department models.Department
	err := r.db.WithContext(ctx).Preload("Faculty").Where("department_id = ?", id).First(&department).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepository) GetAll(ctx context.Context) ([]models.Department, error) {
	var departments []models.Department
	err := r.db.WithContext(ctx).Preload("Faculty").Find(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}

func (r *departmentRepository) Update(ctx context.Context, department *models.Department) error {
	return r.db.WithContext(ctx).Save(department).Error
}

func (r *departmentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Department{}, id).Error
}

func (r *departmentRepository) GetByFacultyID(ctx context.Context, facultyID uint) ([]models.Department, error) {
	var departments []models.Department
	err := r.db.WithContext(ctx).Preload("Faculty").Where("faculty_id = ?", facultyID).Find(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}

func (r *departmentRepository) GetByName(ctx context.Context, name string) (*models.Department, error) {
	var department models.Department
	err := r.db.WithContext(ctx).Preload("Faculty").Where("department_name = ?", name).First(&department).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}
