package models

import "time"

type CommitteeVoteLog struct {
	VoteLogID uint      `gorm:"primaryKey;column:vote_log_id" json:"vote_log_id"`
	FormID    uint      `gorm:"column:form_id;not null;index" json:"form_id"`
	UserID    uint      `gorm:"column:user_id;not null;index" json:"user_id"`
	Operation string    `gorm:"column:operation;type:varchar(50);not null" json:"operation"` // "approve" | "reject"
	VotedAt   time.Time `gorm:"column:voted_at;not null" json:"voted_at"`
}

func (CommitteeVoteLog) TableName() string {
	return "Committee_Vote_Log"
}
