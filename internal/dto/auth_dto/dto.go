package dto

// RegisterRequest ใช้สำหรับรับข้อมูลตอนสมัครสมาชิก (Manual Register)
type RegisterRequest struct {
    Email           string `json:"email" validate:"required,email"`
    Password        string `json:"password" validate:"required,min=6"`
    ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
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
	UserID       uint   `json:"userID"`
	Email        string `json:"email"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	ImagePath    string `json:"imagePath"`
	RoleID       int    `json:"roleID"`
	IsFirstLogin bool   `json:"isFirstLogin"` // แจ้ง Frontend ว่าต้องไปหน้าตั้งค่าโปรไฟล์หรือไม่
}