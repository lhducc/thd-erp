package checkin

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/model/dto"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/pkg"
	"erp/backend/pkg/job"
	"erp/backend/pkg/minIO"
	"erp/backend/pkg/variable"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AttendanceRecordHandler struct {
	biz service_interface.AttendanceRecordService
}

func NewAttendanceRecordHandler(biz service_interface.AttendanceRecordService) *AttendanceRecordHandler {
	return &AttendanceRecordHandler{biz: biz}
}

func (h *AttendanceRecordHandler) CreateAttendanceRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req dto.AttendanceRecordCreate
		if err := c.ShouldBind(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
			return
		}

		if err := req.Validate(); err != nil {
			utils.ResponseMessage(c, "Invalid input: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		employeeID, err := extractEmployeeIDFromContext(c)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		record := dto.ConvertToAttendanceRecordStruct(&req)
		record.EmployeeID = employeeID

		//check exist category
		category, err := h.biz.CheckCategoryExists(ctx, &record)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if record.OfficeID != nil && *record.OfficeID != "" && category.AutoApprove == true {
			if err := h.biz.ValidateAttendanceRecordDistance(c, &record, category); err != nil {
				utils.ResponseMessage(c, "Validation failed: "+err.Error(), http.StatusBadRequest, nil)
				return
			}
		} else {
			record.Status = variable.Pending
		}
		if category.IsCheckLocation == false && category.IsGPS == false {
			record.Longitude = nil
			record.Latitude = nil
		}

		if category.IsCamera == true {
			if err := processImageIfRequired(c, &record, employeeID); err != nil {
				utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
				return
			}
		}

		if err := h.biz.CreateAttendanceRecord(ctx, &record); err != nil {
			utils.ResponseMessage(c, "Failed to create attendance record: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Attendance record created successfully", http.StatusOK, nil)
	}
}

func processImageIfRequired(c *gin.Context, record *model.AttendanceRecord, employeeID string) error {
	const maxUploadSize = 10 << 20 // 10MB

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		return fmt.Errorf("missing image: %w", err)
	}
	defer file.Close()

	fileName := uuid.New().String()
	imageName := utils.GenerateImageName(employeeID, record.Timestamp, fileName)
	record.ImageName = imageName

	job.EnqueueUploadJob(file, employeeID, header, record.Timestamp, imageName)
	return nil
}

func extractEmployeeIDFromContext(c *gin.Context) (string, error) {
	employeeID, err := utils.ExtractFromContext[string](c, "employeeId")
	if err != nil {
		return "", fmt.Errorf("Failed to extract employee ID: %w", err)
	}
	return employeeID, nil
}

func (h *AttendanceRecordHandler) GetAttendanceRecordByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		id := c.Param("id")
		record, err := h.biz.GetAttendanceRecordByID(ctx, id)
		if err != nil {
			utils.ResponseMessage(c, "Attendance record not found", http.StatusNotFound, nil)
			return
		}

		url, err := minIO.GeneratePresignedURL(c, minIO.AttendanceBucket, record.ImageName, 15*time.Minute)
		if err != nil {
			utils.ResponseMessage(c, "Failed to generate presigned URL", http.StatusInternalServerError, nil)
			return
		}
		record.ImageURL = url

		utils.ResponseSuccess(c, "Attendance record found", http.StatusOK, &record)
	}
}

func (h *AttendanceRecordHandler) DeleteAttendanceRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		ctx := c.Request.Context()
		if err := h.biz.DeleteAttendanceRecord(ctx, id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Attendance record deleted successfully", http.StatusOK, nil)
	}
}

func (h *AttendanceRecordHandler) GetRecordsByEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID := c.Param("employeeId")
		ctx := c.Request.Context()

		records, err := h.biz.ListAttendanceRecordsByEmployee(ctx, employeeID)
		if err != nil {
			utils.ResponseMessage(c, "Failed to get records: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		addPresignedURLs(c, records)
		utils.ResponseMessage(c, "List of attendance records", http.StatusOK, &records)
	}
}

func (h *AttendanceRecordHandler) GetRequestsByDateRange() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID := c.Param("employeeId")
		from := c.Query("from")
		to := c.Query("to")
		ctx := c.Request.Context()

		records, err := h.biz.ListAttendanceRequestsByDateRange(ctx, employeeID, from, to)
		if err != nil {
			utils.ResponseMessage(c, "Failed to get records by date range: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		addPresignedURLs(c, records)
		utils.ResponseMessage(c, "Records in date range", http.StatusOK, &records)
	}
}

