package models

type Dean struct {
	DID       uint    `gorm:"primaryKey;column:d_id" json:"d_id"`
	UserID    uint    `gorm:"column:user_id;uniqueIndex" json:"user_id"` // 1 User เป็น 1 Dean
	User      User    `gorm:"foreignKey:UserID"`                         // ความสัมพันธ์กับ User
	// DCode     string  `gorm:"type:varchar(50);column:d_code" json:"d_code"`
	FacultyID uint    `gorm:"column:faculty_id" json:"faculty_id"`
	Faculty   Faculty `gorm:"foreignKey:FacultyID"`
}

func (Dean) TableName() string {
	return "Dean"
}
