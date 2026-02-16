package models

import (
	"time"
)

type User struct {
	UserID         uint   `gorm:"primaryKey;column:user_id" json:"user_id"` // แก้จาก ID
	Prefix         string `gorm:"type:varchar(50);column:prefix" json:"prefix"`
	Firstname      string `gorm:"type:varchar(100);column:firstname" json:"firstname"` // Done
	Lastname       string `gorm:"type:varchar(100);column:lastname" json:"lastname"`   // Done
	Email          string `gorm:"uniqueIndex;not null;column:email" json:"email"`      // Done
	HashedPassword string `gorm:"type:string;column:hashed_password" json:"-"`         // ใส่ "-" เพื่อไม่ให้ส่ง password ออกไปทาง JSON
	// Address        string    `gorm:"type:text;column:address" json:"address"`
	// PhoneNumber    string    `gorm:"type:varchar(20);column:phone_number" json:"phone_number"`
	ImagePath    string    `gorm:"type:text;column:image_path" json:"image_path"`                     // แก้จาก Image
	Provider     string    `gorm:"type:varchar(50);default:'manual';column:provider" json:"provider"` // Done
	RoleID       int       `gorm:"column:role_id" json:"role_id"`
	CampusID     int       `gorm:"column:campus_id" json:"campus_id"`
	IsFirstLogin bool      `gorm:"column:is_first_login;default:true" json:"is_first_login"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"` // Done
	LatestUpdate time.Time `gorm:"column:latest_update" json:"latest_update"`
}

// TableName กำหนดชื่อตารางให้เป็น "User"
func (User) TableName() string {
	return "User"
}
