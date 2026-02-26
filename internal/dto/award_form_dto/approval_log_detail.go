package awardformdto

import "time"

// ApprovalLogDetailResponse is the response DTO for a single approval log detail
// (extends ApprovalLogHistoryResponse with more fields if needed)
type ApprovalLogDetailResponse struct {
	ApprovalLogID  uint      `json:"approval_log_id"`
	FormID         uint      `json:"form_id"`
	UserID         uint      `json:"user_id"`
	ApprovalStatus string    `json:"approval_status"`
	RejectReason   string    `json:"reject_reason,omitempty"`
	ApprovedAt     time.Time `json:"approved_at"`
}
