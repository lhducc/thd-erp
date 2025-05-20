package officehandler

import (
	"bytes"
	"context"
	"encoding/json"
	officemodel "erp/backend/internal/hrm/office/model"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock OfficeBiz implementation
type mockOfficeBiz struct {
	mock.Mock
}

func (m *mockOfficeBiz) CreateOffice(ctx context.Context, data *officemodel.OfficeCreate) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *mockOfficeBiz) GetOffice(ctx context.Context, id string) (*officemodel.Office, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*officemodel.Office), args.Error(1)
}

func (m *mockOfficeBiz) GetAllOffice(ctx context.Context) ([]officemodel.Office, error) {
	args := m.Called(ctx)
	return args.Get(0).([]officemodel.Office), args.Error(1)
}

func (m *mockOfficeBiz) UpdateOffice(ctx context.Context, id string, data *officemodel.OfficeCreate) error {
	args := m.Called(ctx, id, data)
	return args.Error(0)
}

func (m *mockOfficeBiz) DeleteOffice(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper: create test handler and router
func setupRouter(mockBiz OficeBiz) *gin.Engine {
	r := gin.Default()
	h := &OfficeHandler{officeBiz: mockBiz}
	r.POST("/office", h.CreateOffice())
	r.GET("/office/:id", h.GetOffice())
	r.GET("/offices", h.GetAllOffice())
	r.PUT("/office/:id", h.UpdateOffice())
	r.DELETE("/office/:id", h.DeleteOffice())
	return r
}

func TestCreateOffice_Success(t *testing.T) {
	mockBiz := new(mockOfficeBiz)
	router := setupRouter(mockBiz)

	officeInput := officemodel.OfficeCreate{
		Name: "Head Office", Address: "123 Main St",
	}
	mockBiz.On("CreateOffice", mock.Anything, &officeInput).Return(nil)

	jsonValue, _ := json.Marshal(officeInput)
	req, _ := http.NewRequest("POST", "/office", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBiz.AssertExpectations(t)
}

func TestGetOffice_NotFound(t *testing.T) {
	mockBiz := new(mockOfficeBiz)
	router := setupRouter(mockBiz)

	mockBiz.On("GetOffice", mock.Anything, "1").Return(nil, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/office/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockBiz.AssertExpectations(t)
}

func TestGetAllOffice_Success(t *testing.T) {
	mockBiz := new(mockOfficeBiz)
	router := setupRouter(mockBiz)

	mockResult := []officemodel.Office{
		{ID: "1", Name: "Branch A"},
		{ID: "2", Name: "Branch B"},
	}
	mockBiz.On("GetAllOffice", mock.Anything).Return(mockResult, nil)

	req, _ := http.NewRequest("GET", "/offices", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBiz.AssertExpectations(t)
}

func TestUpdateOffice_Success(t *testing.T) {
	mockBiz := new(mockOfficeBiz)
	router := setupRouter(mockBiz)

	data := officemodel.OfficeCreate{Name: "Updated Office", Address: "New Address"}
	mockBiz.On("UpdateOffice", mock.Anything, "1", &data).Return(nil)

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("PUT", "/office/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBiz.AssertExpectations(t)
}

func TestDeleteOffice_Success(t *testing.T) {
	mockBiz := new(mockOfficeBiz)
	router := setupRouter(mockBiz)

	mockBiz.On("DeleteOffice", mock.Anything, "1").Return(nil)

	req, _ := http.NewRequest("DELETE", "/office/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockBiz.AssertExpectations(t)
}
