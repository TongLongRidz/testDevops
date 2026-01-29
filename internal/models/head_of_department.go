package models

type HeadOfDepartment struct {
	HodID        uint       `gorm:"primaryKey;column:hod_id" json:"hod_id"`
	UserID       uint       `gorm:"column:user_id;uniqueIndex" json:"user_id"` // 1 User เป็น 1 Head of Department
	User         User       `gorm:"foreignKey:UserID"`                         // ความสัมพันธ์กับ User
	// HodCode      string     `gorm:"type:varchar(50);column:hod_code" json:"hod_code"`
	FacultyID    uint       `gorm:"column:faculty_id" json:"faculty_id"`
	Faculty      Faculty    `gorm:"foreignKey:FacultyID"`
	DepartmentID uint       `gorm:"column:department_id" json:"department_id"`
	Department   Department `gorm:"foreignKey:DepartmentID"`
}

func (HeadOfDepartment) TableName() string {
	return "Head_Of_Department"
}
