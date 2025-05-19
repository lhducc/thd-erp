package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	"erp/backend/internal/hrm/hr_profile/usecase"
	utils "erp/backend/pkg"
	"erp/backend/pkg/errors"

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
}

type ContractHandler struct {
	contractBiz ContractBiz
	employeeBiz usecase.EmployeeRepo
}

func NewContractHandler(db *gorm.DB) *ContractHandler {
	storeJob := repository.NewContractStore(db)
	storeUser := repository.NewUserStore(db)
	biz := usecase.NewContractBiz(storeJob, storeUser)

	return &ContractHandler{
		contractBiz: biz,
		employeeBiz: storeUser,
	}
}

func (h *ContractHandler) CreateContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.ContractCreate

		if err := c.ShouldBind(&data); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusOK, nil)
			return
		}

		if err := h.contractBiz.CreateContract(c.Request.Context(), &data); err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusOK, nil)
			return
		}

		utils.ResponseMessage(c, errors.MsgCreatedSuccess, http.StatusOK, nil)
	}
}

func (h *ContractHandler) GetContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		result, err := h.contractBiz.GetContract(c.Request.Context(), idParam)
		if err != nil {
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusOK, nil)
			return
		}

		utils.ResponseMessage(c, errors.MsgListData, http.StatusOK, result)
	}
}

func (h *ContractHandler) GetAllContract() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := h.contractBiz.GetAllContract(ctx)
		if err != nil {
			utils.ResponseMessage(ctx, fmt.Sprintf("Lỗi: %s", err.Error()), http.StatusOK, nil)
		}
		utils.ResponseMessage(ctx, errors.MsgListData, http.StatusOK, result)
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

		if err := h.contractBiz.UpdateContract(c.Request.Context(), idParam, &data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		utils.ResponseMessage(c, errors.MsgUpdateSuccess, http.StatusOK, nil)
	}
}

func (h *ContractHandler) DeleteContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		if err := h.contractBiz.DeleteContract(c.Request.Context(), idParam); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			utils.ResponseMessage(c, fmt.Sprintf("Lỗi xóa hợp đồng: %s", err.Error()), http.StatusOK, nil)
			return
		}

		utils.ResponseMessage(c, errors.MsgDeleteSuccess, http.StatusOK, nil)
	}
}

func (biz *ContractHandler) ExportContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := strings.Split(c.Query("fields"), ",")
		data, filename, err := biz.contractBiz.ExportContractTest(c, fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)

	}
}
