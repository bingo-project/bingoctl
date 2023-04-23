package middleware

import (
	"github.com/gin-gonic/gin"
)

func {{.StructName}}() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
