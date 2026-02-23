package models

type Chancellor struct {
	ChancellorID uint `gorm:"primaryKey;column:chancellor_id" json:"chancellor_id"`
	UserID       uint `gorm:"column:user_id;uniqueIndex" json:"user_id"` // 1 User เป็น 1 Chancellor
	User         User `gorm:"foreignKey:UserID"`                         // ความสัมพันธ์กับ User
}

func (Chancellor) TableName() string {
	return "Chancellor"
}
