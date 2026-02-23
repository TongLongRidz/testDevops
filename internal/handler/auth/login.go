package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	authDto "backend/internal/dto/auth_dto"
	"backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req authDto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
	}

	token, user, err := h.AuthService.AuthenticateAndToken(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// ตั้ง cookie (ปรับ Secure/SameSite ตาม environment ของคุณ)
	c.Cookie(&fiber.Cookie{
		Name:     "token",   // ชื่อ cookie
		Value:    token,     // ค่า token ที่ได้มา
		Path:     "/",       // ให้ cookie นี้ใช้ได้ทั่วทั้งเว็บ
		MaxAge:   3600 * 24, // 1 วัน (หน่วยเป็นวินาที)
		HTTPOnly: true,      // ป้องกันการเข้าถึงจาก JavaScript
		Secure:   false,     // ตั้งเป็น true ถ้า deploy แล้วใช้ HTTPS
		SameSite: "Lax",     // "Lax" หรือ "Strict" เพื่อความปลอดภัย
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"user_id":        user.UserID,
			"prefix":         user.Prefix,
			"firstname":      user.Firstname,
			"lastname":       user.Lastname,
			"email":          user.Email,
			"image_path":     user.ImagePath,
			"provider":       user.Provider,
			"role_id":        user.RoleID,
			"campus_id":      user.CampusID,
			"is_first_login": user.IsFirstLogin,
			"created_at":     user.CreatedAt,
			"latest_update":  user.LatestUpdate,
		},
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// ตั้ง cookie ให้หมดอายุ (ล้าง token)

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // ลบ cookie ทันที
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Expires:  time.Now().Add(-time.Hour),
	})

	return c.JSON(fiber.Map{"message": "logged out"})
}

// เพิ่ม /me เพื่อคืนข้อมูลผู้ใช้อย่างปลอดภัย (ต้องผ่าน middleware.RequireAuth ก่อน)
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	u := c.Locals("current_user")
	if u == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	user, ok := u.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user"})
	}

	// ดึงข้อมูล user แบบเต็มจาก database
	fullUser, err := h.AuthService.GetUserByID(c.Context(), user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch user"})
	}

	// สร้าง response object
	response := authDto.MeResponse{
		UserID:       fullUser.UserID,
		Prefix:       fullUser.Prefix,
		Firstname:    fullUser.Firstname,
		Lastname:     fullUser.Lastname,
		Email:        fullUser.Email,
		ImagePath:    fullUser.ImagePath,
		Provider:     fullUser.Provider,
		RoleID:       fullUser.RoleID,
		CampusID:     fullUser.CampusID,
		IsFirstLogin: fullUser.IsFirstLogin,
		CreatedAt:    fullUser.CreatedAt,
		LatestUpdate: fullUser.LatestUpdate,
	}

	// ถ้า RoleID = 1 (Student) ให้ดึงข้อมูล Student ด้วย
	if fullUser.RoleID == 1 && h.StudentService != nil {
		student, err := h.StudentService.GetStudentByUserID(c.Context(), fullUser.UserID)
		if err == nil && student != nil {
			response.StudentData = &authDto.StudentMeData{
				StudentID:     student.StudentID,
				StudentNumber: student.StudentNumber,
				FacultyID:     student.FacultyID,
				DepartmentID:  student.DepartmentID,
			}
		}
	}

	// ถ้า RoleID = 8 (Organization) ให้ดึงข้อมูล Organization ด้วย
	if fullUser.RoleID == 8 && h.OrganizationService != nil {
		org, err := h.OrganizationService.GetByUserID(c.Context(), fullUser.UserID)
		if err == nil && org != nil {
			response.OrganizationData = &authDto.OrganizationMeData{
				OrganizationID:       org.OrganizationID,
				OrganizationName:     org.OrganizationName,
				OrganizationType:     org.OrganizationType,
				OrganizationLocation: org.OrganizationLocation,
				OrganizationPhone:    org.OrganizationPhoneNumber,
			}
		}
	}

	if fullUser.RoleID == 2 {
		hod, err := h.AuthService.GetHeadOfDepartmentByUserID(c.Context(), fullUser.UserID)
		if err == nil && hod != nil {
			response.HeadOfDepartmentData = &authDto.HeadOfDepartmentMeData{
				HodID:        hod.HodID,
				UserID:       hod.UserID,
				FacultyID:    hod.FacultyID,
				DepartmentID: hod.DepartmentID,
			}
		}
	}

	if fullUser.RoleID == 3 {
		ad, err := h.AuthService.GetAssociateDeanByUserID(c.Context(), fullUser.UserID)
		if err == nil && ad != nil {
			response.AssociateDeanData = &authDto.AssociateDeanMeData{
				AdID:      ad.AdID,
				UserID:    ad.UserID,
				AdCode:    ad.AdCode,
				FacultyID: ad.FacultyID,
			}
		}
	}

	if fullUser.RoleID == 4 {
		dean, err := h.AuthService.GetDeanByUserID(c.Context(), fullUser.UserID)
		if err == nil && dean != nil {
			response.DeanData = &authDto.DeanMeData{
				DID:       dean.DID,
				UserID:    dean.UserID,
				FacultyID: dean.FacultyID,
			}
		}
	}

	if fullUser.RoleID == 5 {
		sd, err := h.AuthService.GetStudentDevelopmentByUserID(c.Context(), fullUser.UserID)
		if err == nil && sd != nil {
			response.StudentDevelopmentData = &authDto.StudentDevelopmentMeData{
				SdID:   sd.SdID,
				UserID: sd.UserID,
			}
		}
	}

	if fullUser.RoleID == 6 {
		com, err := h.AuthService.GetCommitteeByUserID(c.Context(), fullUser.UserID)
		if err == nil && com != nil {
			response.CommitteeData = &authDto.CommitteeMeData{
				ComID:      com.ComID,
				UserID:     com.UserID,
				IsChairman: com.IsChairman,
			}
		}
	}

	if fullUser.RoleID == 7 {
		ch, err := h.AuthService.GetChancellorByUserID(c.Context(), fullUser.UserID)
		if err == nil && ch != nil {
			response.ChancellorData = &authDto.ChancellorMeData{
				ChancellorID: ch.ChancellorID,
				UserID:       ch.UserID,
			}
		}
	}

	return c.JSON(fiber.Map{
		"user": response,
	})
}

