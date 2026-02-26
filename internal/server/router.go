package server

import (
	"backend/config"
	academicyear "backend/internal/handler/academic_year"
	"backend/internal/handler/auth"
	"backend/internal/handler/campus"
	"backend/internal/handler/department"
	"backend/internal/handler/faculty"
	formstatus "backend/internal/handler/form_status"
	"backend/internal/handler/role"
	"backend/internal/handler/student"
	"backend/internal/handler/user"

	awardform "backend/internal/handler/award_form"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Middleware พื้นฐาน
	app.Use(logger.New())

	app.Static("/uploads", "./uploads")
	// --- 1. Infrastructure / Config ---
	googleConfig := config.LoadGoogleAuthConfig()

	// --- 2. Repository Layer ---
	// สร้าง User Repository เพื่อใช้จัดการข้อมูลผู้ใช้ในฐานข้อมูล
	userRepo := repository.NewUserRepository(db)
	awardRepo := repository.NewAwardRepository(db)
	academicYearRepo := repository.NewAcademicYearRepository(db)
	facultyRepo := repository.NewFacultyRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	organizationRepo := repository.NewOrganizationRepository(db)
	roleProfileRepo := repository.NewRoleProfileRepository(db)
	campusRepo := repository.NewCampusRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	formStatusRepo := repository.NewFormStatusRepository(db)

	// --- 3. Usecase Layer (Business Logic) ---
	// ส่ง Repository และ Config เข้าไปใน Usecase
	authService := usecase.NewAuthUsecaseWithRepos(userRepo, studentRepo, organizationRepo, roleProfileRepo, googleConfig)
	academicYearService := usecase.NewAcademicYearService(academicYearRepo)
	studentService := usecase.NewStudentService(studentRepo)
	organizationService := usecase.NewOrganizationService(organizationRepo)
	awardService := usecase.NewAwardUseCase(awardRepo, studentService, organizationService, academicYearService)
	userService := usecase.NewUserUsecase(userRepo)
	facultyService := usecase.NewFacultyService(facultyRepo)
	departmentService := usecase.NewDepartmentService(departmentRepo)
	campusService := usecase.NewCampusService(campusRepo)
	roleService := usecase.NewRoleService(roleRepo)
	formStatusService := usecase.NewFormStatusService(formStatusRepo)

	// --- 4. Handler Layer (Controller) ---
	// สร้าง Handler ที่จะรับ HTTP Request
	authHandler := auth.NewAuthHandlerWithServices(authService, studentService, organizationService)
	awardHandler := awardform.NewAwardHandler(awardService, studentService, academicYearService)
	userHandler := user.NewUserHandlerWithAuth(userService, authService)
	academicYearHandler := academicyear.NewAcademicYearHandler(academicYearService)
	facultyHandler := faculty.NewFacultyHandler(facultyService)
	departmentHandler := department.NewDepartmentHandler(departmentService)
	studentHandler := student.NewStudentHandler(studentService)
	campusHandler := campus.NewCampusHandler(campusService)
	roleHandler := role.NewRoleHandler(roleService)
	formStatusHandler := formstatus.NewFormStatusHandler(formStatusService)

	// --- 5. Routing Definition ---
	apiGroup := app.Group("/api")

	// --- Auth Routes ---
	authGroup := apiGroup.Group("/auth")
	authGroup.Get("/google/login", authHandler.GoogleLogin)       // Endpoint สำหรับ Redirect ไปหน้า Login ของ Google
	authGroup.Get("/google/callback", authHandler.GoogleCallback) // Endpoint สำหรับรับ Callback หลังจาก User Login สำเร็จ
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/create-account", authHandler.CreateAccount)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/logout", authHandler.Logout)
	authGroup.Get("/me", middleware.RequireAuth(userRepo), authHandler.Me)
	authGroup.Put("/me", middleware.RequireAuth(userRepo), authHandler.UpdateMe)
	authGroup.Put("/first-login", middleware.RequireAuth(userRepo), authHandler.FirstLogin)
	// เหลือ Change Password, Reset Password

	// --- Academic Year Routes ---
	academicYearGroup := apiGroup.Group("/academic-years")
	academicYearGroup.Get("/all", academicYearHandler.GetAllAcademicYears)    // ส่ง List เฉพาะปี (ไม่ซ้ำ) เอาไป sort
	academicYearGroup.Post("/create", academicYearHandler.CreateAcademicYear) // สร้างปีการศึกษา ()
	academicYearGroup.Get("/current", academicYearHandler.GetCurrentSemester)

	// --- Faculty Routes ---
	facultyGroup := apiGroup.Group("/faculty")
	facultyGroup.Post("/create", facultyHandler.CreateFaculty)
	facultyGroup.Get("/", facultyHandler.GetAllFaculties)
	facultyGroup.Get("/:id", facultyHandler.GetFacultyByID)

	// --- Department Routes ---
	departmentGroup := apiGroup.Group("/department")
	departmentGroup.Post("/create", departmentHandler.CreateDepartment)
	departmentGroup.Get("/", departmentHandler.GetAllDepartments)
	departmentGroup.Get("/:id", departmentHandler.GetDepartmentByID)
	departmentGroup.Get("/faculty/:facultyId", departmentHandler.GetDepartmentsByFacultyID)

	// --- Student Routes ---
	studentGroup := apiGroup.Group("/students")
	// studentGroup.Get("/me", middleware.RequireAuth(userRepo), studentHandler.GetMyStudent)
	studentGroup.Get("/:id", studentHandler.GetStudentByID)
	// studentGroup.Get("/", studentHandler.GetAllStudents)
	// studentGroup.Post("/user/:userId", studentHandler.CreateStudent) // ?
	// studentGroup.Put("/edit/:id", studentHandler.UpdateStudent) // ?
	// studentGroup.Put("/me", middleware.RequireAuth(userRepo), studentHandler.UpdateMyStudent) // ?
	// studentGroup.Delete("/delete/:id", studentHandler.DeleteStudent) // ?

	awardGroup := apiGroup.Group("/awards", middleware.RequireAuth(userRepo))
	awardGroup.Post("/submit", awardHandler.Submit)                                         // POST /awards/submit
	awardGroup.Get("/search", awardHandler.GetByKeyword)                                    // ค้นหาและกรองพร้อม pagination (query: keyword, date, student_year, page, limit)
	awardGroup.Get("/my/submissions", awardHandler.GetMySubmissions)                        // ดูการส่งฟอร์มของตัวเอง (Student/Organization) - sorted by created_at desc (ทั้งหมดที่เคยส่ง)
	awardGroup.Get("/my/submissions/current", awardHandler.GetMyCurrentSemesterSubmissions) // ดูการส่งฟอร์มของตัวเองในภาคเรียนปัจจุบัน (isActive)
	awardGroup.Get("/types", awardHandler.GetAllAwardTypes)
	awardGroup.Get("/details/:formId", awardHandler.GetByFormID) // GET ดูรายละเอียดฟอร์ม
	awardGroup.Get("/my/approval-logs", awardHandler.GetMyApprovalLogs)
	awardGroup.Get("/approval-logs/:formId", awardHandler.GetApprovalLogDetail) // GET /awards/approval-logs/:id

	awardGroup.Post("/committee/vote/:formId", awardHandler.CommitteeVote)

	awardGroup.Put("/form-status/change/:formId", awardHandler.UpdateFormStatus) //PUT อัพเดท formStatus

	awardGroup.Put("/award-type/change/:formId", awardHandler.UpdateAwardType)

	userGroup := apiGroup.Group("/users", middleware.RequireAuth(userRepo))
	userGroup.Get("/", userHandler.GetAllUsersByCampus) // GET /users (ดึง user ตามวิทยาเขตของคนที่ login)
	userGroup.Get("/info/:id", userHandler.GetUserByID) // GET /users/:id
	userGroup.Put("/promote-chairman/:id", userHandler.ChangeCommitteeRole)
	userGroup.Put("/update/:id", userHandler.UpdateUserByID) // PUT /users/:id

	// --- Campus Routes ---
	campusGroup := apiGroup.Group("/campus")
	campusGroup.Get("/", campusHandler.GetAllCampuses)

	// --- Role Routes ---
	roleGroup := apiGroup.Group("/roles")
	roleGroup.Get("/", roleHandler.GetAllRoles)

	// --- Form Status Routes ---
	formStatusGroup := apiGroup.Group("/form-statuses")
	formStatusGroup.Get("/", formStatusHandler.GetAllFormStatuses)
}
