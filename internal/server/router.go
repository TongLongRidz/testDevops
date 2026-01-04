package server

import (
	"backend/config"
	"backend/internal/handler/auth"
	"backend/internal/handler/academic_year"
	"backend/internal/repository"
	"backend/internal/usecase"
	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Middleware พื้นฐาน
	app.Use(logger.New())

	// --- 1. Infrastructure / Config ---
	googleConfig := config.LoadGoogleAuthConfig()

	// --- 2. Repository Layer ---
	// สร้าง User Repository เพื่อใช้จัดการข้อมูลผู้ใช้ในฐานข้อมูล
	userRepo := repository.NewUserRepository(db)
	academicYearRepo := repository.NewAcademicYearRepository(db)

	// --- 3. Usecase Layer (Business Logic) ---
	// ส่ง Repository และ Config เข้าไปใน Usecase
	authService := usecase.NewAuthUsecase(userRepo, googleConfig)
	academicYearService := usecase.NewAcademicYearService(academicYearRepo)

	// --- 4. Handler Layer (Controller) ---
	// สร้าง Handler ที่จะรับ HTTP Request
	authHandler := auth.NewAuthHandler(authService)
	academicYearHandler := academic_year.NewAcademicYearHandler(academicYearService)

	// --- 5. Routing Definition ---

	// --- Auth Routes ---
	authGroup := app.Group("/auth")
	authGroup.Get("/google/login", authHandler.GoogleLogin) // Endpoint สำหรับ Redirect ไปหน้า Login ของ Google
	authGroup.Get("/google/callback", authHandler.GoogleCallback) // Endpoint สำหรับรับ Callback หลังจาก User Login สำเร็จ
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/logout", authHandler.Logout)
	authGroup.Get("/me", middleware.RequireAuth(userRepo), authHandler.Me)
	// เหลือ Update User Profile, Change Password, Reset Password

	// --- Academic Year Routes ---
	// แนะนำให้ใช้ RequireAuth ครอบไว้หากต้องการให้เฉพาะคนที่ Login แล้วจัดการข้อมูลได้ ต้องกรองให้ Role เป็นพวก Admin ด้วย
	academicYearGroup := app.Group("/academic-years", middleware.RequireAuth(userRepo))
	academicYearGroup.Get("/latest", academicYearHandler.GetLatest)   // GET /academic-years/latest
	academicYearGroup.Get("/list", academicYearHandler.GetList)          // GET /academic-years
	academicYearGroup.Get("/:id", academicYearHandler.GetByID)        // GET /academic-years/1
	academicYearGroup.Post("/create", academicYearHandler.Create)          // POST /academic-years
	academicYearGroup.Put("/edit/:id", academicYearHandler.Update)        // PUT /academic-years/1
	academicYearGroup.Delete("/delete/:id", academicYearHandler.Delete)     // DELETE /academic-years/1

}