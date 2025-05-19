package usecase

import (
	"context"
	"erp/backend/internal/hrm/hr_profile/model"
	utils "erp/backend/pkg"
	"errors"
	"fmt"
	"time"
)

type DecisionRepo interface {
	CreateDecision(context context.Context, data *model.Decision) error
	GetDecision(ctx context.Context, id string) (*model.Decision, error)
	GetAllDecision(ctx context.Context) ([]model.Decision, error)
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

func (biz *decisionBiz) CreateDecision(context context.Context, data *model.DecisionCreate) error {
	code, err := utils.GenerateCode("QD", 6, func() (string, error) {
		var last model.Decision
		err := biz.repo.GetLastDecisionByCode(context, &last)
		if err != nil {
			return "", err
		}
		return last.DecisionID, nil
	})

	if err != nil {
		return fmt.Errorf("không thể tạo mã: %w", err)
	}

	decision := &model.Decision{
		DecisionID:     code,
		DecisionName:   data.DecisionName,
		EffectiveDate:  data.EffectiveDate,
		SignDate:       data.SignDate,
		Content:        data.Content,
		AttachedFile:   data.AttachedFile,
		Condition:      data.Condition,
		DecisionTypeID: data.DecisionTypeID,
		EmployeeID:     data.EmployeeID,
		CreatedDate:    time.Now(),
	}
	check, err := biz.repo.CheckExistName(data.DecisionName)
	if err != nil || check {
		return err
	}
	now := time.Now()
	if decision.SignDate.After(now) {
		return errors.New("Ngày ký không thể trong tương lai")
	}

	if decision.EffectiveDate.Before(decision.SignDate) {
		return errors.New("Ngày hiệu lực không thể trước ngày ký")
	}

	_, err = biz.employeeRepo.GetUserById(decision.EmployeeID)
	if err != nil {
		return err
	}
	err = biz.repo.CreateDecision(context, decision)
	if err != nil {
		return err
	}

	return nil
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

func (biz *decisionBiz) GetAllDecision(ctx context.Context) ([]model.Decision, error) {
	positions, err := biz.repo.GetAllDecision(ctx)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

func (d *decisionBiz) ExportDecisionTest(ctx context.Context, selectedFields []string) ([]byte, string, error) {
	exporter := utils.NewExcelExporter("Decisions")

	exporter.RegisterField("decision_id", "Mã quyết định", "decision_id", func(item interface{}) any {
		dec := item.(*model.Decision)
		return dec.DecisionID
	})

	exporter.RegisterField("decision_name", "Tên quyết định", "decision_name", func(item interface{}) any {
		dec := item.(*model.Decision)
		return dec.DecisionName
	})

	exporter.RegisterField("effective_date", "Ngày hiệu lực", "effective_date", func(item interface{}) any {
		dec := item.(*model.Decision)
		return dec.EffectiveDate.Format("2006-01-02")
	})

	exporter.RegisterField("sign_date", "Ngày ký", "sign_date", func(item interface{}) any {
		dec := item.(*model.Decision)
		return dec.SignDate.Format("2006-01-02")
	})

	exporter.RegisterField("content", "Nội dung", "content", func(item interface{}) any {
		dec := item.(*model.Decision)
		return dec.Content
	})

	exporter.RegisterField("condition", "Điều kiện", "condition", func(item interface{}) any {
		dec := item.(*model.Decision)
		return dec.Condition
	})

	exporter.RegisterField("attached_file", "Tệp đính kèm", "attached_file", func(item interface{}) any {
		dec := item.(*model.Decision)
		return dec.AttachedFile
	})

	exporter.RegisterField("created_date", "Ngày tạo", "created_date", func(item interface{}) any {
		dec := item.(*model.Decision)
		return dec.CreatedDate.Format("2006-01-02")
	})

	exporter.RegisterField("employee_id", "Mã nhân viên", "employee_id", func(item interface{}) any {
		dec := item.(*model.Decision)
		return dec.EmployeeID
	})

	exporter.RegisterField("employee_name", "Tên nhân viên", "employee_name", func(item interface{}) any {
		dec := item.(*model.Decision)
		if dec.Employee != nil {
			return dec.Employee.Fullname
		}
		return ""
	})

	exporter.RegisterField("decision_type_id", "Loại quyết định", "decision_type_id", func(item interface{}) any {
		dec := item.(*model.Decision)
		return dec.DecisionTypeID
	})

	decisions, err := d.repo.GetAllDecision(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("không thể lấy dữ liệu quyết định: %w", err)
	}

	decisionPtrs := make([]*model.Decision, len(decisions))
	for i := range decisions {
		decisionPtrs[i] = &decisions[i]
	}

	return exporter.Export(decisionPtrs, selectedFields)
}
