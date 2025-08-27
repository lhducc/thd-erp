package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type DecisionRepo interface {
	CreateDecision(ctx context.Context, data *model.Decision, employeeIDs []string) error
	GetDecision(ctx context.Context, id string) (*model.Decision, error)
	GetAllDecision(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]model.Decision, int64, error)
	GetAllDecisionNoPagination(ctx context.Context) ([]model.Decision, error)
	UpdateDecision(ctx context.Context, id string, data *model.DecisionCreate) error
	DeleteDecision(ctx context.Context, id string) error
	CheckExistName(name string) (bool, error)
	GetLastDecisionByCode(ctx context.Context, decision *model.Decision) error
}

type decisionBiz struct {
	repo         DecisionRepo
	employeeRepo EmployeeRepo
}

func NewDecisionBiz(store DecisionRepo, employeeRepo EmployeeRepo) *decisionBiz {
	return &decisionBiz{
		repo:         store,
		employeeRepo: employeeRepo,
	}
}

func (biz *decisionBiz) CreateDecision(ctx context.Context, data *model.DecisionCreate) (string, error) {
	code, err := utils.GenerateCode("QD", 6, func() (string, error) {
		var last model.Decision
		err := biz.repo.GetLastDecisionByCode(ctx, &last)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
		return last.DecisionID, nil
	})
	if err != nil {
		return "", fmt.Errorf("không thể tạo mã: %w", err)
	}

	exists, err := biz.repo.CheckExistName(data.DecisionName)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("tên quyết định đã tồn tại")
	}

	//now := time.Now()
	//if data.SignDate.After(now) {
	//	return "", errors.New("ngày ký không thể trong tương lai")
	//}
	if data.EffectiveDate.Before(data.SignDate) {
		return "", errors.New("ngày hiệu lực không thể trước ngày ký")
	}

	for _, empID := range data.EmployeeIDs {
		_, err := biz.employeeRepo.GetUserById(empID)
		if err != nil {
			return "", fmt.Errorf("không tìm thấy nhân viên ID: %s", empID)
		}
	}

	decision := &model.Decision{
		DecisionID:     code,
		DecisionName:   data.DecisionName,
		EffectiveDate:  data.EffectiveDate,
		SignDate:       data.SignDate,
		Content:        data.Content,
		Condition:      data.Condition,
		AttachedFile:   data.AttachedFile,
		CreatedDate:    time.Now(),
		DecisionTypeID: data.DecisionTypeID,
	}

	if err := biz.repo.CreateDecision(ctx, decision, data.EmployeeIDs); err != nil {
		return "", fmt.Errorf("lỗi tạo quyết định: %w", err)
	}

	return code, nil
}

func (biz *decisionBiz) GetDecision(ctx context.Context, id string) (*model.Decision, error) {
	Contract, err := biz.repo.GetDecision(ctx, id)
	if err != nil {
		return nil, err
	}
	return Contract, nil
}
func (biz *decisionBiz) UpdateDecision(ctx context.Context, id string, data *model.DecisionCreate) error {
	if err := biz.repo.UpdateDecision(ctx, id, data); err != nil {
		return err
	}
	return nil
}
func (biz *decisionBiz) DeleteDecision(ctx context.Context, id string) error {
	if err := biz.repo.DeleteDecision(ctx, id); err != nil {
		return err
	}
	return nil
}

func (biz *decisionBiz) GetAllDecisionPagination(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]model.Decision, int64, error) {
	decisions, totalRecords, err := biz.repo.GetAllDecision(ctx, page, pageSize, filters)
	if err != nil {
		return nil, 0, err
	}
	return decisions, totalRecords, nil
}

func (d *decisionBiz) ExportDecisionTest(ctx context.Context, selectedFields []string) ([]byte, string, error) {
	exporter := utils.NewExcelExporter("Decisions")

	exporter.RegisterField("decision_id", "Mã quyết định", "decision_id", func(item interface{}) any {
		return item.(*model.Decision).DecisionID
	})

	exporter.RegisterField("decision_name", "Tên quyết định", "decision_name", func(item interface{}) any {
		return item.(*model.Decision).DecisionName
	})

	exporter.RegisterField("effective_date", "Ngày hiệu lực", "effective_date", func(item interface{}) any {
		return item.(*model.Decision).EffectiveDate.Format("2006-01-02")
	})

	exporter.RegisterField("sign_date", "Ngày ký", "sign_date", func(item interface{}) any {
		return item.(*model.Decision).SignDate.Format("2006-01-02")
	})

	exporter.RegisterField("content", "Nội dung", "content", func(item interface{}) any {
		return item.(*model.Decision).Content
	})

	exporter.RegisterField("condition", "Điều kiện", "condition", func(item interface{}) any {
		return item.(*model.Decision).Condition
	})

	exporter.RegisterField("attached_file", "Tệp đính kèm", "attached_file", func(item interface{}) any {
		return item.(*model.Decision).AttachedFile
	})

	exporter.RegisterField("created_date", "Ngày tạo", "created_date", func(item interface{}) any {
		return item.(*model.Decision).CreatedDate.Format("2006-01-02")
	})

	exporter.RegisterField("employee_ids", "Mã nhân viên", "employee_ids", func(item interface{}) any {
		dec := item.(*model.Decision)
		if len(dec.Employees) == 0 {
			return ""
		}
		ids := make([]string, 0, len(dec.Employees))
		for _, emp := range dec.Employees {
			ids = append(ids, emp.EmployeeID)
		}
		return strings.Join(ids, ", ")
	})

	exporter.RegisterField("employee_names", "Tên nhân viên", "employee_names", func(item interface{}) any {
		dec := item.(*model.Decision)
		if len(dec.Employees) == 0 {
			return ""
		}
		names := make([]string, 0, len(dec.Employees))
		for _, emp := range dec.Employees {
			names = append(names, emp.Fullname)
		}
		return strings.Join(names, ", ")
	})

	exporter.RegisterField("decision_type_id", "Loại quyết định", "decision_type_id", func(item interface{}) any {
		return item.(*model.Decision).DecisionTypeID
	})

	// Load data
	decisions, err := d.repo.GetAllDecisionNoPagination(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("không thể lấy dữ liệu quyết định: %w", err)
	}

	// Convert slice
	decisionPtrs := make([]*model.Decision, len(decisions))
	for i := range decisions {
		decisionPtrs[i] = &decisions[i]
	}

	return exporter.Export(decisionPtrs, selectedFields)
}
