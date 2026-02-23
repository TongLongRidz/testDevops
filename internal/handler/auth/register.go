package auth

import (
	authDto "backend/internal/dto/auth_dto"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) Register(c *fiber.Ctx) error {

	fmt.Printf("[Register] raw body: %s\n", string(c.Body()))
	var req authDto.RegisterRequest
	// 1. Parse JSON payload (ใช้ BodyParser ของ Fiber)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	fmt.Printf("[Register] parsed request: %+v\n", req)
	// 2. Simple validation (Email และ Password ต้องมี)
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email and password required",
		})
	}

	// 3. Password matching check (ตรวจสอบ ConfirmPassword)
	if req.Password != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "passwords do not match",
		})
	}

	fmt.Printf("[Register] calling usecase.Register for email=%s\n", req.Email)
	// 4. Call Usecase (ส่ง Request ไปให้ Service จัดการต่อ)
	createdUser, err := h.AuthService.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fmt.Printf("[Register] created user id=%d email=%s provider=%s\n", createdUser.UserID, createdUser.Email, createdUser.Provider)
	// 5. Response format ตามที่คุณต้องการ อิงตาม Model ใหม่
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_id":        createdUser.UserID,
		"email":          createdUser.Email,
		"role_id":        createdUser.RoleID,
		"is_first_login": createdUser.IsFirstLogin,
		"provider":       createdUser.Provider,
	})
}

func (h *AuthHandler) CreateAccount(c *fiber.Ctx) error {
	var req authDto.CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	createdUser, err := h.AuthService.CreateAccount(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_id":        createdUser.UserID,
		"email":          createdUser.Email,
		"role_id":        createdUser.RoleID,
		"campus_id":      createdUser.CampusID,
		"is_first_login": createdUser.IsFirstLogin,
		"provider":       createdUser.Provider,
	})
}
