package repository

import (
    "backend/internal/models"
    "context"

    "gorm.io/gorm"
)

type AcademicYearRepoDB struct {
    db *gorm.DB
}

func NewAcademicYearRepository(db *gorm.DB) *AcademicYearRepoDB {
    return &AcademicYearRepoDB{db: db}
}

func (r *AcademicYearRepoDB) Create(ctx context.Context, ay *models.AcademicYear) error {
    return r.db.WithContext(ctx).Create(ay).Error
}

func (r *AcademicYearRepoDB) Update(ctx context.Context, ay *models.AcademicYear) error {
    return r.db.WithContext(ctx).Save(ay).Error
}

func (r *AcademicYearRepoDB) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Where("academic_year_id = ?", id).Delete(&models.AcademicYear{}).Error
}

func (r *AcademicYearRepoDB) GetByID(ctx context.Context, id uint) (*models.AcademicYear, error) {
    var ay models.AcademicYear
    if err := r.db.WithContext(ctx).First(&ay, "academic_year_id = ?", id).Error; err != nil {
        return nil, err
    }
    return &ay, nil
}

func (r *AcademicYearRepoDB) GetList(ctx context.Context) ([]models.AcademicYear, error) {
	var list []models.AcademicYear
	// เรียงลำดับจากปีล่าสุดและเทอมล่าสุดลงไป
	err := r.db.WithContext(ctx).
		Order("year desc").
		Order("semester desc").
		Find(&list).Error
	
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *AcademicYearRepoDB) GetLatest(ctx context.Context) (*models.AcademicYear, error) {
    var ay models.AcademicYear
    // เลือกโดย year desc, semester desc เพื่อให้ได้ปี+เทอมล่าสุด
    if err := r.db.WithContext(ctx).Order("year desc").Order("semester desc").First(&ay).Error; err != nil {
        return nil, err
    }
    return &ay, nil
}