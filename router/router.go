package router

import (
	"github.com/gin-gonic/gin"
	v1 "material/controller/v1"
	"material/middleware/app_middleware"
	"material/middleware/company"
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
	r.Use(app_middleware.App())
	//注册AuthRouter
	setAuthRouter(r)
	// 接口主节点
	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		setV1Router(apiv1)
		//需要企业授权的
		apiv2 := apiv1.Group("/company")
		apiv2.Use(company.Company())
		{
			apiv2.POST("/project_list", v1.ProjectList)      //获取项目列表
			apiv2.POST("/project_create", v1.ProjecttCreate) //创建项
			apiv2.POST("/project_edit", v1.ProjectEdit)      //编辑项目
		}
	}

	return r
}
