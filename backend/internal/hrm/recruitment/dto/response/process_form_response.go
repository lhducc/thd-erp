package dto

type ProcessFormResponse struct {
	ProcessFormID   string `json:"process_form_id"`
	ProcessFormName string `json:"process_form_name"`
	IsDeleted       bool   `json:"is_deleted"`
}
