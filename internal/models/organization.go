package models

type Organization struct {
	OrganizationID         	uint   `gorm:"primaryKey;column:organization_id" json:"organization_id"`
	UserID        			uint   `gorm:"column:user_id;uniqueIndex" json:"user_id"` // 1 User เป็น 1 Student
	User          			User   `gorm:"foreignKey:UserID"`                         // ความสัมพันธ์กับ User
	OrganizationName        string `gorm:"column:org_name" json:"org_name"`
	OrganizationType        string `gorm:"column:org_type" json:"org_type"`
	OrganizationLocation    string `gorm:"column:org_location" json:"org_location"`
	OrganizationPhoneNumber string `gorm:"column:org_phone_number" json:"org_phone_number"`
}

// TableName กำหนดชื่อตารางให้เป็น "Organization"
func (Organization) TableName() string {
	return "Organization"
}