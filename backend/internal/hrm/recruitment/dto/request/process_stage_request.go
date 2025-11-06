package dto

type CreateStageRequest struct {
	ProcessFormID string `json:"process_form_id" binding:"required"`
	StageName     string `json:"stage_name" binding:"required"`
}

type UpdateStageRequest struct {
	StageName string `json:"stage_name" binding:"required"`
}
