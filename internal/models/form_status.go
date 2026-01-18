package models

type FormStatus struct {
	FormStatusID   uint   `gorm:"primaryKey;column:form_status_id" json:"form_status_id"`
	FormStatusName string `gorm:"column:form_status_name" json:"form_status_name"`
}

// TableName กำหนดชื่อตารางให้เป็น "Form_Status"
func (FormStatus) TableName() string {
	return "Form_Status"
}

// ต้องทำให้ Award Form รองรับ FK กับแก้อันนี้ด้วย และก็ยังไม่ได้ทำใน main.go (AutoMigrate)