package router

import (
	"github.com/gin-gonic/gin"
	"material/controller"
	"material/controller/platform"
)

func setAuthRouter(r *gin.Engine) {
	// 普通接口集
	r.POST("/auth/login", controller.Login)
	r.POST("/auth/auto_login", controller.AutoLogin)
	r.POST("/auth/password_login", controller.PasswordLogin)
	r.POST("/auth/get_auth_login_sms_code", controller.GetAutoLoginSMSCode)
	// 平台接口集
	r.POST("/platform_auth/get_access_token", platform.GetAccessToken)
	// 企业接口
	r.POST("/company_api_auth/get_access_token", controller.GetAccessTokenApi)

}
