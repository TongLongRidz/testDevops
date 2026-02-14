package authdto

import "time"

// RegisterRequest ใช้สำหรับรับข้อมูลตอนสมัครสมาชิก (Manual Register)
type RegisterRequest struct {
    Email           string `json:"email" validate:"required,email"`
    Password        string `json:"password" validate:"required,min=6"`
    ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

// RegisterWithRoleRequest ใช้สำหรับสร้าง account ที่กำหนด role ได้
type RegisterWithRoleRequest struct {
    Email           string `json:"email" validate:"required,email"`
    Password        string `json:"password" validate:"required,min=6"`
    ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
    Firstname       string `json:"firstname"`
    Lastname        string `json:"lastname"`
    RoleID          int    `json:"role_id" validate:"required"`
    CampusID        int    `json:"campus_id" validate:"required"`

    // Optional fields สำหรับ role เฉพาะ
    StudentNumber string `json:"student_number"`
    FacultyID     uint   `json:"faculty_id"`
    DepartmentID  uint   `json:"department_id"`
    AdCode        string `json:"ad_code"`
    IsChairman    *bool  `json:"is_chairman"`
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
    UserID        uint                `json:"user_id"`
    Firstname     string              `json:"firstname"`
    Lastname      string              `json:"lastname"`
    Email         string              `json:"email"`
    ImagePath     string              `json:"image_path"`
    Provider      string              `json:"provider"`
    RoleID        int                 `json:"role_id"`
    CampusID      int                 `json:"campus_id"`
    IsFirstLogin  bool                `json:"is_first_login"`
    CreatedAt     time.Time           `json:"created_at"`
    LatestUpdate  time.Time           `json:"latest_update"`
    StudentData   *StudentMeData      `json:"student_data,omitempty"`
}

type StudentMeData struct {
    StudentID      uint   `json:"student_id"`
    StudentNumber  string `json:"student_number"`
    FacultyID      uint   `json:"faculty_id"`
    DepartmentID   uint   `json:"department_id"`
}

// UpdateUserRequest ใช้สำหรับอัพเดทข้อมูล current user
type UpdateUserRequest struct {
    Firstname    *string `json:"firstname"`
    Lastname     *string `json:"lastname"`
    ImagePath    *string `json:"image_path"`
    CampusID     *int    `json:"campus_id"`
    RoleID       *int    `json:"role_id"`
    IsFirstLogin *bool   `json:"is_first_login"`
}