package repository

import (
	"backend/internal/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type AwardRepository struct {
	db *gorm.DB
}

func NewAwardRepository(db *gorm.DB) *AwardRepository {
	return &AwardRepository{db: db}
}

// ปรับปรุง: เพิ่มพารามิเตอร์ files เพื่อรองรับการบันทึกไฟล์แนบ (ถ้ามี)
func (r *AwardRepository) CreateWithTransaction(ctx context.Context, form *models.AwardForm, detail interface{}, files []models.AwardFileDirectory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. บันทึกตารางหลัก (Award_Form)
		if err := tx.Create(form).Error; err != nil {
			return err
		}

		// 2. บันทึกตารางรายละเอียด (Detail)
		if detail != nil {
			switch d := detail.(type) {
			case *models.ExtracurricularActivity:
				d.FormID = form.FormID
			case *models.GoodBehavior:
				d.FormID = form.FormID
			case *models.CreativityInnovation:
				d.FormID = form.FormID
			}

			if err := tx.Create(detail).Error; err != nil {
				return err
			}
		}

		// 3. บันทึกตารางไฟล์แนบ (ถ้า len > 0 คือมีการแนบไฟล์มา)
		if len(files) > 0 {
			for i := range files {
				// ผูก ID ของไฟล์เข้ากับ FormID ที่เพิ่งสร้างใหม่
				files[i].FormID = form.FormID
				if err := tx.Create(&files[i]).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// GetByKeyword ค้นหาและกรองพร้อม pagination
func (r *AwardRepository) GetByKeyword(ctx context.Context, campusID int, keyword string, date string, studentYear int, page int, limit int, arrangement string) ([]models.AwardForm, int64, error) {
	var list []models.AwardForm
	var total int64

	// สร้าง query พื้นฐาน - กรองตามวิทยาเขตเสมอ
	query := r.db.WithContext(ctx).Model(&models.AwardForm{}).Where("campus_id = ?", campusID)

	// ค้นหาด้วย keyword (firstname, lastname, studentNumber, semester, year, award_type_id)
	if keyword != "" {
		query = query.Where(
			"student_firstname LIKE ? OR student_lastname LIKE ? OR student_number LIKE ? OR CAST(semester AS CHAR) LIKE ? OR CAST(academic_year AS CHAR) LIKE ? OR CAST(award_type_id AS CHAR) LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
		)
	}

	// กรองตามวันที่ (ถ้ามี)
	if date != "" {
		query = query.Where("DATE(created_at) = ?", date)
	}

	// กรองตามชั้นปี (ถ้ามี)
	if studentYear > 0 {
		query = query.Where("student_year = ?", studentYear)
	}

	// นับจำนวนทั้งหมด
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// คำนวณ offset
	offset := (page - 1) * limit

	// ดึงข้อมูลพร้อม pagination และ preload
	orderClause := "created_at desc"
	if arrangement == "asc" {
		orderClause = "created_at asc"
	}

	err := query.
		Preload("Student.User").
		Preload("Student.Faculty").
		Preload("Student.Department").
		Preload("AwardType").
		Preload("Extracurricular").
		Preload("GoodBehavior").
		Preload("Creativity").
		Preload("AwardFiles").
		Order(orderClause).
		Limit(limit).
		Offset(offset).
		Find(&list).Error

	return list, total, err
}

func (r *AwardRepository) GetByType(ctx context.Context, typeID int, campusID int) ([]models.AwardForm, error) {
	var list []models.AwardForm
	query := r.db.WithContext(ctx).
		Where("award_type_id = ? AND campus_id = ?", typeID, campusID).
		Preload("Student.User").
		Preload("Student.Faculty").
		Preload("Student.Department").
		Preload("AwardType").
		Preload("AwardFiles")

	switch typeID {
	case 1:
		query = query.Preload("Extracurricular")
	case 2:
		query = query.Preload("GoodBehavior")
	case 3:
		query = query.Preload("Creativity")
	}

	err := query.Order("created_at desc").Find(&list).Error
	return list, err
}

func (r *AwardRepository) GetByStudentID(ctx context.Context, studentID int) ([]models.AwardForm, error) {
	var list []models.AwardForm
	err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Preload("Student.User").
		Preload("Student.Faculty").
		Preload("Student.Department").
		Preload("AwardType").
		Preload("Extracurricular").
		Preload("GoodBehavior").
		Preload("Creativity").
		Preload("AwardFiles").
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *AwardRepository) GetByFormID(ctx context.Context, formID int) (*models.AwardForm, error) {
	var form models.AwardForm
	err := r.db.WithContext(ctx).
		Where("form_id = ?", formID).
		Preload("Student.User").
		Preload("Student.Faculty").
		Preload("Student.Department").
		Preload("AwardType").
		Preload("Extracurricular").
		Preload("GoodBehavior").
		Preload("Creativity").
		Preload("AwardFiles").
		First(&form).Error
	if err != nil {
		return nil, err
	}
	return &form, nil
}

func (r *AwardRepository) CheckDuplicate(studentID int, year int, semester int) (bool, error) {
	var count int64
	// เช็คในตาราง AwardForm ว่ามีข้อมูลที่ student_id, academic_year, semester ตรงกันไหม
	err := r.db.Model(&models.AwardForm{}).
		Where("student_id = ? AND academic_year = ? AND semester = ?", studentID, year, semester).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *AwardRepository) UpdateAwardType(ctx context.Context, formID uint, awardTypeID int) error {
	return r.db.WithContext(ctx).
		Model(&models.AwardForm{}).
		Where("form_id = ?", formID).
		Updates(map[string]interface{}{
			"award_type_id": awardTypeID,
			"latest_update": time.Now(),
		}).Error
}

func (r *AwardRepository) UpdateFormStatus(ctx context.Context, formID uint, formStatusID int) error {
	return r.db.WithContext(ctx).
		Model(&models.AwardForm{}).
		Where("form_id = ?", formID).
		Updates(map[string]interface{}{
			"form_status_id": formStatusID,
			"latest_update":  time.Now(),
		}).Error
}
