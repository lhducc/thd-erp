package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"erp/backend/pkg/errors"
	errpkg "erp/backend/pkg/errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"erp/backend/pkg/minIO"

	"github.com/gin-gonic/gin"
)

type ContractBiz interface {
	CreateContract(context context.Context, data *model.ContractCreate) error
	GetContract(ctx context.Context, id string) (*model.Contract, error)
	GetAllContract(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]model.Contract, int64, error)
	UpdateContract(ctx context.Context, id string, data *model.ContractCreate) error
	DeleteContract(ctx context.Context, id string) error
	ExportContractTest(c context.Context, selectedFields []string) ([]byte, string, error)
	UpdateApproveStatus(ctx context.Context, id string, status string) error
	GetContractByEmployeeID(ctx context.Context, employeeId string) ([]model.ContractBasicInfo, error)
}

type ContractHandler struct {
	ContractBiz ContractBiz
}

func NewContractHandler(biz ContractBiz) *ContractHandler {
	return &ContractHandler{
		ContractBiz: biz,
	}
}

func (h *ContractHandler) CreateContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.ContractCreate

		// Support multipart/form-data so client can upload a file; fall back to JSON
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			if err := c.ShouldBindJSON(&data); err != nil {
				utils.ResponseError(c, "Dữ liệu đầu vào không hợp lệ", err, http.StatusBadRequest)
				return
			}
		} else {
			if err := c.ShouldBind(&data); err != nil {
				utils.ResponseError(c, "Dữ liệu đầu vào không hợp lệ", err, http.StatusBadRequest)
				return
			}

			// Handle attached file if provided (field name: attached_file)
			file, header, err := c.Request.FormFile("attached_file")
			if err == nil && file != nil {
				defer file.Close()
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

				contentType := header.Header.Get("Content-Type")
				if contentType == "" {
					switch ext {
					case ".pdf":
						contentType = "application/pdf"
					case ".docx":
						contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
					//case ".jpg", ".jpeg":
					//	contentType = "image/jpeg"
					//case ".png":
					//	contentType = "image/png"
					default:
						contentType = "application/octet-stream"
					}
				}

				objectName := uuid.New().String() + ext
				if err := minIO.UploadImageToMinIO(c.Request.Context(), minIO.ContractBucket, objectName, file, header.Size, contentType, 30); err != nil {
					utils.ResponseMessage(c, fmt.Sprintf("Upload file failed: %v", err), http.StatusInternalServerError, nil)
					return
				}
				data.AttachedFile = objectName
			}
		}

		if err := h.ContractBiz.CreateContract(c.Request.Context(), &data); err != nil {
			utils.ResponseError(c, "Không thể tạo hợp đồng", err, http.StatusInternalServerError)
			return
		}

		utils.ResponseMessage(c, "Tạo hợp đồng thành công", http.StatusOK, gin.H{"data": data})
	}
}
func (h *ContractHandler) GetContractByEmployeeID() gin.HandlerFunc {
	return func(c *gin.Context) {
		employeeID := c.Param("id")
		if employeeID == "" {
			utils.ResponseMessage(c, "Thiếu mã nhân viên", http.StatusBadRequest, nil)
			return
		}

		contracts, err := h.ContractBiz.GetContractByEmployeeID(c.Request.Context(), employeeID)
		if err != nil {
			utils.ResponseMessage(c, "Lỗi khi lấy hợp đồng", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách hợp đồng của nhân viên", http.StatusOK, contracts)
	}
}

func (h *ContractHandler) GetContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		result, err := h.ContractBiz.GetContract(c.Request.Context(), idParam)
		if err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			return
		}

		var res = model.ContractResponse{
			ContractID:    result.ContractId,
			EffectiveDate: result.EffectiveDate,
			ExpiredDate:   result.ExpiredDate,
			SignDate:      result.SignDate,
			Note:          result.Note,
			AttachedFile:  result.AttachedFile,
			Condition:     result.Condition,
			CreatedDate:   result.CreatedDate,
			ContractType:  result.ContractTypeId,
			ApproveStatus: result.ApproveStatus,
			Employee: model.EmployeeSimple{
				EmployeeID: result.Employee.EmployeeID,
				FullName:   result.Employee.Fullname,
				Department: model.DepartmentSimple{
					DepartmentID:   result.Employee.Department.ID,
					DepartmentName: result.Employee.Department.Name,
					Office: model.OfficeSimple{
						OfficeID:   result.Employee.Department.Office.ID,
						OfficeName: result.Employee.Department.Office.Name,
					},
				},
			},
			ContractTypeObj: *result.ContractType,
		}
		// If AttachedFile is present but not a full URL, try to generate presigned/fallback URL
		if res.AttachedFile != "" && !strings.HasPrefix(res.AttachedFile, "http") {
			if url, err := minIO.GeneratePresignedURL(c.Request.Context(), minIO.ContractBucket, res.AttachedFile, 15*time.Minute); err == nil {
				res.AttachedFile = url
			} else {
				log.Printf("Warning: failed to generate contract attached file URL for %s: %v", res.AttachedFile, err)
			}
		}

		utils.ResponseMessage(c, errors.MsgListData, http.StatusOK, []model.ContractResponse{res})
	}
}

