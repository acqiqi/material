package product_service

import (
	"errors"
	"github.com/astaxie/beego/validation"
	"log"
	"material/lib/app"
	"material/lib/utils"
	"material/models"
)

type ProductAdd struct {
	Id                 int64   `json:"id"`
	MaterialName       string  `json:"material_name"`       // 材料名称
	BlankingAttachment string  `json:"blanking_attachment"` // 下料附件信息(与码里公装关联)
	Attachment         string  `json:"attachment"`          // 附件
	InstallMap         string  `json:"install_map"`         // 安装示意图
	Price              float64 `json:"price"`               // 价格
	Count              float64 `json:"count"`               // 数量
	ContractCount      float64 `json:"contract_count"`      // 与供应商签的合同数量(来源码里公装)
	PackCount          float64 `json:"pack_count"`          // 打包数量
	SendCount          float64 `json:"send_count"`          // 发货数量
	ReturnCount        float64 `json:"return_count"`        // 退货数量
	ReceiveCount       float64 `json:"receive_count"`       // 签收数量
	Unit               string  `json:"unit"`                // 单位
	ProjectId          int     `json:"project_id"`

	ProjectName         string  `json:"project_name"`
	ReplenishmentFlag   int     `json:"replenishment_flag"`    // 是否补货产品
	ProductSubFlag      int     `json:"product_sub_flag"`      // 是否有子部件
	ConfigData          string  `json:"config_data"`           // 自定义字段信息
	AppendAttachment    string  `json:"append_attachment"`     // 附加的资源库信息
	ProjectMaterialId   int     `json:"project_material_id"`   // 码里公装对应下料材料id
	AdminMaterialInfoId int     `json:"admin_materialInfo_id"` // 码里公装对应合同材料id，统计累计数量需要
	ProjectAdditional   string  `json:"project_additional"`    // 项目补充信息
	Remark              string  `json:"remark"`                // 备注
	Length              float64 `json:"length"`                // 长
	Width               float64 `json:"width"`                 // 宽
	Hight               float64 `json:"hight"`                 // 厚
	Location            string  `json:"location"`              // 安装位置
	Standard            string  `json:"standard"`              // 规格
	ArriveDate          int     `json:"arrive_date"`           // 到货时间
	Cuid                int     `json:"cuid"`
	CompanyId           int64   `json:"company_id"`
	SupplyCycle         int     `json:"supply_cycle"` // 供货周期
}

//新增产品
func Add(data *ProductAdd) (*models.Product, error) {
	// 表单验证
	log.Println(utils.JsonEncode(data))
	valid := validation.Validation{}
	valid.Required(data.Cuid, "Cuid").Message("CUID不能为空！")
	valid.Required(data.MaterialName, "MaterialName").Message("请选择产品名称")
	valid.Required(data.CompanyId, "CompanyId").Message("请选择项目")
	valid.Required(data.Price, "Price").Message("请输入单价")
	valid.Required(data.Count, "Count").Message("请输入生产数量")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}

	//检测Company是否存在
	company, err := models.CompanyUsersGetInfoOrCompanyId(data.CompanyId)
	if err != nil {
		return nil, err
	}
	// 检测有用户权限
	if company.Cuid != data.Cuid {
		return nil, errors.New("非法请求")
	}

	log.Println("???")
	model := models.Product{
		MaterialName:        data.MaterialName,
		BlankingAttachment:  "",
		Attachment:          data.Attachment,
		InstallMap:          data.InstallMap,
		Price:               data.Price,
		Count:               data.Count,
		ContractCount:       data.ContractCount,
		PackCount:           0,
		SendCount:           0,
		ReturnCount:         0,
		ReceiveCount:        0,
		Unit:                data.Unit,
		ProjectId:           data.ProjectId,
		ProjectName:         data.ProjectName,
		ReplenishmentFlag:   0,
		ProductSubFlag:      0,
		ConfigData:          data.ConfigData,
		AppendAttachment:    data.AppendAttachment,
		ProjectMaterialId:   0,
		AdminMaterialInfoId: 0,
		ProjectAdditional:   data.ProjectAdditional,
		Remark:              data.Remark,
		Length:              data.Length,
		Width:               data.Width,
		Hight:               data.Hight,
		Location:            data.Location,
		Standard:            data.Standard,
		ArriveDate:          data.ArriveDate,
		Cuid:                data.Cuid,
		CompanyId:           data.CompanyId,
		SupplyCycle:         data.SupplyCycle,
	}
	if err := models.ProductAdd(&model); err != nil {
		return nil, err
	}
	return &model, nil
}

// 编辑项目
func Edit(data *ProductAdd) (*models.Product, error) {
	// 表单验证
	log.Println(utils.JsonEncode(data))
	valid := validation.Validation{}
	valid.Required(data.MaterialName, "MaterialName").Message("请选择产品名称")
	valid.Required(data.CompanyId, "CompanyId").Message("请选择项目")
	valid.Required(data.Price, "Price").Message("请输入单价")
	valid.Required(data.Count, "Count").Message("请输入生产数量")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}

	log.Println("???")
	model := models.Product{}
	model.MaterialName = data.MaterialName
	model.Attachment = data.Attachment
	model.InstallMap = data.InstallMap
	model.Price = data.Price
	model.Count = data.Count
	model.ContractCount = data.ContractCount
	model.Unit = data.Unit
	model.ConfigData = data.ConfigData
	model.AppendAttachment = data.AppendAttachment
	model.ProjectAdditional = data.ProjectAdditional
	model.Remark = data.Remark
	model.Length = data.Length
	model.Width = data.Width
	model.Hight = data.Hight
	model.Location = data.Location
	model.Standard = data.Standard
	model.ArriveDate = data.ArriveDate
	model.SupplyCycle = data.SupplyCycle

	if err := models.ProductEdit(model.Id, model); err != nil {
		return nil, err
	}

	return &model, nil
}

// 获取Api列表
func ApiLists(page int, limit int, maps string) ([]*models.Project, error) {
	offset := (page - 1) * limit
	return models.ProjectGetLists(offset, limit, maps)
}
