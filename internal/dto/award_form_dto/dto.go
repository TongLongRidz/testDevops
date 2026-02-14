package awardformdto

import (
	"time"
)

// --- Request DTOs ---
type SubmitAwardRequest struct {
	// ข้อมูลเหล่านี้จะดึงจาก Token: StudentID, StudentFirstname, StudentLastname, Email, StudentNumber, FacultyID, DepartmentID, CampusID
	// ข้อมูลเหล่านี้จะดึงจาก Academic Year Service: AcademicYear, Semester
	AwardTypeID int       `json:"award_type_id" binding:"required"`
	StudentYear int       `json:"student_year" binding:"required"`
	AdvisorName string    `json:"advisor_name" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Address     string    `json:"address" binding:"required"`
	GPA         float64   `json:"gpa" binding:"required"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`

	Extracurricular *ExtracurricularRequest `json:"extracurricular_detail,omitempty"`
	Creativity      *CreativityRequest      `json:"creativity_detail,omitempty"`
	GoodBehavior    *GoodBehaviorRequest    `json:"good_behavior_detail,omitempty"`
}

type ExtracurricularRequest struct {
	QualificationType string    `json:"qualification_type"`
	DateReceived      time.Time `json:"date_received"`
	TeamName          string    `json:"team_name"`
	ProjectTitle      string    `json:"project_title"`
	Prize             string    `json:"prize"`
	OrganizedBy       string    `json:"organized_by"`
	CompetitionLevel  string    `json:"competition_level"`
	ActivityCategory  string    `json:"activity_category"`
	CompetitionName   string    `json:"competition_name"`
}

type CreativityRequest struct {
	DateReceived     time.Time `json:"date_received"`
	TeamName         string    `json:"team_name"`
	ProjectTitle     string    `json:"project_title"`
	Prize            string    `json:"prize"`
	OrganizedBy      string    `json:"organized_by"`
	CompetitionLevel string    `json:"competition_level"`
	ActivityCategory string    `json:"activity_category"`
	CompetitionName  string    `json:"competition_name"`
}

type GoodBehaviorRequest struct {
	// ตอนนี้ไม่มีฟิลด์เพิ่มเติม
}

// --- Response DTOs ---
type AwardFormResponse struct {
	FormID           uint      `json:"form_id"`
	StudentID        int       `json:"student_id"`
	StudentFirstname string    `json:"student_firstname"`
	StudentLastname  string    `json:"student_lastname"`
	Email            string    `json:"email"`
	StudentNumber    string    `json:"student_number"`
	FacultyID        int       `json:"faculty_id"`
	DepartmentID     int       `json:"department_id"`
	CampusID         int       `json:"campus_id"`
	AcademicYear     int       `json:"academic_year"`
	Semester         int       `json:"semester"`
	FormStatusID     int       `json:"form_status_id"`
	AwardTypeID      int       `json:"award_type_id"`
	AwardTypeName    string    `json:"award_type_name"`
	CreatedAt        time.Time `json:"created_at"`
	LatestUpdate     time.Time `json:"latest_update"`
	StudentYear      int       `json:"student_year"`
	AdvisorName      string    `json:"advisor_name"`
	PhoneNumber      string    `json:"phone_number"`
	Address          string    `json:"address"`
	GPA              float64   `json:"gpa"`
	DateOfBirth      time.Time `json:"date_of_birth"`

	// ข้อมูลรายละเอียด (จะถูกเติมเฉพาะประเภทที่ตรงกัน)
	Detail interface{} `json:"detail,omitempty"`

	// ข้อมูลไฟล์แนบ
	Files []FileResponse `json:"files,omitempty"`
}

type FileResponse struct {
	FileDirID uint   `json:"file_dir_id"`
	FileName  string `json:"file_name"`
	FileType  string `json:"file_type"`
	FileSize  int64  `json:"file_size"`
	FilePath  string `json:"file_path"`
}

// --- Search & Pagination DTOs ---
type SearchAwardRequest struct {
	Keyword     string `query:"keyword"`      // ค้นหาใน firstname, lastname, studentNumber, semester, year, type
	Date        string `query:"date"`         // กรองตามวันที่ (format: YYYY-MM-DD)
	StudentYear int    `query:"student_year"` // กรองตามชั้นปี
	Page        int    `query:"page"`         // หน้าปัจจุบัน (default: 1)
	Limit       int    `query:"limit"`        // จำนวนต่อหน้า (default: 10)
	Arrangement string `query:"arrangement"` // เรียงลำดับ: asc หรือ desc (default: desc)
}

type PaginatedAwardResponse struct {
	Data       []AwardFormResponse `json:"data"`
	Pagination PaginationMeta      `json:"pagination"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	Limit       int   `json:"limit"`
}

// --- Award Form Log DTOs ---
type CreateAwardFormLogRequest struct {
	FormID    uint   `json:"form_id" binding:"required"`
	FieldName string `json:"field_name" binding:"required"`
	OldValue  string `json:"old_value"`
	NewValue  string `json:"new_value"`
}

type UpdateAwardTypeRequest struct {
	AwardTypeID int `json:"award_type_id" binding:"required"`
}

type UpdateFormStatusRequest struct {
	FormStatusID int `json:"form_status_id" binding:"required"`
}

type AwardFormLogResponse struct {
	LogID      uint      `json:"log_id"`
	FormID     uint      `json:"form_id"`
	FieldName  string    `json:"field_name"`
	OldValue   string    `json:"old_value"`
	NewValue   string    `json:"new_value"`
	ChangedBy  *int      `json:"changed_by"`
	CreatedAt  time.Time `json:"created_at"`
	LatestEdit time.Time `json:"latest_edit"`
}
