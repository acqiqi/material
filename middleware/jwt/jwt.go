package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"material/models"

	"material/lib/e"
	"material/lib/utils"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authentication")
		if token == "" {
			e.ApiOpt(c, e.INVALID_PARAMS, e.GetMsg(e.INVALID_PARAMS), e.GetEmptyStruct())
			return
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					e.ApiOpt(c, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT), e.GetEmptyStruct())
					return
				default:
					e.ApiOpt(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL), e.GetEmptyStruct())
					return
				}
			}
			log.Print(claims.Id)
			user, err := models.GetAdminUsers(claims.Id)
			if err != nil {
				e.ApiOpt(c, e.ERROR_AUTH, e.GetMsg(e.ERROR_AUTH), e.GetEmptyStruct())
				return
			}

			c.Set("user_info", *user)
		}

		c.Next()
	}
}
