package middleware

import (
	"erp/backend/internal/auth/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil || strings.TrimSpace(accessToken) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "access token not found"})
			return
		}

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("accountId", payload.ID)
		c.Set("role", payload.Roles)

		c.Next()
	}
}
