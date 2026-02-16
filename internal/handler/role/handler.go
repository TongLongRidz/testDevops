package role

import (
	roleDTO "backend/internal/dto/role_dto"
	"backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	service usecase.RoleService
}

func NewRoleHandler(service usecase.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

// GetAllRoles ดึงข้อมูล role ทั้งหมด
func (h *RoleHandler) GetAllRoles(c *fiber.Ctx) error {
	roles, err := h.service.GetAllRoles(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responses := make([]roleDTO.RoleResponse, 0, len(roles))
	for _, r := range roles {
		responses = append(responses, roleDTO.RoleResponse{
			RoleID:     r.RoleID,
			RoleName:   r.RoleName,
			RoleNameTH: r.RoleNameTH,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Roles retrieved successfully",
		"data":    responses,
	})
}
