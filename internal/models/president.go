package models

type President struct {
	PresID uint `gorm:"primaryKey;column:pres_id" json:"pres_id"`
	UserID uint `gorm:"column:user_id;uniqueIndex" json:"user_id"` // 1 User เป็น 1 President
	User   User `gorm:"foreignKey:UserID"`                         // ความสัมพันธ์กับ User
}

func (President) TableName() string {
	return "President"
}
