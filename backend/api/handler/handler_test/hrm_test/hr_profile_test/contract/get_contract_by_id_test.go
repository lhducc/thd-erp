package handler

import (
	_ "context"
	mockBiz "erp/backend/api/handler/handler_test/hrm_test/mock_test_biz"
	handler "erp/backend/api/handler/hrm/hr_profile"
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetContract_Success(t *testing.T) {
	mockBiz := new(mockBiz.MockContractBiz)
	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.GET("/contracts/:id", h.GetContract())

	contract := &model.Contract{
		ContractId:     "HD000003",
		EffectiveDate:  time.Now(),
		ExpiredDate:    time.Now().AddDate(1, 0, 0),
		SignDate:       time.Now(),
		Note:           "Test",
		AttachedFile:   "file.pdf",
		Condition:      "Đang hiệu lực",
		CreatedDate:    time.Now(),
		ContractTypeId: "CT001",
		ApproveStatus:  "Đang duyệt",
		EmployeeID:     "EMP001",
		Employee: &model.Employee{
			EmployeeID: "EMP001",
			Fullname:   "Nguyễn Văn A",
			Department: &model.Department{
				ID:   "DEPT01",
				Name: "Phòng Kỹ thuật",
				Office: &model.Office{
					ID:   "OFF001",
					Name: "Trụ sở chính",
				},
			},
		},
	}

	mockBiz.On("GetContract", mock.Anything, "HD000003").Return(contract, nil)

	req := httptest.NewRequest(http.MethodGet, "/contracts/HD000003", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "HD000003")
	assert.Contains(t, resp.Body.String(), "EMP001")
	assert.Contains(t, resp.Body.String(), "Test")
	mockBiz.AssertExpectations(t)
}

func TestGetContract_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBiz := new(mockBiz.MockContractBiz)
	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.GET("/contracts/:id", h.GetContract())

	mockBiz.On("GetContract", mock.Anything, "INVALID_ID").
		Return(nil, errors.New("contract not found"))

	req := httptest.NewRequest(http.MethodGet, "/contracts/INVALID_ID", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "contract not found")
	mockBiz.AssertExpectations(t)
}
