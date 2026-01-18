package models

import (
	"time"
)

type AwardFileDirectory struct {
    FileDirID  uint      `gorm:"primaryKey;column:file_dir_id" json:"file_dir_id"`
    FormID     uint       `gorm:"column:form_id" json:"form_id"`
    FileType   string    `gorm:"column:file_type" json:"file_type"`
    FileSize   int64     `gorm:"column:file_size" json:"file_size"` // หน่วยเป็น Bytes
    FilePath   string    `gorm:"column:file_path" json:"file_path"`
    UploadedAt time.Time `gorm:"column:uploaded_at" json:"uploaded_at"`

	// Relationship
	AwardForm *AwardForm `gorm:"foreignKey:FormID" json:"-"`
}

// TableName กำหนดชื่อตารางให้เป็น "Award_File_Directory"
func (AwardFileDirectory) TableName() string {
	return "Award_File_Directory"
}