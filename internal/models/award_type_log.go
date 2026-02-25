package models

import "time"

type AwardTypeLog struct {
	TypeLogID    uint      `gorm:"primaryKey;column:type_log_id" json:"type_log_id"`
	FormID       uint      `gorm:"column:form_id;not null;index" json:"form_id"`
	UserID       uint      `gorm:"column:user_id;not null;index" json:"user_id"`
	LogType      string    `gorm:"column:log_type;type:varchar(50)" json:"log_type"` // "award_type_change" | "approval" | "rejection"
	OldValue     string    `gorm:"column:old_value;type:text" json:"old_value"`
	NewValue     string    `gorm:"column:new_value;type:text" json:"new_value"`
	Status       string    `gorm:"column:status;type:varchar(50)" json:"status"` // "approve" | "reject"
	RejectReason string    `gorm:"column:reject_reason;type:text" json:"reject_reason"`
	ChangedAt    time.Time `gorm:"column:changed_at;not null" json:"changed_at"`
}

func (AwardTypeLog) TableName() string {
	return "Award_Type_Log"
}
