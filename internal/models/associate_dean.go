package models

type AssociateDean struct {
	AdID      uint    `gorm:"primaryKey;column:ad_id" json:"ad_id"`
	UserID    uint    `gorm:"column:user_id;uniqueIndex" json:"user_id"` // 1 User เป็น 1 Associate Dean
	User      User    `gorm:"foreignKey:UserID"`                         // ความสัมพันธ์กับ User
	AdCode    string  `gorm:"type:varchar(50);column:ad_code" json:"ad_code"`
	FacultyID uint    `gorm:"column:faculty_id" json:"faculty_id"`
	Faculty   Faculty `gorm:"foreignKey:FacultyID"`
}

func (AssociateDean) TableName() string {
	return "Associate_Dean"
}
