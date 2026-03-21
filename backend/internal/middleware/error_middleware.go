package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			last := c.Errors.Last()
			c.JSON(http.StatusBadRequest, gin.H{"error": last.Error()})
		}
	}
}
