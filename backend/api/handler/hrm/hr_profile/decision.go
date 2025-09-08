package handler

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"

	// ...existing code...

	"github.com/google/uuid"

	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"erp/backend/pkg/minIO"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DecisionBiz interface {
	CreateDecision(ctx context.Context, data *model.DecisionCreate) (string, error)
	GetDecision(ctx context.Context, id string) (*model.Decision, error)
	GetAllDecisionPagination(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]model.Decision, int64, error)
	UpdateDecision(ctx context.Context, id string, data *model.DecisionCreate) error
	DeleteDecision(ctx context.Context, id string) error
	ExportDecisionTest(ctx context.Context, selectedFields []string) ([]byte, string, error)
}

type DecisionHandler struct {
	decisionBiz DecisionBiz
}

func NewDecisionHandler(biz DecisionBiz) *DecisionHandler {
	return &DecisionHandler{
		decisionBiz: biz,
	}
}

func (h *DecisionHandler) CreateDecision() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.DecisionCreate

		// Support multipart/form-data so client can upload a file
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			// if not multipart, fall back to JSON bind
			if err := c.ShouldBindJSON(&data); err != nil {
				utils.ResponseMessage(c, fmt.Sprintf("Invalid input: %v", err), http.StatusBadRequest, nil)
				return
			}
		} else {
			// Bind form fields to struct
			if err := c.ShouldBind(&data); err != nil {
				utils.ResponseMessage(c, fmt.Sprintf("Invalid form input: %v", err), http.StatusBadRequest, nil)
				return
			}

			// Handle attached file if provided
			file, header, err := c.Request.FormFile("file")
			if err == nil && file != nil {
				defer file.Close()
				// Validate extension and allow only .pdf and .docx
				ext := strings.ToLower(filepath.Ext(header.Filename))
				allowed := map[string]bool{
					".pdf":  true,
					".docx": true,
				}
				if !allowed[ext] {
					utils.ResponseMessage(c, "Định dạng file không được hỗ trợ. Chỉ cho phép pdf, docx, jpg, jpeg, png", http.StatusBadRequest, nil)
					return
				}

				// Determine content type; fall back based on extension if header is missing
				contentType := header.Header.Get("Content-Type")
				if contentType == "" {
					switch ext {
					case ".pdf":
						contentType = "application/pdf"
					case ".docx":
						contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
					default:
						contentType = "application/octet-stream"
					}
				}

				// generate unique object name (uuid + original ext)
				objectName := uuid.New().String() + ext
				// use UploadImageToMinIO helper
				if err := minIO.UploadImageToMinIO(c.Request.Context(), minIO.DecisionBucket, objectName, file, header.Size, contentType, 0); err != nil {
					utils.ResponseMessage(c, fmt.Sprintf("Upload file failed: %v", err), http.StatusInternalServerError, nil)
					return
				}
				// store object name in DB; GET will generate presigned URL when serving
				data.AttachedFile = objectName
			}
		}
		fmt.Println(data)
		code, err := h.decisionBiz.CreateDecision(c.Request.Context(), &data)
		if err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Tạo thành công", http.StatusOK, gin.H{
			"decision_id": code,
		})

	}
}

func (h *DecisionHandler) GetDecision() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		d, err := h.decisionBiz.GetDecision(c.Request.Context(), idParam)
		if err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}
		var decisionTypeName string
		if d.DecisionType != nil {
			decisionTypeName = d.DecisionType.DecisionType
		} else {
			decisionTypeName = ""
		}

		response := model.DecisionResponse{
			DecisionID:       d.DecisionID,
			DecisionName:     d.DecisionName,
			DecisionTypeID:   d.DecisionTypeID,
			DecisionTypeName: decisionTypeName,
			EffectiveDate:    d.EffectiveDate.Format("2006-01-02"),
			SignDate:         d.SignDate.Format("2006-01-02"),
			Condition:        d.Condition,
			Content:          d.Content,
			AttachedFile:     d.AttachedFile,
			CreatedDate:      d.CreatedDate,
		}
		response.Employees = make([]model.EmployeeShort, 0)
		for _, emp := range d.Employees {
			response.Employees = append(response.Employees, model.EmployeeShort{
				EmployeeID: emp.EmployeeID,
				Fullname:   emp.Fullname,
			})
		}
		if d.DecisionType != nil {
			response.DecisionTypeName = d.DecisionType.DecisionType
		}
		// If AttachedFile is present but not a full URL, try to generate presigned/fallback URL
		if response.AttachedFile != "" && !strings.HasPrefix(response.AttachedFile, "http") {
			if url, err := minIO.GeneratePresignedURL(c.Request.Context(), minIO.DecisionBucket, response.AttachedFile, 15*time.Minute); err == nil {
				response.AttachedFile = url
			} else {
				log.Printf("Warning: failed to generate decision attached file URL for %s: %v", response.AttachedFile, err)
				// keep original object name or set empty; here we keep original name
			}
		}

		utils.ResponseMessage(c, "Thông tin quyết định", http.StatusOK, response)
	}
}

