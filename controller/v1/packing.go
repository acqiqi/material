package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/packing_service"
	"material/service/product_service"
)

// 仓库列表
func PackingList(c *gin.Context) {
	user_info, _ := c.Get("user_info")
	log.Print(user_info.(models.Users))
	company, _ := c.Get("company")
	log.Println(company.(models.CompanyUsers))

	data := e.ApiPageLists{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	e.CheckApiPageListDefault(&data) //处理页数据

	maps := utils.WhereToMap(data.Map)
	maps["company_id"] = company.(models.CompanyUsers).Company.Id
	maps["flag"] = 1

	lists, _ := packing_service.ApiLists(data.Page, data.Limit, utils.BuildWhere(maps))
	e.ApiOk(c, "获取成功", e.ApiPageLists{
		Page:  data.Page,
		Limit: data.Limit,
		Lists: lists,
		Total: models.PackingGetListsCount(utils.BuildWhere(maps)),
		Map:   data.Map,
	})
}

// 创建仓库
func PackingCreate(c *gin.Context) {
	data := struct {
		Packing packing_service.PackingAdd          `json:"packing"`
		Links   []packing_service.PackingProductAdd `json:"links"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	//查询合同是否存在
	contract, err := models.ContractInfo(data.Packing.ContractId)
	if err != nil {
		e.ApiErr(c, "合同不存在")
		return
	}
	company, _ := c.Get("company")
	if contract.Id != company.(models.Company).Id {
		e.ApiErr(c, "非法请求")
		return
	}
	data.Packing.CompanyId = contract.CompanyId

	if len(data.Links) > 0 {
		//先查询数据
		products_in := make([]int64, len(data.Links))
		for i, v := range data.Links {
			products_in[i] = v.Id
		}
		maps := utils.WhereToMap(nil)
		maps["product_id__in"] = products_in
		maps["company_id"] = company.(models.CompanyUsers).Company.Id
		maps["contract_id"] = contract.Id
		maps["flag"] = 1
		products, err := product_service.Select(utils.BuildWhere(maps))
		if err != nil {
			e.ApiErr(c, "获取产品列表有误")
			return
		}

		if len(products) != len(data.Links) {
			e.ApiErr(c, "材料数据有误")
			return
		}
		for i, v := range data.Links {
			in_flag := false
			for _, p := range products {
				if p.Id == v.ProductId {
					in_flag = true
					//判断库存是否满足
					if (p.Count-p.PackCount)-v.Count < 0 {
						e.ApiErr(c, p.ProjectName+"库存不足")
						return
					}
					data.Links[i].CompanyId = contract.CompanyId
					data.Links[i].ContractId = contract.Id
					data.Links[i].MaterialId = p.MaterialId
					data.Links[i].MaterialName = p.MaterialName
				}
			}
			if !in_flag {
				e.ApiErr(c, "材料数据有误")
				return
			}
		}
		cb, err := packing_service.Add(data.Packing, data.Links)
		if err != nil {
			e.ApiErr(c, "打包失败"+err.Error())
			return
		}
		e.ApiOk(c, "打包成功", cb)
	} else {
		e.ApiErr(c, "请选择要打包的材料")
		return
	}
}
