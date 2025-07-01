package utils

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// ExtractFilterArrays extracts and sanitizes query parameters like "job_title_id", "department_id", etc.
func ExtractFilterArrays(ctx *gin.Context, fields []string) map[string]interface{} {
	filters := make(map[string]interface{})

	for _, field := range fields {
		values := ctx.QueryArray(field)
		cleaned := make([]string, 0)

		for _, v := range values {
			for _, part := range strings.Split(v, ",") {
				if trimmed := strings.TrimSpace(part); trimmed != "" {
					cleaned = append(cleaned, trimmed)
				}
			}
		}

		if len(cleaned) > 0 {
			filters[field] = cleaned
		}
	}

	return filters
}
