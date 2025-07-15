package checkin

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/service/service_interface"
	"erp/backend/pkg"
	"erp/backend/pkg/job"
	"erp/backend/pkg/minIO"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
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
		req, err := bindAndValidateAttendanceRequest(c)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		employeeID, err := extractEmployeeIDFromContext(c)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		record := model.ConvertToAttendanceRecordStruct(req)
		record.EmployeeID = employeeID

		if err := record.Validate(); err != nil {
			utils.ResponseMessage(c, "Invalid input: "+err.Error(), http.StatusBadRequest, nil)
			return
		}

		//check exist category
		category, err := h.biz.CheckCatrgoryExists(c, &record)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if record.OfficeID != nil && *record.OfficeID != "" && category.AutoApprove == true {
			if err := h.biz.ValidateAttendanceRecordDistance(c, &record, category); err != nil {
				utils.ResponseMessage(c, "Validation failed: "+err.Error(), http.StatusBadRequest, nil)
				return
			}
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

		if err := h.biz.CreateAttendanceRecord(c, &record); err != nil {
			utils.ResponseMessage(c, "Failed to create attendance record: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Attendance record created successfully", http.StatusOK, nil)
	}
}

func processImageIfRequired(c *gin.Context, record *model.AttendanceRecord, employeeID string) error {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		return fmt.Errorf("missing image: %w", err)
	}
	defer file.Close()

	fileName := uuid.New().String()
	imageName := utils.GenerateImageName(employeeID, record.Timestamp, header.Filename)
	record.ImageName = imageName

	job.EnqueueUploadJob(file, employeeID, header, record.Timestamp, fileName)
	return nil
}

func bindAndValidateAttendanceRequest(c *gin.Context) (*model.AttendanceRecordCreate, error) {
	var req model.AttendanceRecordCreate
	if err := c.ShouldBind(&req); err != nil {
		return nil, errors.New("Invalid input data")
	}

	fmt.Printf("%+v\n", req)

	return &req, nil
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
		id := c.Param("id")
		record, err := h.biz.GetAttendanceRecordByID(c, id)
		if err != nil {
			utils.ResponseMessage(c, "Attendance record not found", http.StatusNotFound, nil)
			return
		}

		url, err := minIO.GeneratePresignedURL(c, minIO.AttendanceBucket, record.ImageName, 15*time.Minute)
		record.ImageURL = url

		utils.ResponseSuccess(c, "Attendance record found", http.StatusOK, &record)
	}
}

func (h *AttendanceRecordHandler) DeleteAttendanceRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.biz.DeleteAttendanceRecord(c, id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Attendance record deleted successfully", http.StatusOK, nil)
	}
}

func (h *AttendanceRecordHandler) GetRecordsByEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID := c.Param("employeeId")

		records, err := h.biz.ListAttendanceRecordsByEmployee(c, employeeID)
		if err != nil {
			utils.ResponseMessage(c, "Failed to get records: "+err.Error(), http.StatusInternalServerError, nil)
			return
		}

		addPresignedURLs(c, records)
		utils.ResponseMessage(c, "List of attendance records", http.StatusOK, &records)
	}
}

func (h *AttendanceRecordHandler) GetRecordsByDateRange() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID := c.Param("employeeId")
		from := c.Query("from")
		to := c.Query("to")

		records, err := h.biz.ListAttendanceRecordsByDateRange(c, employeeID, from, to)
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
		employeeID := c.Param("employeeId")
		total, err := h.biz.GetTotalReqOfEmployee(c, employeeID)
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
		id := c.Param("id")
		var updateReq model.AttendanceRecordUpdate
		err := c.ShouldBindJSON(&updateReq)
		if err != nil {
			utils.ResponseMessage(c, "Failed to parse request body", http.StatusBadRequest, nil)
			return
		}

		err = h.biz.UpdateAttendanceRecord(c, &updateReq, id)
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Update thành công", http.StatusOK, nil)
	}
}
