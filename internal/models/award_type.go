package models

type AwardType struct {
    AwardTypeID uint   `gorm:"primaryKey;column:award_type_id" json:"award_type_id"`
    AwardName   string `gorm:"column:award_name" json:"award_name"`
    Description string `gorm:"column:description" json:"description"`

	// Relationship
	AwardForms []AwardForm `gorm:"foreignKey:AwardTypeID" json:"award_forms,omitempty"`
}

// TableName กำหนดชื่อตารางให้เป็น "Award_Type"
func (AwardType) TableName() string {
	return "Award_Type"
}