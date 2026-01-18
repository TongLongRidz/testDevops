package student

import (
	studentDTO "backend/internal/dto/student_dto"
	"backend/internal/models"
	"backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type StudentHandler struct {
	service usecase.StudentService
}

func NewStudentHandler(service usecase.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

// CreateStudent สร้าง student ใหม่ (สำหรับ admin)
func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("userId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	req := new(studentDTO.CreateStudentRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	student, err := h.service.CreateStudent(c.Context(), uint(userID), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &studentDTO.StudentWithUserResponse{
		StudentID:      student.StudentID,
		UserID:         student.UserID,
		StudentNumber:  student.StudentNumber,
		FacultyID:      student.FacultyID,
		DepartmentID:   student.DepartmentID,
		Firstname:      student.User.Firstname,
		Lastname:       student.User.Lastname,
		Email:          student.User.Email,
		ImagePath:      student.User.ImagePath,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Student created successfully",
		"data":    response,
	})
}

// GetAllStudents ดึงข้อมูล student ทั้งหมด
func (h *StudentHandler) GetAllStudents(c *fiber.Ctx) error {
	students, err := h.service.GetAllStudents(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var responses []studentDTO.StudentWithUserResponse
	for _, s := range students {
		responses = append(responses, studentDTO.StudentWithUserResponse{
			StudentID:      s.StudentID,
			UserID:         s.UserID,
			StudentNumber:  s.StudentNumber,
			FacultyID:      s.FacultyID,
			DepartmentID:   s.DepartmentID,
			Firstname:      s.User.Firstname,
			Lastname:       s.User.Lastname,
			Email:          s.User.Email,
			ImagePath:      s.User.ImagePath,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Students retrieved successfully",
		"data":    responses,
	})
}

// GetStudentByID ดึงข้อมูล student ตามรหัส
func (h *StudentHandler) GetStudentByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid student ID",
		})
	}

	student, err := h.service.GetStudentByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	response := &studentDTO.StudentWithUserResponse{
		StudentID:      student.StudentID,
		UserID:         student.UserID,
		StudentNumber:  student.StudentNumber,
		FacultyID:      student.FacultyID,
		DepartmentID:   student.DepartmentID,
		Firstname:      student.User.Firstname,
		Lastname:       student.User.Lastname,
		Email:          student.User.Email,
		ImagePath:      student.User.ImagePath,
	}

	return c.JSON(fiber.Map{
		"message": "Student retrieved successfully",
		"data":    response,
	})
}

// GetMyStudent ดึงข้อมูล student ของตัวเอง (สำหรับ student ที่ login)
func (h *StudentHandler) GetMyStudent(c *fiber.Ctx) error {
	u := c.Locals("current_user")
	if u == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	user, ok := u.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user"})
	}

	student, err := h.service.GetStudentByUserID(c.Context(), user.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	response := &studentDTO.StudentWithUserResponse{
		StudentID:      student.StudentID,
		UserID:         student.UserID,
		StudentNumber:  student.StudentNumber,
		FacultyID:      student.FacultyID,
		DepartmentID:   student.DepartmentID,
		Firstname:      student.User.Firstname,
		Lastname:       student.User.Lastname,
		Email:          student.User.Email,
		ImagePath:      student.User.ImagePath,
	}

	return c.JSON(fiber.Map{
		"message": "Student retrieved successfully",
		"data":    response,
	})
}

// UpdateStudent แก้ไข student
func (h *StudentHandler) UpdateStudent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid student ID",
		})
	}

	req := new(studentDTO.CreateStudentRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	student, err := h.service.UpdateStudent(c.Context(), uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &studentDTO.StudentWithUserResponse{
		StudentID:      student.StudentID,
		UserID:         student.UserID,
		StudentNumber:  student.StudentNumber,
		FacultyID:      student.FacultyID,
		DepartmentID:   student.DepartmentID,
		Firstname:      student.User.Firstname,
		Lastname:       student.User.Lastname,
		Email:          student.User.Email,
		ImagePath:      student.User.ImagePath,
	}

	return c.JSON(fiber.Map{
		"message": "Student updated successfully",
		"data":    response,
	})
}

// UpdateMyStudent แก้ไข student ของตัวเอง
func (h *StudentHandler) UpdateMyStudent(c *fiber.Ctx) error {
	u := c.Locals("current_user")
	if u == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	user, ok := u.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user"})
	}

	// ดึง student ตาม UserID
	student, err := h.service.GetStudentByUserID(c.Context(), user.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	req := new(studentDTO.CreateStudentRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	updatedStudent, err := h.service.UpdateStudent(c.Context(), student.StudentID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &studentDTO.StudentWithUserResponse{
		StudentID:      updatedStudent.StudentID,
		UserID:         updatedStudent.UserID,
		StudentNumber:  updatedStudent.StudentNumber,
		FacultyID:      updatedStudent.FacultyID,
		DepartmentID:   updatedStudent.DepartmentID,
		Firstname:      updatedStudent.User.Firstname,
		Lastname:       updatedStudent.User.Lastname,
		Email:          updatedStudent.User.Email,
		ImagePath:      updatedStudent.User.ImagePath,
	}

	return c.JSON(fiber.Map{
		"message": "Student updated successfully",
		"data":    response,
	})
}

// DeleteStudent ลบ student
func (h *StudentHandler) DeleteStudent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid student ID",
		})
	}

	if err := h.service.DeleteStudent(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Student deleted successfully",
	})
}
