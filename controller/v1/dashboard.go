package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/utils"
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
		"company":  company,
	})
}

func DashboardAccount(c *gin.Context) {
	data := struct {
		BeginTime utils.Time `json:"begin_time"`
		EndTime   utils.Time `json:"end_time"`
		ProjectId int64      `json:"project_id"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	//data.ProjectId = 23

	//查询当前项目
	project, err := models.ProjectGetInfo(data.ProjectId)
	if err != nil {
		e.ApiErr(c, "项目不存在")
		return
	}
	product_list, err := models.ProductGetAccountList(project.Id)
	if err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	//查询项目下对应产品列表 全部

	m_count := 0
	//算出年差月
	y := data.EndTime.Year() - data.BeginTime.Year()
	m_count = 12 * y
	//算出月总差
	m := int(data.EndTime.Month()) - int(data.BeginTime.Month())
	m_count = m_count + m + 1

	product_maps := make([]map[string]interface{}, len(product_list))

	str_arr := make([]string, m_count)

	for i, v := range product_list {
		product_map := make(map[string]interface{})
		product_map["id"] = v.Id
		product_map["material_name"] = v.MaterialName
		product_map["standard"] = v.Standard
		product_map["length"] = v.Length
		product_map["width"] = v.Width
		product_map["height"] = v.Height
		product_map["unit"] = v.Unit
		product_map["count"] = v.Count
		product_map["price"] = v.Price
		product_map["packing_product"] = v.PackingProduct
		product_maps[i] = product_map
	}
	for key, value := range product_maps {
		t := data.BeginTime.AddDate(0, 0, 0)
		for i := 0; i < m_count; i++ {
			str_arr[i] = t.Format("20060102")
			product_maps[key][t.Format("20060102")+"_count"] = float64(0)
			product_maps[key][t.Format("20060102")+"_price"] = float64(0)
			for _, pp := range value["packing_product"].([]models.PackingProduct) {
				log.Println(pp.ReceiveTime.Year(), t.Year())
				log.Println(pp.ReceiveTime.Month(), t.Month())

				if pp.ReceiveTime.Year() == t.Year() &&
					pp.ReceiveTime.Month() == t.Month() {
					product_maps[key][t.Format("20060102")+"_count"] = e.ToFloat64(product_maps[key][t.Format("20060102")+"_count"]) + pp.ReceiveCount
					product_maps[key][t.Format("20060102")+"_price"] = e.ToFloat64(product_maps[key][t.Format("20060102")+"_price"]) + value["price"].(float64)
				}
			}
			t = t.AddDate(0, 1, 0)
		}
	}
	e.ApiOk(c, "获取成功", struct {
		Table  interface{} `json:"table"`
		Header interface{} `json:"header"`
	}{
		Table:  product_maps,
		Header: str_arr,
	})
}
