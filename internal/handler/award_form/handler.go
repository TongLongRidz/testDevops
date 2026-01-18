package awardform

import (
	awardformdto "backend/internal/dto/award_form_dto"
	"backend/internal/models"
	"backend/internal/usecase"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AwardHandler struct {
	useCase             usecase.AwardUseCase
	studentService      usecase.StudentService
	academicYearService usecase.AcademicYearService
}

func NewAwardHandler(u usecase.AwardUseCase, s usecase.StudentService, ays usecase.AcademicYearService) *AwardHandler {
	return &AwardHandler{useCase: u, studentService: s, academicYearService: ays}
}

func (h *AwardHandler) Submit(c *fiber.Ctx) error {
	var req awardformdto.SubmitAwardRequest

	// สร้าง uploads folder อัตโนมัติ
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create uploads directory",
		})
	}

	// ดึงข้อมูลผู้ใช้ที่ login อยู่จาก middleware
	currentUser := c.Locals("current_user")
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: User not found",
		})
	}
	user, ok := currentUser.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user data",
		})
	}

	// 1. รับข้อมูลที่เป็น Text/JSON จาก Form
	// Debug: ดูค่าที่ได้รับทั้งหมด
	fmt.Println("=== DEBUG FORM VALUES ===")
	fmt.Printf("award_type_id: '%s'\n", c.FormValue("award_type_id"))
	fmt.Printf("student_year: '%s'\n", c.FormValue("student_year"))
	fmt.Printf("advisor_name: '%s'\n", c.FormValue("advisor_name"))

	awardTypeIDStr := c.FormValue("award_type_id")
	if awardTypeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "award_type_id is required",
		})
	}
	awardTypeID, err := strconv.Atoi(awardTypeIDStr)
	if err != nil || awardTypeID == 0 {
		fmt.Printf("❌ Parse error: awardTypeIDStr='%s', awardTypeID=%d, err=%v\n", awardTypeIDStr, awardTypeID, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("award_type_id must be valid (1, 2, or 3), got: '%s'", awardTypeIDStr),
		})
	}
	req.AwardTypeID = awardTypeID

	// ตรวจสอบ duplicate จะทำใน usecase หลังจากดึงข้อมูล student และ academic year แล้ว
	// (ย้ายไป usecase เพราะต้องรู้ student_id และ academic_year ก่อน)

	// ดึงข้อมูลทั่วไป (ที่ไม่ได้มาจาก token)
	studentYear, _ := strconv.Atoi(c.FormValue("student_year"))
	req.StudentYear = studentYear
	req.AdvisorName = c.FormValue("advisor_name")
	req.PhoneNumber = c.FormValue("phone_number")
	req.Address = c.FormValue("address")
	gpa, _ := strconv.ParseFloat(c.FormValue("gpa"), 64)
	req.GPA = gpa
	if dobStr := c.FormValue("date_of_birth"); dobStr != "" {
		dob, _ := time.Parse("2006-01-02", dobStr)
		req.DateOfBirth = dob
	}

	// ตรวจสอบประเภทรางวัลเพื่อดึงข้อมูลรายละเอียด
	switch req.AwardTypeID {
	case 1: // Extracurricular Activity
		req.Extracurricular = &awardformdto.ExtracurricularRequest{
			QualificationType: c.FormValue("qualification_type"),
			TeamName:          c.FormValue("team_name"),
			ProjectTitle:      c.FormValue("project_title"),
			Prize:             c.FormValue("prize"),
			OrganizedBy:       c.FormValue("organized_by"),
			CompetitionLevel:  c.FormValue("competition_level"),
			ActivityCategory:  c.FormValue("activity_category"),
		}
		if dateStr := c.FormValue("date_received"); dateStr != "" {
			t, _ := time.Parse("2006-01-02", dateStr) // จัด Layout เฉยๆ
			req.Extracurricular.DateReceived = t
		}

	case 2: // Good Behavior
		// หากในอนาคตมีฟิลด์ของ Good Behavior ให้มาเพิ่มที่นี่
		req.GoodBehavior = &awardformdto.GoodBehaviorRequest{}

	case 3: // Creativity & Innovation
		req.Creativity = &awardformdto.CreativityRequest{
			TeamName:         c.FormValue("team_name"),
			ProjectTitle:     c.FormValue("project_title"),
			Prize:            c.FormValue("prize"),
			OrganizedBy:      c.FormValue("organized_by"),
			CompetitionLevel: c.FormValue("competition_level"),
			ActivityCategory: c.FormValue("activity_category"),
		}
		if dateStr := c.FormValue("date_received"); dateStr != "" {
			t, _ := time.Parse("2006-01-02", dateStr)
			req.Creativity.DateReceived = t
		}
	}

	// จัดการกับไฟล์แนบ (ถ้ามี)
	var awardFiles []models.AwardFileDirectory

	form, err := c.MultipartForm()
	if err == nil {
		// Debug: แสดง field names ทั้งหมดที่มีในฟอร์ม
		fmt.Println("🔍 Form fields ที่มีในฟอร์ม:")
		for fieldName := range form.File {
			fmt.Printf("  - %s: %d files\n", fieldName, len(form.File[fieldName]))
		}

		files := form.File["files"]

		fmt.Printf("📁 จำนวนไฟล์ที่ได้รับ (field 'files'): %d\n", len(files))

		// --- STEP 1: VALIDATION LOOP ---
		// เช็คไฟล์ทั้งหมดก่อนว่ามีอันไหนไม่ valid ไหม
		allowedExtensions := map[string]bool{".pdf": true}
		maxTotalSize := int64(10 * 1024 * 1024) // 10 MB
		var totalSize int64

		for _, file := range files {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if !allowedExtensions[ext] {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": fmt.Sprintf("ไม่อนุญาตให้อัปโหลดไฟล์ประเภท %s (รองรับเฉพาะ PDF)", ext),
				})
			}
			totalSize += file.Size
		}

		// เช็คขนาดไฟล์รวม
		if totalSize > maxTotalSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("ขนาดไฟล์รวมเกิน 10 MB (ได้รับ %.2f MB)", float64(totalSize)/(1024*1024)),
			})
		}

		// --- STEP 2: PROCESSING & SAVING LOOP ---
		// ถ้าผ่านการเช็คด้านบนมาได้ แสดงว่าไฟล์ทุกไฟล์ valid พร้อมบันทึก
		for _, file := range files {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			subDir := "pdf" // รองรับเฉพาะ PDF

			targetDir := filepath.Join(uploadDir, subDir)
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": "error", "message": "Failed to create directory",
				})
			}

			newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
			savePath := filepath.Join(targetDir, newFileName)

			if err := c.SaveFile(file, savePath); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to save file: " + err.Error(),
				})
			}

			fmt.Printf("✅ บันทึกไฟล์สำเร็จ: %s (ขนาด: %d bytes)\n", savePath, file.Size)

			cleanPath := filepath.ToSlash(savePath)
			awardFiles = append(awardFiles, models.AwardFileDirectory{
				FilePath:   cleanPath,
				FileType:   strings.TrimPrefix(ext, "."),
				FileSize:   file.Size,
				UploadedAt: time.Now(),
			})
		}
	}

	// 3. ส่งข้อมูลไปยัง UseCase พร้อม userID
	if err := h.useCase.SubmitAward(c.UserContext(), user.UserID, req, awardFiles); err != nil {

		// --- ส่วนที่เพิ่มเข้ามา: ลบไฟล์ทิ้งถ้า DB บันทึกไม่สำเร็จ ---
		for _, f := range awardFiles {
			// f.FilePath เก็บค่าเช่น "uploads/pdf/xxx.pdf"
			if removeErr := os.Remove(f.FilePath); removeErr != nil {
				fmt.Printf("Failed to cleanup file %s: %v\n", f.FilePath, removeErr)
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "บันทึกข้อมูลไม่สำเร็จ (อาจมีการส่งข้อมูลในปีการศึกษานี้ไปแล้ว): " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Award form submitted successfully",
	})
}

func (h *AwardHandler) GetAll(c *fiber.Ctx) error {
	results, err := h.useCase.GetAllAwards(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": results})
}

func (h *AwardHandler) GetByType(c *fiber.Ctx) error {
	typeID, _ := c.ParamsInt("type_id")
	results, err := h.useCase.GetAwardsByType(c.UserContext(), typeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": results})
}
