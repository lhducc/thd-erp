package router

import (
	authHandler "erp/backend/api/handler/hrm/auth"
	"erp/backend/api/handler/hrm/checkin"
	checking_handler "erp/backend/api/handler/hrm/checkin"
	handler "erp/backend/api/handler/hrm/hr_profile"
	"erp/backend/api/middleware"
	"erp/backend/config"
	authRepo "erp/backend/internal/auth/repository"
	"erp/backend/internal/auth/token"
	authUsecase "erp/backend/internal/auth/usecase"
	checkinRepo "erp/backend/internal/hrm/checkin/repository"
	checkinService "erp/backend/internal/hrm/checkin/service"

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

	// hrm_repo
	accountAuthRepo := authRepo.NewAccountRepository(db)
	contractTypeRepo := repository.NewContractTypeRepository(db)
	decisionTypeRepo := repository.NewDecisionTypeRepository(db)
	insuranceRepo := repository.NewInsuranceRepository(db)
	officeRepo := repository.NewOfficeStore(db)
	departmentRepo := repository.NewDepartmentStore(db)
	jobTitleRepo := repository.NewJobTitleStore(db)
	positionRepo := repository.NewPositionStore(db)
	userRepo := repository.NewUserStore(db)
	accountManagementRepo := repository.NewAccountStore(db)
	contractRepo := repository.NewContractStore(db)
	decisionRepo := repository.DecisionStore(db)
	documentTypeRepo := repository.NewDocumentType(db)
	employeeDocumentRepo := repository.NewEmployeeDocument(db)
	hierarchyLevelRepo := repository.NewhierarchyLevelStore(db)
	allowanceRepo := repository.NewAllowanceRepo(db)

	//usecase hrm
	accountUsecase := authUsecase.NewAuthentication(accountAuthRepo, jwtMaker)
	contractTypeUsecase := usecase.NewContractTypeUsecase(contractTypeRepo)
	decisionTypeUsecase := usecase.NewDecisionTypeUsecase(decisionTypeRepo)
	insuranceUsecase := usecase.NewInsuranceUsecase(insuranceRepo)
	departmentUsecase := usecase.NewDepartmentBiz(departmentRepo)
	officeUsecase := usecase.NewOfficeBiz(officeRepo)
	jobTitleUsecase := usecase.NewJobTitleBiz(jobTitleRepo, hierarchyLevelRepo)
	employeeUsecase := usecase.NewEmployeeBiz(userRepo, accountManagementRepo)
	contractUsecase := usecase.NewContractBiz(contractRepo, userRepo)
	employeeDocumentUsecase := usecase.NewEmployeeDocumentBiz(employeeDocumentRepo)
	positionUsecase := usecase.NewPositionBiz(positionRepo)
	decisionUsecase := usecase.NewDecisionBiz(decisionRepo, userRepo)
	documentTypUsecase := usecase.NewDocumentTypeBiz(documentTypeRepo)
	hierarchyLevelUsecase := usecase.NewHierarchyLevelBiz(hierarchyLevelRepo)
	allowanceUsecase := usecase.NewAllowanceBiz(allowanceRepo)

	//hanlder hrm
	officeHandler := handler.NewOficeHandler(officeUsecase)
	departmentHandler := handler.NewDepartmentHandler(departmentUsecase)
	jobtitleHandler := handler.NewJobTitleHandler(jobTitleUsecase)
	hierarchyLevel := handler.NewHierarchyLevelHandler(hierarchyLevelUsecase)
	positionHandler := handler.NewPositionHandler(positionUsecase)
	employeeHandler := handler.NewEmployeeHandler(employeeUsecase, departmentUsecase)
	contractHandler := handler.NewContractHandler(contractUsecase)
	decisionHandler := handler.NewDecisionHandler(decisionUsecase)
	accountHandler := authHandler.NewAuthenticationHandler(accountUsecase)
	documentTypeHandler := handler.NewDocumentTypeHandler(documentTypUsecase)
	contractTypeHandler := handler.NewContractTypeHandler(contractTypeUsecase)
	decisionTypeHandler := handler.NewDecisionTypeHandler(decisionTypeUsecase)
	insuranceHandler := handler.NewInsuranceHandler(insuranceUsecase)
	employeeDocumentHandler := handler.NewEmployeeDocumentHandler(employeeDocumentUsecase)
	allowanceHandler := handler.NewAllowanceHandler(allowanceUsecase)

	//checkin Repo
	workShiftRepo := checkinRepo.NewWorkShiftStore(db)
	attendanceRecordRepo := checkinRepo.NewAttendanceRecordRepository(db)
	categoryRepo := checkinRepo.NewAttendanceCategoryRepository(db)
	employeeWorkShiftRepo := checkinRepo.NewEmployeeWorkShiftRepo(db)
	workScheduleRepo := checkinRepo.NewWorkScheduleRepo(db)
	timesheetListRepo := checkinRepo.NewTimesheetListRepo(db)
	timesheetRepo := checkinRepo.NewTimesheetRepo(db)
	timesheetDetailRepo := checkinRepo.NewTimesheetDetailRepo(db)

	//checkin service
	attendanceRecordService := checkinService.NewAttendanceRecordService(attendanceRecordRepo, categoryRepo, officeRepo)
	attendanceCategoryService := checkinService.NewAttendanceCategoryService(categoryRepo)
	workShiftService := checkinService.NewWorkShiftService(workShiftRepo)
	employeeWorkShiftService := checkinService.NewEmployeeWorkshiftService(employeeWorkShiftRepo, userRepo, workShiftRepo, workScheduleRepo)
	workScheduleService := checkinService.NewWorkScheduleService(workScheduleRepo, userRepo)
	shiftAllocationService := checkinService.NewShiftAllocationService(workScheduleRepo, userRepo, employeeWorkShiftRepo)
	timesheetListService := checkinService.NewTimesheetListSerivce(timesheetListRepo, userRepo, timesheetRepo, timesheetDetailRepo)
	timesheetService := checkinService.NewTimesheetService(timesheetRepo, timesheetListRepo, attendanceRecordRepo, employeeWorkShiftRepo)

	//CheckIn handler
	workShiftHandler := checkin.NewWorkShiftHandler(workShiftService)
	attandanceRecord := checkin.NewAttendanceRecordHandler(attendanceRecordService)
	attendanceCategoryHandler := checkin.NewAttendanceCategoryHandler(attendanceCategoryService)
	employeeWorkShiftHandler := checkin.NewEmployeeWorkshift(employeeWorkShiftService)
	workScheduleHandler := checkin.NewWorkScheduleHandler(workScheduleService)
	shiftAllocationHandler := checkin.NewShiftAllocationHandler(shiftAllocationService)
	timesheetListHandler := checkin.NewTimesheetListHandler(timesheetListService)
	timesheetHandler := checkin.NewTimesheetHandler(timesheetService)

	//setup routes
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
	setupAllowanceRoutes(hrmRouter, allowanceHandler)
	setupEmployeeDocumentRoutes(hrmRouter, employeeDocumentHandler)

	//CheckIn
	setupWorkShiftRoutes(hrmRouter, workShiftHandler)
	setupAttandanceRecordRoutes(hrmRouter, attandanceRecord)
	setupAttendanceCategory(hrmRouter, attendanceCategoryHandler)
	setupEmployeeWorkshiftRoutes(hrmRouter, employeeWorkShiftHandler)
	setupWorkSchedule(hrmRouter, workScheduleHandler)
	setupShiftAllocation(hrmRouter, shiftAllocationHandler)
	setupTimesheetList(hrmRouter, timesheetListHandler)
	setupTimesheet(hrmRouter, timesheetHandler)

}

