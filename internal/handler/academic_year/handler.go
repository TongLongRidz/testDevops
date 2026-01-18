package academicyear

import (
	academicYearDTO "backend/internal/dto/academic_year_dto"
	"backend/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AcademicYearHandler struct {
	service usecase.AcademicYearService
}

func NewAcademicYearHandler(service usecase.AcademicYearService) *AcademicYearHandler {
	return &AcademicYearHandler{service: service}
}

// CreateAcademicYear สร้าง academic year ใหม่
func (h *AcademicYearHandler) CreateAcademicYear(c *fiber.Ctx) error {
	req := new(academicYearDTO.CreateAcademicYear)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	academicYear, err := h.service.CreateAcademicYear(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &academicYearDTO.AcademicYearResponse{
		AcademicYearID: academicYear.AcademicYearID,
		Year:           academicYear.Year,
		Semester:       academicYear.Semester,
		StartDate:      academicYear.StartDate,
		EndDate:        academicYear.EndDate,
		IsCurrent:      academicYear.IsCurrent,
		IsOpenRegister: academicYear.IsOpenRegister,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Academic year created successfully",
		"data":    response,
	})
}

// GetAllAcademicYears ดึงข้อมูล academic year ทั้งหมด
func (h *AcademicYearHandler) GetAllAcademicYears(c *fiber.Ctx) error {
	academicYears, err := h.service.GetAllAcademicYears(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var responses []academicYearDTO.AcademicYearResponse
	for _, ay := range academicYears {
		responses = append(responses, academicYearDTO.AcademicYearResponse{
			AcademicYearID: ay.AcademicYearID,
			Year:           ay.Year,
			Semester:       ay.Semester,
			StartDate:      ay.StartDate,
			EndDate:        ay.EndDate,
			IsCurrent:      ay.IsCurrent,
			IsOpenRegister: ay.IsOpenRegister,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Academic years retrieved successfully",
		"data":    responses,
	})
}

// GetAcademicYearByID ดึงข้อมูล academic year ตามรหัส
func (h *AcademicYearHandler) GetAcademicYearByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid academic year ID",
		})
	}

	academicYear, err := h.service.GetAcademicYearByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Academic year not found",
		})
	}

	response := &academicYearDTO.AcademicYearResponse{
		AcademicYearID: academicYear.AcademicYearID,
		Year:           academicYear.Year,
		Semester:       academicYear.Semester,
		StartDate:      academicYear.StartDate,
		EndDate:        academicYear.EndDate,
		IsCurrent:      academicYear.IsCurrent,
		IsOpenRegister: academicYear.IsOpenRegister,
	}

	return c.JSON(fiber.Map{
		"message": "Academic year retrieved successfully",
		"data":    response,
	})
}

// UpdateAcademicYear แก้ไข academic year
func (h *AcademicYearHandler) UpdateAcademicYear(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid academic year ID",
		})
	}

	req := new(academicYearDTO.UpdateAcademicYear)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	academicYear, err := h.service.UpdateAcademicYear(c.Context(), uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &academicYearDTO.AcademicYearResponse{
		AcademicYearID: academicYear.AcademicYearID,
		Year:           academicYear.Year,
		Semester:       academicYear.Semester,
		StartDate:      academicYear.StartDate,
		EndDate:        academicYear.EndDate,
		IsCurrent:      academicYear.IsCurrent,
		IsOpenRegister: academicYear.IsOpenRegister,
	}

	return c.JSON(fiber.Map{
		"message": "Academic year updated successfully",
		"data":    response,
	})
}

// DeleteAcademicYear ลบ academic year
func (h *AcademicYearHandler) DeleteAcademicYear(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid academic year ID",
		})
	}

	if err := h.service.DeleteAcademicYear(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Academic year deleted successfully",
	})
}

// ToggleCurrent เปิด/ปิด isCurrent (มีได้แค่อันเดียว)
func (h *AcademicYearHandler) ToggleCurrent(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid academic year ID",
		})
	}

	academicYear, err := h.service.ToggleCurrent(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &academicYearDTO.AcademicYearResponse{
		AcademicYearID: academicYear.AcademicYearID,
		Year:           academicYear.Year,
		Semester:       academicYear.Semester,
		StartDate:      academicYear.StartDate,
		EndDate:        academicYear.EndDate,
		IsCurrent:      academicYear.IsCurrent,
		IsOpenRegister: academicYear.IsOpenRegister,
	}

	return c.JSON(fiber.Map{
		"message": "isCurrent toggled successfully",
		"data":    response,
	})
}

// ToggleOpenRegister เปิด/ปิด isOpenRegister (มีได้แค่อันเดียว)
func (h *AcademicYearHandler) ToggleOpenRegister(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid academic year ID",
		})
	}

	academicYear, err := h.service.ToggleOpenRegister(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &academicYearDTO.AcademicYearResponse{
		AcademicYearID: academicYear.AcademicYearID,
		Year:           academicYear.Year,
		Semester:       academicYear.Semester,
		StartDate:      academicYear.StartDate,
		EndDate:        academicYear.EndDate,
		IsCurrent:      academicYear.IsCurrent,
		IsOpenRegister: academicYear.IsOpenRegister,
	}

	return c.JSON(fiber.Map{
		"message": "isOpenRegister toggled successfully",
		"data":    response,
	})
}

// GetCurrentSemester ดึงข้อมูล academic year ที่เป็น current
func (h *AcademicYearHandler) GetCurrentSemester(c *fiber.Ctx) error {
	academicYear, err := h.service.GetCurrentSemester(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No current academic year found",
		})
	}

	response := &academicYearDTO.AcademicYearResponse{
		AcademicYearID: academicYear.AcademicYearID,
		Year:           academicYear.Year,
		Semester:       academicYear.Semester,
		StartDate:      academicYear.StartDate,
		EndDate:        academicYear.EndDate,
		IsCurrent:      academicYear.IsCurrent,
		IsOpenRegister: academicYear.IsOpenRegister,
	}

	return c.JSON(fiber.Map{
		"message": "Current academic year retrieved successfully",
		"data":    response,
	})
}

// GetLatestAbleRegister ดึงข้อมูล academic year ที่ current = true และ isOpenRegister = true
func (h *AcademicYearHandler) GetLatestAbleRegister(c *fiber.Ctx) error {
	academicYear, err := h.service.GetLatestAbleRegister(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No academic year available for registration",
		})
	}

	response := &academicYearDTO.AcademicYearResponse{
		AcademicYearID: academicYear.AcademicYearID,
		Year:           academicYear.Year,
		Semester:       academicYear.Semester,
		StartDate:      academicYear.StartDate,
		EndDate:        academicYear.EndDate,
		IsCurrent:      academicYear.IsCurrent,
		IsOpenRegister: academicYear.IsOpenRegister,
	}

	return c.JSON(fiber.Map{
		"message": "Available academic year for registration retrieved successfully",
		"data":    response,
	})
}
