package usecase

import (
	"backend/config"
	authDto "backend/internal/dto/auth_dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	jwt "github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	GetGoogleLoginURL() string
	ProcessGoogleLogin(code string) (*models.User, error)
	IssueToken(user *models.User) (string, error)
	Register(req *authDto.RegisterRequest) (*models.User, error)
	AuthenticateAndToken(ctx context.Context, email, password string) (string, *models.User, error)

	GetUserByID(ctx context.Context, userID uint) (*models.User, error)
	UpdateUser(ctx context.Context, userID uint, req *authDto.UpdateUserRequest) (*models.User, error)
	CompleteFirstLogin(ctx context.Context, userID uint, req *authDto.FirstLoginRequest, imagePath string) (*models.User, *models.Student, error)
}

var ErrInvalidCredentials = errors.New("invalid credentials")

type authService struct {
	repo        repository.UserRepository
	studentRepo repository.StudentRepository
	googleCfg   *config.GoogleOAuthConfig
}

func NewAuthUsecase(repo repository.UserRepository, cfg *config.GoogleOAuthConfig) AuthService {
	return &authService{repo: repo, googleCfg: cfg}
}

func NewAuthUseWithStudent(repo repository.UserRepository, studentRepo repository.StudentRepository, cfg *config.GoogleOAuthConfig) AuthService {
	return &authService{repo: repo, studentRepo: studentRepo, googleCfg: cfg}
}

func (u *authService) GetGoogleLoginURL() string {
	return u.googleCfg.Config.AuthCodeURL("state-token") // ในโปรดักชั่นควรสุ่ม state
}