func setupEmployeeWorkshiftRoutes(router *gin.RouterGroup, handler *checking_handler.EmployeeWorkshiftHandler) {
	employeeWorkshift := router.Group("/employee-workshifts")
	{
		employeeWorkshift.GET("/:employeeID", handler.GetAllByEmployeeID())
		employeeWorkshift.GET("", handler.GetAll())
		employeeWorkshift.POST("", handler.Register())
		employeeWorkshift.PUT("/:id", handler.Update())
		employeeWorkshift.GET("/register-shift", handler.GetListShiftAllowRegister())
		employeeWorkshift.DELETE("/:id", handler.Delete())
	}
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
	contractGroup.PUT("/:id/reapprove", handler.ReapproveContract())
	contractGroup.DELETE("/:id", handler.DeleteContract())
	contractGroup.GET("/:id", handler.GetContract())
	contractGroup.GET("", handler.GetAllContract())
	contractGroup.GET("/export", handler.ExportContract())
	contractGroup.GET("/employee/:id", handler.GetContractByEmployeeID())
}

func setupDocumentTypeRoutes(r *gin.RouterGroup, h *handler.DocumentTypeHandler) {
	group := r.Group("/documenttype")
	{
		group.POST("", h.CreateDocumentType())
		group.PUT("/:id", h.UpdateDocumentType())
		group.DELETE("/:id", h.DeleteDocumentType())
		group.GET("/:id", h.GetDocumentTypeById())
		group.GET("", h.GetAllDocumentType())
		group.GET("/enums", h.GetEnumDocumentType())
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

func setupAttandanceRecordRoutes(router *gin.RouterGroup, handler *checkin.AttendanceRecordHandler) {
	attandanceRecordGroup := router.Group("/attendance-record")
	attandanceRecordGroup.POST("", handler.CreateAttendanceRecord())
	attandanceRecordGroup.GET("/:id", handler.GetAttendanceRecordByID())
	attandanceRecordGroup.GET("/employee/:employeeId", handler.GetRecordsByEmployee())
	attandanceRecordGroup.GET("/employee/:employeeId/date-range", handler.GetRequestsByDateRange())
	attandanceRecordGroup.GET("/employee/:employeeId/total-request", handler.GetTotalReqOfOneEmployee())
	attandanceRecordGroup.PUT("/status/:id", handler.UpdateStatusRecord())
	attandanceRecordGroup.DELETE("/:id", handler.DeleteAttendanceRecord())
	attandanceRecordGroup.GET("/history-record/:employee-id", handler.GetHistoryRecordByEmployee())
	attandanceRecordGroup.POST("/history-record/manual", handler.CreateAttendanceRecordByAdmin())
}

func setupAllowanceRoutes(router *gin.RouterGroup, handler *handler.AllowanceHandler) {
	group := router.Group("/allowance")
	{
		group.POST("", handler.Create())
		group.GET("", handler.GetAll())
		group.DELETE("/:id", handler.Delete())
		group.PUT("/:id", handler.Update())
		group.GET("/:id", handler.GetByID())
	}
}

func setupEmployeeDocumentRoutes(router *gin.RouterGroup, employeeDocumentHandler *handler.EmployeeDocumentHandler) {
	employeeDoc := router.Group("/employee-document")
	{
		employeeDoc.POST("", employeeDocumentHandler.CreateEmployeeDocument())
		employeeDoc.GET("", employeeDocumentHandler.GetAllEmployeeDocuments())
		employeeDoc.GET("/:id", employeeDocumentHandler.GetEmployeeDocumentById())
		employeeDoc.DELETE("/:id", employeeDocumentHandler.DeleteEmployeeDocument())
		employeeDoc.PUT("/:id", employeeDocumentHandler.UpdateEmployeeDocument())
		employeeDoc.GET("/employee/:id", employeeDocumentHandler.GetEmployeeDocumentByEmployeeId())
	}
}

func setupAttendanceCategory(router *gin.RouterGroup, attendanceCategoryHandler *checkin.AttendanceCategoryHandler) {
	attendanceCategory := router.Group("/attendance-category")
	{
		attendanceCategory.POST("", attendanceCategoryHandler.CreateAttendanceCategory())
		attendanceCategory.GET("/:id", attendanceCategoryHandler.GetAttendanceCategory())
		attendanceCategory.GET("", attendanceCategoryHandler.GetAllAttendanceCategories())
		attendanceCategory.PUT("/:id", attendanceCategoryHandler.UpdateAttendanceCategory())
		attendanceCategory.DELETE("/:id", attendanceCategoryHandler.DeleteAttendanceCategory())
		attendanceCategory.GET("/office/:officeId", attendanceCategoryHandler.GetAttendanceCategoryByOfficeID())
	}
}

func setupWorkSchedule(router *gin.RouterGroup, workScheduleHandler *checkin.WorkScheduleHandler) {
	workSchedule := router.Group("/work-schedule")
	{
		workSchedule.POST("", workScheduleHandler.CreateWorkSchedule())
		workSchedule.PUT("/:id", workScheduleHandler.UpdateWorkScheduleAuto())
		workSchedule.DELETE("/:id", workScheduleHandler.DeleteWorkScheduleAuto())
		workSchedule.POST("/assign/:work-schedule-id", workScheduleHandler.AddManagerToWorkScheduleAuto())
		workSchedule.GET("", workScheduleHandler.FetchListWorkScheduleAuto())
		workSchedule.GET("/:id", workScheduleHandler.FetchWorkScheduleByID())
		workSchedule.DELETE("/manager/:schedule-id", workScheduleHandler.DeleteManagerFromWorkScheduleAuto())
		workSchedule.GET("/register", workScheduleHandler.FetchListWorkScheduleRegister())
	}
	workScheduleRegister := router.Group("/work-schedule-register")
	{
		workScheduleRegister.POST("", workScheduleHandler.CreateWorkSchedule())
		workScheduleRegister.PUT("/:id", workScheduleHandler.UpdateWorkScheduleRegister())
		workScheduleRegister.DELETE("/:id", workScheduleHandler.DeleteWorkScheduleRegister())
		workScheduleRegister.POST("/assign/:work-schedule-id", workScheduleHandler.AddManagerToWorkScheduleRegister())
		workScheduleRegister.GET("", workScheduleHandler.FetchListWorkScheduleRegister())
		workScheduleRegister.GET("/:id", workScheduleHandler.FetchWorkScheduleByID())
		workScheduleRegister.DELETE("/manager/:schedule-id", workScheduleHandler.DeleteManagerFromWorkScheduleRegister())
		workScheduleRegister.GET("/register", workScheduleHandler.FetchListWorkScheduleRegister())
	}
}

func setupShiftAllocation(router *gin.RouterGroup, shiftAllocationHandler *checkin.ShiftAllocationHandler) {
	shiftAllocation := router.Group("/shift-allocation")
	{
		shiftAllocation.GET("", shiftAllocationHandler.GetAllByHR())
	}
}

func setupTimesheetList(router *gin.RouterGroup, timesheetListHandler *checkin.TimesheetListHandler) {
	timesheetList := router.Group("/timesheet-list")
	{
		timesheetList.GET("/:id", timesheetListHandler.GetByID())
		timesheetList.GET("", timesheetListHandler.List())
		timesheetList.DELETE("/:id", timesheetListHandler.Delete())
		timesheetList.POST("", timesheetListHandler.Create())
		timesheetList.PUT("/:id", timesheetListHandler.Update())
		timesheetList.PUT("/:id/locked", timesheetListHandler.LockedTimeSheet())
	}
}

func setupTimesheet(router *gin.RouterGroup, timesheetHandler *checkin.TimesheetHandler) {
	timesheet := router.Group("/timesheet")
	{
		timesheet.GET("/personal", timesheetHandler.GetPersonalTimesheet())
		timesheet.GET("/:id", timesheetHandler.CalculationTimeSheet())
	}
}
