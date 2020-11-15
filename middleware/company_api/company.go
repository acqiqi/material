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
		log.Println("Company_Middleware")
		token := c.GetHeader("Authentication")
		if token == "" {
			e.ApiOpt(c, e.INVALID_PARAMS, e.GetMsg(e.INVALID_PARAMS), e.GetEmptyStruct())
			return
		} else {
			//检测Token
			comapny_ak := gredis.GetCacheString("COMPANY_API" + token)
			if comapny_ak == "" {
				e.ApiOpt(c, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT), e.GetEmptyStruct())
				return
			}
			log.Println(comapny_ak)
			company, err := models.CompanyGetInfoOrAk(comapny_ak)
			if err != nil {
				log.Println(err)
				e.ApiOpt(c, e.ERROR_AUTH, e.GetMsg(e.ERROR_AUTH), e.GetEmptyStruct())
				return
			}

			user_company, err := models.CompanyUsersGetInfoOrIdUid(company.Id, int64(company.Cuid))
			if err != nil {
				log.Println(err)
				e.ApiOpt(c, e.ERROR_AUTH, "用户数据有误！", e.GetEmptyStruct())
				return
			}
			c.Set("company", *user_company)
			user, err := models.GetUsersInfoCuid(int64(company.Cuid))
			c.Set("user_info", *user)
		}
		c.Next()
	}
}
