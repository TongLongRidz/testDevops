package campus

import (
	campusDTO "backend/internal/dto/campus_dto"
	"backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type CampusHandler struct {
	service usecase.CampusService
}

func NewCampusHandler(service usecase.CampusService) *CampusHandler {
	return &CampusHandler{service: service}
}

// GetAllCampuses ดึงข้อมูลวิทยาเขตทั้งหมด
func (h *CampusHandler) GetAllCampuses(c *fiber.Ctx) error {
	campuses, err := h.service.GetAllCampuses(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responses := make([]campusDTO.CampusResponse, 0, len(campuses))
	for _, campus := range campuses {
		responses = append(responses, campusDTO.CampusResponse{
			CampusID:   campus.CampusID,
			CampusCode: campus.CampusCode,
			CampusName: campus.CampusName,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Campuses retrieved successfully",
		"data":    responses,
	})
}
