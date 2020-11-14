package company

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/gredis"
	"material/models"
)

func Company() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authentication")
		if token == "" {
			e.ApiOpt(c, e.INVALID_PARAMS, e.GetMsg(e.INVALID_PARAMS), e.GetEmptyStruct())
			return
		} else {
			platform_key := ""
			if token == "MLGZ_ZJSZ" {
				platform_key = "MLGZ_ZJSZ"
			} else {
				//检测Token
				platform_key = gredis.GetCacheString("PLATFORM" + token)
				if platform_key == "" {
					e.ApiOpt(c, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT), e.GetEmptyStruct())
					return
				}
			}
			log.Println(platform_key)
			platform, err := models.PlatformGetInfoOrKey(platform_key)
			if err != nil {
				log.Println(err)
				e.ApiOpt(c, e.ERROR_AUTH, e.GetMsg(e.ERROR_AUTH), e.GetEmptyStruct())
				return
			}
			c.Set("platform", *platform)
			c.Set("platform_key", platform_key)
		}

		c.Next()
	}
}
