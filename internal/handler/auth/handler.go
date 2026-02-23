package auth

import (
	"backend/internal/usecase"
	"os"

	// "time"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AuthService         usecase.AuthService
	StudentService      usecase.StudentService
	OrganizationService usecase.OrganizationService
}

func NewAuthHandler(u usecase.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: u}
}

func NewAuthHandlerWithStudent(u usecase.AuthService, s usecase.StudentService) *AuthHandler {
	return &AuthHandler{AuthService: u, StudentService: s}
}

func NewAuthHandlerWithServices(u usecase.AuthService, s usecase.StudentService, o usecase.OrganizationService) *AuthHandler {
	return &AuthHandler{AuthService: u, StudentService: s, OrganizationService: o}
}

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	url := h.AuthService.GetGoogleLoginURL()
	return c.Redirect(url)
}

// func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
// 	code := c.Query("code")
// 	if code == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Code not found"})
// 	}

// 	user, err := h.AuthService.ProcessGoogleLogin(code)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	token, err := h.AuthService.IssueToken(user)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	c.Cookie(&fiber.Cookie{
// 		Name:     "token",
// 		Value:    token,
// 		Path:     "/",
// 		MaxAge:   3600 * 24,
// 		HTTPOnly: true,
// 		Secure:   false,
// 		SameSite: "Lax",
// 		Expires:  time.Now().Add(24 * time.Hour),
// 	})

// 	frontendBase := os.Getenv("FRONTEND_BASE_URL")
// 	if frontendBase == "" {
// 		frontendBase = "http://localhost:3000"
// 	}
// 	redirectPath := "/student/main/student-nomination-form"
// 	if user.IsFirstLogin {
// 		redirectPath = "/student/auth/first-login"
// 	}
// 	return c.Redirect(frontendBase + redirectPath)
// }

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Code not found"})
	}

	// 1. Process Google Login
	user, err := h.AuthService.ProcessGoogleLogin(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 2. Issue Token
	token, err := h.AuthService.IssueToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 3. Prepare Data for Frontend
	isFirstLoginStr := "false"
	if user.IsFirstLogin {
		isFirstLoginStr = "true"
	}

	// Map Role ID to Role Name
	roleName := "student"
	switch user.RoleID {
	case 1:
		roleName = "student"
	case 2:
		roleName = "head_of_department"
	case 3:
		roleName = "associate_dean"
	case 4:
		roleName = "dean"
	case 5:
		roleName = "student_development"
	case 6:
		roleName = "committee"
	case 7:
		roleName = "chancellor"
	case 8:
		roleName = "organization"
	default:
		roleName = "student"
	}

	// 4. Get Frontend Base URL
	frontendBase := os.Getenv("FRONTEND_BASE_URL")
	if frontendBase == "" {
		frontendBase = "http://localhost:3000"
	}

	// 5. ✅ [จุดสำคัญ] Redirect ไปที่หน้า google-callback ของ Frontend
	// พร้อมแนบข้อมูลที่ Frontend จำเป็นต้องใช้ในการ Login
	redirectURL := fmt.Sprintf(
		"%s/google-callback?token=%s&role=%s&first_login=%s&firstname=%s",
		frontendBase,
		token,
		roleName,
		isFirstLoginStr,
		user.Firstname,
	)

	return c.Redirect(redirectURL)
}
