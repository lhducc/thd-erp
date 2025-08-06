package employee_test_test

import (
	"bytes"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg" // Đảm bảo đường dẫn này đúng tới gói `pkg` chứa ExcelExporter và ExcelReader của bạn
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// MockEmployeeRepo là một triển khai giả lập của EmployeeRepository interface
type MockEmployeeRepo struct {
	mock.Mock
}

func (m *MockEmployeeRepo) GetAllEmployees() ([]model.Employee, error) {
	args := m.Called()
	return args.Get(0).([]model.Employee), args.Error(1)
}

type employeeBiz struct {
	repo *MockEmployeeRepo
}

func (e *employeeBiz) ExportEmployeeTest(selectedFields []string) ([]byte, string, error) {
	exporter := utils.NewExcelExporter("Employees")

	// Đăng ký tất cả các trường bạn muốn có khả năng export
	exporter.RegisterField("employee_id", "Mã Nhân Viên", "employee_id", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.EmployeeID
	})

	exporter.RegisterField("fullname", "Họ và tên", "fullname", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.Fullname
	})

	exporter.RegisterField("birthday", "Ngày sinh", "birthday", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.Birthday
	})

	exporter.RegisterField("gender", "Giới tính", "gender", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.Gender
	})

	exporter.RegisterField("phone", "Số điện thoại", "phone", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.PhoneNumber
	})

	exporter.RegisterField("email", "Email", "email", func(item interface{}) any {
		emp := item.(*model.Employee)
		if emp.Account != nil {
			return emp.Account.LoginMail
		}
		return ""
	})

	exporter.RegisterField("manager", "Quản lý", "manager", func(item interface{}) any {
		emp := item.(*model.Employee)
		if emp.Manager != nil {
			return emp.Manager.Fullname
		}
		return ""
	})

	exporter.RegisterField("created_date", "Ngày tạo", "created_date", func(item interface{}) any {
		emp := item.(*model.Employee)
		return emp.CreatedDate.Format("2006-01-02")
	})

	// Lấy dữ liệu nhân viên từ mock repository
	employees, err := e.repo.GetAllEmployees()
	if err != nil {
		return nil, "", fmt.Errorf("không thể lấy dữ liệu nhân viên: %w", err)
	}

	// Nếu không có nhân viên nào, trả về file rỗng
	if len(employees) == 0 {
		return exporter.Export([]interface{}{}, selectedFields)
	}

	// Chuyển đổi slice of model.Employee sang slice of interface{}
	itemsToExport := make([]interface{}, len(employees))
	for i := range employees {
		itemsToExport[i] = &employees[i]
	}

	// Nếu selectedFields rỗng, sử dụng tất cả các trường đã đăng ký
	if len(selectedFields) == 0 {
		selectedFields = []string{"employee_id", "fullname", "birthday", "gender", "phone", "email", "manager", "created_date"}
	}

	// Gọi hàm ExportTimeSheet từ gói utils
	return exporter.Export(itemsToExport, selectedFields)
}

// --- Phần Test Cases ---
func TestExportEmployeeTest(t *testing.T) {
	// Dữ liệu mẫu để test
	mockEmployees := []model.Employee{
		{
			EmployeeID:  "EMP001",
			Fullname:    "Nguyen Van A",
			Birthday:    "1990-01-15",
			Gender:      "Male",
			PhoneNumber: "0912345678",
			Account: &model.Account{
				LoginMail: "nguyenvana@example.com",
			},
			Manager: &model.ManagerResponse{
				Fullname: "Tran Thi B (Manager)",
			},
			CreatedDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			EmployeeID:  "EMP002",
			Fullname:    "Le Thi C",
			Birthday:    "1992-03-20",
			Gender:      "Female",
			PhoneNumber: "0987654321",
			Account:     nil, // Test case cho nil account
			Manager:     nil, // Test case cho nil manager
			CreatedDate: time.Date(2023, 2, 10, 0, 0, 0, 0, time.UTC),
		},
	}

	tests := []struct {
		name              string
		selectedFields    []string
		mockRepoSetup     func(*MockEmployeeRepo)
		expectedError     error
		expectedFileName  string
		expectedHeaderRow []string
		expectedDataRows  [][]any
	}{
		{
			name:           "ExportTimeSheet thành công - Các trường chọn lọc (Họ và tên, Email)",
			selectedFields: []string{"fullname", "email"},
			mockRepoSetup: func(m *MockEmployeeRepo) {
				m.On("GetAllEmployees").Return(mockEmployees, nil)
			},
			expectedError:    nil,
			expectedFileName: "Employees",
			expectedHeaderRow: []string{
				"Họ và tên", "Email",
			},
			expectedDataRows: [][]any{
				{"Nguyen Van A", "nguyenvana@example.com"},
				{"Le Thi C", ""},
			},
		},
		// (các trường hợp khác giữ nguyên)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockEmployeeRepo)
			tt.mockRepoSetup(mockRepo)

			biz := &employeeBiz{repo: mockRepo}

			// Gọi hàm cần test
			excelBytes, fileName, err := biz.ExportEmployeeTest(tt.selectedFields)

			// Kiểm tra lỗi
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, excelBytes)
				assert.Empty(t, fileName)
			} else {
				assert.NoError(t, err)

				// Kiểm tra tên file
				assert.Contains(t, fileName, tt.expectedFileName)
				assert.Contains(t, fileName, ".xlsx")

				assert.NotNil(t, excelBytes)

				// Sử dụng ExcelReader để đọc file Excel đã tạo và xác minh nội dung
				file, err := utils.NewExcelReader(bytes.NewReader(excelBytes))
				assert.NoError(t, err)

				// Lấy tên sheet đang hoạt động
				sheetName := file.GetSheetName(0)
				assert.NotEmpty(t, sheetName)

				// Đọc tất cả các hàng từ sheet
				rows, err := file.GetRows(sheetName)
				assert.NoError(t, err)

				// Xác minh hàng tiêu đề
				assert.NotNil(t, rows)
				if len(tt.expectedHeaderRow) > 0 {
					assert.True(t, len(rows) > 0, "Mong đợi ít nhất một hàng cho tiêu đề")
					assert.Equal(t, tt.expectedHeaderRow, rows[0])
				} else {
					if len(rows) > 0 {
						assert.NotEqual(t, tt.expectedHeaderRow, rows[0])
					}
				}

				// Xác minh các hàng dữ liệu
				if len(tt.expectedDataRows) > 0 {
					assert.Equal(t, len(tt.expectedDataRows)+1, len(rows), "Không khớp số lượng hàng (bao gồm tiêu đề)")
					for i, expectedRow := range tt.expectedDataRows {
						actualRow := rows[i+1]
						for j := 0; j < len(expectedRow); j++ {
							var actualCell string
							if j < len(actualRow) {
								actualCell = actualRow[j]
							} else {
								actualCell = ""
							}

							expectedVal := fmt.Sprintf("%v", expectedRow[j])
							assert.Equal(t, expectedVal, actualCell, "Sai giá trị tại hàng %d, cột %d", i+1, j)
						}
					}
				} else {
					if len(tt.expectedHeaderRow) > 0 {
						assert.Equal(t, 1, len(rows), "Chỉ mong đợi hàng tiêu đề khi không có hàng dữ liệu nào")
					} else {
						assert.Equal(t, 0, len(rows), "Không mong đợi hàng nào khi không có dữ liệu và không có tiêu đề")
					}
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
