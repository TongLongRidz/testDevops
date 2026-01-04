package middleware

import (
    "fmt"
    "os"
    "strings"

    "backend/internal/repository"
    "github.com/gofiber/fiber/v2"
    jwt "github.com/golang-jwt/jwt/v5"
)

func RequireAuth(userRepo repository.UserRepository) fiber.Handler {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "dev-secret"
    }

    return func(c *fiber.Ctx) error {
        // อ่าน token จาก cookie หรือ Authorization header
        tokenStr := c.Cookies("token")
        if tokenStr == "" {
            auth := c.Get("Authorization")
            if strings.HasPrefix(auth, "Bearer ") {
                tokenStr = strings.TrimPrefix(auth, "Bearer ")
            }
        }
        if tokenStr == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
        }

        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
        }

        emailI, _ := claims["email"]
        email, _ := emailI.(string)
        if email == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token payload"})
        }

        user, err := userRepo.GetUserByEmail(email)
        if err != nil || user == nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
        }

        c.Locals("current_user", user)
        return c.Next()
    }
}