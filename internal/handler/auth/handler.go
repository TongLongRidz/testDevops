package auth

import (
	"backend/internal/usecase"
	"os"
	"time"

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

	token, err := h.AuthService.IssueToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   3600 * 24,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	frontendBase := os.Getenv("FRONTEND_BASE_URL")
	if frontendBase == "" {
		frontendBase = "http://localhost:3000"
	}
	redirectPath := "/student/main/student-nomination-form"
	if user.IsFirstLogin {
		redirectPath = "/student/auth/first-login"
	}
	return c.Redirect(frontendBase + redirectPath)
}
