package models

type Campus struct {
    CampusID   uint   `gorm:"primaryKey;column:campus_id" json:"campus_id"`
    CampusCode string `gorm:"type:varchar(50);column:campus_code" json:"campus_code"`
    CampusName string `gorm:"type:varchar(255);column:campus_name" json:"campus_name"`
}

// TableName กำหนดชื่อตารางให้เป็น "Campuses" ตามในรูปภาพ
func (Campus) TableName() string {
    return "Campuses"
}