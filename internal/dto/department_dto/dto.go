package departmentdto

type CreateDepartment struct {
	FacultyID      uint   `json:"faculty_id" binding:"required"`
	DepartmentName string `json:"department_name" binding:"required,min=1,max=255"`
}

type UpdateDepartment struct {
	FacultyID      uint   `json:"faculty_id" binding:"required"`
	DepartmentName string `json:"department_name" binding:"required,min=1,max=255"`
}

type DepartmentResponse struct {
	DepartmentID   uint   `json:"department_id"`
	FacultyID      uint   `json:"faculty_id"`
	DepartmentName string `json:"department_name"`
}
