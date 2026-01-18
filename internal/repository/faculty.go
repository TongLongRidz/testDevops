package repository

import (
	"backend/internal/models"
	"context"
	"gorm.io/gorm"
)

type FacultyRepository interface {
	Create(ctx context.Context, faculty *models.Faculty) error
	GetByID(ctx context.Context, id uint) (*models.Faculty, error)
	GetAll(ctx context.Context) ([]models.Faculty, error)
	Update(ctx context.Context, faculty *models.Faculty) error
	Delete(ctx context.Context, id uint) error
	GetByName(ctx context.Context, name string) (*models.Faculty, error)
}

type facultyRepository struct {
	db *gorm.DB
}

func NewFacultyRepository(db *gorm.DB) FacultyRepository {
	return &facultyRepository{db: db}
}

func (r *facultyRepository) Create(ctx context.Context, faculty *models.Faculty) error {
	return r.db.WithContext(ctx).Create(faculty).Error
}

func (r *facultyRepository) GetByID(ctx context.Context, id uint) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.db.WithContext(ctx).Where("faculty_id = ?", id).First(&faculty).Error
	if err != nil {
		return nil, err
	}
	return &faculty, nil
}

func (r *facultyRepository) GetAll(ctx context.Context) ([]models.Faculty, error) {
	var faculties []models.Faculty
	err := r.db.WithContext(ctx).Find(&faculties).Error
	if err != nil {
		return nil, err
	}
	return faculties, nil
}

func (r *facultyRepository) Update(ctx context.Context, faculty *models.Faculty) error {
	return r.db.WithContext(ctx).Save(faculty).Error
}

func (r *facultyRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Faculty{}, id).Error
}

func (r *facultyRepository) GetByName(ctx context.Context, name string) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.db.WithContext(ctx).Where("faculty_name = ?", name).First(&faculty).Error
	if err != nil {
		return nil, err
	}
	return &faculty, nil
}
