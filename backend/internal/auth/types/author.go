package types

import "github.com/gin-gonic/gin"

type Author interface {
	GetUserID(c *gin.Context) string
	IsAllowed(subject string, object string, action string) (bool, error)
	GetRoles(c *gin.Context) []string
}
