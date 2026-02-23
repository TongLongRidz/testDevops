package userdto

import "time"

// ใช้ pointer เพื่อทำ Partial update (เฉพาะฟิลด์ที่ส่งมา)
type EditUserRequest struct {
    Prefix       *string `json:"prefix,omitempty"`
    Firstname    *string `json:"firstname,omitempty"`
    Lastname     *string `json:"lastname,omitempty"`
    Email        *string `json:"email,omitempty"`
    ImagePath    *string `json:"image_path,omitempty"`
    Provider     *string `json:"provider,omitempty"`
    RoleID       *int    `json:"role_id,omitempty"`
    CampusID     *int    `json:"campus_id,omitempty"`
    IsFirstLogin *bool   `json:"is_first_login,omitempty"`
}

type ChangeCommitteeRoleRequest struct {
    IsChairman bool `json:"is_chairman"`
}

// UserResponse ส่วนข้อมูล user ที่ปลอดภัยในการส่งออก
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