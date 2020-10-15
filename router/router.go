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
			apiv2.POST("/contract_select", v1.ContractSelect) //合同select

			//产品 材料
			apiv2.POST("/product_class_list", v1.ProductClassList)     //材料类型列表
			apiv2.POST("/product_class_create", v1.ProductClassCreate) //新增材料类型
			apiv2.POST("/product_class_edit", v1.ProductClassEdit)     //编辑材料类型
			apiv2.POST("/product_class_delete", v1.ProductClassDelete) //删除材料类型

			apiv2.POST("/product_list", v1.ProductList)   //材料列表
			apiv2.POST("/product_table", v1.ProductTable) //材料表格

			//仓库
			apiv2.POST("/depository_create", v1.DepositoryCreate) //创建昂库
			apiv2.POST("/depository_edit", v1.DepositoryEdit)     //编辑仓库
			apiv2.POST("/depository_list", v1.DepositoryList)
			apiv2.POST("/depository_select", v1.DepositorySelect)

			// 打包
			apiv2.POST("/packing_create", v1.PackingCreate) //打包入库
			apiv2.POST("/packing_list", v1.PackingList)
			apiv2.POST("/packing_delete", v1.PackingDelete) //拆包
			apiv2.POST("/packing_table", v1.PackingTable)   //表格

			//发货相关
			apiv2.POST("/send_create", v1.SendCreate) //发货
			apiv2.POST("/send_list", v1.SendList)     //发货列表

			// 请款相关
			apiv2.POST("/pr_type_select", v1.PrTypeSelect) //获取请款type
			apiv2.POST("/pr_check_price", v1.PrCheckPrice) //获取请款type
			apiv2.POST("/pr_create", v1.PrCreate)          //请款
			apiv2.POST("/pr_list", v1.PrList)              //请款列表
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
		// 合同同步
		api_platformv1.POST("/contract_sync", platform_v1.ContractSync) //同步项目
		// 下料单
		api_platformv1.POST("/product_cats_list", platform_v1.ProductCatsList) //获取材料大类
		api_platformv1.POST("/material_sync", platform_v1.MaterialSync)        //同步下料单

		// 用户部分
		api_platformv1.POST("/users_dd_sync_qrcode", platform_v1.UsersDDSyncQrcode) //获取同步用户二维码
		api_platformv1.POST("/users_sync", platform_v1.UsersSync)                   //同步用户
	}

	return r
}
