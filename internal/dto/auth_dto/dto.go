package authdto

import "time"

// RegisterRequest ใช้สำหรับรับข้อมูลตอนสมัครสมาชิก (Manual Register)
type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

// CreateAccountRequest ใช้สำหรับสร้าง account ที่กำหนด role ได้
type CreateAccountRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
	Prefix          string `json:"prefix"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	RoleID          int    `json:"role_id" validate:"required"`
	CampusID        int    `json:"campus_id" validate:"required"`

	// Optional fields สำหรับ role เฉพาะ
	StudentNumber string `json:"student_number"`
	FacultyID     uint   `json:"faculty_id"`
	DepartmentID  uint   `json:"department_id"`
	IsChairman    bool   `json:"is_chairman"`

	OrganizationName     string `json:"organization_name"`
	OrganizationType     string `json:"organization_type"`
	OrganizationLocation string `json:"organization_location"`
	OrganizationPhone    string `json:"organization_phone"`
}

// LoginRequest ใช้สำหรับรับข้อมูลตอนเข้าสู่ระบบด้วย Email/Password
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse คือข้อมูลที่จะส่งกลับไปให้ Frontend เมื่อ Login สำเร็จ
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse คือรายละเอียดของผู้ใช้ที่อนุญาตให้ส่งออกไปภายนอก (Safe Data)
type UserResponse struct {
	UserID       uint      `json:"user_id"`
	Prefix       string    `json:"prefix"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	ImagePath    string    `json:"image_path"`
	Provider     string    `json:"provider"`
	RoleID       int       `json:"role_id"`
	CampusID     int       `json:"campus_id"`
	IsFirstLogin bool      `json:"is_first_login"`
	CreatedAt    time.Time `json:"created_at"`
	LatestUpdate time.Time `json:"latest_update"`
}

// MeResponse คือข้อมูล /me endpoint
type MeResponse struct {
	UserID       uint      `json:"user_id"`
	Prefix       string    `json:"prefix"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	ImagePath    string    `json:"image_path"`
	Provider     string    `json:"provider"`
	RoleID       int       `json:"role_id"`
	CampusID     int       `json:"campus_id"`
	IsFirstLogin bool      `json:"is_first_login"`
	CreatedAt    time.Time `json:"created_at"`
	LatestUpdate time.Time `json:"latest_update"`

	StudentData            *StudentMeData            `json:"student_data,omitempty"`
	OrganizationData       *OrganizationMeData       `json:"organization_data,omitempty"`
	HeadOfDepartmentData   *HeadOfDepartmentMeData   `json:"head_of_department_data,omitempty"`
	AssociateDeanData      *AssociateDeanMeData      `json:"associate_dean_data,omitempty"`
	DeanData               *DeanMeData               `json:"dean_data,omitempty"`
	StudentDevelopmentData *StudentDevelopmentMeData `json:"student_development_data,omitempty"`
	CommitteeData          *CommitteeMeData          `json:"committee_data,omitempty"`
	ChancellorData         *ChancellorMeData         `json:"chancellor_data,omitempty"`
}

type StudentMeData struct {
	StudentID     uint   `json:"student_id"`
	StudentNumber string `json:"student_number"`
	FacultyID     uint   `json:"faculty_id"`
	DepartmentID  uint   `json:"department_id"`
}

type OrganizationMeData struct {
	OrganizationID       uint   `json:"organization_id"`
	OrganizationName     string `json:"organization_name"`
	OrganizationType     string `json:"organization_type"`
	OrganizationLocation string `json:"organization_location"`
	OrganizationPhone    string `json:"organization_phone"`
}

type HeadOfDepartmentMeData struct {
	HodID        uint `json:"hod_id"`
	UserID       uint `json:"user_id"`
	FacultyID    uint `json:"faculty_id"`
	DepartmentID uint `json:"department_id"`
}

type AssociateDeanMeData struct {
	AdID      uint   `json:"ad_id"`
	UserID    uint   `json:"user_id"`
	AdCode    string `json:"ad_code"`
	FacultyID uint   `json:"faculty_id"`
}

type DeanMeData struct {
	DID       uint `json:"d_id"`
	UserID    uint `json:"user_id"`
	FacultyID uint `json:"faculty_id"`
}

type StudentDevelopmentMeData struct {
	SdID   uint `json:"sd_id"`
	UserID uint `json:"user_id"`
}

type CommitteeMeData struct {
	ComID      uint `json:"com_id"`
	UserID     uint `json:"user_id"`
	IsChairman bool `json:"is_chairman"`
}

type ChancellorMeData struct {
	ChancellorID uint `json:"chancellor_id"`
	UserID       uint `json:"user_id"`
}

// UpdateUserRequest ใช้สำหรับอัพเดทข้อมูล current user
type UpdateUserRequest struct {
	Prefix       *string `json:"prefix"`
	Firstname    *string `json:"firstname"`
	Lastname     *string `json:"lastname"`
	ImagePath    *string `json:"image_path"`
	CampusID     *int    `json:"campus_id"`
	RoleID       *int    `json:"role_id"`
	IsFirstLogin *bool   `json:"is_first_login"`
}

// FirstLoginRequest ใช้สำหรับตั้งค่าข้อมูลครั้งแรกของนักศึกษา/องค์กร
type FirstLoginRequest struct {
	Prefix    string `json:"prefix"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	CampusID  int    `json:"campus_id"`

	// Student fields (สำหรับเมล @ku.th)
	StudentNumber string `json:"student_number,omitempty"`
	FacultyID     uint   `json:"faculty_id,omitempty"`
	DepartmentID  uint   `json:"department_id,omitempty"`

	// Organization fields (สำหรับเมลอื่นๆ)
	OrganizationName     string `json:"organization_name,omitempty"`
	OrganizationType     string `json:"organization_type,omitempty"`
	OrganizationLocation string `json:"organization_location,omitempty"`
	OrganizationPhone    string `json:"organization_phone,omitempty"`
}