// UpdateMe - อัพเดทข้อมูล current user
func (h *AuthHandler) UpdateMe(c *fiber.Ctx) error {
	u := c.Locals("current_user")
	if u == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	user, ok := u.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user"})
	}

	var req authDto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	// เรียกใช้ service เพื่ออัพเดทข้อมูล
	updatedUser, err := h.AuthService.UpdateUser(c.Context(), user.UserID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "user updated successfully",
		"user": fiber.Map{
			"user_id":        updatedUser.UserID,
			"prefix":         updatedUser.Prefix,
			"firstname":      updatedUser.Firstname,
			"lastname":       updatedUser.Lastname,
			"email":          updatedUser.Email,
			"image_path":     updatedUser.ImagePath,
			"provider":       updatedUser.Provider,
			"role_id":        updatedUser.RoleID,
			"campus_id":      updatedUser.CampusID,
			"is_first_login": updatedUser.IsFirstLogin,
			"created_at":     updatedUser.CreatedAt,
			"latest_update":  updatedUser.LatestUpdate,
		},
	})
}

// FirstLogin - ตั้งค่าข้อมูลครั้งแรกสำหรับนักศึกษา/องค์กร
func (h *AuthHandler) FirstLogin(c *fiber.Ctx) error {
	u := c.Locals("current_user")
	if u == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	user, ok := u.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid user"})
	}

	// รับข้อมูลจาก form-data
	prefix := strings.TrimSpace(c.FormValue("prefix"))
	firstname := strings.TrimSpace(c.FormValue("firstname"))
	lastname := strings.TrimSpace(c.FormValue("lastname"))
	campusIDStr := strings.TrimSpace(c.FormValue("campus_id"))

	// Validation common fields
	if prefix == "" || firstname == "" || lastname == "" || campusIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "prefix, firstname, lastname, and campus_id are required"})
	}

	campusID, err := strconv.Atoi(campusIDStr)
	if err != nil || campusID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid campus_id"})
	}

	// Prepare request object
	req := authDto.FirstLoginRequest{
		Prefix:    prefix,
		Firstname: firstname,
		Lastname:  lastname,
		CampusID:  campusID,
	}

	var imagePath string

	// Handle ตามประเภท Role
	switch user.RoleID {
	case 1: // Student
		studentNumber := strings.TrimSpace(c.FormValue("student_number"))
		facultyIDStr := strings.TrimSpace(c.FormValue("faculty_id"))
		departmentIDStr := strings.TrimSpace(c.FormValue("department_id"))

		if studentNumber == "" || facultyIDStr == "" || departmentIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "student_number, faculty_id, and department_id are required for student"})
		}

		facultyID, err := strconv.ParseUint(facultyIDStr, 10, 32)
		if err != nil || facultyID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid faculty_id"})
		}

		departmentID, err := strconv.ParseUint(departmentIDStr, 10, 32)
		if err != nil || departmentID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid department_id"})
		}

		req.StudentNumber = studentNumber
		req.FacultyID = uint(facultyID)
		req.DepartmentID = uint(departmentID)

		// Handle profile image for Student
		file, err := c.FormFile("profile_image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "profile_image is required for student"})
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == "" {
			ext = ".jpg"
		}

		uploadDir := filepath.Join("uploads", "user-profile")
		if err := os.MkdirAll(uploadDir, 0o755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to prepare upload directory"})
		}

		fileName := fmt.Sprintf("%d%s", user.UserID, ext)
		savePath := filepath.Join(uploadDir, fileName)
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save profile image"})
		}

		imagePath = "/" + filepath.ToSlash(savePath)

	case 8: // Organization
		orgName := strings.TrimSpace(c.FormValue("organization_name"))
		if orgName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "organization_name is required for organization"})
		}

		req.OrganizationName = orgName
		req.OrganizationType = strings.TrimSpace(c.FormValue("organization_type"))
		req.OrganizationLocation = strings.TrimSpace(c.FormValue("organization_location"))
		req.OrganizationPhone = strings.TrimSpace(c.FormValue("organization_phone"))

		// Optional: profile image for Organization
		file, err := c.FormFile("profile_image")
		if err == nil {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext == "" {
				ext = ".jpg"
			}

			uploadDir := filepath.Join("uploads", "user-profile")
			if err := os.MkdirAll(uploadDir, 0o755); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to prepare upload directory"})
			}

			fileName := fmt.Sprintf("%d%s", user.UserID, ext)
			savePath := filepath.Join(uploadDir, fileName)
			if err := c.SaveFile(file, savePath); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save profile image"})
			}

			imagePath = "/" + filepath.ToSlash(savePath)
		}

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid role_id"})
	}

	fmt.Printf("[FirstLogin] userID=%d roleID=%d imagePath=%s\n", user.UserID, user.RoleID, imagePath)

	updatedUser, _, err := h.AuthService.CompleteFirstLogin(c.Context(), user.UserID, &req, imagePath)
	if err != nil {
		fmt.Printf("[FirstLogin] error: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// สร้าง response
	response := authDto.MeResponse{
		UserID:       updatedUser.UserID,
		Prefix:       updatedUser.Prefix,
		Firstname:    updatedUser.Firstname,
		Lastname:     updatedUser.Lastname,
		Email:        updatedUser.Email,
		ImagePath:    updatedUser.ImagePath,
		Provider:     updatedUser.Provider,
		RoleID:       updatedUser.RoleID,
		CampusID:     updatedUser.CampusID,
		IsFirstLogin: updatedUser.IsFirstLogin,
		CreatedAt:    updatedUser.CreatedAt,
		LatestUpdate: updatedUser.LatestUpdate,
	}

	// ถ้า RoleID = 1 (Student) ให้ดึงข้อมูล Student ด้วย
	if updatedUser.RoleID == 1 && h.StudentService != nil {
		student, err := h.StudentService.GetStudentByUserID(c.Context(), updatedUser.UserID)
		if err == nil && student != nil {
			response.StudentData = &authDto.StudentMeData{
				StudentID:     student.StudentID,
				StudentNumber: student.StudentNumber,
				FacultyID:     student.FacultyID,
				DepartmentID:  student.DepartmentID,
			}
		}
	}

	// ถ้า RoleID = 8 (Organization) ให้ดึงข้อมูล Organization ด้วย
	if updatedUser.RoleID == 8 && h.OrganizationService != nil {
		org, err := h.OrganizationService.GetByUserID(c.Context(), updatedUser.UserID)
		if err == nil && org != nil {
			response.OrganizationData = &authDto.OrganizationMeData{
				OrganizationID:       org.OrganizationID,
				OrganizationName:     org.OrganizationName,
				OrganizationType:     org.OrganizationType,
				OrganizationLocation: org.OrganizationLocation,
				OrganizationPhone:    org.OrganizationPhoneNumber,
			}
		}
	}

	return c.JSON(fiber.Map{
		"message": "first login completed",
		"user":    response,
	})
}
