package models

type Committee struct {
	ComID      uint `gorm:"primaryKey;column:com_id" json:"com_id"`
	UserID     uint `gorm:"column:user_id;uniqueIndex" json:"user_id"` // 1 User เป็น 1 Committee
	User       User `gorm:"foreignKey:UserID"`                         // ความสัมพันธ์กับ User
	IsChairman bool `gorm:"type:boolean;column:is_chairman;default:false" json:"is_chairman"`
	// ComCode     string `gorm:"type:varchar(50);column:com_code" json:"com_code"`
}

func (Committee) TableName() string {
	return "Committee"
}
