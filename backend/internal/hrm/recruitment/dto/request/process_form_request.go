package dto

type CreateProcessFormRequest struct {
	ProcessFormName string                    `json:"process_form_name" validate:"required"`
	Stages          []CreateProcessStageInput `json:"stages" validate:"required,dive"`
}

type CreateProcessStageInput struct {
	StageName string `json:"stage_name" validate:"required"`
}

type UpdateProcessFormRequest struct {
	ProcessFormID   string `json:"process_form_id" validate:"required"`
	ProcessFormName string `json:"process_form_name" validate:"required"`
}

type DeleteProcessFormRequest struct {
	ProcessFormID string `json:"process_form_id" validate:"required"`
}
