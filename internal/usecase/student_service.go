package usecase

import (
	studentDTO "backend/internal/dto/student_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"fmt"
)

type StudentService interface {
	CreateStudent(ctx context.Context, userID uint, req *studentDTO.CreateStudentRequest) (*models.Student, error)
	GetStudentByID(ctx context.Context, id uint) (*models.Student, error)
	GetStudentByUserID(ctx context.Context, userID uint) (*models.Student, error)
	GetAllStudents(ctx context.Context) ([]models.Student, error)
	UpdateStudent(ctx context.Context, id uint, req *studentDTO.CreateStudentRequest) (*models.Student, error)
	DeleteStudent(ctx context.Context, id uint) error
}

type studentService struct {
	repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) StudentService {
	return &studentService{repo: repo}
}

func (s *studentService) CreateStudent(ctx context.Context, userID uint, req *studentDTO.CreateStudentRequest) (*models.Student, error) {
	if err := validateStudentNumber(req.StudentNumber); err != nil {
		return nil, err
	}

	student := &models.Student{
		UserID:        userID,
		StudentNumber: req.StudentNumber,
		FacultyID:     req.FacultyID,
		DepartmentID:  req.DepartmentID,
	}

	if err := s.repo.Create(ctx, student); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, student.StudentID)
}

func (s *studentService) GetStudentByID(ctx context.Context, id uint) (*models.Student, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *studentService) GetStudentByUserID(ctx context.Context, userID uint) (*models.Student, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *studentService) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	return s.repo.GetAll(ctx)
}

func (s *studentService) UpdateStudent(ctx context.Context, id uint, req *studentDTO.CreateStudentRequest) (*models.Student, error) {
	student, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := validateStudentNumber(req.StudentNumber); err != nil {
		return nil, err
	}

	student.StudentNumber = req.StudentNumber
	student.FacultyID = req.FacultyID
	student.DepartmentID = req.DepartmentID

	if err := s.repo.Update(ctx, student); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *studentService) DeleteStudent(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func validateStudentNumber(studentNumber string) error {
	if len(studentNumber) != 10 {
		return fmt.Errorf("student_number must be exactly 10 digits")
	}
	for _, r := range studentNumber {
		if r < '0' || r > '9' {
			return fmt.Errorf("student_number must be exactly 10 digits")
		}
	}
	return nil
}
