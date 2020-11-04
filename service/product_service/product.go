package product_service

import (
	"errors"
	"github.com/astaxie/beego/validation"
	"log"
	"material/lib/app"
	"material/lib/e"
	"material/lib/utils"
	"material/models"
	"strconv"
	"time"
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

	BeginTime utils.Time `json:"begin_time"` //同步開始時間

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
	SupplyCycle int64      `json:"supply_cycle"` // 供货周期

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

type MaterialSelectData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"` //
}

//同步下料单
func ProductSync(m_data *MaterialAdd, data []map[string]interface{}, platform models.Platform) (cb interface{}, err error) {
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
	//contract, err := models.ContractInfo(m_data.ContractId)
	//if err != nil {
	//	return nil, errors.New("合同查询有误")
	//}
	//if contract.ProjectId != project.Id {
	//	return nil, errors.New("合同和项目不对应")
	//}

	// 创建
	t := *models.NewTransaction()
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
		BeginTime:      m_data.BeginTime,
	}
	models.MaterialAddT(&mm, &t)

	//处理材料数据
	validate_item := validation.Validation{}
	if len(data) > 0 {
		for i, v := range data {
			log.Println(i, v)
			//检测数据是否正常
			validate_item.Required(e.ToString(v["material_name"]), "MaterialName").Message("请输入材料名称")
			validate_item.Required(e.ToFloat64(v["price"]), "TotalAmount").Message("请输入正确的材料单价")
			validate_item.Required(e.ToFloat64(v["count"]), "Count").Message("请输入正确的数量")
			validate_item.Required(e.ToFloat64(v["company_id"]), "ContractCount").Message("请输入正确的合同数量")
			validate_item.Required(e.ToInt64(v["project_id"]), "ProjectId").Message("项目有误")
			validate_item.Required(e.ToInt64(v["platform_id"]), "PlatformId").Message("请输入id")
			validate_item.Required(e.ToString(v["unit"]), "Unit").Message("请输入单位")
			//第一层表单验证
			if valid.HasErrors() {
				app.MarkErrors(valid.Errors)
				log.Println(valid.Errors)
				t.Rollback() //回滚
				return nil, app.ErrorsGetOne(valid.Errors)
			}
			data[i]["platform_key"] = platform.PlatformKey
			data[i]["project_id"] = project.Id
			data[i]["project_name"] = project.ProjectName
			data[i]["company_id"] = project.CompanyId

			//判断产品是否存在
			p_item, err := models.ProductCheckInfo(e.ToString(data[i]["material_name"]), e.ToString(data[i]["standard"]), project.Id)
			productModel := models.Product{}
			if err != nil {
				// 处理新增数据结构
				productModel = models.Product{
					Model:              models.Model{},
					MaterialName:       e.ToString(data[i]["material_name"]),
					BlankingAttachment: e.ToString(data[i]["blanking_attachment"]),
					Attachment:         e.ToString(data[i]["attachment"]),
					InstallMap:         e.ToString(data[i]["install_map"]),
					Price:              e.ToFloat64(data[i]["price"]),
					Count:              e.ToFloat64(data[i]["count"]),
					ContractCount:      e.ToFloat64(data[i]["contract_count"]),
					PackCount:          0,
					SendCount:          0,
					ReturnCount:        0,
					ReceiveCount:       0,
					Unit:               e.ToString(data[i]["unit"]),
					ProjectId:          project.Id,
					ProjectName:        project.ProjectName,
					ReplenishmentFlag:  0,
					ProductSubFlag:     0,
					ConfigData:         "[]",
					AppendAttachment:   e.ToString(data[i]["append_attachment"]),
					//ProjectMaterialId:   data[i].ProjectMaterialId,
					//AdminMaterialInfoId: "",
					ProjectAdditional: e.ToString(data[i]["project_additional"]),
					Remark:            e.ToString(data[i]["remark"]),
					Length:            e.ToFloat64(data[i]["length"]),
					Width:             e.ToFloat64(data[i]["width"]),
					Height:            e.ToFloat64(data[i]["height"]),
					Location:          e.ToString(data[i]["location"]),
					Standard:          e.ToString(data[i]["standard"]),
					ArriveDate:        utils.Time{},
					Cuid:              0,
					CompanyId:         project.CompanyId,
					SupplyCycle:       e.ToInt64(data[i]["supply_cycle"]),
					MaterialId:        mm.Id,
					PlatformKey:       platform.PlatformKey,
					PlatformUid:       mm.PlatformUid,
					PlatformId:        e.ToString(data[i]["platform_id"]),
					ContractId:        mm.ContractId,
				}
				err = models.ProductAddT(&productModel, &t)
				if err != nil {
					t.Rollback()
					return nil, err
				}
				if mm.Type == models.P内装材料 {

				}
				switch mm.Type {
				case models.P内装材料:
					break
				case models.P幕墙面材:
					lm := models.ProductLinkSurface{
						W1:               e.ToFloat64(data[i]["w1"]),
						W2:               e.ToFloat64(data[i]["w2"]),
						W3:               e.ToFloat64(data[i]["w3"]),
						W4:               e.ToFloat64(data[i]["w4"]),
						W5:               e.ToFloat64(data[i]["w5"]),
						W6:               e.ToFloat64(data[i]["w6"]),
						W7:               e.ToFloat64(data[i]["w7"]),
						W8:               e.ToFloat64(data[i]["w8"]),
						W9:               e.ToFloat64(data[i]["w9"]),
						H1:               e.ToFloat64(data[i]["h1"]),
						H2:               e.ToFloat64(data[i]["h2"]),
						H3:               e.ToFloat64(data[i]["h3"]),
						H4:               e.ToFloat64(data[i]["h4"]),
						H5:               e.ToFloat64(data[i]["h5"]),
						H6:               e.ToFloat64(data[i]["h6"]),
						H7:               e.ToFloat64(data[i]["h7"]),
						H8:               e.ToFloat64(data[i]["h8"]),
						H9:               e.ToFloat64(data[i]["h9"]),
						L1:               e.ToFloat64(data[i]["l1"]),
						L2:               e.ToFloat64(data[i]["l2"]),
						L3:               e.ToFloat64(data[i]["l3"]),
						L4:               e.ToFloat64(data[i]["l4"]),
						L5:               e.ToFloat64(data[i]["l5"]),
						L6:               e.ToFloat64(data[i]["l6"]),
						L7:               e.ToFloat64(data[i]["l7"]),
						L8:               e.ToFloat64(data[i]["l8"]),
						L9:               e.ToFloat64(data[i]["l9"]),
						WSize:            e.ToInt64(data[i]["w_size"]),
						HSize:            e.ToInt64(data[i]["h_size"]),
						LSize:            e.ToInt64(data[i]["l_size"]),
						SurfaceTreatment: e.ToString(data[i]["surface_treatment"]),
						Color:            e.ToString(data[i]["color"]),
						Area:             e.ToString(data[i]["area"]),
						TotalCount:       e.ToString(data[i]["total_count"]),
						ProductId:        productModel.Id,
					}
					models.ProductLinkSurfaceAddT(&lm, &t)
					break
				case models.P幕墙辅材:
					lm := models.ProductLinkAuxiliary{
						MaterialStatus: e.ToString(data[i]["material_status"]),
						Weight:         e.ToString(data[i]["weight"]),
						TotalArea:      e.ToString(data[i]["total_area"]),
						ProductId:      productModel.Id,
					}
					models.ProductLinkAuxiliaryAddT(&lm, &t)
					break
				case models.P幕墙线材:
					lm := models.ProductLinkWire{
						SurfaceTreatment: e.ToString(data[i]["surface_treatment"]),
						Color:            e.ToString(data[i]["color"]),
						Area:             e.ToString(data[i]["area"]),
						TotalCount:       e.ToString(data[i]["total_count"]),
						ProductId:        productModel.Id,
					}
					models.ProductLinkWireAddT(&lm, &t)
					break
				}
			} else {
				productModel = *p_item
				models.ProductEditT(productModel.Id, map[string]interface{}{
					"count": productModel.Count + e.ToFloat64(data[i]["count"]),
					"price": e.ToFloat64(data[i]["price"]),
				}, &t)
			}
			ml := models.MaterialLink{
				MaterialId:  mm.Id,
				ProductId:   productModel.Id,
				CompanyId:   project.CompanyId,
				Count:       e.ToFloat64(data[i]["count"]),
				ProjectId:   project.Id,
				Price:       e.ToFloat64(data[i]["price"]),
				SupplyCycle: e.ToInt64(data[i]["supply_cycle"]),
				ReceiveTime: utils.Time{Time: time.Now()},
				Status:      0,
			}
			models.MaterialLinkAddT(&ml, &t)
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
		MaterialName:       data.MaterialName,
		BlankingAttachment: "",
		Attachment:         data.Attachment,
		InstallMap:         data.InstallMap,
		Price:              data.Price,
		Count:              data.Count,
		ContractCount:      data.ContractCount,
		PackCount:          0,
		SendCount:          0,
		ReturnCount:        0,
		ReceiveCount:       0,
		Unit:               data.Unit,
		ProjectId:          data.ProjectId,
		ProjectName:        data.ProjectName,
		ReplenishmentFlag:  0,
		ProductSubFlag:     0,
		ConfigData:         data.ConfigData,
		AppendAttachment:   data.AppendAttachment,
		//ProjectMaterialId:   0,
		//AdminMaterialInfoId: "",
		ProjectAdditional: data.ProjectAdditional,
		Remark:            data.Remark,
		Length:            data.Length,
		Width:             data.Width,
		Height:            data.Height,
		Location:          data.Location,
		Standard:          data.Standard,
		ArriveDate:        data.ArriveDate,
		Cuid:              data.Cuid,
		CompanyId:         data.CompanyId,
		SupplyCycle:       data.SupplyCycle,
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

// 下料单列表
func MaterialApiLists(page int, limit int, maps string) ([]*models.Material, error) {
	offset := (page - 1) * limit
	return models.MaterialGetLists(offset, limit, maps)
}

// 获取Api列表
func ApiLists(page int, limit int, maps string) ([]*models.Product, error) {
	offset := (page - 1) * limit
	return models.ProductGetLists(offset, limit, maps)
}

// 获取Select列表
func Select(maps string) ([]*models.Product, error) {
	return models.ProductGetSelect(maps)
}

func Tables(project_id, material_id, company_id int64) ([]*models.MaterialLink, error) {
	maps := utils.WhereToMap(nil)
	maps["flag"] = 1
	maps["company_id"] = company_id

	maps["material_id"] = material_id

	//if material_id != 0 {
	//	maps["material_id"] = material_id
	//}

	maps["project_id"] = project_id
	return models.MaterialLinkGetAllLists(utils.BuildWhere(maps))
}

func SelectMaterialLists(company_id, project_id int64) ([]MaterialSelectData, error) {
	lists, err := models.MaterialGetSelect("flag =1 AND company_id = " + strconv.Itoa(int(company_id)) + " AND project_id = " + strconv.Itoa(int(project_id)))
	if err != nil {
		return nil, err
	}
	cb := make([]MaterialSelectData, len(lists))
	for i := 0; i < len(lists); i++ {
		cb[i] = MaterialSelectData{
			Id:   lists[i].Id,
			Name: lists[i].Name,
		}
	}
	return cb, nil
}
