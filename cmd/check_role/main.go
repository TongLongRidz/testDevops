package main

import (
	"backend/config"
	"backend/internal/models"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not found")
	}

	db := config.ConnectDB()

	var user models.User
	if err := db.Where("email = ?", "dev@test.com").First(&user).Error; err != nil {
		log.Fatal("User not found:", err)
	}

	fmt.Printf("User: %s, RoleID: %d\n", user.Email, user.RoleID)
}
