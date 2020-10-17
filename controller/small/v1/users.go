package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/models"
)

func UserGetInfo(c *gin.Context) {
	user_info, _ := c.Get("user_info")
	e.ApiOk(c, "获取成功", user_info.(models.Users))
}
