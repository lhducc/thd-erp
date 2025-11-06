package dto

type ProcessStageResponse struct {
	ProcessStageID int    `json:"process_stage_id"`
	StageName      string `json:"stage_name"`
	StageOrder     int    `json:"stage_order"`
}
