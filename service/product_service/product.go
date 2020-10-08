package product_service

import (
	"errors"
	"github.com/astaxie/beego/validation"
	"log"
	"material/lib/app"
	"material/lib/utils"
	"material/models"
)

// 下单新增结构体
type MaterialAdd struct {
	Id             int64   `json:"id"`
	Name           string  `json:"name"`             // 材料单名称
	TotalAmount    float64 `json:"total_amount"`     // 下料总额（不含税）
	TotalTaxAmount float64 `json:"total_tax_amount"` // 下料总额（含税）
	DesignNo       string  `json:"design_no"`        // 设计订单号
	ApplyNo        string  `json:"apply_no"`         // 下料单号
	Remark         string  `json:"remark"`           // 备注
	CreateType     int     `json:"create_type"`      // 创建类型 新建,    采购计划生成
	Type           int     `json:"type"`             // 类型    内装材料,    幕墙面材,    幕墙辅材,    幕墙线材
	PlatformKey    string  `json:"platform_key"`     // 平台key
	PlatformUid    string  `json:"platform_uid"`     // 平台用户id
	PlatformId     string  `json:"platform_id"`      // 平台id  或者对照订单号

	ProjectId  int64           `json:"project_id"`
	Project    models.Project  `gorm:"ForeignKey:ProjectId" json:"project"`
	CompanyId  int64           `json:"company_id"`
	Company    models.Company  `gorm:"ForeignKey:CompanyId" json:"company"`
	ContractId int64           `json:"contract_id"` //合同
	Contract   models.Contract `gorm:"ForeignKey:ContractId" json:"contract"`
}

// 产品新增结构体
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
	ProjectId          int64   `json:"project_id"`

	ProjectName       string `json:"project_name"`
	ReplenishmentFlag int    `json:"replenishment_flag"` // 是否补货产品
	ProductSubFlag    int    `json:"product_sub_flag"`   // 是否有子部件
	ConfigData        string `json:"config_data"`        // 自定义字段信息
	AppendAttachment  string `json:"append_attachment"`  // 附加的资源库信息
	//ProjectMaterialId   int     `json:"project_material_id"`   // 码里公装对应下料材料id
	//AdminMaterialInfoId int     `json:"admin_materialInfo_id"` // 码里公装对应合同材料id，统计累计数量需要
	ProjectAdditional string  `json:"project_additional"` // 项目补充信息
	Remark            string  `json:"remark"`             // 备注
	Length            float64 `json:"length"`             // 长
	Width             float64 `json:"width"`              // 宽
	Height            float64 `json:"height"`             // 厚

	Location    string     `json:"location"`    // 安装位置
	Standard    string     `json:"standard"`    // 规格
	ArriveDate  utils.Time `json:"arrive_date"` // 到货时间
	Cuid        int        `json:"cuid"`
	CompanyId   int64      `json:"company_id"`
	SupplyCycle int        `json:"supply_cycle"` // 供货周期

	PlatformKey string `json:"platform_key"` //平台key
	PlatformUid string `json:"platform_uid"` //平台uid
	PlatformId  string `json:"platform_id"`  //平台id

	ContractId int64 `json:"contract_id"` //合同
}

// 用于Sync 接口的回调
type CBProjectSync struct {
	Project  []models.Project `json:"project"`
	Material models.Material
}

