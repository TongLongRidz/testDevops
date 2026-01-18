package user

import (
    userdto "backend/internal/dto/user_dto"
    // "backend/internal/models"
    "backend/internal/usecase"
    "strconv"

    "github.com/gofiber/fiber/v2"
)

type UserHandler struct {
    UserService usecase.UserService
}

func NewUserHandler(us usecase.UserService) *UserHandler {
    return &UserHandler{UserService: us}
}

// GET /users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id64, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
    }

    user, err := h.UserService.GetUserByID(c.Context(), uint(id64))
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
    }

    return c.JSON(userdto.UserResponse{
        UserID:       user.UserID,
        Firstname:    user.Firstname,
        Lastname:     user.Lastname,
        Email:        user.Email,
        ImagePath:    user.ImagePath,
        Provider:     user.Provider,
        RoleID:       user.RoleID,
        CampusID:     user.CampusID,
        IsFirstLogin: user.IsFirstLogin,
        CreatedAt:    user.CreatedAt,
        LatestUpdate: user.LatestUpdate,
    })
}

// PUT /users/:id
func (h *UserHandler) UpdateUserByID(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id64, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
    }

    var req userdto.EditUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
    }

    updated, err := h.UserService.UpdateUserByID(c.Context(), uint(id64), &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(userdto.UserResponse{
        UserID:       updated.UserID,
        Firstname:    updated.Firstname,
        Lastname:     updated.Lastname,
        Email:        updated.Email,
        ImagePath:    updated.ImagePath,
        Provider:     updated.Provider,
        RoleID:       updated.RoleID,
        CampusID:     updated.CampusID,
        IsFirstLogin: updated.IsFirstLogin,
        CreatedAt:    updated.CreatedAt,
        LatestUpdate: updated.LatestUpdate,
    })
}