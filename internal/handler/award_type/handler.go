package awardtype

import (
	awardTypeDTO "backend/internal/dto/award_type_dto"
	"backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type AwardTypeHandler struct {
	service usecase.AwardTypeService
}

func NewAwardTypeHandler(service usecase.AwardTypeService) *AwardTypeHandler {
	return &AwardTypeHandler{service: service}
}

// GetAllAwardTypes ดึงข้อมูล award type ทั้งหมด
func (h *AwardTypeHandler) GetAllAwardTypes(c *fiber.Ctx) error {
	types, err := h.service.GetAllAwardTypes(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responses := make([]awardTypeDTO.AwardTypeResponse, 0, len(types))
	for _, t := range types {
		responses = append(responses, awardTypeDTO.AwardTypeResponse{
			AwardTypeID: t.AwardTypeID,
			AwardName:   t.AwardName,
			Description: t.Description,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Award types retrieved successfully",
		"data":    responses,
	})
}
