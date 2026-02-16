package awardtypedto

type AwardTypeResponse struct {
	AwardTypeID uint   `json:"award_type_id"`
	AwardName   string `json:"award_name"`
	Description string `json:"description"`
}
