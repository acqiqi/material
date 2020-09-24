package router

import (
	"github.com/gin-gonic/gin"
	v1 "material/controller/v1"
	"net/http"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"material/lib/export"
	"material/lib/qrcode"
	"material/middleware/jwt"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//注册AuthRouter
	setAuthRouter(r)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		//获取项目列表
		apiv1.POST("/project_list", v1.ProjectList)

	}

	return r
}
