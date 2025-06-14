package officehandler

import (
	"context"
	officemodel "erp/backend/internal/hrm/office/model"
	hrmrepository "erp/backend/internal/hrm/office/repository"
	hrmbiz "erp/backend/internal/hrm/office/usecase"
	utils "erp/backend/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OficeBiz interface {
	CreateOffice(context context.Context, data *officemodel.OfficeCreate) error
	GetOffice(ctx context.Context, id string) (*officemodel.Office, error)
	GetAllOffice(ctx context.Context) ([]officemodel.Office, error)
	UpdateOffice(ctx context.Context, id string, data *officemodel.OfficeCreate) error
	DeleteOffice(ctx context.Context, id string) error
}

type OfficeHandler struct {
	officeBiz OficeBiz
}

func NewOficeHandler(db *gorm.DB) *OfficeHandler {
	storeJob := hrmrepository.NewOfficeStore(db)
	biz := hrmbiz.NewOfficeBiz(storeJob)

	return &OfficeHandler{
		officeBiz: biz,
	}
}

func (h *OfficeHandler) CreateOffice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data officemodel.OfficeCreate

		if err := c.ShouldBind(&data); err != nil {
			utils.ResponseMessage(c, "Xóa bộ phận thất bại ", http.StatusBadRequest, nil)

			return
		}

		if err := h.officeBiz.CreateOffice(c.Request.Context(), &data); err != nil {
			utils.ResponseMessage(c, "Tạo văn phòng thất bại", http.StatusInternalServerError, nil)
			return
		}
		utils.ResponseMessage(c, "Tạo văn phòng thành công", http.StatusCreated, nil)
	}
}

// GetOffice xử lý request lấy Office theo ID
func (h *OfficeHandler) GetOffice() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		result, err := h.officeBiz.GetOffice(c.Request.Context(), idParam)
		if err != nil {
			utils.ResponseMessage(c, "Lấy thông tin thất bại", http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": result})
	}
}

func (h *OfficeHandler) GetAllOffice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := h.officeBiz.GetAllOffice(ctx)
		if err != nil {
			utils.ResponseMessage(ctx, "Lấy thông tin thất bại", http.StatusNotFound, nil)
		}
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	}
}

func (h *OfficeHandler) UpdateOffice() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		var data officemodel.OfficeCreate
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, "Cập nhật văn phòng thất bại", http.StatusBadRequest, nil)
			return
		}

		if err := h.officeBiz.UpdateOffice(c.Request.Context(), idParam, &data); err != nil {
			utils.ResponseMessage(c, "Cập nhật văn phòng thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật văn phòng thành công", http.StatusOK, nil)
	}
}

// DeleteOffice xử lý request xóa Office
func (h *OfficeHandler) DeleteOffice() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		if err := h.officeBiz.DeleteOffice(c.Request.Context(), idParam); err != nil {
			utils.ResponseMessage(c, "Xóa văn phòng thất bại", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa văn phòng thành công", http.StatusOK, nil)
	}
}
