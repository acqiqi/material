package router

import (
	"github.com/gin-gonic/gin"
	v1 "material/controller/v1"
)

func setV1Router(r *gin.RouterGroup) {
	r.POST("/company_create", v1.CompanyCreate)  //创建公司
	r.POST("/company_my_info", v1.CompanyMyInfo) //我的公司详情
	r.POST("/company_my_list", v1.CompanyMyList) //我的公司列表

}
