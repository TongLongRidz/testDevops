package models

type Committee struct {
	ComID  uint `gorm:"primaryKey;column:com_id" json:"com_id"`
	UserID uint `gorm:"column:user_id;uniqueIndex" json:"user_id"` // 1 User เป็น 1 Committee
	User   User `gorm:"foreignKey:UserID"`                         // ความสัมพันธ์กับ User
	// ComCode     string `gorm:"type:varchar(50);column:com_code" json:"com_code"`
	IsChairman bool `gorm:"type:boolean;column:is_chairman" json:"is_chairman"`
}

func (Committee) TableName() string {
	return "Committee"
}
