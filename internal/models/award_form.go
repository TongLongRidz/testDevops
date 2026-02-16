package models

import (
	"time"
)

type AwardForm struct {
	FormID           uint      `gorm:"primaryKey;column:form_id" json:"form_id"`
	StudentID        int       `gorm:"uniqueIndex:idx_student_semester;column:student_id" json:"student_id"`
	StudentFirstname string    `gorm:"column:student_firstname" json:"student_firstname"`
	StudentLastname  string    `gorm:"column:student_lastname" json:"student_lastname"`
	Email            string    `gorm:"column:email" json:"email"`
	StudentNumber    string    `gorm:"column:student_number" json:"student_number"`
	FacultyID        int       `gorm:"column:faculty_id" json:"faculty_id"`
	DepartmentID     int       `gorm:"column:department_id" json:"department_id"`
	CampusID         int       `gorm:"column:campus_id" json:"campus_id"`
	AcademicYear     int       `gorm:"uniqueIndex:idx_student_semester;column:academic_year" json:"academic_year"`
	Semester         int       `gorm:"uniqueIndex:idx_student_semester;column:semester" json:"semester"`
	FormStatusID     int       `gorm:"column:form_status_id" json:"form_status_id"`
	AwardTypeID      int       `gorm:"column:award_type_id" json:"award_type_id"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	LatestUpdate     time.Time `gorm:"column:latest_update" json:"latest_update"`
	StudentYear      int       `gorm:"column:student_year" json:"student_year"`
	AdvisorName      string    `gorm:"column:advisor_name" json:"advisor_name"`
	PhoneNumber      string    `gorm:"column:phone_number" json:"phone_number"`
	Address          string    `gorm:"column:address" json:"address"`
	GPA              float64   `gorm:"column:gpa" json:"gpa"`
	DateOfBirth      time.Time `gorm:"column:date_of_birth;type:date" json:"date_of_birth"`
	RejectReason     string    `gorm:"column:reject_reason" json:"reject_reason"`

	// Relationship
	Student    Student              `gorm:"foreignKey:StudentID" json:"student"`
	AwardType  AwardType            `gorm:"foreignKey:AwardTypeID" json:"award_type"`
	AwardFiles []AwardFileDirectory `gorm:"foreignKey:FormID" json:"award_files"`

	Extracurricular *ExtracurricularActivity `gorm:"foreignKey:FormID" json:"extracurricular,omitempty"`
	GoodBehavior    *GoodBehavior            `gorm:"foreignKey:FormID" json:"good_behavior,omitempty"`
	Creativity      *CreativityInnovation    `gorm:"foreignKey:FormID" json:"creativity,omitempty"`
}

// TableName กำหนดชื่อตารางให้เป็น "Award_Form"
func (AwardForm) TableName() string {
	return "Award_Form"
}
