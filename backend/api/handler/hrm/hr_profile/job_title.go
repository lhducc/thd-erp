package handler

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	"erp/backend/internal/hrm/hr_profile/repository"
	"erp/backend/internal/hrm/hr_profile/usecase"
	utils "erp/backend/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobTitleBiz interface {
	CreateJobTitle(ctx context.Context, data *model.JobTitleCreate) error
	GetJobTitleById(ctx context.Context, id string) (*model.JobTitle, error)
	UpdateJobTitleById(ctx context.Context, id string, data *model.JobTitleCreate) error
	DeleteJobTitleById(ctx context.Context, id string) error
	GetAllJobTitle(ctx context.Context) ([]model.JobTitle, error)
}

type JobTitleHandler struct {
	jobTitleBiz JobTitleBiz
}

func NewJobTitleHandler(db *gorm.DB) *JobTitleHandler {
	storeJob := repository.NewJobTitleStore(db)
	storeHie := repository.NewhierarchyLevelStore(db)
	biz := usecase.NewJobTitleBiz(storeJob, storeHie)

	return &JobTitleHandler{
		jobTitleBiz: biz,
	}
}

func (h *JobTitleHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data model.JobTitleCreate

		if err := c.ShouldBind(&data); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := h.jobTitleBiz.CreateJobTitle(c.Request.Context(), &data); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		utils.ResponseMessage(c, "Tao chuc vu thanh cong", http.StatusCreated, gin.H{
			"data": data,
		})
	}
}

func (h *JobTitleHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseID(c)
		if err != nil {
			return
		}

		result, err := h.jobTitleBiz.GetJobTitleById(c.Request.Context(), string(id))
		if err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusNotFound, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (h *JobTitleHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.jobTitleBiz.GetAllJobTitle(c)
		if err != nil {
			utils.ResponseMessage(c, "Không thể lấy danh sách vị trí", http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Danh sách dữ liệu", http.StatusOK, result)
	}
}

func (h *JobTitleHandler) UpdateByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseID(c)
		if err != nil {
			return
		}

		var data model.JobTitleCreate
		if err := c.ShouldBindJSON(&data); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusBadRequest, nil)
			return
		}

		if err := h.jobTitleBiz.UpdateJobTitleById(c.Request.Context(), string(id), &data); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Cập nhật thành công", http.StatusOK, nil)
	}
}

func (h *JobTitleHandler) DeleteByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := parseID(c)
		if err != nil {
			return
		}

		if err := h.jobTitleBiz.DeleteJobTitleById(c.Request.Context(), id); err != nil {
			utils.ResponseMessage(c, err.Error(), http.StatusInternalServerError, nil)
			return
		}

		utils.ResponseMessage(c, "Xóa thành công", http.StatusOK, nil)
	}
}

func parseID(c *gin.Context) (string, error) {
	idParam := c.Param("id")

	return idParam, nil
}
