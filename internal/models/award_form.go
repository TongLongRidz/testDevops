package models

import (
	"time"
)

type AwardForm struct {
	FormID       uint      `gorm:"primaryKey;column:form_id" json:"form_id"`
	StudentID    int       `gorm:"column:student_id" json:"student_id"` // FK int
	AcademicYear int       `gorm:"column:academic_year" json:"academic_year"` // เก็บเป็นปี พ.ศ. หรือ ค.ศ.
	Semester     int       `gorm:"column:semester" json:"semester"`
	FormStatus   int       `gorm:"column:form_status" json:"form_status"`
	AwardTypeID  int       `gorm:"column:award_type_id" json:"award_type_id"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	CampusID     int       `gorm:"column:campus_id" json:"campus_id"`
}

// TableName กำหนดชื่อตารางให้เป็น "Award_Form"
func (AwardForm) TableName() string {
	return "Award_Form"
}