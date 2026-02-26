package repository

import (
	"backend/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	// "golang.org/x/crypto/bcrypt"
	// "time"
	// "fmt"
)

const (
	committeeRoleID = 6
)

type UserRepository interface {
	UpsertUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserList() ([]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUserFields(ctx context.Context, userID uint, updates map[string]interface{}) (*models.User, error)
	GetUserListByCampus(ctx context.Context, campusID int, page int, limit int) ([]models.User, error)
	SetCommitteeChairman(ctx context.Context, targetUserID uint, isChairman bool) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) UpsertUser(user *models.User) error {
	// ค้นหาด้วย Email ถ้าเจอให้ Update ถ้าไม่เจอให้ Create (Upsert)
	return r.db.Where("email = ?", user.Email).Assign(user).FirstOrCreate(user).Error
}

// GetUserByID
func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserList() ([]models.User, error) {
	var users []models.User
	// ใช้ .Find เพื่อดึงข้อมูลทั้งหมด
	// หากไม่พบข้อมูล Find จะส่งกลับเป็น Slice ว่าง และไม่ถือว่าเป็น Error
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) UpdateUserFields(ctx context.Context, userID uint, updates map[string]interface{}) (*models.User, error) {
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetUserByID(userID)
}

func (r *userRepository) GetUserListByCampus(ctx context.Context, campusID int, page int, limit int) ([]models.User, error) {
	var users []models.User
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).
		Where("campus_id = ?", campusID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) ensureCommitteeProfile(tx *gorm.DB, userID uint) error {
	return tx.Where("user_id = ?", userID).FirstOrCreate(&models.Committee{UserID: userID}).Error
}

func (r *userRepository) SetCommitteeChairman(ctx context.Context, targetUserID uint, isChairman bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var target models.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", targetUserID).
			First(&target).Error; err != nil {
			return err
		}

		if target.RoleID != committeeRoleID {
			return errors.New("role change allowed only for committee role")
		}

		if err := r.ensureCommitteeProfile(tx, targetUserID); err != nil {
			return err
		}

		if !isChairman {
			return tx.Model(&models.Committee{}).
				Where("user_id = ?", targetUserID).
				Update("is_chairman", false).Error
		}

		if err := tx.Model(&models.Committee{}).
			Where("is_chairman = ? AND user_id <> ?", true, targetUserID).
			Update("is_chairman", false).Error; err != nil {
			return err
		}

		return tx.Model(&models.Committee{}).
			Where("user_id = ?", targetUserID).
			Update("is_chairman", true).Error
	})
}

// func (r *userRepository) GetUserListSortedByCampus() ([]models.User, error) {
//     var users []models.User

//     // ใช้ .Order("campus ASC") เพื่อเรียงลำดับจาก ก-ฮ หรือ A-Z
//     // สมมติว่าใน struct models.User มี field ชื่อ Campus
//     err := r.db.Order("campus ASC").Find(&users).Error

//     if err != nil {
//         return nil, err
//     }
//     return users, nil
// }

// func (r *userRepository) GetUserListSortedByProvider() ([]models.User, error) {
//     var users []models.User

//     // เรียงลำดับตาม Provider จาก A-Z
//     // หากต้องการจาก Z-A ให้ใช้ "provider DESC"
//     err := r.db.Order("provider ASC").Find(&users).Error

//     if err != nil {
//         return nil, err
//     }
//     return users, nil
// }

// func (r *userRepository) SearchUsers(keyword string, sortBy string) ([]models.User, error) {
//     var users []models.User

//     // สร้างรูปแบบการค้นหา เช่น %keyword% เพื่อหาคำที่อยู่ส่วนไหนของประโยคก็ได้
//     searchQuery := "%" + keyword + "%"

//     // ค้นหาในฟิลด์ Email หรือ Campus (หรือฟิลด์อื่นๆ ที่ต้องการ)
//     // จากนั้นเรียงลำดับตาม sortBy ที่ส่งมา
//     err := r.db.Where("email LIKE ? OR campus LIKE ?", searchQuery, searchQuery).
//         Order(sortBy + " ASC").
//         Find(&users).Error

//     if err != nil {
//         return nil, err
//     }
//     return users, nil
// }
