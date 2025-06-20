package handler

import (
	mock_test "erp/backend/api/handler/handler_test/hrm_test/mock_test_biz"
	handler "erp/backend/api/handler/hrm/hr_profile"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExportContract_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBiz := new(mock_test.MockContractBiz)
	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.GET("/contracts/export", h.ExportContract())

	// Mock dữ liệu giả
	expectedData := []byte("file content")
	expectedFilename := "contracts.xlsx"
	mockBiz.On("ExportContractTest", mock.Anything, mock.Anything).Return(expectedData, expectedFilename, nil)

	req := httptest.NewRequest(http.MethodGet, "/contracts/export", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expectedData, resp.Body.Bytes())
	assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", resp.Header().Get("Content-Type"))
	assert.Contains(t, resp.Header().Get("Content-Disposition"), expectedFilename)
	mockBiz.AssertExpectations(t)
}

func TestExportContract_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockBiz := new(mock_test.MockContractBiz)
	router := gin.New()
	h := &handler.ContractHandler{ContractBiz: mockBiz}
	router.GET("/contracts/export", h.ExportContract())

	mockBiz.On("ExportContractTest", mock.Anything, mock.Anything).Return(nil, "", fmt.Errorf("export failed"))

	req := httptest.NewRequest(http.MethodGet, "/contracts/export", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "export failed")
	mockBiz.AssertExpectations(t)
}
