package models

type Faculty struct {
	FacultyID   uint   `gorm:"primaryKey;column:faculty_id" json:"faculty_id"`
	// FacultyCode string `gorm:"type:text;column:faculty_code" json:"faculty_code"`
	FacultyName string `gorm:"type:text;column:faculty_name" json:"faculty_name"`
}

func (Faculty) TableName() string {
	return "Faculty"
}