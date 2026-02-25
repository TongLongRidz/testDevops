package models

import (
	"time"
)

type AcademicYear struct {
	AcademicYearID uint      `gorm:"primaryKey;column:academic_year_id" json:"academic_year_id"`
	Year           int       `gorm:"column:year" json:"year"`
	Semester       int       `gorm:"column:semester" json:"semester"`
	StartDate      time.Time `gorm:"type:date;column:start_date" json:"start_date"` // เก็บแค่วันที่
	EndDate        time.Time `gorm:"type:date;column:end_date" json:"end_date"`     // เก็บแค่วันที่
	// CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	// UpdatedAt      time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (AcademicYear) TableName() string {
	return "Academic_Year"
}