func (h *ContractHandler) GetAllContract() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		filterFields := []string{"department_id", "contract_type_id", "condition"}
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

		// get data
		contract, totalRecords, err := h.ContractBiz.GetAllContract(ctx, page, pageSize, filters)
		if err != nil {
			utils.ResponseMessage(ctx, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			log.Printf("Lỗi lấy hợp đồng: %+v", err)
		}

		var contractResponses []model.ContractResponse
		for _, v := range contract {
			empSimple := model.EmployeeSimple{}

			if v.Employee != nil {
				empSimple.EmployeeID = v.Employee.EmployeeID
				empSimple.FullName = v.Employee.Fullname

				if v.Employee.Department != nil {
					empSimple.Department.DepartmentID = v.Employee.Department.ID
					empSimple.Department.DepartmentName = v.Employee.Department.Name

					if v.Employee.Department.Office != nil {
						empSimple.Department.Office.OfficeID = v.Employee.Department.Office.ID
						empSimple.Department.Office.OfficeName = v.Employee.Department.Office.Name
					}
				}
			} else {
				fmt.Println("Contract", v.ContractId, "không có employee")
			}

			res := model.ContractResponse{
				ContractID:    v.ContractId,
				EffectiveDate: v.EffectiveDate,
				ExpiredDate:   v.ExpiredDate,
				SignDate:      v.SignDate,
				Note:          v.Note,
				AttachedFile:  v.AttachedFile,
				Condition:     v.Condition,
				CreatedDate:   v.CreatedDate,
				ContractType:  v.ContractTypeId,
				ApproveStatus: v.ApproveStatus,
				Employee:      empSimple,
				Allowances:    v.Allowances,
			}

			contractResponses = append(contractResponses, res)
		}

		// total page
		totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)

		// response data
		utils.ResponseSuccess(ctx, errors.MsgListData, http.StatusOK, gin.H{
			"data":         contractResponses,
			"totalRecords": totalRecords,
			"page":         page,
			"pageSize":     pageSize,
			"totalPages":   totalPages,
		})
	}
}

func (h *ContractHandler) UpdateContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		var data model.ContractCreate
		// Support multipart/form-data for update as well
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			if err := c.ShouldBindJSON(&data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		} else {
			if err := c.ShouldBind(&data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			file, header, err := c.Request.FormFile("attached_file")
			if err == nil && file != nil {
				defer file.Close()
				ext := strings.ToLower(filepath.Ext(header.Filename))
				// Allow only .pdf and .docx for contract attachments
				allowed := map[string]bool{
					".pdf":  true,
					".docx": true,
				}
				if !allowed[ext] {
					utils.ResponseMessage(c, "Định dạng file không được hỗ trợ. Chỉ cho phép pdf, docx, jpg, jpeg, png", http.StatusBadRequest, nil)
					return
				}

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

				objectName := uuid.New().String() + ext
				if err := minIO.UploadImageToMinIO(c.Request.Context(), minIO.ContractBucket, objectName, file, header.Size, contentType, 30); err != nil {
					utils.ResponseMessage(c, fmt.Sprintf("Upload file failed: %v", err), http.StatusInternalServerError, nil)
					return
				}
				data.AttachedFile = objectName
			}
		}

		err := h.ContractBiz.UpdateContract(c.Request.Context(), idParam, &data)
		if err != nil {
			switch err {
			case errpkg.ErrApprovedContractCannotEdit,
				errpkg.ErrOnlyExpiredContractCanBeLiquidated:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống: " + err.Error()})
			}
			return
		}

		utils.ResponseMessage(c, errors.MsgUpdateSuccess, http.StatusOK, nil)
	}
}

func (h *ContractHandler) ReapproveContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu mã hợp đồng"})
			return
		}

		err := h.ContractBiz.UpdateApproveStatus(c.Request.Context(), id, "Chờ duyệt")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Hợp đồng đã gửi yêu cầu duyệt lại"})
	}
}

func (h *ContractHandler) DeleteContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		if err := h.ContractBiz.DeleteContract(c.Request.Context(), idParam); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi xóa hợp đồng: %s", err.Error()), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, errors.MsgDeleteSuccess, http.StatusOK, nil)
	}
}

func (biz *ContractHandler) ExportContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := strings.Split(c.Query("fields"), ",")
		data, filename, err := biz.ContractBiz.ExportContractTest(c, fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)

	}
}
