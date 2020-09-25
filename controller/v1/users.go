package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/models"
)

// 获取用户信息
func UsersGetInfo(c *gin.Context) {
	user_info, _ := c.Get("user_info")
	company_list, _ := models.CompanyUsersGetMyList(user_info.(models.Users).Id)
	e.ApiOk(c, "获取成功", struct {
		UserInfo    models.Users          `json:"user_info"`
		CompanyList []models.CompanyUsers `json:"company_list"`
	}{
		UserInfo:    user_info.(models.Users),
		CompanyList: company_list,
	})
}
