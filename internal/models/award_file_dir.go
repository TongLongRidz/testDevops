package models

import (
	"time"
)

type AwardFileDirectory struct {
	FileDirID  uint      `gorm:"primaryKey;column:file_dir_id" json:"file_dir_id"`
	FormID     int       `gorm:"column:form_id" json:"form_id"` // FK int
	FileType   string    `gorm:"type:varchar(50);column:file_type" json:"file_type"`
	FileSize   int64     `gorm:"column:file_size" json:"file_size"` // ใช้ int64 สำหรับเก็บขนาดไฟล์ (bytes)
	FilePath   string    `gorm:"type:text;column:file_path" json:"file_path"`
	UploadedAt time.Time `gorm:"column:uploaded_at" json:"uploaded_at"`
}

// TableName กำหนดชื่อตารางให้เป็น "Award_File_Directory"
func (AwardFileDirectory) TableName() string {
	return "Award_File_Directory"
}