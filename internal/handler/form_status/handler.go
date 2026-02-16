package formstatus

import (
	formStatusDTO "backend/internal/dto/form_status_dto"
	"backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type FormStatusHandler struct {
	service usecase.FormStatusService
}

func NewFormStatusHandler(service usecase.FormStatusService) *FormStatusHandler {
	return &FormStatusHandler{service: service}
}

// GetAllFormStatuses ดึงข้อมูล form status ทั้งหมด
func (h *FormStatusHandler) GetAllFormStatuses(c *fiber.Ctx) error {
	statuses, err := h.service.GetAllFormStatuses(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responses := make([]formStatusDTO.FormStatusResponse, 0, len(statuses))
	for _, s := range statuses {
		responses = append(responses, formStatusDTO.FormStatusResponse{
			FormStatusID:   s.FormStatusID,
			FormStatusName: s.FormStatusName,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Form statuses retrieved successfully",
		"data":    responses,
	})
}
