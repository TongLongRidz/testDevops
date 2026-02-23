package migration

import (
	"backend/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func SeedFormStatus(db *gorm.DB) error {
	formStatuses := []models.FormStatus{
		{FormStatusName: "ฟอร์มใหม่"},
		{FormStatusName: "อนุมัติโดยหัวหน้าภาควิชา"}, // ส่งต่อให้รองคณบดี
		{FormStatusName: "ปฏิเสธโดยหัวหน้าภาควิชา"},
		{FormStatusName: "อนุมัติโดยรองคณบดี"}, // ส่งต่อให้คณบดี
		{FormStatusName: "ปฏิเสธโดยรองคณบดี"},
		{FormStatusName: "อนุมัติโดยคณบดี"}, // ส่งต่อให้กองพัฒนานิสิต
		{FormStatusName: "ปฏิเสธโดยคณบดี"},
		{FormStatusName: "อนุมัติโดยกองพัฒนานิสิต"},  // ส่งต่อให้คณะกรรมการ
		{FormStatusName: "ตีกลับโดยกองพัฒนานิสิต"},      // ให้นิสิตแก้ให้เสร็จแล้วส่งใหม่ กลับไป status ก่อนหน้า [ลูปก่อนหน้ากับอันนี้จนกว่าจะผ่าน]
		{FormStatusName: "อนุมัติโดยคณะกรรมการ"},     // โหวตผ่านเกินครึ่ง
		{FormStatusName: "ปฏิเสธโดยคณะกรรมการ"},      // โหวตไม่ผ่านเกินครึ่ง (หรือไม่ถึงในเวลา?)
		{FormStatusName: "ลงนามโดยประธานคณะกรรมการ"}, // ส่งต่อให้อธิการบดี
		{FormStatusName: "ลงนามโดยอธิการบดี"},
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
		{RoleName: "Chancellor", RoleNameTH: "อธิการบดี"},
		{RoleName: "Organization", RoleNameTH: "หน่วยงานภายนอก"},
	}

	var existingRoles []models.Role
	if err := db.Find(&existingRoles).Error; err != nil {
		return err
	}

	existingByName := make(map[string]struct{}, len(existingRoles))
	for _, role := range existingRoles {
		existingByName[role.RoleName] = struct{}{}
	}

	var newRoles []models.Role
	for _, role := range roles {
		if _, exists := existingByName[role.RoleName]; !exists {
			newRoles = append(newRoles, role)
		}
	}

	if len(newRoles) == 0 {
		return nil
	}

	return db.CreateInBatches(newRoles, 100).Error
}

