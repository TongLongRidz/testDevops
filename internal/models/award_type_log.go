package models

import "time"

type AwardTypeLog struct {
	TypeLogID uint      `gorm:"primaryKey;column:type_log_id" json:"type_log_id"`
	FormID    uint      `gorm:"column:form_id;not null;index" json:"form_id"`
	UserID    uint      `gorm:"column:user_id;not null;index" json:"user_id"`
	OldValue  string    `gorm:"column:old_value;type:text;not null" json:"old_value"`
	NewValue  string    `gorm:"column:new_value;type:text;not null" json:"new_value"`
	ChangedAt time.Time `gorm:"column:changed_at;not null" json:"changed_at"`
}

func (AwardTypeLog) TableName() string {
	return "Award_Type_Log"
}
