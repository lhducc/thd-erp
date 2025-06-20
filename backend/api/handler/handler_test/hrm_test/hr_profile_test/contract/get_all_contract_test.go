package handler

import (
	mockBizz "erp/backend/api/handler/handler_test/hrm_test/mock_test_biz"
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

func TestGetAllContract_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBiz := new(mockBizz.MockContractBiz)
	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.GET("/contracts", h.GetAllContract())

	contracts := []model.Contract{
		{
			ContractId:     "HD001",
			EffectiveDate:  time.Now(),
			ExpiredDate:    time.Now().AddDate(1, 0, 0),
			SignDate:       time.Now(),
			Note:           "Note 1",
			AttachedFile:   "file1.pdf",
			Condition:      "Điều kiện 1",
			CreatedDate:    time.Now(),
			ContractTypeId: "CT001",
			ApproveStatus:  "Đã duyệt",
			EmployeeID:     "EMP001",
			Employee: &model.Employee{
				EmployeeID: "EMP001",
				Fullname:   "Nguyễn Văn A",
				Department: &model.Department{
					ID:   "D001",
					Name: "Phòng A",
					Office: &model.Office{
						ID:   "O001",
						Name: "Trụ sở",
					},
				},
			},
		},
	}

	mockBiz.On("GetAllContract", mock.Anything).Return(contracts, nil)

	req := httptest.NewRequest(http.MethodGet, "/contracts", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Nguyễn Văn A")
	assert.Contains(t, resp.Body.String(), "HD001")
	mockBiz.AssertExpectations(t)
}

func TestGetAllContract_Fail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBiz := new(mockBizz.MockContractBiz)
	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.GET("/contracts", h.GetAllContract())

	mockBiz.On("GetAllContract", mock.Anything).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest(http.MethodGet, "/contracts", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "database error")
	mockBiz.AssertExpectations(t)
}
