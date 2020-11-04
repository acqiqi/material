package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"material/service/packing_service"
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

	//查询项目是否存在
	project, err := models.ProjectGetInfo(data.Packing.ProjectId)
	if err != nil {
		e.ApiErr(c, "项目不存在")
		return
	}

	//查询仓库
	depository, err := models.DepositoryGetInfo(data.Packing.DepositoryId)
	if err != nil {
		e.ApiErr(c, "请选择打包仓库")
		return
	}

	material, err := models.MaterialGetInfo(data.Packing.MaterialId)
	if err != nil {
		e.ApiErr(c, "下料单有误")
		return
	}

	company, _ := c.Get("company")
	if project.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		return
	}
	data.Packing.CompanyId = project.CompanyId

	if len(data.Links) > 0 {
		//先查询数据
		products_in := make([]int64, len(data.Links))
		for i, v := range data.Links {
			products_in[i] = v.MaterialLinkId
		}
		maps := utils.WhereToMap(nil)
		maps["id__in"] = products_in
		maps["company_id"] = company.(models.CompanyUsers).Company.Id
		maps["project_id"] = project.Id
		maps["flag"] = 1
		products, err := models.MaterialLinkGetAllLists(utils.BuildWhere(maps))
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
				if p.Id == v.MaterialLinkId {
					in_flag = true
					//判断库存是否满足
					if (p.Product.Count-p.Product.PackCount)-v.Count < 0 {
						e.ApiErr(c, p.Product.ProjectName+"库存不足")
						return
					}
					if v.Count <= 0 {
						e.ApiErr(c, p.Product.ProjectName+"请输入正确的打包数量")
						return
					}
					//v.CompanyId = project.CompanyId
					//v.ProjectId = project.Id
					//v.MaterialId = p.MaterialId
					//v.MaterialName = p.MaterialName
					data.Links[i].CompanyId = project.CompanyId
					data.Links[i].ProjectId = project.Id
					data.Links[i].MaterialId = material.Id
					data.Links[i].MaterialName = p.Product.MaterialName
					data.Links[i].ProductId = p.Product.Id
					data.Links[i].DepositoryId = depository.Id
					data.Links[i].MaterialLinkId = p.Id
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

// 拆包
func PackingDelete(c *gin.Context) {
	data := struct {
		Id int64 `json:"id"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}
	company, _ := c.Get("company")
	if err := packing_service.Delete(data.Id, company.(models.CompanyUsers).Company.Id); err != nil {
		e.ApiErr(c, err.Error())
	} else {
		e.ApiOk(c, "拆包成功", e.GetEmptyStruct())
	}
}

// 打包表格
func PackingTable(c *gin.Context) {
	data := struct {
		Id int64 `json:"id"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	company, _ := c.Get("company")
	list, err := packing_service.Tables(data.Id, company.(models.CompanyUsers).Company.Id)
	if err != nil {
		e.ApiErr(c, "非法请求")
		return
	}

	e.ApiOk(c, "获取成功", struct {
		Table interface{} `json:"table"`
	}{
		Table: list,
	})
}

//查看二維碼
func PackingLookQrcode(c *gin.Context) {
	data := struct {
		Id int64 `json:"id"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	//查询打包
	packing, err := models.PackingGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "打包信息不存在")
		return
	}

	company, _ := c.Get("company")
	if packing.CompanyId != company.(models.CompanyUsers).Company.Id {
		e.ApiErr(c, "非法请求")
		return
	}

	cb_url, err := packing_service.QrcodeBuild(*packing)
	if err != nil {
		e.ApiErr(c, "获取失败"+err.Error())
	} else {
		e.ApiOk(c, "获取成功", struct {
			Url string `json:"url"`
		}{
			Url: cb_url,
		})
	}
}

// 打包详情
func PackingInfo(c *gin.Context) {
	data := e.ApiId{}
	if err := c.BindJSON(&data); err != nil {
		e.ApiErr(c, err.Error())
		return
	}

	info, err := models.PackingGetInfo(data.Id)
	if err != nil {
		e.ApiErr(c, "下料单不存在")
		return
	}

	ps, err := models.PackingProductGetLists(0, 999999, utils.BuildWhere(map[string]interface{}{
		"packing_id": info.Id,
	}))

	e.ApiOk(c, "获取成功", struct {
		Info     models.Packing `json:"info"`
		Products interface{}    `json:"products"`
	}{
		Info:     *info,
		Products: ps,
	})
}
