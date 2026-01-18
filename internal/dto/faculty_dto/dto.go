package facultydto

type CreateFaculty struct {
	FacultyName string `json:"faculty_name" binding:"required,min=1,max=255"`
}

type UpdateFaculty struct {
	FacultyName string `json:"faculty_name" binding:"required,min=1,max=255"`
}

type FacultyResponse struct {
	FacultyID   uint   `json:"faculty_id"`
	FacultyName string `json:"faculty_name"`
}
