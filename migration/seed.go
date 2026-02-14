package migration

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

func SeedFormStatus(db *gorm.DB) error {
	formStatuses := []models.FormStatus{
		{FormStatusName: "Send_to_HoD"},
		{FormStatusName: "Send_to_AsD"},
		{FormStatusName: "Rejected_by_HoD"},
		{FormStatusName: "Send_to_Dean"},
		{FormStatusName: "Rejected_by_AsD"},
		{FormStatusName: "Send_to_StD"},
		{FormStatusName: "Rejected_by_Dean"},
		{FormStatusName: "Send_to_Com"},
		{FormStatusName: "Rejected_by_StD"},
		{FormStatusName: "Send_to_ComPres"},
		{FormStatusName: "Rejected_by_Com"},
		{FormStatusName: "Send_to_Chan"},
		{FormStatusName: "Rejected_by_ComPres"},
		{FormStatusName: "Accepted_By_Chan"},
		{FormStatusName: "Rejected_by_Chan"},
	}

	// ตรวจสอบว่า FormStatus มีข้อมูลอยู่แล้วหรือไม่
	var count int64
	db.Model(&models.FormStatus{}).Count(&count)
	if count > 0 {
		return nil // ข้อมูล FormStatus มีอยู่แล้ว ไม่ต้องสร้างใหม่
	}

	// บันทึก FormStatus ลงฐานข้อมูล
	return db.CreateInBatches(formStatuses, 100).Error
}

func SeedCampus(db *gorm.DB) error {
	campuses := []models.Campus{
		{CampusName: "บางเขน", CampusCode: "KU"},
		{CampusName: "กำแพงแสน", CampusCode: "KU-KPS"},
		{CampusName: "ศรีราชา", CampusCode: "KU-SR"},
		{CampusName: "เฉลิมพระเกียรติ จังหวัดสกลนคร", CampusCode: "KU-CSC"},
		{CampusName: "สุพรรณบุรี", CampusCode: "KU-SLA"},
	}

	// ตรวจสอบว่า Campus มีข้อมูลอยู่แล้วหรือไม่
	var count int64
	db.Model(&models.Campus{}).Count(&count)
	if count > 0 {
		return nil // ข้อมูล Campus มีอยู่แล้ว ไม่ต้องสร้างใหม่
	}

	// บันทึก Campus ลงฐานข้อมูล
	return db.CreateInBatches(campuses, 100).Error
}

func SeedRole(db *gorm.DB) error {
	roles := []models.Role{
		{RoleName: "Student", RoleNameTH: "นักศึกษา"},
		{RoleName: "Head of Department", RoleNameTH: "หัวหน้าภาควิชา"},
		{RoleName: "Associate Dean", RoleNameTH: "รองคณบดี"},
		{RoleName: "Dean", RoleNameTH: "คณบดี"},
		{RoleName: "Student Development", RoleNameTH: "กองพัฒนานิสิต"},
		{RoleName: "Committee", RoleNameTH: "คณะกรรมการ"},
		{RoleName: "Committee President", RoleNameTH: "ประธานคณะกรรมการ"},
		{RoleName: "Chancellor", RoleNameTH: "อธิการบดี"},
	}

	// ตรวจสอบว่า Role มีข้อมูลอยู่แล้วหรือไม่
	var count int64
	db.Model(&models.Role{}).Count(&count)
	if count > 0 {
		return nil // ข้อมูล Role มีอยู่แล้ว ไม่ต้องสร้างใหม่
	}

	// บันทึก Role ลงฐานข้อมูล
	return db.CreateInBatches(roles, 100).Error
}

func SeedAdmin(db *gorm.DB) error { return nil }
