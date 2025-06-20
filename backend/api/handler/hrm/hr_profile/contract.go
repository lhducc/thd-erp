package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	"erp/backend/internal/hrm/hr_profile/usecase"
	utils "erp/backend/pkg"
	"erp/backend/pkg/errors"
	"log"

	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContractBiz interface {
	CreateContract(context context.Context, data *model.ContractCreate) error
	GetContract(ctx context.Context, id string) (*model.Contract, error)
	GetAllContract(ctx context.Context) ([]model.Contract, error)
	UpdateContract(ctx context.Context, id string, data *model.ContractCreate) error
	DeleteContract(ctx context.Context, id string) error
	ExportContractTest(c context.Context, selectedFields []string) ([]byte, string, error)
	UpdateApproveStatus(ctx context.Context, id string, status string) error
}

type ContractHandler struct {
	ContractBiz ContractBiz
	employeeBiz usecase.EmployeeRepo
}

func NewContractHandler(db *gorm.DB) *ContractHandler {
	storeJob := repository.NewContractStore(db)
	storeUser := repository.NewUserStore(db)
	biz := usecase.NewContractBiz(storeJob, storeUser)

	return &ContractHandler{
		ContractBiz: biz,
		employeeBiz: storeUser,
	}
}

func (h *ContractHandler) CreateContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.ContractCreate

		err := c.ShouldBindJSON(&data)
		if err != nil {
			utils.ResponseError(c, "Dữ liệu đầu vào không hợp lệ", err, http.StatusBadRequest)
			return
		}

		if err := h.ContractBiz.CreateContract(c.Request.Context(), &data); err != nil {
			utils.ResponseError(c, "Không thể tạo hợp đồng", err, http.StatusInternalServerError)
			return
		}

		utils.ResponseMessage(c, "Tạo hợp đồng thành công", http.StatusOK, gin.H{"data": data})
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
		}
		utils.ResponseMessage(c, errors.MsgListData, http.StatusOK, []model.ContractResponse{res})
	}
}

func (h *ContractHandler) GetAllContract() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := h.ContractBiz.GetAllContract(ctx)
		if err != nil {
			utils.ResponseMessage(ctx, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusBadRequest, nil)
			log.Printf("Lỗi lấy hợp đồng: %+v", err)
		}
		var responseList []model.ContractResponse
		for _, v := range result {
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
			}

			responseList = append(responseList, res)
		}
		utils.ResponseMessage(ctx, errors.MsgListData, http.StatusOK, responseList)
	}
}

func (h *ContractHandler) UpdateContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		var data model.ContractCreate
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := h.ContractBiz.UpdateContract(c.Request.Context(), idParam, &data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

		err := h.contractBiz.UpdateApproveStatus(c.Request.Context(), id, "Chờ duyệt")
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
