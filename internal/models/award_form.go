package models

import (
	"time"
)

type AwardForm struct {
	FormID             uint      `gorm:"primaryKey;column:form_id" json:"form_id"`
	UserID             uint      `gorm:"uniqueIndex:idx_user_semester;column:user_id" json:"user_id"`
	StudentFirstname   string    `gorm:"column:student_firstname" json:"student_firstname"`
	StudentLastname    string    `gorm:"column:student_lastname" json:"student_lastname"`
	StudentEmail       string    `gorm:"column:student_email" json:"student_email"`
	StudentNumber      string    `gorm:"column:student_number" json:"student_number"`
	FacultyID          int       `gorm:"column:faculty_id" json:"faculty_id"`
	DepartmentID       int       `gorm:"column:department_id" json:"department_id"`
	CampusID           int       `gorm:"column:campus_id" json:"campus_id"`
	AcademicYear       int       `gorm:"uniqueIndex:idx_user_semester;column:academic_year" json:"academic_year"`
	Semester           int       `gorm:"uniqueIndex:idx_user_semester;column:semester" json:"semester"`
	FormStatusID       int       `gorm:"column:form_status_id" json:"form_status"`
	AwardType          string    `gorm:"column:award_type" json:"award_type"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
	LatestUpdate       time.Time `gorm:"column:latest_update" json:"latest_update"`
	StudentYear        int       `gorm:"column:student_year" json:"student_year"`
	AdvisorName        string    `gorm:"column:advisor_name" json:"advisor_name"`
	StudentPhoneNumber string    `gorm:"column:student_phone_number" json:"student_phone_number"`
	StudentAddress     string    `gorm:"column:student_address" json:"student_address"`
	GPA                float64   `gorm:"column:gpa" json:"gpa"`
	StudentDateOfBirth time.Time `gorm:"column:student_date_of_birth;type:date" json:"student_date_of_birth"`
	OrgName            string    `gorm:"column:org_name" json:"org_name"`
	OrgType            string    `gorm:"column:org_type" json:"org_type"`
	OrgLocation        string    `gorm:"column:org_location" json:"org_location"`
	OrgPhoneNumber     string    `gorm:"column:org_phone_number" json:"org_phone_number"`
	FormDetail         string    `gorm:"column:form_detail" json:"form_detail"`
	RejectReason       string    `gorm:"column:reject_reason" json:"reject_reason"`

	// Relationships
	AwardFiles []AwardFileDirectory `gorm:"foreignKey:FormID" json:"award_files"`
}

// TableName กำหนดชื่อตารางให้เป็น "Award_Form"
func (AwardForm) TableName() string {
	return "Award_Form"
}
