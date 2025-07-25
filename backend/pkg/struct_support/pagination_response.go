package struct_support

type PaginatedResponse struct {
	Data       any   `json:"data"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
}
