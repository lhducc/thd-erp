package checkin

//
//import (
//	"context"
//	"encoding/json"
//	handler "erp/backend/api/handler/hrm/checkin"
//	"erp/backend/internal/hrm/checkin/model/dto"
//	"errors"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//)
//
//// Mock service
//type MockAttendanceCategoryService struct {
//	mock.Mock
//}
//
//func (m *MockAttendanceCategoryService) CreateAttendanceCategory(ctx context.Context, category *dto.AttendanceCategoryRequest) error {
//	args := m.Called(ctx, category)
//	return args.Error(0)
//}
//
//func (m *MockAttendanceCategoryService) GetAttendanceCategoryByID(ctx context.Context, id string) (*dto.AttendanceCategoryRequest, error) {
//	args := m.Called(ctx, id)
//	return args.Get(0).(*dto.AttendanceCategoryRequest), args.Error(1)
//}
//
//func TestCreateAttendanceCategory_Success(t *testing.T) {
//	// Setup
//	mockService := new(MockAttendanceCategoryService)
//	handler := handler.NewAttendanceCategoryHandler(mockService)
//
//	// Mock data
//	reqBody := `{
//		"name": "Test Category",
//		"code": "TEST001",
//		"description": "Test description"
//	}`
//
//	expectedCategory := dto.AttendanceCategoryRequest{
//		AttendanceCategoryName: "Test Category",
//		Code:                   "TEST001",
//		Description:            "Test description",
//		CreatedBy:              "THD001",
//	}
//
//	// Mock expectations
//	mockService.On("CreateAttendanceCategory", mock.Anything, mock.AnythingOfType("*dto.AttendanceCategory")).
//		Return(nil).
//		Run(func(args mock.Arguments) {
//			category := args.Get(1).(*dto.AttendanceCategory)
//			assert.Equal(t, expectedCategory.Name, category.Name)
//			assert.Equal(t, expectedCategory.Code, category.Code)
//		})
//
//	// Create test context
//	w := httptest.NewRecorder()
//	c, _ := gin.CreateTestContext(w)
//	c.Request = httptest.NewRequest("POST", "/categories", strings.NewReader(reqBody))
//	c.Request.Header.Set("Content-Type", "application/json")
//
//	// Mock context values
//	c.Set("employeeId", "emp123")
//
//	// Execute
//	handlerFunc := handler.CreateAttendanceCategory()
//	handlerFunc(c)
//
//	// Assertions
//	assert.Equal(t, http.StatusOK, w.Code)
//	mockService.AssertExpectations(t)
//
//	var response pkg.Response
//	json.Unmarshal(w.Body.Bytes(), &response)
//	assert.Equal(t, "Attendance category created successfully", response.Message)
//}
//
//func TestCreateAttendanceCategory_InvalidJSON(t *testing.T) {
//	mockService := new(MockAttendanceCategoryService)
//	handler := NewAttendanceCategoryHandler(mockService)
//
//	w := httptest.NewRecorder()
//	c, _ := gin.CreateTestContext(w)
//	c.Request = httptest.NewRequest("POST", "/categories", strings.NewReader(`invalid json`))
//	c.Request.Header.Set("Content-Type", "application/json")
//
//	handlerFunc := handler.CreateAttendanceCategory()
//	handlerFunc(c)
//
//	assert.Equal(t, http.StatusBadRequest, w.Code)
//}
//
//func TestCreateAttendanceCategory_ValidationError(t *testing.T) {
//	mockService := new(MockAttendanceCategoryService)
//	handler := NewAttendanceCategoryHandler(mockService)
//
//	reqBody := `{"name": ""}` // Invalid - empty name
//
//	w := httptest.NewRecorder()
//	c, _ := gin.CreateTestContext(w)
//	c.Request = httptest.NewRequest("POST", "/categories", strings.NewReader(reqBody))
//	c.Request.Header.Set("Content-Type", "application/json")
//
//	handlerFunc := handler.CreateAttendanceCategory()
//	handlerFunc(c)
//
//	assert.Equal(t, http.StatusBadRequest, w.Code)
//}
//
//func TestCreateAttendanceCategory_ServiceError(t *testing.T) {
//	mockService := new(MockAttendanceCategoryService)
//	handler := NewAttendanceCategoryHandler(mockService)
//
//	reqBody := `{
//		"name": "Test",
//		"code": "TEST001"
//	}`
//
//	mockService.On("CreateAttendanceCategory", mock.Anything, mock.Anything).
//		Return(errors.New("database error"))
//
//	w := httptest.NewRecorder()
//	c, _ := gin.CreateTestContext(w)
//	c.Request = httptest.NewRequest("POST", "/categories", strings.NewReader(reqBody))
//	c.Request.Header.Set("Content-Type", "application/json")
//	c.Set("employeeId", "emp123")
//
//	handlerFunc := handler.CreateAttendanceCategory()
//	handlerFunc(c)
//
//	assert.Equal(t, http.StatusInternalServerError, w.Code)
//}
//
//func TestGetAttendanceCategory_Success(t *testing.T) {
//	mockService := new(MockAttendanceCategoryService)
//	handler := NewAttendanceCategoryHandler(mockService)
//
//	expectedCategory := &dto.AttendanceCategory{
//		AttendanceCategoryID: "cat123",
//		Name:                 "Test Category",
//		Code:                 "TEST001",
//	}
//
//	mockService.On("GetAttendanceCategoryByID", mock.Anything, "cat123").
//		Return(expectedCategory, nil)
//
//	w := httptest.NewRecorder()
//	c, _ := gin.CreateTestContext(w)
//	c.Request = httptest.NewRequest("GET", "/categories/cat123", nil)
//	c.Params = []gin.Param{{Key: "id", Value: "cat123"}}
//
//	handlerFunc := handler.GetAttendanceCategory()
//	handlerFunc(c)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var response pkg.Response
//	json.Unmarshal(w.Body.Bytes(), &response)
//	assert.Equal(t, expectedCategory, response.Data)
//}
//
//func TestGetAttendanceCategory_NotFound(t *testing.T) {
//	mockService := new(MockAttendanceCategoryService)
//	handler := NewAttendanceCategoryHandler(mockService)
//
//	mockService.On("GetAttendanceCategoryByID", mock.Anything, "nonexistent").
//		Return(&dto.AttendanceCategory{}, errors.New("not found"))
//
//	w := httptest.NewRecorder()
//	c, _ := gin.CreateTestContext(w)
//	c.Request = httptest.NewRequest("GET", "/categories/nonexistent", nil)
//	c.Params = []gin.Param{{Key: "id", Value: "nonexistent"}}
//
//	handlerFunc := handler.GetAttendanceCategory()
//	handlerFunc(c)
//
//	assert.Equal(t, http.StatusNotFound, w.Code)
//}
//
//// Test helper function
//func createTestContext(method, url, body string) (*gin.Context, *httptest.ResponseRecorder) {
//	w := httptest.NewRecorder()
//	c, _ := gin.CreateTestContext(w)
//	c.Request = httptest.NewRequest(method, url, strings.NewReader(body))
//	c.Request.Header.Set("Content-Type", "application/json")
//	return c, w
//}
//
//// Test main
//func TestMain(m *testing.M) {
//	gin.SetMode(gin.TestMode)
//	m.Run()
//}
