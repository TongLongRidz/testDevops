package models

type Role struct {
	RoleID uint `gorm:"primaryKey;column:role_id" json:"role_id"`
	// RoleCode    string `gorm:"type:varchar(50);unique;column:role_code" json:"role_code"`
	RoleName   string `gorm:"type:text;column:role_name" json:"role_name"`
	RoleNameTH string `gorm:"type:text;column:role_name_th" json:"role_name_th"`
	// Description string `gorm:"type:text;column:description" json:"description"`
}

func (Role) TableName() string {
	return "Role"
}
