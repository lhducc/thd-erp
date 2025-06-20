package handler

import (
	"bytes"
	mock_test "erp/backend/api/handler/handler_test/hrm_test/mock_test_biz"
	handler "erp/backend/api/handler/hrm/hr_profile"
	"erp/backend/internal/hrm/hr_profile/model"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestUpdateContract_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBiz := new(mock_test.MockContractBiz)
	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.PUT("/contracts/:id", h.UpdateContract())

	// Tạo fixed time để không bị lệch m=+xxx
	now := time.Date(2025, 6, 17, 15, 0, 0, 0, time.Local)

	contractData := model.ContractCreate{
		ContractId:     "HD001",
		EffectiveDate:  now,
		ExpiredDate:    now.AddDate(1, 0, 0),
		SignDate:       now,
		Note:           "Ghi chú",
		AttachedFile:   "file.pdf",
		Condition:      "Điều kiện",
		CreatedDate:    now,
		ContractTypeId: "CT001",
		ApproveStatus:  "Đã duyệt",
		Manager:        "EMP001",
	}

	// Dùng MatchedBy để match struct
	mockBiz.On("UpdateContract", mock.Anything, "HD001", mock.MatchedBy(func(input *model.ContractCreate) bool {
		return reflect.DeepEqual(*input, contractData)
	})).Return(nil)

	jsonValue, _ := json.Marshal(contractData)

	req := httptest.NewRequest(http.MethodPut, "/contracts/HD001", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Cập nhật thành công")
	mockBiz.AssertExpectations(t)
}

func TestUpdateContract_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBiz := new(mock_test.MockContractBiz)
	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.PUT("/contracts/:id", h.UpdateContract())

	invalidJSON := `{"contract_id":}` // JSON lỗi cú pháp

	req := httptest.NewRequest(http.MethodPut, "/contracts/HD001", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
}

func TestUpdateContract_BizError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBiz := new(mock_test.MockContractBiz)
	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.PUT("/contracts/:id", h.UpdateContract())

	contractData := model.ContractCreate{
		ContractId:    "HD001",
		Manager:       "EMP001",
		ApproveStatus: "Đang duyệt",
	}

	jsonValue, _ := json.Marshal(contractData)
	mockBiz.On("UpdateContract", mock.Anything, "HD001", &contractData).
		Return(errors.New("update failed"))

	req := httptest.NewRequest(http.MethodPut, "/contracts/HD001", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "update failed")
	mockBiz.AssertExpectations(t)
}
