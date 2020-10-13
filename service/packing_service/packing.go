package packing_service

import (
	"errors"
	"github.com/astaxie/beego/validation"
	uuid "github.com/satori/go.uuid"
	"log"
	"material/lib/app"
	"material/lib/utils"
	"material/models"
)

//打包
type PackingAdd struct {
	Id          int64   `json:"id"`
	PackingName string  `json:"packing_name"` // 包装名称
	SerialNo    string  `json:"serial_no"`    // 包装编号
	Count       float64 `json:"count"`        // 产品总数
	ReturnCount float64 `json:"return_count"` // 包装下退货数量
	Remark      string  `json:"remark"`       // 描述
	CompanyId   int64   `json:"company_id"`
	ProductId   int64   `json:"product_id"`
	MaterialId  int64   `json:"material_id"`

	ContractId int64           `json:"contract_id"` //合同
	Contract   models.Contract `gorm:"ForeignKey:ContractId" json:"contract"`

	ProjectId int64          `json:"project_id"`
	Project   models.Project `gorm:"ForeignKey:ProjectId" json:"project"`

	DepositoryId int64             `json:"depository_id"`
	Depository   models.Depository `gorm:"ForeignKey:DepositoryId" json:"depository"`

	Status int `json:"status"` //0已打包 1已发货 4已收货 已验收
}

// 获取Api列表
func ApiLists(page int, limit int, maps string) ([]*models.Packing, error) {
	offset := (page - 1) * limit
	return models.PackingGetLists(offset, limit, maps)
}

//新增
func Add(data PackingAdd, links []PackingProductAdd) (*models.Packing, error) {
	// 表单验证
	valid := validation.Validation{}
	valid.Required(data.PackingName, "PackingName").Message("请输入打包名称")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}
	log.Println("???")

	serial_no := uuid.NewV4().String()

	//查询总打包量
	total_count := float64(0)
	for _, v := range links {
		total_count = total_count + v.Count
	}

	model := models.Packing{
		PackingName:  data.PackingName,
		SerialNo:     serial_no,
		Count:        total_count,
		ReturnCount:  data.ReturnCount,
		Remark:       data.Remark,
		CompanyId:    data.CompanyId,
		ProductId:    data.ProductId,
		MaterialId:   data.MaterialId,
		ProjectId:    data.ProjectId,
		DepositoryId: data.DepositoryId,
	}
	//创建事务
	t := *models.NewTransaction()

	if err := models.PackingAddT(&model, &t); err != nil {
		return nil, err
	}

	//处理链接
	for _, v := range links {
		v.PackingId = model.Id

		link_model := models.PackingProduct{
			PackingId:     v.PackingId,
			CompanyId:     v.CompanyId,
			OrderReturnid: v.OrderReturnid,
			ProductId:     v.ProductId,
			MaterialId:    v.MaterialId,
			Count:         v.Count,
			ReturnCount:   0,
			MaterialName:  v.MaterialName,
			ContractId:    v.ContractId,
			ProjectId:     v.ProductId,
			DepositoryId:  v.DepositoryId,
		}
		models.PackingProductAddT(&link_model, &t)

		//减少库存
		product, err := models.ProductGetInfo(link_model.ProductId)
		if err != nil {
			t.Rollback()
			return nil, errors.New("材料有误")
		}
		product.PackCount = product.PackCount + link_model.Count
		if err := models.ProductEditT(product.Id, product, &t); err != nil {
			t.Rollback()
			return nil, errors.New("材料有误1")
		}
	}
	t.Commit()
	return &model, nil
}

// 编辑项目
func Edit(data *PackingAdd) error {
	// 表单验证
	log.Println(utils.JsonEncode(data))
	valid := validation.Validation{}
	valid.Required(data.PackingName, "PackingName").Message("请输入打包名称")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return app.ErrorsGetOne(valid.Errors)
	}

	log.Println("???")
	c, err := models.PackingGetInfo(data.Id)
	if err != nil {
		return err
	}
	//model.Id = data.CompanyId
	model := make(map[string]interface{})
	model["PackingName"] = data.PackingName
	//model["SerialNo"] = data.SerialNo
	model["Count"] = data.Count
	model["ReturnCount"] = data.ReturnCount
	model["Remark"] = data.Remark
	model["CompanyId"] = data.CompanyId
	model["ProductId"] = data.ProductId
	model["MaterialId"] = data.MaterialId

	log.Println(model)
	if err := models.PackingEdit(c.Id, model); err != nil {
		return err
	}
	return nil
}

// 解包
func Delete(id int64, company_id int64) error {

	//查询packing
	packing, err := models.PackingGetInfo(id)
	if err != nil {
		return errors.New("包装不存在")
	}

	if packing.CompanyId != company_id {
		return errors.New("非法请求")
	}

	//查询对应的材料
	maps := utils.WhereToMap(nil)
	maps["flag"] = 1
	maps["company_id"] = company_id
	maps["packing_id"] = packing.Id
	pp, err := ApiListsPP(utils.BuildWhere(maps))
	if err != nil {
		return errors.New("材料数据有误")
	}

	//创建事务
	t := *models.NewTransaction()

	for _, v := range pp {
		//查询product
		product, err := models.ProductGetInfo(v.ProductId)
		if err != nil {
			t.Rollback()
			return errors.New("产品数据有误")
		}
		psave := make(map[string]interface{})
		psave["PackCount"] = product.PackCount - v.Count
		if err := models.ProductEditT(product.Id, psave, &t); err != nil {
			t.Rollback()
			return errors.New("保存数据失败")
		}
		vsave := make(map[string]interface{})
		vsave["Flag"] = -1
		if err := models.PackingProductEditT(v.Id, vsave, &t); err != nil {
			t.Rollback()
			return errors.New("包装产品数据有误")
		}
	}
	packsave := make(map[string]interface{})
	packsave["Flag"] = -1
	if err := models.PackingEditT(packing.Id, packsave, &t); err != nil {
		t.Rollback()
		return errors.New("打包数据有误")
	}
	t.Commit()
	return nil
}

func Tables(project_id int64, company_id int64) ([]*models.Packing, error) {
	maps := utils.WhereToMap(nil)
	maps["flag"] = 1
	maps["company_id"] = company_id
	maps["project_id"] = project_id
	maps["status"] = 0
	return models.PackingGetSelect(utils.BuildWhere(maps))
}

// 获取Select列表
func Select(maps string) ([]*models.Packing, error) {
	return models.PackingGetSelect(maps)
}
