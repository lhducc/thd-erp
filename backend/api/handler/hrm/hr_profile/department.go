package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type DepartmentBiz interface {
	CreateDepartment(context context.Context, data *model.DepartmentCreate) error
	GetDepartment(ctx context.Context, id string) (*model.Department, error)
	UpdateDepartment(ctx context.Context, id string, data *model.DepartmentCreate) error
	DeleteDepartment(ctx context.Context, id string) error
	GetAllDepartment(ctx context.Context) ([]model.Department, error)
	GetDepartmentByOfficeID(officeID string) ([]model.Department, error)
}

type DepartmentHandler struct {
	departmentBiz DepartmentBiz
}

func NewDepartmentHandler(biz DepartmentBiz) *DepartmentHandler {
	return &DepartmentHandler{
		departmentBiz: biz,
	}
}

func (handler *DepartmentHandler) CreateDepartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.DepartmentCreate

		if err := c.ShouldBind(&data); err != nil {
			utils.ResponseMessage(c, "Không thể đọc dữ liệu", http.StatusBadRequest, nil)
			return
		}

		if err := handler.departmentBiz.CreateDepartment(c.Request.Context(), &data); err != nil {
			utils.ResponseMessage(c, "Tạo bộ phận thất bại", http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Thêm mới bộ phận thành công", http.StatusOK, nil)
	}
}

func (handler *DepartmentHandler) GetDepartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		result, err := handler.departmentBiz.GetDepartment(c.Request.Context(), idParam)
		if err != nil {
			utils.ResponseMessage(c, "Không tìm thấy dữ liệu văn phòng", http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách dữ liệu", http.StatusOK, result)
	}
}
func (handler *DepartmentHandler) UpdateDepartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		var data model.DepartmentCreate
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Cập nhật thất bại", http.StatusBadRequest, nil)
			return
		}

		if err := handler.departmentBiz.UpdateDepartment(c.Request.Context(), idParam, &data); err != nil {
			utils.ResponseMessage(c, "Cập nhật thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật bộ phận thành công", http.StatusOK, nil)
	}
}

func (h *DepartmentHandler) GetAllDepartment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := h.departmentBiz.GetAllDepartment(ctx)
		if err != nil {
			utils.ResponseMessage(ctx, "Lấy thông tin thất bại", http.StatusNotFound, nil)
		}
		utils.ResponseMessage(ctx, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (handler *DepartmentHandler) DeleteDepartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		if err := handler.departmentBiz.DeleteDepartment(c.Request.Context(), idParam); err != nil {
			utils.ResponseMessage(c, "Xóa bộ phận thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa bộ phận thành công", http.StatusOK, nil)
	}
}

func (handler *DepartmentHandler) GetDepartmentByOfficeID() gin.HandlerFunc {
	return func(c *gin.Context) {
		officeID := c.Param("office-id")
		if strings.TrimSpace(officeID) == "" {
			utils.ResponseMessage(c, "Mã văn phòng không hợp lệ", http.StatusInternalServerError, nil)
			return
		}
		result, err := handler.departmentBiz.GetDepartmentByOfficeID(officeID)
		if err != nil {
			utils.ResponseMessage(c, "Lấy dữ liệu thất bại", http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Lấy Danh sách dữ liệu thành công", http.StatusOK, result)

	}
}
