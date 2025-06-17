package router

import (
	authHandler "erp/backend/api/handler/hrm/auth"
	"erp/backend/api/handler/hrm/checkin"
	handler "erp/backend/api/handler/hrm/hr_profile"
	"erp/backend/api/middleware"
	"erp/backend/config"
	authRepo "erp/backend/internal/auth/repository"
	"erp/backend/internal/auth/token"
	authUsecase "erp/backend/internal/auth/usecase"
	"erp/backend/internal/hrm/hr_profile/repository"
	"erp/backend/internal/hrm/hr_profile/usecase"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api/v1")

	setupHRMRoutes(api, db)

}

func setupHRMRoutes(router *gin.RouterGroup, db *gorm.DB) {
	RegisterRoutes(router, db)
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {

	config.LoadEnv()
	jwtMaker, err := token.NewJWTMaker(config.GetSecretKey())
	if err != nil {
		log.Fatal("Error creating JWT maker:", err)
	}

	officeHandler := handler.NewOficeHandler(db)
	departmentHandler := handler.NewDepartmentHandler(db)
	jobtitleHandler := handler.NewJobTitleHandler(db)
	hierarchyLevel := handler.NewHierarchyLevelHandler(db)
	positionHandler := handler.NewPositionHandler(db)
	employeeHandler := handler.NewEmployeeHandler(db)
	contractHandler := handler.NewContractHandler(db)
	decisionHandler := handler.NewDecisionHandler(db)
	accountRepo := authRepo.NewAccountRepository(db)
	accountUsecase := authUsecase.NewAuthentication(accountRepo, jwtMaker)
	accountHandler := authHandler.NewAuthenticationHandler(accountUsecase)
	documentTypeHandler := handler.NewDocumentTypeHandler(db)
	contractTypeRepo := repository.NewContractTypeRepository(db)
	contractTypeUsecase := usecase.NewContractTypeUsecase(contractTypeRepo)
	contractTypeHandler := handler.NewContractTypeHandler(contractTypeUsecase)
	decisionTypeRepo := repository.NewDecisionTypeRepository(db)
	decisionTypeUsecase := usecase.NewDecisionTypeUsecase(decisionTypeRepo)
	decisionTypeHandler := handler.NewDecisionTypeHandler(decisionTypeUsecase)
	insuranceRepo := repository.NewInsuranceRepository(db)
	insuranceUsecase := usecase.NewInsuranceUsecase(insuranceRepo)
	insuranceHandler := handler.NewInsuranceHandler(insuranceUsecase)
	employeeDocumentHandler := handler.NewEmployeeDocumentHandler(db)

	//CheckIn
	workShiftHandler := checkin.NewWorkShiftHandler(db)
	allow := checkin.NewAllowedWorkingScheduleHandler(db)
	holiday := checkin.NewHolidayHandler(db)
	attandanceForm := checkin.NewAttendanceFormHandler(db)

	public := router.Group("/auth")
	public.POST("/login", accountHandler.SignIn)

	hrmRouter := router.Group("")
	hrmRouter.Use(middleware.AuthMiddleware(jwtMaker))
	hrmRouter.PUT("/auth/password", accountHandler.ChangePassword)

	setupOfficeRoutes(hrmRouter, officeHandler)
	setupDepartmentRoutes(hrmRouter, departmentHandler)
	setupJobTitleRoutes(hrmRouter, jobtitleHandler)
	setupPositionRoutes(hrmRouter, positionHandler)
	setuphierarchyLevelRoutes(hrmRouter, hierarchyLevel)
	setupEmployeeRouters(hrmRouter, employeeHandler)
	setupContractRoutes(hrmRouter, contractHandler)
	setupDecisionRoutes(hrmRouter, decisionHandler)
	setupDocumentTypeRoutes(hrmRouter, documentTypeHandler)
	setupInsuranceRoutes(hrmRouter, insuranceHandler)
	setupDecisionTypeRoutes(hrmRouter, decisionTypeHandler)
	setupContractTypeRoutes(hrmRouter, contractTypeHandler)
	setupEmployeeDocumentRoutes(hrmRouter, employeeDocumentHandler)

	//CheckIn
	setupAllowedWorkingScheduleRoutes(hrmRouter, allow)
	setupWorkShiftRoutes(hrmRouter, workShiftHandler)
	setupHolidayRoutes(hrmRouter, holiday)
	setupAttandanceFormRoutes(hrmRouter, attandanceForm)

}

func setupOfficeRoutes(router *gin.RouterGroup, handler *handler.OfficeHandler) {
	officeRouter := router.Group("/office")
	{
		officeRouter.GET("", handler.GetAllOffice())
		officeRouter.POST("", handler.CreateOffice())
		officeRouter.GET("/:id", handler.GetOffice())
		officeRouter.PUT("/:id", handler.UpdateOffice())
		officeRouter.DELETE("/:id", handler.DeleteOffice())
	}
}

func setupDepartmentRoutes(router *gin.RouterGroup, handler *handler.DepartmentHandler) {
	departmentRouter := router.Group("/department")
	{
		departmentRouter.GET("", handler.GetAllDepartment())
		departmentRouter.POST("", handler.CreateDepartment())
		departmentRouter.GET("/:id", handler.GetDepartment())
		departmentRouter.PUT("/:id", handler.UpdateDepartment())
		departmentRouter.DELETE("/:id", handler.DeleteDepartment())
	}
}

func setupPositionRoutes(router *gin.RouterGroup, handler *handler.PositionHandler) {
	positionRouter := router.Group("/position")
	{
		positionRouter.GET("", handler.GetAll())
		positionRouter.POST("", handler.Create())
		positionRouter.GET("/:id", handler.GetByID())
		positionRouter.PUT("/:id", handler.UpdateByID())
		positionRouter.DELETE("/:id", handler.DeleteByID())
	}
}

func setupJobTitleRoutes(router *gin.RouterGroup, handler *handler.JobTitleHandler) {
	jobTitleRouter := router.Group("/jobtitle")
	{
		jobTitleRouter.GET("", handler.GetAll())
		jobTitleRouter.POST("", handler.Create())
		jobTitleRouter.GET("/:id", handler.GetByID())
		jobTitleRouter.PUT("/:id", handler.UpdateByID())
		jobTitleRouter.DELETE("/:id", handler.DeleteByID())
	}
}

func setuphierarchyLevelRoutes(router *gin.RouterGroup, handler *handler.HierarchyLevelHandler) {
	hierarchyLevelRouter := router.Group("/hierarchylevel")
	{
		hierarchyLevelRouter.GET("", handler.GetAll())
		hierarchyLevelRouter.POST("", handler.Create())
		hierarchyLevelRouter.GET("/:id", handler.GetByID())
		hierarchyLevelRouter.PUT("/:id", handler.UpdateByID())
		hierarchyLevelRouter.DELETE("/:id", handler.DeleteByID())
	}
}

func setupEmployeeRouters(router *gin.RouterGroup, handler *handler.EmployeeHandler) {
	employeeRouter := router.Group("/employee")
	{
		employeeRouter.GET("", handler.GetAllEmployees())
		employeeRouter.POST("", handler.CreateEmployee())
		employeeRouter.GET("/:id", handler.GetUserById())
		employeeRouter.PUT("/:id", handler.UpdateEmployee())
		employeeRouter.DELETE("/:id", handler.DeleteEmployee())
		employeeRouter.GET("/export", handler.ExportEmployees())
	}
}

func setupDecisionRoutes(router *gin.RouterGroup, handler *handler.DecisionHandler) {
	decisionGroup := router.Group("/decision")
	decisionGroup.POST("", handler.CreateDecision())
	decisionGroup.PUT("/:id", handler.UpdateDecision())
	decisionGroup.DELETE("/:id", handler.DeleteDecision())
	decisionGroup.GET("/:id", handler.GetDecision())
	decisionGroup.GET("", handler.GetAllDecision())
	decisionGroup.GET("/export", handler.ExportDecision())
}

func setupContractRoutes(router *gin.RouterGroup, handler *handler.ContractHandler) {
	contractGroup := router.Group("/contract")
	contractGroup.POST("", handler.CreateContract())
	contractGroup.PUT("/:id", handler.UpdateContract())
	contractGroup.DELETE("/:id", handler.DeleteContract())
	contractGroup.GET("/:id", handler.GetContract())
	contractGroup.GET("", handler.GetAllContract())
	contractGroup.GET("/export", handler.ExportContract())
}

func setupDocumentTypeRoutes(r *gin.RouterGroup, h *handler.DocumentTypeHandler) {
	group := r.Group("/documenttype")
	{
		group.POST("", h.CreateDocumentType())
		group.PUT("/:id", h.UpdateDocumentType())
		group.DELETE("/:id", h.DeleteDocumentType())
		group.GET("/:id", h.GetDocumentTypeById())
		group.GET("", h.GetAllDocumentType())
	}
}

func setupInsuranceRoutes(r *gin.RouterGroup, h *handler.InsuranceHandler) {
	group := r.Group("/insurance")
	{
		group.POST("", h.CreateInsurance())
		group.PUT("/:id", h.UpdateInsurance())
		group.DELETE("/:id", h.DeleteInsurance())
		group.GET("/:id", h.GetInsuranceById())
		group.GET("", h.GetInsurances())
	}
}

func setupDecisionTypeRoutes(r *gin.RouterGroup, h *handler.DecisionTypeHandler) {
	group := r.Group("/decisiontype")
	{
		group.POST("", h.CreateDecisionType())
		group.PUT("/:id", h.UpdateDecisionType())
		group.DELETE("/:id", h.DeleteDecisionType())
		group.GET("/:id", h.GetDecisionTypeById())
		group.GET("", h.GetDecisionTypes())
	}
}

func setupContractTypeRoutes(r *gin.RouterGroup, h *handler.ContractTypeHandler) {
	group := r.Group("/contracttype")
	{
		group.POST("", h.CreateContractType())
		group.PUT("/:id", h.UpdateContractType())
		group.DELETE("/:id", h.DeleteContractType())
		group.GET("/:id", h.GetContractTypeById())
		group.GET("", h.GetContractTypes())
	}
}

func setupWorkShiftRoutes(router *gin.RouterGroup, workShiftHandler *checkin.WorkShiftHandler) {
	workshift := router.Group("/workshifts")
	{
		workshift.POST("", workShiftHandler.CreateWorkShift())
		workshift.GET("/:id", workShiftHandler.GetWorkShift())
		workshift.GET("", workShiftHandler.GetAllWorkShift())
		workshift.PUT("/:id", workShiftHandler.UpdateWorkShift())
		workshift.DELETE("/:id", workShiftHandler.DeleteWorkShift())
	}
}
func setupHolidayRoutes(router *gin.RouterGroup, holidayHandler *checkin.HolidayHandler) {
	holidays := router.Group("/holiday")
	{
		holidays.POST("", holidayHandler.CreateHoliday())
		holidays.GET("/:id", holidayHandler.GetHoliday())
		holidays.GET("", holidayHandler.GetAllHoliday())
		holidays.PUT("/:id", holidayHandler.UpdateHoliday())
		holidays.DELETE("/:id", holidayHandler.DeleteHoliday())
	}
}

func setupAllowedWorkingScheduleRoutes(router *gin.RouterGroup, scheduleHandler *checkin.AllowedWorkingScheduleHandler) {
	schedule := router.Group("/allowed-working-schedule")
	{
		schedule.POST("", scheduleHandler.CreateAllowedWorkingSchedule())
		schedule.GET("/:id", scheduleHandler.GetAllowedWorkingSchedule())
		schedule.GET("", scheduleHandler.GetAllAllowedWorkingSchedule())
		schedule.PUT("/:id", scheduleHandler.UpdateAllowedWorkingSchedule())
		schedule.DELETE("/:id", scheduleHandler.DeleteAllowedWorkingSchedule())
	}
}

func setupAttandanceFormRoutes(router *gin.RouterGroup, handler *checkin.AttendanceFormHandler) {
	attandanceFormGroup := router.Group("/attandance-form")
	attandanceFormGroup.POST("", handler.CreateAttendanceForm())
	attandanceFormGroup.GET("/:id", handler.GetAttendanceForm())
	attandanceFormGroup.GET("", handler.GetAllAttendanceForm())
	attandanceFormGroup.PUT("/:id", handler.UpdateAttendanceForm())
	attandanceFormGroup.DELETE("/:id", handler.DisableAttendanceForm())
}

func setupEmployeeDocumentRoutes(router *gin.RouterGroup, employeeDocumentHandler *handler.EmployeeDocumentHandler) {
	employeeDoc := router.Group("/employee-document")
	{
		employeeDoc.POST("", employeeDocumentHandler.CreateEmployeeDocument())
		employeeDoc.GET("", employeeDocumentHandler.GetAllEmployeeDocuments())
		employeeDoc.GET("/:id", employeeDocumentHandler.GetEmployeeDocumentById())
		employeeDoc.DELETE("/:id", employeeDocumentHandler.DeleteEmployeeDocument())
		employeeDoc.PUT("/:id", employeeDocumentHandler.UpdateEmployeeDocument())
	}
}