func (u *authService) ProcessGoogleLogin(code string) (*models.User, error) {
	fmt.Println("--- 1. เข้าสู่ ProcessGoogleLogin แล้ว ---")

	token, err := u.googleCfg.Config.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("--- 2. แลก Token ไม่สำเร็จ:", err, " ---")
		return nil, err
	}

	// เรียกดึงข้อมูลจาก Google API
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("--- 3. ดึงข้อมูล User ไม่สำเร็จ ---")
		return nil, err
	}
	defer resp.Body.Close()

	var googleUser struct {
		Email      string `json:"email"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture    string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		fmt.Println("--- 4. Decode JSON ไม่สำเร็จ ---")
		return nil, err
	}

	now := time.Now()

	existing, err := u.repo.GetUserByEmail(googleUser.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("--- 5. Repository Error:", err, " ---")
		return nil, err
	}

	if existing == nil {
		// --- สร้างผู้ใช้ใหม่ (เหมือน register สำหรับ student) ---
		user := &models.User{
			Email:        googleUser.Email,
			Firstname:    googleUser.GivenName,  // ใช้ Firstname (n ตัวเล็ก) ตามที่คุณกำหนด
			Lastname:     googleUser.FamilyName, // ใช้ Lastname (n ตัวเล็ก) ตามที่คุณกำหนด
			ImagePath:    googleUser.Picture,    // ใช้ ImagePath ตามที่คุณกำหนด
			Provider:     "google",              // ระบุเป็น google เพื่อแยกกับ 'manual'
			RoleID:       1,
			IsFirstLogin: true,
			CreatedAt:    now,
			LatestUpdate: now,
		}

		fmt.Printf("--- 5. กำลังจะส่งไป Repository: %+v ---\n", user)

		if err := u.repo.UpsertUser(user); err != nil {
			fmt.Println("--- 6. Repository Error:", err, " ---")
			return nil, err
		}

		if u.studentRepo != nil {
			student := &models.Student{
				UserID:        user.UserID,
				StudentNumber: "",
				FacultyID:     0,
				DepartmentID:  0,
			}
			if err := u.studentRepo.Create(context.Background(), student); err != nil {
				return nil, err
			}
		}

		fmt.Println("--- 7. บันทึกสำเร็จ! ---")
		return user, nil
	}

	updates := map[string]interface{}{
		"firstname":     googleUser.GivenName,
		"lastname":      googleUser.FamilyName,
		"image_path":    googleUser.Picture,
		"provider":      "google",
		"latest_update": now,
	}

	updatedUser, err := u.repo.UpdateUserFields(context.Background(), existing.UserID, updates)
	if err != nil {
		return nil, err
	}

	fmt.Println("--- 7. บันทึกสำเร็จ! ---")
	return updatedUser, nil
}

// Register สำหรับ Manual Sign-up (Basic validation, password hash, duplicate check)
func (u *authService) Register(req *authDto.RegisterRequest) (*models.User, error) {
	// ตรวจสอบว่ามี Email อยู่แล้วหรือไม่
	existing, err := u.repo.GetUserByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("email already registered")
	}
	if err != nil {
		// ถ้า error ที่ได้ไม่ใช่ RecordNotFound ให้ส่ง error กลับ
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// ตรวจสอบ Password และ ConfirmPassword อีกครั้ง
	if req.Password != req.ConfirmPassword {
		return nil, fmt.Errorf("passwords do not match")
	}

	// Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:          req.Email,
		HashedPassword: string(hashed),
		Provider:       "manual",
		RoleID:         1, // Default: Student
		IsFirstLogin:   true,
		CreatedAt:      time.Now(),
		LatestUpdate:   time.Now(),
	}

	// บันทึกลง DB
	if err := u.repo.UpsertUser(user); err != nil {
		return nil, err
	}

	// สร้าง Student record ถ้ามี studentRepo
	if u.studentRepo != nil {
		student := &models.Student{
			UserID:        user.UserID,
			StudentNumber: "", // ค่าว่างเพราะยังไม่ได้กรอก
			FacultyID:     0,  // ค่าว่างเพราะยังไม่ได้กรอก
			DepartmentID:  0,  // ค่าว่างเพราะยังไม่ได้กรอก
		}
		// ไม่ return error ถ้าสร้าง student ไม่สำเร็จ เพื่อให้ register สำเร็จ
		_ = u.studentRepo.Create(context.Background(), student)
	}

	return user, nil
}

// AuthenticateAndToken ตรวจสอบรหัสผ่านและออก JWT
func (u *authService) AuthenticateAndToken(ctx context.Context, email, password string) (string, *models.User, error) {
	// ดึง user จาก repository
	user, err := u.repo.GetUserByEmail(email)
	if err != nil || user == nil {
		return "", nil, ErrInvalidCredentials
	}

	// ถ้าเป็น account จาก provider อื่น (เช่น google) ให้ไม่ยอมรับ password แบบ manual
	if user.Provider != "" && user.Provider != "manual" {
		return "", nil, ErrInvalidCredentials
	}

	// เปรียบเทียบ bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	signed, err := u.IssueToken(user)
	if err != nil {
		return "", nil, err
	}

	return signed, user, nil
}

func (u *authService) IssueToken(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	ttl := 24 * time.Hour
	now := time.Now()

	claims := jwt.MapClaims{
		"sub":    fmt.Sprint(user.UserID),
		"email":  user.Email,
		"roleID": user.RoleID,
		"iat":    now.Unix(),
		"exp":    now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (u *authService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	user, err := u.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser อัพเดทข้อมูล user
func (u *authService) UpdateUser(ctx context.Context, userID uint, req *authDto.UpdateUserRequest) (*models.User, error) {
	updates := make(map[string]interface{})

	if req.Firstname != nil {
		updates["firstname"] = *req.Firstname
	}
	if req.Prefix != nil {
		updates["prefix"] = strings.TrimSpace(*req.Prefix)
	}
	if req.Lastname != nil {
		updates["lastname"] = *req.Lastname
	}
	if req.ImagePath != nil {
		updates["image_path"] = *req.ImagePath
	}
	if req.CampusID != nil {
		updates["campus_id"] = *req.CampusID
	}
	if req.RoleID != nil {
		updates["role_id"] = *req.RoleID
	}
	if req.IsFirstLogin != nil {
		updates["is_first_login"] = *req.IsFirstLogin
	}

	updates["latest_update"] = time.Now()

	if len(updates) == 1 { // เฉพาะ latest_update
		return u.repo.GetUserByID(userID)
	}

	return u.repo.UpdateUserFields(ctx, userID, updates)
}

// CompleteFirstLogin ตั้งค่าข้อมูลครั้งแรกสำหรับนักศึกษา
func (u *authService) CompleteFirstLogin(ctx context.Context, userID uint, req *authDto.FirstLoginRequest, imagePath string) (*models.User, *models.Student, error) {
	if u.studentRepo == nil {
		return nil, nil, errors.New("student repository not configured")
	}

	prefix := strings.TrimSpace(req.Prefix)
	firstname := strings.TrimSpace(req.Firstname)
	lastname := strings.TrimSpace(req.Lastname)
	imagePath = strings.TrimSpace(imagePath)
	studentNumber := strings.TrimSpace(req.StudentNumber)

	if prefix == "" || firstname == "" || lastname == "" || imagePath == "" || studentNumber == "" {
		return nil, nil, errors.New("missing required fields")
	}
	if req.CampusID <= 0 {
		return nil, nil, errors.New("invalid campus id")
	}
	if req.FacultyID == 0 {
		return nil, nil, errors.New("invalid faculty id")
	}
	if req.DepartmentID == 0 {
		return nil, nil, errors.New("invalid department id")
	}
	if err := validateStudentNumber(studentNumber); err != nil {
		return nil, nil, err
	}

	if existing, err := u.studentRepo.GetByStudentNumber(ctx, studentNumber); err == nil && existing != nil && existing.UserID != userID {
		return nil, nil, errors.New("student_number already in use")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	updates := map[string]interface{}{
		"prefix":         prefix,
		"firstname":      firstname,
		"lastname":       lastname,
		"image_path":     imagePath,
		"campus_id":      req.CampusID,
		"is_first_login": false,
		"latest_update":  time.Now(),
	}

	updatedUser, err := u.repo.UpdateUserFields(ctx, userID, updates)
	if err != nil {
		return nil, nil, err
	}

	student, err := u.studentRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			student = &models.Student{
				UserID:        userID,
				StudentNumber: studentNumber,
				FacultyID:     req.FacultyID,
				DepartmentID:  req.DepartmentID,
			}
			if err := u.studentRepo.Create(ctx, student); err != nil {
				return updatedUser, nil, err
			}
			student, err = u.studentRepo.GetByUserID(ctx, userID)
			if err != nil {
				return updatedUser, nil, err
			}
		} else {
			return updatedUser, nil, err
		}
	} else {
		student.StudentNumber = studentNumber
		student.FacultyID = req.FacultyID
		student.DepartmentID = req.DepartmentID
		if err := u.studentRepo.Update(ctx, student); err != nil {
			return updatedUser, nil, err
		}
	}

	return updatedUser, student, nil
}
// func validateStudentNumber(studentNumber string) error {
// 	if len(studentNumber) != 10 {
// 		return fmt.Errorf("student_number must be exactly 10 digits")
// 	}
// 	for _, r := range studentNumber {
// 		if r < '0' || r > '9' {
// 			return fmt.Errorf("student_number must be exactly 10 digits")
// 		}
// 	}
// 	return nil
// }
