package repository

import (
	"backend/internal/models"
	"context"
	"gorm.io/gorm"
)

type StudentRepository interface {
	Create(ctx context.Context, student *models.Student) error
	GetByID(ctx context.Context, id uint) (*models.Student, error)
	GetByUserID(ctx context.Context, userID uint) (*models.Student, error)
	GetAll(ctx context.Context) ([]models.Student, error)
	Update(ctx context.Context, student *models.Student) error
	Delete(ctx context.Context, id uint) error
	GetByStudentNumber(ctx context.Context, studentNumber string) (*models.Student, error)
	GetByFacultyID(ctx context.Context, facultyID uint) ([]models.Student, error)
	GetByDepartmentID(ctx context.Context, departmentID uint) ([]models.Student, error)
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(ctx context.Context, student *models.Student) error {
	return r.db.WithContext(ctx).Create(student).Error
}

func (r *studentRepository) GetByID(ctx context.Context, id uint) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Faculty").
		Preload("Department").
		Where("student_id = ?", id).
		First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByUserID(ctx context.Context, userID uint) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Faculty").
		Preload("Department").
		Where("user_id = ?", userID).
		First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetAll(ctx context.Context) ([]models.Student, error) {
	var students []models.Student
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Faculty").
		Preload("Department").
		Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *studentRepository) Update(ctx context.Context, student *models.Student) error {
	return r.db.WithContext(ctx).Save(student).Error
}

func (r *studentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Student{}, id).Error
}

func (r *studentRepository) GetByStudentNumber(ctx context.Context, studentNumber string) (*models.Student, error) {
	var student models.Student
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Faculty").
		Preload("Department").
		Where("student_number = ?", studentNumber).
		First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByFacultyID(ctx context.Context, facultyID uint) ([]models.Student, error) {
	var students []models.Student
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Faculty").
		Preload("Department").
		Where("faculty_id = ?", facultyID).
		Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *studentRepository) GetByDepartmentID(ctx context.Context, departmentID uint) ([]models.Student, error) {
	var students []models.Student
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Faculty").
		Preload("Department").
		Where("department_id = ?", departmentID).
		Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}
