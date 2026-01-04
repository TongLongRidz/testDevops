package models

type AwardType struct {
	AwardTypeID uint   `gorm:"primaryKey;column:award_type_id" json:"award_type_id"`
	AwardName   string `gorm:"type:varchar(255);column:award_name" json:"award_name"`
	Description string `gorm:"type:text;column:description" json:"description"`
}

// TableName กำหนดชื่อตารางให้เป็น "Award_Type"
func (AwardType) TableName() string {
	return "Award_Type"
}