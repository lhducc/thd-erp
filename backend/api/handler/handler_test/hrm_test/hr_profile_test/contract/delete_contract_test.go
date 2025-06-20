package handler

import (
	mock_test "erp/backend/api/handler/handler_test/hrm_test/mock_test_biz"
	handler "erp/backend/api/handler/hrm/hr_profile"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteContract_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBiz := new(mock_test.MockContractBiz)
	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.DELETE("/contracts/:id", h.DeleteContract())

	mockBiz.On("DeleteContract", mock.Anything, "HD001").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/contracts/HD001", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var body struct {
		Message    string `json:"message"`
		StatusCode int    `json:"statuscode"`
	}
	_ = json.Unmarshal(resp.Body.Bytes(), &body)
	assert.Equal(t, http.StatusOK, body.StatusCode)
	assert.Contains(t, body.Message, "Xóa dữ liệu thành công")

	mockBiz.AssertExpectations(t)
}

func TestDeleteContract_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockBiz := new(mock_test.MockContractBiz)
	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.DELETE("/contracts/:id", h.DeleteContract())

	mockBiz.On("DeleteContract", mock.Anything, "HD999").Return(fmt.Errorf("lỗi xóa db"))

	req := httptest.NewRequest(http.MethodDelete, "/contracts/HD999", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var body map[string]interface{}
	_ = json.Unmarshal(resp.Body.Bytes(), &body)
	assert.Contains(t, body["message"], "Lỗi xóa hợp đồng")
	mockBiz.AssertExpectations(t)
}
