package auth

import (
	"fmt"
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
        "user":  fiber.Map{"id": user.UserID, "email": user.Email, "role": user.RoleID},
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

    user, ok := u.(*models.User) // ตรวจว่าตรงกับ type ของ project
    if !ok {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user"})
    }

    return c.JSON(fiber.Map{
        "user": fiber.Map{
            "id":    fmt.Sprint(user.UserID),
            "email": user.Email,
            "roleID":  user.RoleID,
        },
    })
}