package awardformdto

import (
	"time"
)

// --- Request DTOs ---
type SubmitAwardRequest struct {
	// === Role: STUDENT (RoleID = 1) ===
	// System Auto-Fill จาก Token: UserID, StudentFirstname, StudentLastname, StudentEmail, StudentNumber, FacultyID, DepartmentID, CampusID
	// System Auto-Fill จาก Academic Year Service: AcademicYear, Semester
	// Student กรอก:
	AwardType          string    `json:"award_type"`
	StudentYear        int       `json:"student_year"`
	AdvisorName        string    `json:"advisor_name"`
	StudentPhoneNumber string    `json:"student_phone_number"`
	StudentAddress     string    `json:"student_address"`
	GPA                float64   `json:"gpa"`
	StudentDateOfBirth time.Time `json:"student_date_of_birth"`
	FormDetail         string    `json:"form_detail"`

	// === Role: ORGANIZATION (RoleID = 9) ===
	// System Auto-Fill จาก Token: UserID, CampusID
	// System Auto-Fill จาก Organization & Academic Year Service: OrgName, OrgType, OrgLocation, OrgPhoneNumber, AcademicYear, Semester
	// Organization กรอก:
	StudentFirstname string `json:"student_firstname"`
	StudentLastname  string `json:"student_lastname"`
	StudentEmail     string `json:"student_email"`
	StudentNumber    string `json:"student_number"`
	FacultyID        int    `json:"faculty_id"`
	DepartmentID     int    `json:"department_id"`
}

// --- Response DTOs ---
type AwardFormResponse struct {
	FormID             uint      `json:"form_id"`
	UserID             uint      `json:"user_id"`
	StudentFirstname   string    `json:"student_firstname"`
	StudentLastname    string    `json:"student_lastname"`
	StudentEmail       string    `json:"student_email"`
	StudentNumber      string    `json:"student_number"`
	FacultyID          int       `json:"faculty_id"`
	DepartmentID       int       `json:"department_id"`
	CampusID           int       `json:"campus_id"`
	AcademicYear       int       `json:"academic_year"`
	Semester           int       `json:"semester"`
	FormStatusID       int       `json:"form_status"`
	AwardType          string    `json:"award_type"`
	CreatedAt          time.Time `json:"created_at"`
	LatestUpdate       time.Time `json:"latest_update"`
	StudentYear        int       `json:"student_year"`
	AdvisorName        string    `json:"advisor_name"`
	StudentPhoneNumber string    `json:"student_phone_number"`
	StudentAddress     string    `json:"student_address"`
	GPA                float64   `json:"gpa"`
	StudentDateOfBirth time.Time `json:"student_date_of_birth"`
	OrgName            string    `json:"org_name"`
	OrgType            string    `json:"org_type"`
	OrgLocation        string    `json:"org_location"`
	OrgPhoneNumber     string    `json:"org_phone_number"`
	FormDetail         string    `json:"form_detail"`
	RejectReason       string    `json:"reject_reason"`

	// ข้อมูลไฟล์แนบ
	Files []FileResponse `json:"files,omitempty"`
}

type FileResponse struct {
	FileDirID uint   `json:"file_dir_id"`
	FileType  string `json:"file_type"`
	FileSize  int64  `json:"file_size"`
	FilePath  string `json:"file_path"`
}

// --- Search & Pagination DTOs ---
type SearchAwardRequest struct {
	Keyword     string `query:"keyword"`      // ค้นหาใน firstname, lastname, studentNumber, semester, year, award_type
	Date        string `query:"date"`         // กรองตามวันที่ (format: YYYY-MM-DD)
	StudentYear int    `query:"student_year"` // กรองตามชั้นปี
	AwardType   string `query:"award_type"`   // กรองตามประเภทรางวัล
	Page        int    `query:"page"`         // หน้าปัจจุบัน (default: 1)
	Limit       int    `query:"limit"`        // จำนวนต่อหน้า (default: 10)
	Arrangement string `query:"arrangement"`  // เรียงลำดับ: asc หรือ desc (default: desc)
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
	AwardType string `json:"award_type" binding:"required"`
}

type UpdateFormStatusRequest struct {
	FormStatusID int `json:"form_status" binding:"required"`
}

type AwardFormLogResponse struct {
	LogID     uint      `json:"log_id"`
	FormID    uint      `json:"form_id"`
	FieldName string    `json:"field_name"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	ChangedBy int       `json:"changed_by"`
	CreatedAt time.Time `json:"created_at"`
}