func (h *AttendanceRecordHandler) GetTotalReqOfOneEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		employeeID := c.Param("employeeId")
		total, err := h.biz.GetTotalReqOfEmployee(ctx, employeeID)
		if err != nil {
			utils.ResponseMessage(c, "Failed: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Total request of employee", http.StatusOK, total)
	}
}

func addPresignedURLs(c *gin.Context, records []model.AttendanceRecord) {
	for i := range records {
		if records[i].ImageName == "" {
			continue
		}

		url, err := minIO.GeneratePresignedURL(
			c,
			minIO.AttendanceBucket,
			records[i].ImageName,
			15*time.Minute,
		)
		if err != nil {
			log.Printf("Error generating URL for image %s: %v", records[i].ImageName, err)
			continue
		}
		records[i].ImageURL = url
	}
}

func (h *AttendanceRecordHandler) UpdateStatusRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		id := c.Param("id")
		var updateReq dto.AttendanceRecordUpdate
		err := c.ShouldBindJSON(&updateReq)
		if err != nil {
			utils.ResponseMessage(c, "Failed to parse request body", http.StatusBadRequest, nil)
			return
		}

		if err := updateReq.Validate(); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		err = h.biz.UpdateAttendanceRecord(ctx, &updateReq, id)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Update thành công", http.StatusOK, nil)
	}
}

func (h *AttendanceRecordHandler) GetHistoryRecordByEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID := c.Param("employee-id")
		ctx := c.Request.Context()

		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			utils.ResponseMessage(c, "Giá trị 'page' không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			utils.ResponseMessage(c, "Giá trị 'limit' không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		records, err := h.biz.GetHistoryRecordByEmployee(ctx, employeeID, page, limit)
		if err != nil {
			utils.ResponseMessage(c, "Không thể lấy lịch sử chấm công: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		addPresignedURLs(c, records)
		utils.ResponseMessage(c, "Danh sách chấm công theo nhân viên", http.StatusOK, records)
	}
}

func (h *AttendanceRecordHandler) GetPersonalHistoryRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID := c.GetString("employeeId")
		ctx := c.Request.Context()

		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			utils.ResponseMessage(c, "Giá trị 'page' không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			utils.ResponseMessage(c, "Giá trị 'limit' không hợp lệ", http.StatusBadRequest, nil)
			return
		}

		records, err := h.biz.GetHistoryRecordByEmployee(ctx, employeeID, page, limit)
		if err != nil {
			utils.ResponseMessage(c, "Không thể lấy lịch sử chấm công: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		addPresignedURLs(c, records)
		utils.ResponseMessage(c, "Danh sách chấm công cá nhân", http.StatusOK, records)
	}
}

func (h *AttendanceRecordHandler) CreateAttendanceRecordByAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		role := c.GetString("role")
		if role != "admin" {
			utils.ResponseMessage(c, "Access deni, must be admin role", http.StatusBadRequest, nil)
			return
		}
		var req dto.AttendanceRecordCreateByAdmin
		if err := c.ShouldBind(&req); err != nil {
			utils.ResponseMessage(c, "Invalid input data", http.StatusBadRequest, nil)
			return
		}

		if err := req.Validate(); err != nil {
			utils.ResponseMessage(c, "Invalid input: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		createrID, err := extractEmployeeIDFromContext(c)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		var record model.AttendanceRecord
		record.EmployeeID = req.EmployeeID
		record.Timestamp = req.Timestamp.UTC()
		record.CreatedBy = &createrID
		record.Status = variable.Approved

		if err := h.biz.CreateAttendanceRecord(ctx, &record); err != nil {
			utils.ResponseMessage(c, "Failed to create attendance record: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Attendance record created successfully", http.StatusOK, nil)
	}
}

func (h *AttendanceRecordHandler) GetPersonalRecordDetailById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		id := c.Param("id")
		employeeID := c.GetString("employeeId")
		record, err := h.biz.GetAttendanceRecordByIDPersonal(ctx, id, employeeID)
		if err != nil {
			utils.ResponseMessage(c, "Attendance record not found", http.StatusNotFound, nil)
			return
		}

		url, err := minIO.GeneratePresignedURL(c, minIO.AttendanceBucket, record.ImageName, 15*time.Minute)
		record.ImageURL = url

		utils.ResponseSuccess(c, "Attendance record found", http.StatusOK, &record)
	}
}
