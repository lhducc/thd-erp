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

	// transaction
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
	roleRepo := repository.NewRoleRepo(db)

	//checkin Repo
	workShiftRepo := checkinRepo.NewWorkShiftStore(db)
	attendanceRecordRepo := checkinRepo.NewAttendanceRecordRepository(db)
	categoryRepo := checkinRepo.NewAttendanceCategoryRepository(db)
	employeeWorkShiftRepo := checkinRepo.NewEmployeeWorkShiftRepo(db)
	workScheduleRepo := checkinRepo.NewWorkScheduleRepo(db)
	timesheetListRepo := checkinRepo.NewTimesheetListRepo(db)
	timesheetRepo := checkinRepo.NewTimesheetRepo(db)
	timesheetDetailRepo := checkinRepo.NewTimesheetDetailRepo(db)

	//usecase hrm
	accountUsecase := authUsecase.NewAuthentication(accountAuthRepo, jwtMaker)
	contractTypeUsecase := usecase.NewContractTypeUsecase(contractTypeRepo)
	decisionTypeUsecase := usecase.NewDecisionTypeUsecase(decisionTypeRepo)
	insuranceUsecase := usecase.NewInsuranceUsecase(insuranceRepo)
	departmentUsecase := usecase.NewDepartmentBiz(departmentRepo)
	officeUsecase := usecase.NewOfficeBiz(officeRepo)
	jobTitleUsecase := usecase.NewJobTitleBiz(jobTitleRepo, hierarchyLevelRepo)
	employeeUsecase := usecase.NewEmployeeBiz(db, userRepo, accountManagementRepo, timesheetRepo, timesheetListRepo, departmentRepo)
	contractUsecase := usecase.NewContractBiz(contractRepo, userRepo, contractTypeRepo)
	employeeDocumentUsecase := usecase.NewEmployeeDocumentBiz(employeeDocumentRepo)
	positionUsecase := usecase.NewPositionBiz(positionRepo)
	decisionUsecase := usecase.NewDecisionBiz(decisionRepo, userRepo)
	documentTypUsecase := usecase.NewDocumentTypeBiz(documentTypeRepo)
	hierarchyLevelUsecase := usecase.NewHierarchyLevelBiz(hierarchyLevelRepo)
	allowanceUsecase := usecase.NewAllowanceBiz(allowanceRepo)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)

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
	roleHanlder := handler.NewRoleHandler(roleUsecase)

	//checkin service
	attendanceRecordService := checkinService.NewAttendanceRecordService(attendanceRecordRepo, categoryRepo, officeRepo)
	attendanceCategoryService := checkinService.NewAttendanceCategoryService(categoryRepo, userRepo)
	workShiftService := checkinService.NewWorkShiftService(workShiftRepo, userRepo, workScheduleRepo)
	employeeWorkShiftService := checkinService.NewEmployeeWorkshiftService(employeeWorkShiftRepo, userRepo, workShiftRepo, workScheduleRepo)
	workScheduleService := checkinService.NewWorkScheduleService(workScheduleRepo, userRepo)
	shiftAllocationService := checkinService.NewShiftAllocationService(workScheduleRepo, userRepo, employeeWorkShiftRepo)
	timesheetListService := checkinService.NewTimesheetListSerivce(timesheetListRepo, userRepo, timesheetRepo, timesheetDetailRepo)
	timesheetService := checkinService.NewTimesheetService(timesheetRepo, timesheetListRepo, attendanceRecordRepo, employeeWorkShiftRepo, timesheetDetailRepo)

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

	managerRouter := hrmRouter.Group("")
	managerRouter.Use(middleware.RoleMiddleware("manager", "admin"))
	adminRouter := managerRouter.Group("")
	adminRouter.Use(middleware.RoleMiddleware("admin"))

	setupOfficeRoutes(adminRouter, hrmRouter, officeHandler)
	setupDepartmentRoutes(adminRouter, departmentHandler)
	setupJobTitleRoutes(adminRouter, jobtitleHandler)
	setupPositionRoutes(adminRouter, positionHandler)
	setuphierarchyLevelRoutes(adminRouter, hierarchyLevel)
	setupEmployeeRouters(adminRouter, hrmRouter, employeeHandler)
	setupContractRoutes(adminRouter, contractHandler)
	setupDecisionRoutes(adminRouter, decisionHandler)
	setupDocumentTypeRoutes(adminRouter, documentTypeHandler)
	setupInsuranceRoutes(adminRouter, insuranceHandler)
	setupDecisionTypeRoutes(adminRouter, decisionTypeHandler)
	setupContractTypeRoutes(adminRouter, contractTypeHandler)
	setupAllowanceRoutes(adminRouter, allowanceHandler)
	setupEmployeeDocumentRoutes(adminRouter, employeeDocumentHandler)
	setupRoleRoutes(adminRouter, roleHanlder)

	//CheckIn
	setupWorkShiftRoutes(adminRouter, hrmRouter, workShiftHandler)
	setupAttandanceRecordRoutes(adminRouter, hrmRouter, attandanceRecord)
	setupAttendanceCategory(adminRouter, hrmRouter, attendanceCategoryHandler)
	setupEmployeeWorkshiftRoutes(adminRouter, managerRouter, hrmRouter, employeeWorkShiftHandler)
	setupWorkScheduleRoutes(adminRouter, workScheduleHandler)
	setupShiftAllocationRoutes(managerRouter, shiftAllocationHandler)
	setupTimesheetListRoutes(adminRouter, timesheetListHandler)
	setupTimesheetRoutes(adminRouter, hrmRouter, timesheetHandler)

}

