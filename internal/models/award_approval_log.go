package models

import "time"

type AwardApprovalLog struct {
	ApprovalLogID  uint      `gorm:"primaryKey;column:approval_log_id" json:"approval_log_id"`
	FormID         uint      `gorm:"column:form_id;not null;index" json:"form_id"`
	UserID         uint      `gorm:"column:user_id;not null;index" json:"user_id"`
	ApprovalStatus string    `gorm:"column:approval_status;type:varchar(10);not null" json:"approval_status"`
	RejectReason   string    `gorm:"column:reject_reason;type:text" json:"reject_reason,omitempty"`
	ApprovedAt     time.Time `gorm:"column:approved_at;not null" json:"approved_at"`
}

func (AwardApprovalLog) TableName() string {
	return "Award_Approval_Log"
}
