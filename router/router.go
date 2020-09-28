package router

import (
	"github.com/gin-gonic/gin"
	platform_v1 "material/controller/platform/v1"
	v1 "material/controller/v1"
	"material/middleware/app_middleware"
	"material/middleware/company"
	"material/middleware/platform"
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
			// 项目
			apiv2.POST("/project_list", v1.ProjectList)       //获取项目列表
			apiv2.POST("/project_create", v1.ProjectCreate)   //创建项
			apiv2.POST("/project_edit", v1.ProjectEdit)       //编辑项目
			apiv2.POST("/project_select", v1.ProjectSelect)   //获取Select
			apiv2.POST("/project_receive", v1.ProjectReceive) //接收项目

			// 合同
			apiv2.POST("/contract_create", v1.ContractCreate) //创建合同
			apiv2.POST("/contract_edit", v1.ContractEdit)     //编辑合同
			apiv2.POST("/contract_list", v1.ContractList)     //合同列表

		}
	}
	api_platformv1 := r.Group("/platform_api/v1")
	api_platformv1.Use(platform.Platform())
	{
		api_platformv1.POST("/company_get_info", platform_v1.CompanyGetInfo)   //获取企业详情
		api_platformv1.POST("/company_bind_link", platform_v1.CompanyBindLink) //绑定企业
		api_platformv1.POST("/company_list", platform_v1.CompanyList)          //绑定企业
		api_platformv1.POST("/company_delete", platform_v1.CompanyDelete)      //绑定企业
		// 项目
		api_platformv1.POST("/project_sync", platform_v1.ProjectSync) //同步项目

	}

	return r
}