func (h *DecisionHandler) GetAllDecision() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		// Filter processing if necessary
		filterFields := []string{"decision_type_id", "condition"}
		filters := make(map[string]interface{})

		for _, field := range filterFields {
			values := ctx.QueryArray(field)
			cleaned := make([]string, 0)
			for _, v := range values {
				for _, part := range strings.Split(v, ",") {
					if trimmed := strings.TrimSpace(part); trimmed != "" {
						cleaned = append(cleaned, trimmed)
					}
				}
			}
			if len(cleaned) > 0 {
				filters[field] = cleaned
			}
		}

		decisions, totalRecords, err := h.decisionBiz.GetAllDecisionPagination(ctx, page, pageSize, filters)
		if err != nil {
			utils.ResponseMessage(ctx, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			return
		}

		var decisionResponses []model.DecisionResponse
		for _, d := range decisions {
			response := model.DecisionResponse{
				DecisionID:       d.DecisionID,
				DecisionName:     d.DecisionName,
				DecisionTypeID:   d.DecisionTypeID,
				DecisionTypeName: "",
				EffectiveDate:    d.EffectiveDate.Format("2006-01-02"),
				SignDate:         d.SignDate.Format("2006-01-02"),
				Condition:        d.Condition,
				Content:          d.Content,
				AttachedFile:     d.AttachedFile,
				CreatedDate:      d.CreatedDate,
			}

			response.Employees = make([]model.EmployeeShort, 0)
			for _, emp := range d.Employees {
				response.Employees = append(response.Employees, model.EmployeeShort{
					EmployeeID: emp.EmployeeID,
					Fullname:   emp.Fullname,
				})
			}
			if d.DecisionType != nil {
				response.DecisionTypeName = d.DecisionType.DecisionType
			}
			// convert attached file object name to presigned/fallback URL like attendance handler
			// if response.AttachedFile != "" && !strings.HasPrefix(response.AttachedFile, "http") {
			// 	if url, err := minIO.GeneratePresignedURL(ctx.Request.Context(), minIO.DecisionBucket, response.AttachedFile, 15*time.Minute); err == nil {
			// 		response.AttachedFile = url
			// 	} else {
			// 		log.Printf("Warning: failed to generate decision attached file URL for %s: %v", response.AttachedFile, err)
			// 	}
			// }
			decisionResponses = append(decisionResponses, response)
		}
		// Calculate the total number of pages
		totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

		utils.ResponseSuccess(ctx, "Danh sách dữ liệu", http.StatusOK, gin.H{
			"data":         decisionResponses,
			"totalRecords": totalRecords,
			"page":         page,
			"pageSize":     pageSize,
			"totalPages":   totalPages,
		})
	}
}

func (h *DecisionHandler) UpdateDecision() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		var data model.DecisionCreate
		// Support multipart/form-data for update as well
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			if err := c.ShouldBindJSON(&data); err != nil {
				utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
				return
			}
		} else {
			if err := c.ShouldBind(&data); err != nil {
				utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
				return
			}

			file, header, err := c.Request.FormFile("attached_file")
			if err == nil && file != nil {
				defer file.Close()
				// Validate extension and allow .pdf and .docx (and common images)
				ext := strings.ToLower(filepath.Ext(header.Filename))
				allowed := map[string]bool{
					".pdf":  true,
					".docx": true,
					// ".jpg":  true,
					// ".jpeg": true,
					// ".png":  true,
				}
				if !allowed[ext] {
					utils.ResponseMessage(c, "Định dạng file không được hỗ trợ. Chỉ cho phép pdf, docx, jpg, jpeg, png", http.StatusBadRequest, nil)
					return
				}

				// Determine content type; fall back based on extension if header is missing
				contentType := header.Header.Get("Content-Type")
				if contentType == "" {
					switch ext {
					case ".pdf":
						contentType = "application/pdf"
					case ".docx":
						contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
					// case ".jpg", ".jpeg":
					// 	contentType = "image/jpeg"
					// case ".png":
					// 	contentType = "image/png"
					default:
						contentType = "application/octet-stream"
					}
				}

				objectName := uuid.New().String() + ext
				if err := minIO.UploadImageToMinIO(c.Request.Context(), minIO.DecisionBucket, objectName, file, header.Size, contentType, 30); err != nil {
					utils.ResponseMessage(c, fmt.Sprintf("Upload file failed: %v", err), http.StatusInternalServerError, nil)
					return
				}
				// store object name in DB
				data.AttachedFile = objectName
			}
		}

		if err := h.decisionBiz.UpdateDecision(c.Request.Context(), idParam, &data); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Đã cập nhật", http.StatusOK, nil)
	}
}

func (h *DecisionHandler) DeleteDecision() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		if err := h.decisionBiz.DeleteDecision(c.Request.Context(), idParam); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Đã xóa", http.StatusOK, nil)
	}
}

func (biz *DecisionHandler) ExportDecision() gin.HandlerFunc {
	return func(c *gin.Context) {
		fieldStr := c.Query("fields")
		if strings.TrimSpace(fieldStr) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu tham số 'fields'"})
			return
		}

		fields := strings.Split(fieldStr, ",")
		data, filename, err := biz.decisionBiz.ExportDecisionTest(c, fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
	}
}
