package auth

import (
	"strings"
	"time"

	authDto "backend/internal/dto/auth_dto"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req authDto.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
    }

    req.Email = strings.TrimSpace(strings.ToLower(req.Email))
    if req.Email == "" || req.Password == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
    }

    token, user, err := h.AuthService.AuthenticateAndToken(c.Context(), req.Email, req.Password)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
    }

    // ตั้ง cookie (ปรับ Secure/SameSite ตาม environment ของคุณ)
    c.Cookie(&fiber.Cookie{
        Name:     "token",     // ชื่อ cookie
		Value:    token,       // ค่า token ที่ได้มา
		Path:     "/",         // ให้ cookie นี้ใช้ได้ทั่วทั้งเว็บ
		MaxAge:   3600 * 24, // 1 วัน (หน่วยเป็นวินาที)
		HTTPOnly: true,      // ป้องกันการเข้าถึงจาก JavaScript
		Secure:   false,     // ตั้งเป็น true ถ้า deploy แล้วใช้ HTTPS
		SameSite: "Lax",     // "Lax" หรือ "Strict" เพื่อความปลอดภัย
		Expires: time.Now().Add(24 * time.Hour),
    })

    return c.JSON(fiber.Map{
        "token": token,
        "user": fiber.Map{
            "user_id":        user.UserID,
            "firstname":      user.Firstname,
            "lastname":       user.Lastname,
            "email":          user.Email,
            "image_path":     user.ImagePath,
            "provider":       user.Provider,
            "role_id":        user.RoleID,
            "campus_id":      user.CampusID,
            "is_first_login": user.IsFirstLogin,
            "created_at":     user.CreatedAt,
            "latest_update":  user.LatestUpdate,
        },
    })
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
    // ตั้ง cookie ให้หมดอายุ (ล้าง token)

    c.Cookie(&fiber.Cookie{
        Name:     "token",
        Value:    "",
        Path:     "/",
        MaxAge:   -1,                  // ลบ cookie ทันที
        HTTPOnly: true,
        Secure:   false,
        SameSite: "Lax",
        Expires:  time.Now().Add(-time.Hour),
    })

    return c.JSON(fiber.Map{"message": "logged out"})
}

// เพิ่ม /me เพื่อคืนข้อมูลผู้ใช้อย่างปลอดภัย (ต้องผ่าน middleware.RequireAuth ก่อน)
func (h *AuthHandler) Me(c *fiber.Ctx) error {
    u := c.Locals("current_user")
    if u == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
    }

    user, ok := u.(*models.User)
    if !ok {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user"})
    }

    // ดึงข้อมูล user แบบเต็มจาก database
    fullUser, err := h.AuthService.GetUserByID(c.Context(), user.UserID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch user"})
    }

    // สร้าง response object
    response := authDto.MeResponse{
        UserID:       fullUser.UserID,
        Firstname:    fullUser.Firstname,
        Lastname:     fullUser.Lastname,
        Email:        fullUser.Email,
        ImagePath:    fullUser.ImagePath,
        Provider:     fullUser.Provider,
        RoleID:       fullUser.RoleID,
        CampusID:     fullUser.CampusID,
        IsFirstLogin: fullUser.IsFirstLogin,
        CreatedAt:    fullUser.CreatedAt,
        LatestUpdate: fullUser.LatestUpdate,
    }

    // ถ้า RoleID = 1 (Student) ให้ดึงข้อมูล Student ด้วย
    if fullUser.RoleID == 1 && h.StudentService != nil {
        student, err := h.StudentService.GetStudentByUserID(c.Context(), fullUser.UserID)
        if err == nil && student != nil {
            response.StudentData = &authDto.StudentMeData{
                StudentID:      student.StudentID,
                StudentNumber:  student.StudentNumber,
                FacultyID:      student.FacultyID,
                DepartmentID:   student.DepartmentID,
            }
        }
    }

    return c.JSON(fiber.Map{
        "user": response,
    })
}

// UpdateMe - อัพเดทข้อมูล current user
func (h *AuthHandler) UpdateMe(c *fiber.Ctx) error {
    u := c.Locals("current_user")
    if u == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
    }

    user, ok := u.(*models.User)
    if !ok {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user"})
    }

    var req authDto.UpdateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
    }

    // เรียกใช้ service เพื่ออัพเดทข้อมูล
    updatedUser, err := h.AuthService.UpdateUser(c.Context(), user.UserID, &req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "message": "user updated successfully",
        "user": fiber.Map{
            "user_id":        updatedUser.UserID,
            "firstname":      updatedUser.Firstname,
            "lastname":       updatedUser.Lastname,
            "email":          updatedUser.Email,
            "image_path":     updatedUser.ImagePath,
            "provider":       updatedUser.Provider,
            "role_id":        updatedUser.RoleID,
            "campus_id":      updatedUser.CampusID,
            "is_first_login": updatedUser.IsFirstLogin,
            "created_at":     updatedUser.CreatedAt,
            "latest_update":  updatedUser.LatestUpdate,
        },
    })
}