package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ExtractAccountIDFromContext(c *gin.Context) (int64, error) {
	raw, exists := c.Get("accountId")
	if !exists {
		return 0, errors.New("accountId not found in context")
	}

	switch v := raw.(type) {
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, errors.New("accountId is a string but not a valid number")
		}
		return id, nil
	default:
		return 0, errors.New("accountId has an unsupported type")
	}
}
