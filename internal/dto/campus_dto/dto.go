package campusdto

// CampusResponse คือรายละเอียดของวิทยาเขตที่ส่งออกไปภายนอก
type CampusResponse struct {
	CampusID   uint   `json:"campusID"`
	CampusCode string `json:"campusCode"`
	CampusName string `json:"campusName"`
}

// CreateCampusRequest ใช้สำหรับรับข้อมูลตอนสร้างวิทยาเขตใหม่
type CreateCampusRequest struct {
	CampusCode string `json:"campusCode" validate:"required"`
	CampusName string `json:"campusName" validate:"required"`
}

// UpdateCampusRequest ใช้สำหรับรับข้อมูลตอนแก้ไขวิทยาเขต
type UpdateCampusRequest struct {
	CampusCode string `json:"campusCode"`
	CampusName string `json:"campusName"`
}