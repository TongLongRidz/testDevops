package academicyear

import (
	academicYearDTO "backend/internal/dto/academic_year_dto"
	"backend/internal/usecase"

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
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Academic year created successfully",
		"data":    response,
	})
}

// GetAllAcademicYears ดึงเฉพาะปี (list ของ years ที่ไม่ซ้ำ)
func (h *AcademicYearHandler) GetAllAcademicYears(c *fiber.Ctx) error {
	academicYears, err := h.service.GetAllAcademicYears(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Extract unique years
	yearMap := make(map[int]bool)
	var years []int
	for _, ay := range academicYears {
		if !yearMap[ay.Year] {
			yearMap[ay.Year] = true
			years = append(years, ay.Year)
		}
	}

	return c.JSON(fiber.Map{
		"message": "All academic years retrieved successfully",
		"data":    years,
	})
}

// GetCurrentSemester ดึงข้อมูล academic year ล่าสุด
func (h *AcademicYearHandler) GetCurrentSemester(c *fiber.Ctx) error {
	academicYear, err := h.service.GetCurrentSemester(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No academic year found",
		})
	}

	response := &academicYearDTO.AcademicYearResponse{
		AcademicYearID: academicYear.AcademicYearID,
		Year:           academicYear.Year,
		Semester:       academicYear.Semester,
		StartDate:      academicYear.StartDate,
		EndDate:        academicYear.EndDate,
	}

	return c.JSON(fiber.Map{
		"message": "Latest academic year retrieved successfully",
		"data":    response,
	})
}
