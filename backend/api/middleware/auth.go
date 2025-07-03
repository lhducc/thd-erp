package middleware

import (
	"erp/backend/internal/auth/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		accessToken := strings.Replace(authorization, "Bearer ", "", 1)

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("employeeId", payload.UserId)
		c.Set("role", payload.Roles)

		c.Next()
	}
}
