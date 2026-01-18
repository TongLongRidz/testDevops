package config

import (
	"fmt"
	"os"
	// "strings"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // เพิ่มอันนี้
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	
	// เพิ่ม Logger: Info เพื่อดู SQL ทุกคำสั่งที่ส่งไป DB
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        DisableForeignKeyConstraintWhenMigrating: true, 
    })

	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	return db
}
