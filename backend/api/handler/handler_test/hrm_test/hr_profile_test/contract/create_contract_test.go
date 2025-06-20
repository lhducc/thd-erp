package handler

import (
	"bytes"
	_ "context"
	mockBiz "erp/backend/api/handler/handler_test/hrm_test/mock_test_biz"
	handler "erp/backend/api/handler/hrm/hr_profile"
	"erp/backend/internal/hrm/hr_profile/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateContract_Success(t *testing.T) {
	mockBiz := new(mockBiz.MockContractBiz)

	payload := `{
		"contract_id": "HD001",
		"effective_date": "2025-06-17T00:00:00Z",
		"expired_date": "2026-06-17T00:00:00Z",
		"sign_date": "2025-06-16T00:00:00Z",
		"note": "Hợp đồng thử việc",
		"attached_file": "file.pdf",
		"condition": "Làm việc toàn thời gian",
		"created_date": "2025-06-15T00:00:00Z",
		"contract_type": "CT001",
		"approve_status": "Đang duyệt",
		"employee_id": "EMP001"
	}`

	expected := &model.ContractCreate{
		ContractId:     "HD001",
		EffectiveDate:  time.Date(2025, 6, 17, 0, 0, 0, 0, time.UTC),
		ExpiredDate:    time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC),
		SignDate:       time.Date(2025, 6, 16, 0, 0, 0, 0, time.UTC),
		Note:           "Hợp đồng thử việc",
		AttachedFile:   "file.pdf",
		Condition:      "Làm việc toàn thời gian",
		CreatedDate:    time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		ContractTypeId: "CT001",
		ApproveStatus:  "Đang duyệt",
		Manager:        "EMP001",
	}

	mockBiz.On("CreateContract", mock.Anything, mock.MatchedBy(func(input *model.ContractCreate) bool {
		return input.ContractId == expected.ContractId &&
			input.Manager == expected.Manager &&
			input.Note == expected.Note
	})).Return(nil)

	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.POST("/contracts", h.CreateContract())

	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Tạo hợp đồng thành công")
}

func TestCreateContract_InvalidJSON(t *testing.T) {
	mockBiz := new(mockBiz.MockContractBiz)

	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.POST("/contracts", h.CreateContract())

	invalidJSON := `{"contract_id": "HD001", "effective_date": "not-a-date"}` // sai format ngày

	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Dữ liệu đầu vào không hợp lệ")
	mockBiz.AssertNotCalled(t, "CreateContract")
}

func TestCreateContract_BizError(t *testing.T) {
	mockBiz := new(mockBiz.MockContractBiz)

	payload := `{
		"contract_id": "HD002",
		"effective_date": "2025-06-17T00:00:00Z",
		"expired_date": "2026-06-17T00:00:00Z",
		"sign_date": "2025-06-16T00:00:00Z",
		"note": "Lỗi nghiệp vụ",
		"attached_file": "file.pdf",
		"condition": "test",
		"created_date": "2025-06-15T00:00:00Z",
		"contract_type": "CT002",
		"approve_status": "Đang duyệt",
		"employee_id": "EMP002"
	}`

	mockBiz.On("CreateContract", mock.Anything, mock.Anything).Return(fmt.Errorf("lỗi nghiệp vụ"))

	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.POST("/contracts", h.CreateContract())

	req := httptest.NewRequest(http.MethodPost, "/contracts", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "Không thể tạo hợp đồng")
	mockBiz.AssertExpectations(t)
}
