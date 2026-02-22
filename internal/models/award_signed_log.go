package models

import "time"

type AwardSignedLog struct {
	SignedLogID uint      `gorm:"primaryKey;column:signed_log_id" json:"signed_log_id"`
	FormID      uint      `gorm:"column:form_id;not null;index" json:"form_id"`
	UserID      uint      `gorm:"column:user_id;not null;index" json:"user_id"`
	SignedAt    time.Time `gorm:"column:signed_at;not null" json:"signed_at"`
}

func (AwardSignedLog) TableName() string {
	return "Award_Signed_Log"
}
