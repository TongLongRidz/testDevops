package auth

import (
	"backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AuthService    usecase.AuthService
	StudentService usecase.StudentService
}

func NewAuthHandler(u usecase.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: u}
}

func NewAuthHandlerWithStudent(u usecase.AuthService, s usecase.StudentService) *AuthHandler {
	return &AuthHandler{AuthService: u, StudentService: s}
}

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	url := h.AuthService.GetGoogleLoginURL()
	return c.Redirect(url)
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Code not found"})
	}

	user, err := h.AuthService.ProcessGoogleLogin(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// เมื่อสำเร็จ อาจจะส่ง User Data กลับไป หรือออก JWT Token
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
	})
}