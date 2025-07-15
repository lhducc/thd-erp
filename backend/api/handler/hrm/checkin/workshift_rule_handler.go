package checkin

import (
	"erp/backend/internal/hrm/checkin/model"
	"erp/backend/internal/hrm/checkin/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WorkshiftRuleHandler struct {
	service *service.WorkshiftRuleService
}

func NewWorkshiftRuleHandler(service *service.WorkshiftRuleService) *WorkshiftRuleHandler {
	return &WorkshiftRuleHandler{
		service: service,
	}
}

// POST /workshift-rules
func (h *WorkshiftRuleHandler) Create(c *gin.Context) {
	var rule model.WorkshiftRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(rule)

	if err := h.service.Create(rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

// GET /workshift-rules
func (h *WorkshiftRuleHandler) GetAll(c *gin.Context) {
	rules, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

// GET /workshift-rules/user/:id
func (h *WorkshiftRuleHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("employeeID")

	rules, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rules)
}

// PUT /workshift-rules/:id
func (h *WorkshiftRuleHandler) Update(c *gin.Context) {
	var rule model.WorkshiftRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Update(id, rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// DELETE /workshift-rules/:id
func (h *WorkshiftRuleHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
