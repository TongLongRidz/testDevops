package models

type StudentDevelopment struct {
	SdID   uint   `gorm:"primaryKey;column:sd_id" json:"sd_id"`
	UserID uint   `gorm:"column:user_id;uniqueIndex" json:"user_id"` // 1 User เป็น 1 Student Development
	User   User   `gorm:"foreignKey:UserID"`                         // ความสัมพันธ์กับ User
	// SdCode string `gorm:"type:varchar(50);column:sd_code" json:"sd_code"`
}

func (StudentDevelopment) TableName() string {
	return "Student_Development"
}
