package models

type Department struct {
	DepartmentID   uint    `gorm:"primaryKey;column:department_id" json:"department_id"`
	FacultyID      uint    `gorm:"column:faculty_id" json:"faculty_id"` // Foreign Key ไปยัง Faculty
	Faculty        Faculty `gorm:"foreignKey:FacultyID"`               // ความสัมพันธ์กับ Faculty
	// DepartmentCode string  `gorm:"type:varchar(50);unique;column:department_code" json:"department_code"`
	DepartmentName string  `gorm:"type:text;column:department_name" json:"department_name"`
}

func (Department) TableName() string {
	return "Department"
}