package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/models"
)

func ProjectList(c *gin.Context) {
	user_info, _ := c.Get("user_info")
	log.Print(user_info.(models.AdminUsers))
	e.ApiOk(c, "登录成功", struct {
		Token string `json:"token"`
	}{Token: "token"})
}