func SeedFacultyAndDepartments(db *gorm.DB) error {
	facultiesToSeed := []string{
		"คณะเกษตร",
		"คณะประมง",
		"คณะวนศาสตร์",
		"คณะวิทยาศาสตร์",
		"คณะวิศวกรรมศาสตร์",
		"คณะบริหารธุรกิจ",
		"คณะเศรษฐศาสตร์",
		"คณะสังคมศาสตร์",
		"คณะมนุษยศาสตร์",
		"คณะศึกษาศาสตร์",
		"คณะสัตวแพทยศาสตร์",
		"คณะเทคนิคการสัตวแพทย์",
		"คณะสถาปัตยกรรมศาสตร์",
		"คณะสิ่งแวดล้อม",
	}

	var existingFaculties []models.Faculty
	if err := db.Find(&existingFaculties).Error; err != nil {
		return err
	}

	facultyByName := make(map[string]uint, len(existingFaculties))
	for _, faculty := range existingFaculties {
		facultyByName[faculty.FacultyName] = faculty.FacultyID
	}

	var newFaculties []models.Faculty
	for _, name := range facultiesToSeed {
		if _, exists := facultyByName[name]; !exists {
			newFaculties = append(newFaculties, models.Faculty{FacultyName: name})
		}
	}

	if len(newFaculties) > 0 {
		if err := db.CreateInBatches(newFaculties, 100).Error; err != nil {
			return err
		}
	}

	var refreshedFaculties []models.Faculty
	if err := db.Find(&refreshedFaculties).Error; err != nil {
		return err
	}

	facultyByName = make(map[string]uint, len(refreshedFaculties))
	for _, faculty := range refreshedFaculties {
		facultyByName[faculty.FacultyName] = faculty.FacultyID
	}

	var departmentCount int64
	db.Model(&models.Department{}).Count(&departmentCount)
	if departmentCount > 0 {
		return nil
	}

	departmentsByFaculty := map[string][]string{
		"คณะเกษตร": {
			"กีฏวิทยา",
			"โรคพืช",
			"สัตวบาล",
			"คหกรรมศาสตร์",
		},
		"คณะประมง": {
			"จัดการประมง",
			"ชีววิทยาประมง",
			"ผลิตภัณฑ์ประมง",
			"เพาะเลี้ยงสัตว์น้ำ",
		},
		"คณะวนศาสตร์": {
			"การจัดการลุ่มน้ำและสิ่งแวดล้อม",
			"วนวัฒนวิทยา",
			"วิศวกรรมป่าไม้",
		},
		"คณะวิทยาศาสตร์": {
			"คณิตศาสตร์",
			"เคมี",
			"ชีววิทยา",
			"วิทยาการคอมพิวเตอร์",
			"สถิติ",
			"ฟิสิกส์",
		},
		"คณะวิศวกรรมศาสตร์": {
			"วิศวกรรมโยธา",
			"เครื่องกล",
			"ไฟฟ้า",
			"คอมพิวเตอร์",
			"เคมี",
			"สิ่งแวดล้อม",
			"การบินและอวกาศ",
		},
		"คณะบริหารธุรกิจ": {
			"การเงิน",
			"การจัดการ",
			"การตลาด",
			"การบัญชี",
		},
		"คณะเศรษฐศาสตร์": {
			"เศรษฐศาสตร์",
			"เศรษฐศาสตร์เกษตรและทรัพยากร",
		},
		"คณะสังคมศาสตร์": {
			"นิติศาสตร์",
			"รัฐศาสตร์",
			"จิตวิทยา",
			"ภูมิศาสตร์",
			"สังคมวิทยาและมานุษยวิทยา",
		},
		"คณะมนุษยศาสตร์": {
			"ภาษาไทย",
			"ภาษาอังกฤษ",
			"ภาษาญี่ปุ่น",
			"การท่องเที่ยวและโรงแรม",
			"ดนตรี",
		},
		"คณะศึกษาศาสตร์": {
			"การสอนคณิตศาสตร์",
			"พลศึกษา",
			"คหกรรมศาสตรศึกษา",
		},
		"คณะสัตวแพทยศาสตร์": {
			"สัตวแพทยศาสตรบัณฑิต",
		},
		"คณะเทคนิคการสัตวแพทย์": {
			"เทคนิคการสัตวแพทย์",
			"การพยาบาลสัตว์",
		},
		"คณะสถาปัตยกรรมศาสตร์": {
			"สถาปัตยกรรมศาสตร์",
			"ภูมิสถาปัตยกรรม",
		},
		"คณะสิ่งแวดล้อม": {
			"วิทยาศาสตร์สิ่งแวดล้อม",
		},
	}

	var departments []models.Department
	for facultyName, departmentNames := range departmentsByFaculty {
		facultyID, ok := facultyByName[facultyName]
		if !ok {
			return fmt.Errorf("faculty not found for department seeding: %s", facultyName)
		}
		for _, departmentName := range departmentNames {
			departments = append(departments, models.Department{
				FacultyID:      facultyID,
				DepartmentName: departmentName,
			})
		}
	}

	return db.CreateInBatches(departments, 100).Error
}

func SeedAdmin(db *gorm.DB) error { return nil }
