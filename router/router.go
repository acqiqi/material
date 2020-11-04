package router

import (
	"github.com/gin-gonic/gin"
	platform_v1 "material/controller/platform/v1"
	small_v1 "material/controller/small/v1"
	v1 "material/controller/v1"
	"material/lib/setting"
	"material/middleware/app_middleware"
	"material/middleware/company"
	"material/middleware/platform"
	"material/middleware/small"
	"net/http"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"fmt"
	"github.com/unrolled/secure"
	"material/lib/export"
	"material/lib/qrcode"
	"material/middleware/jwt"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()

	//r.Use(TlsHandler())
	//r.RunTLS(fmt.Sprintf(":%d", setting.ServerSetting.HttpsPort), setting.ServerSetting.CertFile, setting.ServerSetting.KeyFile)

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
	apiv1.POST("/uploader", v1.Uploader)
	apiv1.Use(jwt.JWT())
	{
		setV1Router(apiv1)
		//需要企业授权的
		apiv2 := apiv1.Group("/company")
		apiv2.Use(company.Company())
		{
			// 控制台
			apiv2.POST("/dashboard_home", v1.DashboardHome)       //控制台首页
			apiv2.POST("/dashboard_account", v1.DashboardAccount) //控制台首页

			// 项目
			apiv2.POST("/project_list", v1.ProjectList)       //获取项目列表
			apiv2.POST("/project_create", v1.ProjectCreate)   //创建项
			apiv2.POST("/project_edit", v1.ProjectEdit)       //编辑项目
			apiv2.POST("/project_select", v1.ProjectSelect)   //获取Select
			apiv2.POST("/project_receive", v1.ProjectReceive) //接收项目
			apiv2.POST("/project_info", v1.ProjectInfo)       //项目详情

			// 合同
			apiv2.POST("/contract_create", v1.ContractCreate) //创建合同
			apiv2.POST("/contract_edit", v1.ContractEdit)     //编辑合同
			apiv2.POST("/contract_list", v1.ContractList)     //合同列表
			apiv2.POST("/contract_select", v1.ContractSelect) //合同select
			apiv2.POST("/contract_info", v1.ContractInfo)     //合同详情

			//产品 材料
			apiv2.POST("/product_class_list", v1.ProductClassList)     //材料类型列表
			apiv2.POST("/product_class_create", v1.ProductClassCreate) //新增材料类型
			apiv2.POST("/product_class_edit", v1.ProductClassEdit)     //编辑材料类型
			apiv2.POST("/product_class_delete", v1.ProductClassDelete) //删除材料类型

			apiv2.POST("/material_list", v1.MaterialList)       //下料单列表
			apiv2.POST("/material_info", v1.MaterialInfo)       //下料单详情
			apiv2.POST("/material_receive", v1.MaterialReceive) //接收下料单
			apiv2.POST("/material_select", v1.MaterialSelect)   //下料单Select

			apiv2.POST("/product_list", v1.ProductList)              //材料列表
			apiv2.POST("/return_product_list", v1.ReturnProductList) //退货列表
			apiv2.POST("/product_info", v1.ProductInfo)              //材料详情
			apiv2.POST("/product_table", v1.ProductTable)            //材料表格

			//仓库
			apiv2.POST("/depository_create", v1.DepositoryCreate) //创建昂库
			apiv2.POST("/depository_edit", v1.DepositoryEdit)     //编辑仓库
			apiv2.POST("/depository_list", v1.DepositoryList)
			apiv2.POST("/depository_select", v1.DepositorySelect)

			// 打包
			apiv2.POST("/packing_create", v1.PackingCreate) //打包入库
			apiv2.POST("/packing_list", v1.PackingList)
			apiv2.POST("/packing_delete", v1.PackingDelete)          //拆包
			apiv2.POST("/packing_table", v1.PackingTable)            //表格
			apiv2.POST("/packing_look_qrcode", v1.PackingLookQrcode) //查看二维码
			apiv2.POST("/packing_info", v1.PackingInfo)              //打包详情

			//发货相关
			apiv2.POST("/send_create", v1.SendCreate)                    //发货
			apiv2.POST("/send_list", v1.SendList)                        //发货列表
			apiv2.POST("/send_info", v1.SendInfo)                        //发货列表
			apiv2.POST("/send_look_qrcode", v1.SendLookQrcode)           //发货列表
			apiv2.POST("/send_return_info", v1.SendReturnInfo)           //退货详情
			apiv2.POST("/send_return_use", v1.SendReturnUse)             //接收退货
			apiv2.POST("/send_return_replenish", v1.SendReturnReplenish) //补货

			// 请款相关
			apiv2.POST("/pr_type_select", v1.PrTypeSelect) //获取请款type
			apiv2.POST("/pr_check_price", v1.PrCheckPrice) //获取请款type
			apiv2.POST("/pr_create", v1.PrCreate)          //请款
			apiv2.POST("/pr_list", v1.PrList)              //请款列表
		}
	}
	apiPlatformV1 := r.Group("/platform_api/v1")
	apiPlatformV1.Use(platform.Platform())
	{
		apiPlatformV1.POST("/company_get_info", platform_v1.CompanyGetInfo)   //获取企业详情
		apiPlatformV1.POST("/company_bind_link", platform_v1.CompanyBindLink) //绑定企业
		apiPlatformV1.POST("/company_list", platform_v1.CompanyList)          //绑定企业
		apiPlatformV1.POST("/company_delete", platform_v1.CompanyDelete)      //绑定企业
		// 项目
		apiPlatformV1.POST("/project_sync", platform_v1.ProjectSync) //同步项目
		// 合同同步
		apiPlatformV1.POST("/contract_sync", platform_v1.ContractSync) //同步项目
		// 下料单
		apiPlatformV1.POST("/product_cats_list", platform_v1.ProductCatsList) //获取材料大类
		apiPlatformV1.POST("/material_sync", platform_v1.MaterialSync)        //同步下料单
		apiPlatformV1.POST("/material_dd_sync", platform_v1.MaterialDDSync)   //同步下料单DD

		// 用户部分
		apiPlatformV1.POST("/users_dd_sync_qrcode", platform_v1.UsersDDSyncQrcode) //获取同步用户二维码
		apiPlatformV1.POST("/users_sync", platform_v1.UsersSync)                   //同步用户
	}

	apiSmallV1 := r.Group("/small/v1")
	apiSmallV1.Use(small.Small())
	{
		apiSmallV1.POST("/user_get_info", small_v1.UserGetInfo) //获取用户详情

		apiSmallV1.POST("/packing_info", small_v1.PackingInfo)                             //打包详情
		apiSmallV1.POST("/product_info", small_v1.ProductInfo)                             //产品详情
		apiSmallV1.POST("/send_info", small_v1.SendInfo)                                   //产品详情
		apiSmallV1.POST("/send_receiver_info", small_v1.SendReceiverInfo)                  //确认收货详情
		apiSmallV1.POST("/packing_product_return_info", small_v1.PackingProductReturnInfo) //退货详情
		apiSmallV1.POST("/packing_product_info", small_v1.PackingProductInfo)              //打包产品详情
		apiSmallV1.POST("/send_return", small_v1.SendReturn)                               //退货详情
		apiSmallV1.POST("/send_receiver", small_v1.SendReceiver)                           //确认收货

	}
	return r
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:" + fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
