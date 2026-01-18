package models

type GoodBehavior struct {
    GoodBeID uint `gorm:"primaryKey;column:good_be_id" json:"good_be_id"`
    FormID   uint `gorm:"column:form_id" json:"form_id"`
}

func (GoodBehavior) TableName() string {
    return "Good_Behavior"
}