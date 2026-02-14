package models

import "time"

type AwardFormLog struct {
	LogID      uint      `gorm:"primaryKey;column:log_id" json:"log_id"`
	FormID     uint      `gorm:"index;column:form_id" json:"form_id"`
	FieldName  string    `gorm:"column:field_name" json:"field_name"`
	OldValue   string    `gorm:"column:old_value" json:"old_value"`
	NewValue   string    `gorm:"column:new_value" json:"new_value"`
	ChangedBy  *int      `gorm:"column:changed_by" json:"changed_by"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	LatestEdit time.Time `gorm:"column:latest_edit" json:"latest_edit"`

	// Relationship
	AwardForm AwardForm `gorm:"foreignKey:FormID" json:"award_form"`
}

func (AwardFormLog) TableName() string {
	return "Award_Form_Log"
}
