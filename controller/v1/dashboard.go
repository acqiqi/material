package v1

import (
	"github.com/gin-gonic/gin"
	"material/lib/e"
	"material/models"
	"time"
)

func DashboardHome(c *gin.Context) {
	company, _ := c.Get("company")
	company_id := company.(models.CompanyUsers).Company.Id
	//获取两天前的时间
	thisTime := time.Now()
	onTime := thisTime.AddDate(0, 0, 0)      //当前时间
	dayTime := thisTime.AddDate(0, 0, -7)    //7天
	monthTime := thisTime.AddDate(0, 0, -30) //30天
	allTime := thisTime.AddDate(0, 0, -9999) //全部

	//查询项目总数
	project_count := models.ProjectGetCount(company_id, allTime, onTime)
	project_day_count := models.ProjectGetCount(company_id, dayTime, onTime)     //7
	project_month_count := models.ProjectGetCount(company_id, monthTime, onTime) //30
	//查询产品数
	product_count := models.ProductGetCount(company_id, allTime, onTime)
	product_day_count := models.ProductGetCount(company_id, dayTime, onTime)     //7
	product_month_count := models.ProductGetCount(company_id, monthTime, onTime) //30

	//查詢請款次數
	pr_count := models.PrGetCount(company_id, allTime, onTime)
	pr_price := models.PrGetSum(company_id, allTime, onTime)

	e.ApiOk(c, "获取成功", map[string]interface{}{
		"project_count":       project_count,
		"project_day_count":   project_day_count,
		"project_month_count": project_month_count,

		"product_count":       product_count,
		"product_day_count":   product_day_count,
		"product_month_count": product_month_count,

		"pr_count": pr_count,
		"pr_price": pr_price,
	})
}
