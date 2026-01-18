package repository

import (
	"backend/internal/models"
	"context"

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

// GetAll และ GetByType คงเดิมตามที่คุณเขียนไว้ (ซึ่งถูกต้องแล้วในการ Preload AwardFiles)
func (r *AwardRepository) GetAll(ctx context.Context) ([]models.AwardForm, error) {
	var list []models.AwardForm
	err := r.db.WithContext(ctx).
		Preload("Student.User").
		Preload("Student.Faculty").
		Preload("Student.Department").
		Preload("AwardType").
		Preload("Extracurricular").
		Preload("GoodBehavior").
		Preload("Creativity").
		Preload("AwardFiles"). // ดึงข้อมูลไฟล์แนบมาแสดงผลด้วย
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *AwardRepository) GetByType(ctx context.Context, typeID int) ([]models.AwardForm, error) {
	var list []models.AwardForm
	query := r.db.WithContext(ctx).
		Where("award_type_id = ?", typeID).
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
