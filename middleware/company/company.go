package company

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/models"
	"strconv"
)

func Company() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_company_id_str := c.GetHeader("CompanyId")
		user_company_id, err := strconv.Atoi(user_company_id_str)
		if err != nil {
			e.ApiOpt(c, e.ERROR_COMPANY_HEADER, e.GetMsg(e.ERROR_COMPANY_HEADER), e.GetEmptyStruct())
			return
		}

		user_company, err := models.CompanyUsersGetInfo(int64(user_company_id))
		if err != nil {
			e.ApiOpt(c, e.ERROR_COMPANY_NOT, e.GetMsg(e.ERROR_COMPANY_NOT), e.GetEmptyStruct())
			return
		}
		c.Set("company", *user_company)
		c.Next()
	}
}
