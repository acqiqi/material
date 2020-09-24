package router

import (
	"github.com/gin-gonic/gin"
	"material/controller"
)

func setAuthRouter(r *gin.Engine) {
	r.POST("/auth/login", controller.Login)
	r.POST("/auth/auto_login", controller.AutoLogin)
	r.POST("/auth/get_auth_login_sms_code", controller.GetAutoLoginSMSCode)

}
