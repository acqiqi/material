package company

import (
	"github.com/gin-gonic/gin"
)

func Company() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}
