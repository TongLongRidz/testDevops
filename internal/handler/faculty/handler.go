package faculty

import (
	facultyDTO "backend/internal/dto/faculty_dto"
	"backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type FacultyHandler struct {
	service usecase.FacultyService
}

func NewFacultyHandler(service usecase.FacultyService) *FacultyHandler {
	return &FacultyHandler{service: service}
}

// CreateFaculty สร้าง faculty ใหม่
func (h *FacultyHandler) CreateFaculty(c *fiber.Ctx) error {
	req := new(facultyDTO.CreateFaculty)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	faculty, err := h.service.CreateFaculty(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &facultyDTO.FacultyResponse{
		FacultyID:   faculty.FacultyID,
		FacultyName: faculty.FacultyName,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Faculty created successfully",
		"data":    response,
	})
}

// GetAllFaculties ดึงข้อมูล faculty ทั้งหมด
func (h *FacultyHandler) GetAllFaculties(c *fiber.Ctx) error {
	faculties, err := h.service.GetAllFaculties(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var responses []facultyDTO.FacultyResponse
	for _, f := range faculties {
		responses = append(responses, facultyDTO.FacultyResponse{
			FacultyID:   f.FacultyID,
			FacultyName: f.FacultyName,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Faculties retrieved successfully",
		"data":    responses,
	})
}

// GetFacultyByID ดึงข้อมูล faculty ตามรหัส
func (h *FacultyHandler) GetFacultyByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid faculty ID",
		})
	}

	faculty, err := h.service.GetFacultyByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Faculty not found",
		})
	}

	response := &facultyDTO.FacultyResponse{
		FacultyID:   faculty.FacultyID,
		FacultyName: faculty.FacultyName,
	}

	return c.JSON(fiber.Map{
		"message": "Faculty retrieved successfully",
		"data":    response,
	})
}

// UpdateFaculty แก้ไข faculty
func (h *FacultyHandler) UpdateFaculty(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid faculty ID",
		})
	}

	req := new(facultyDTO.UpdateFaculty)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	faculty, err := h.service.UpdateFaculty(c.Context(), uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &facultyDTO.FacultyResponse{
		FacultyID:   faculty.FacultyID,
		FacultyName: faculty.FacultyName,
	}

	return c.JSON(fiber.Map{
		"message": "Faculty updated successfully",
		"data":    response,
	})
}

// DeleteFaculty ลบ faculty
func (h *FacultyHandler) DeleteFaculty(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid faculty ID",
		})
	}

	if err := h.service.DeleteFaculty(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Faculty deleted successfully",
	})
}