func setupEmployeeWorkshiftRoutes(adminRouter, managerRouter, userRouter *gin.RouterGroup, handler *checking_handler.EmployeeWorkshiftHandler) {
	adminGr := adminRouter.Group("/employee-workshifts")
	{
		adminGr.GET("", handler.GetAll())
		adminGr.PUT("/:id", handler.Update())
		adminGr.GET("/:employeeID", handler.GetAllByEmployeeID())
		adminGr.POST("/assign", handler.AssignEmpShifts())
		adminGr.POST("/register-many", handler.RegisterMany())
	}
	userGr := userRouter.Group("/employee-workshifts")
	{
		userGr.POST("", handler.RegisterPersonal())
		userGr.GET("/register-shift", handler.GetListShiftAllowRegister())
		userGr.GET("/all", handler.GetAllEmployeeWorkshift())
		userGr.GET("/personal", handler.GetAllPersonal())
		userGr.DELETE("/personal/:id", handler.DeletePersonalShift())

	}
	managerGr := managerRouter.Group("/employee-workshifts")
	{
		managerGr.DELETE("/:id", handler.DeleteByManager())
		managerGr.POST("/manager", handler.Register())
	}
}

func setupOfficeRoutes(adminRouter, userRouter *gin.RouterGroup, handler *handler.OfficeHandler) {
	officeRouter := userRouter.Group("/office")
	{
		officeRouter.GET("", handler.GetAllOffice())
		officeRouter.GET("/:id", handler.GetOffice())
	}
	adminOnly := adminRouter.Group("/office")
	{
		adminOnly.POST("", handler.CreateOffice())
		adminOnly.PUT("/:id", handler.UpdateOffice())
		adminOnly.DELETE("/:id", handler.DeleteOffice())
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
		departmentRouter.GET("/office/:office-id", handler.GetDepartmentByOfficeID())
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

func setupEmployeeRouters(adminRouter, userRouter *gin.RouterGroup, handler *handler.EmployeeHandler) {
	adminEmployeeRouter := adminRouter.Group("/employee")
	{
		adminEmployeeRouter.GET("", handler.GetAllEmployees())
		adminEmployeeRouter.POST("", handler.CreateEmployee())
		adminEmployeeRouter.GET("/:id", handler.GetUserById())
		adminEmployeeRouter.PUT("/:id", handler.UpdateEmployee())
		adminEmployeeRouter.DELETE("/:id", handler.DeleteEmployee())
		adminEmployeeRouter.GET("/export", handler.ExportEmployees())
		adminEmployeeRouter.GET("/user", handler.GetEmployeesByRoleID())
		adminEmployeeRouter.PUT("/:id/status", handler.UpdateStatusEmp())
	}

	userEmployeeRouter := userRouter.Group("/employee")
	{
		userEmployeeRouter.GET("/personal", handler.GetPersonalInfById())
	}
}

func setupDecisionRoutes(adminRouter *gin.RouterGroup, handler *handler.DecisionHandler) {
	decisionGroup := adminRouter.Group("/decision")
	decisionGroup.POST("", handler.CreateDecision())
	decisionGroup.PUT("/:id", handler.UpdateDecision())
	decisionGroup.DELETE("/:id", handler.DeleteDecision())
	decisionGroup.GET("/:id", handler.GetDecision())
	decisionGroup.GET("", handler.GetAllDecision())
	decisionGroup.GET("/export", handler.ExportDecision())
}

func setupContractRoutes(adminRouter *gin.RouterGroup, handler *handler.ContractHandler) {
	contractGroup := adminRouter.Group("/contract")
	contractGroup.POST("", handler.CreateContract())
	contractGroup.PUT("/:id", handler.UpdateContract())
	contractGroup.PUT("/:id/reapprove", handler.ReapproveContract())
	contractGroup.DELETE("/:id", handler.DeleteContract())
	contractGroup.GET("/:id", handler.GetContract())
	contractGroup.GET("", handler.GetAllContract())
	contractGroup.GET("/export", handler.ExportContract())
	contractGroup.GET("/employee/:id", handler.GetContractByEmployeeID())
}

func setupDocumentTypeRoutes(adminRouter *gin.RouterGroup, h *handler.DocumentTypeHandler) {
	group := adminRouter.Group("/documenttype")
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

func setupWorkShiftRoutes(adminRouter, userRouter *gin.RouterGroup, workShiftHandler *checkin.WorkShiftHandler) {
	workshiftAdmin := adminRouter.Group("/workshifts")
	{
		workshiftAdmin.POST("", workShiftHandler.CreateWorkShift())
		workshiftAdmin.GET("/:id", workShiftHandler.GetWorkShift())
		workshiftAdmin.GET("", workShiftHandler.GetAllWorkShift())
		workshiftAdmin.PUT("/:id", workShiftHandler.UpdateWorkShift())
		workshiftAdmin.DELETE("/:id", workShiftHandler.DeleteWorkShift())
	}
	workshiftUser := userRouter.Group("/workshifts")
	{
		workshiftUser.GET("/allow-register", workShiftHandler.GetWorkshiftForRegister())
	}

}

func setupAttandanceRecordRoutes(amdinRouter, userRouter *gin.RouterGroup, handler *checkin.AttendanceRecordHandler) {
	adminGroup := amdinRouter.Group("/attendance-record")
	{
		adminGroup.GET("/:id", handler.GetAttendanceRecordByID())
		adminGroup.GET("/employee/:employeeId", handler.GetRecordsByEmployee())
		adminGroup.GET("/employee/:employeeId/date-range", handler.GetRequestsByDateRange())
		adminGroup.GET("/employee/:employeeId/total-request", handler.GetTotalReqOfOneEmployee())
		adminGroup.PUT("/status/:id", handler.UpdateStatusRecord())
		adminGroup.DELETE("/:id", handler.DeleteAttendanceRecord())
		adminGroup.GET("/history-record/:employee-id", handler.GetHistoryRecordByEmployee())
		adminGroup.POST("/history-record/manual", handler.CreateAttendanceRecordByAdmin())
		adminGroup.GET("/history/by-date", handler.GetHistoryByDate())
		adminGroup.GET("/attendance/export", handler.ExportExcelByDate())
	}
	userGr := userRouter.Group("/attendance-record")
	{
		userGr.POST("", handler.CreateAttendanceRecord())
		userGr.GET("/personal", handler.GetPersonalHistoryRecord())
		userGr.GET("/:id/personal", handler.GetPersonalRecordDetailById())
	}

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

func setupAttendanceCategory(adminRouter, userRouter *gin.RouterGroup, attendanceCategoryHandler *checkin.AttendanceCategoryHandler) {
	adminGr := adminRouter.Group("/attendance-category")
	{
		adminGr.POST("", attendanceCategoryHandler.CreateAttendanceCategory())
		adminGr.GET("/:id", attendanceCategoryHandler.GetAttendanceCategory())
		adminGr.GET("", attendanceCategoryHandler.GetAllAttendanceCategories())
		adminGr.PUT("/:id", attendanceCategoryHandler.UpdateAttendanceCategory())
		adminGr.DELETE("/:id", attendanceCategoryHandler.DeleteAttendanceCategory())
	}
	userGr := userRouter.Group("/attendance-category")
	{
		userGr.GET("/office", attendanceCategoryHandler.GetAttendanceCategoryByOfficeID())
	}
}

func setupWorkScheduleRoutes(router *gin.RouterGroup, workScheduleHandler *checkin.WorkScheduleHandler) {
	workSchedule := router.Group("/work-schedule")
	{
		workSchedule.POST("", workScheduleHandler.CreateWorkSchedule())
		workSchedule.PUT("/:id", workScheduleHandler.UpdateWorkScheduleAuto())
		workSchedule.DELETE("/:id", workScheduleHandler.DeleteWorkScheduleAuto())
		workSchedule.POST("/assign/:work-schedule-id", workScheduleHandler.AddManagerToWorkScheduleAuto())
		workSchedule.GET("", workScheduleHandler.FetchListWorkScheduleAuto())
		workSchedule.GET("/:id", workScheduleHandler.FetchWorkScheduleByID())
		workSchedule.DELETE("/manager/:schedule-id", workScheduleHandler.DeleteManagerFromWorkScheduleAuto())
		//workSchedule.PUT("/:id/recurring", workScheduleHandler.UpdateStatusAutoRecurring())
	}
	workScheduleRegister := router.Group("/work-schedule-register")
	{
		workScheduleRegister.POST("", workScheduleHandler.CreateWorkSchedule())
		workScheduleRegister.PUT("/:id", workScheduleHandler.UpdateWorkScheduleRegister())
		workScheduleRegister.DELETE("/:id", workScheduleHandler.DeleteWorkScheduleRegister())
		workScheduleRegister.POST("/assign/:work-schedule-id", workScheduleHandler.AddManagerToWorkScheduleRegister())
		workScheduleRegister.GET("/:id", workScheduleHandler.FetchWorkScheduleByID())
		workScheduleRegister.DELETE("/manager/:schedule-id", workScheduleHandler.DeleteManagerFromWorkScheduleRegister())
		workScheduleRegister.GET("/register", workScheduleHandler.FetchListWorkScheduleRegister())
	}
}

func setupShiftAllocationRoutes(managerRouter *gin.RouterGroup, shiftAllocationHandler *checkin.ShiftAllocationHandler) {
	shiftAllocation := managerRouter.Group("/shift-allocation")
	{
		shiftAllocation.GET("", shiftAllocationHandler.GetAll())
	}
}

func setupTimesheetListRoutes(adminRouter *gin.RouterGroup, timesheetListHandler *checkin.TimesheetListHandler) {
	adminGr := adminRouter.Group("/timesheet-list")
	{
		adminGr.GET("/:id", timesheetListHandler.GetByID())
		adminGr.GET("", timesheetListHandler.List())
		adminGr.DELETE("/:id", timesheetListHandler.Delete())
		adminGr.POST("", timesheetListHandler.Create())
		adminGr.PUT("/:id", timesheetListHandler.Update())
		adminGr.PUT("/:id/locked", timesheetListHandler.LockedTimeSheet())
	}
}

func setupTimesheetRoutes(adminRouter, userRouter *gin.RouterGroup, timesheetHandler *checkin.TimesheetHandler) {
	adminGr := adminRouter.Group("/timesheet")
	{
		adminGr.GET("/timesheet-list/:id", timesheetHandler.CalculationTimeSheet())
		adminGr.PUT("/:id", timesheetHandler.AdjustWorkDayManual())
		adminGr.GET("/:id/reset", timesheetHandler.ResetWorkDayAdjustment())
		adminGr.GET("/export/:timesheet-list-id", timesheetHandler.ExportTimesheetListByMonth())
	}
	userGr := userRouter.Group("/timesheet")
	{
		userGr.GET("/personal", timesheetHandler.GetPersonalTimesheet())
	}

}

func setupRoleRoutes(router *gin.RouterGroup, roleHandler *handler.RoleHandler) {
	role := router.Group("/role")
	{
		role.GET("", roleHandler.GetAllRole())
	}
}
