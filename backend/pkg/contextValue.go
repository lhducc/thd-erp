package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ExtractFromContext[T any](c *gin.Context, key string) (T, error) {
	var zero T

	value, exists := c.Get(key)
	if !exists {
		return zero, fmt.Errorf("key '%s' not found in context", key)
	}

	result, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("key '%s' is not of expected type", key)
	}

	return result, nil
}