//同步下料单
func ProductSync(m_data *MaterialAdd, data []ProductAdd, platform models.Platform) (cb interface{}, err error) {
	// 查询是否有同步过的材料
	_, err = models.MaterialCheck(m_data.PlatformId, m_data.PlatformKey, m_data.PlatformUid)
	if err == nil {
		return nil, errors.New("已经同步过，请勿重复同步")
	}

	valid := validation.Validation{}
	valid.Required(m_data.Name, "Name").Message("请输入下料单名称")
	valid.Required(m_data.TotalAmount, "TotalAmount").Message("请输入下料总额")
	valid.Required(m_data.TotalTaxAmount, "TotalTaxAmount").Message("请输入下料含税总额")
	valid.Required(m_data.CompanyId, "CompanyId").Message("材料商id有误")
	valid.Required(m_data.ProjectId, "ProjectId").Message("项目有误")
	valid.Required(m_data.PlatformId, "PlatformId").Message("请输入id")
	//第一层表单验证
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}
	//查询企业是否正确
	comapny, err := models.CompanyGetInfo(m_data.CompanyId)
	if err != nil {
		return nil, errors.New("企业查询有误")
	}
	m_data.CompanyId = comapny.Id
	// 查询项目
	project, err := models.ProjectGetInfo(m_data.ProjectId)
	if err != nil {
		return nil, errors.New("项目查询有误")
	}
	if project.PlatformKey != platform.PlatformKey {
		return nil, errors.New("非法请求")
	}
	//查询合同
	contract, err := models.ContractInfo(m_data.ContractId)
	if err != nil {
		return nil, errors.New("合同查询有误")
	}
	if contract.ProjectId != project.Id {
		return nil, errors.New("合同和项目不对应")
	}

	// 创建
	t := models.NewTransaction()
	mm := models.Material{
		Model:          models.Model{},
		Name:           m_data.Name,
		TotalAmount:    m_data.TotalAmount,
		TotalTaxAmount: m_data.TotalTaxAmount,
		DesignNo:       m_data.DesignNo,
		ApplyNo:        m_data.ApplyNo,
		Remark:         m_data.Remark,
		CreateType:     m_data.CreateType,
		Type:           m_data.Type, //  内装材料,    幕墙面材,    幕墙辅材,    幕墙线材
		PlatformKey:    platform.PlatformKey,
		PlatformUid:    m_data.PlatformUid,
		PlatformId:     m_data.PlatformId, //处理对应id
		ProjectId:      project.Id,
		CompanyId:      project.CompanyId,
		ContractId:     m_data.ContractId, //合同id
	}
	models.MaterialAddT(&mm, t)

	//处理材料数据
	validate_item := validation.Validation{}
	if len(data) > 0 {
		for i, v := range data {
			log.Println(i, v)
			//检测数据是否正常
			validate_item.Required(v.MaterialName, "MaterialName").Message("请输入材料名称")
			validate_item.Required(v.Price, "TotalAmount").Message("请输入正确的材料单价")
			validate_item.Required(v.Count, "Count").Message("请输入正确的数量")
			validate_item.Required(v.CompanyId, "ContractCount").Message("请输入正确的合同数量")
			validate_item.Required(v.ProjectId, "ProjectId").Message("项目有误")
			validate_item.Required(v.PlatformId, "PlatformId").Message("请输入id")
			validate_item.Required(v.Unit, "Unit").Message("请输入单位")
			//第一层表单验证
			if valid.HasErrors() {
				app.MarkErrors(valid.Errors)
				log.Println(valid.Errors)
				t.Rollback() //回滚
				return nil, app.ErrorsGetOne(valid.Errors)
			}
			data[i].PlatformKey = platform.PlatformKey
			data[i].ProjectId = project.Id
			data[i].ProjectName = project.ProjectName
			data[i].CompanyId = project.CompanyId
			// 处理新增数据结构

			product_model := models.Product{
				Model:              models.Model{},
				MaterialName:       data[i].MaterialName,
				BlankingAttachment: data[i].BlankingAttachment,
				Attachment:         data[i].Attachment,
				InstallMap:         data[i].InstallMap,
				Price:              data[i].Price,
				Count:              data[i].Count,
				ContractCount:      data[i].ContractCount,
				PackCount:          0,
				SendCount:          0,
				ReturnCount:        0,
				ReceiveCount:       0,
				Unit:               data[i].Unit,
				ProjectId:          project.Id,
				ProjectName:        project.ProjectName,
				ReplenishmentFlag:  0,
				ProductSubFlag:     0,
				ConfigData:         "[]",
				AppendAttachment:   data[i].AppendAttachment,
				//ProjectMaterialId:   data[i].ProjectMaterialId,
				AdminMaterialInfoId: 0,
				ProjectAdditional:   data[i].ProjectAdditional,
				Remark:              data[i].Remark,
				Length:              data[i].Length,
				Width:               data[i].Width,
				Height:              data[i].Height,
				Location:            data[i].Location,
				Standard:            data[i].Standard,
				ArriveDate:          data[i].ArriveDate,
				Cuid:                0,
				CompanyId:           project.CompanyId,
				Company:             models.Company{},
				SupplyCycle:         data[i].SupplyCycle,
				MaterialId:          mm.Id,
				PlatformKey:         platform.PlatformKey,
				PlatformUid:         mm.PlatformUid,
				PlatformId:          data[i].PlatformId,
				ContractId:          mm.ContractId,
			}
			err := models.ProductAddT(&product_model, t)
			if err != nil {
				t.Rollback()
				return nil, err
			}
		}
	} else {
		// 材料没有的情况
		//return nil,errors.New("材料为空")
	}
	t.Commit()
	return struct {
		MData models.Material `json:"m_data"`
	}{
		MData: mm,
	}, nil
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
		Height:              data.Height,
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
	model.Height = data.Height
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
func ApiLists(page int, limit int, maps string) ([]*models.Product, error) {
	offset := (page - 1) * limit
	return models.ProductGetLists(offset, limit, maps)
}
