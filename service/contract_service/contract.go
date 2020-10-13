package contract_service

import (
	"errors"
	"github.com/astaxie/beego/validation"
	"log"
	"material/lib/app"
	"material/lib/utils"
	"material/models"
	"strconv"
)

type ContractAdd struct {
	Id                    int64          `json:"id"`
	ContractName          string         `json:"contract_name"`  // 合同名
	ContractNo            string         `json:"contract_no"`    // 合同编号
	UseTime               utils.Time     `json:"use_time"`       // 签订时间
	UseAddress            string         `json:"use_address"`    // 签约地点
	Price                 float64        `json:"price"`          // 全部总金额
	AName                 string         `json:"a_name"`         // 甲方名
	ATel                  string         `json:"a_tel"`          // 甲方电话
	AEmail                string         `json:"a_email"`        // 甲方email
	BName                 string         `json:"b_name"`         // 乙方名
	BTel                  string         `json:"b_tel"`          // 乙方电话
	BEmail                string         `json:"b_email"`        // 乙方email
	ContractPrice         float64        `json:"contract_price"` // 合同金额
	Attachment            []string       `json:"attachment"`     // 合同附件
	ContractType          string         `json:"contract_type"`  // 合同类型 供应商合同 框架协议
	ProjectId             int64          `json:"project_id"`     // 项目id
	Project               models.Project `gorm:"ForeignKey:ProjectId" json:"project"`
	StartDate             utils.Time     `json:"start_date"`               // 合同开始时间
	EndDate               utils.Time     `json:"end_date"`                 // 合同结束时间
	PayWay                string         `json:"pay_way"`                  // 付款方式
	BreachItem            string         `json:"breach_item"`              // 违约条款
	TotalContractTaxPrice float64        `json:"total_contract_tax_price"` // 合同含税总价
	Remark                string         `json:"remark"`                   // 备注
	ItemReceiptAmount     float64        `json:"item_receipt_amount"`      // 已开进项发票总额
	InStorageAmount       float64        `json:"in_storage_amount"`        // 合同入库材料总金额
	RequestAccount        float64        `json:"request_account"`          // 总请款金额
	ReceiptAccount        float64        `json:"receipt_account"`          // 已收发票金额
	PayAccount            float64        `json:"pay_account"`              // 付款总金额
	HasR                  float64        `json:"has_r"`                    // 已请总金额
	CompanyId             int64          `json:"company_id"`               // 公司id
	Company               models.Company `gorm:"ForeignKey:CompanyId" json:"company"`
	Cuid                  int            `json:"cuid"`
	CreatedAt             utils.Time     `json:"created_at"`
	UpdatedAt             utils.Time     `json:"updated_at"`
	PlatformKey           string         `json:"platform_key"` // 平台key
	PlatformUid           string         `json:"platform_uid"` // 平台用户id
	PlatformId            string         `json:"platform_id"`  // 平台用户id
	IsPlatform            int            `json:"is_platform"`  // 是否三方平台同步
	BindState             int            `json:"bind_state"`   //是否绑定 0否 1是

}

type ContractSelectData struct {
	Id           int64  `json:"id"`
	ContractName string `json:"name"`
	ContractNo   string `json:"contract_no"` // 合同编号
}

//新增合同
func Add(data *ContractAdd) (*models.Contract, error) {
	// 表单验证
	valid := validation.Validation{}
	valid.Required(data.Cuid, "Cuid").Message("CUID不能为空！")
	valid.Required(data.ContractName, "ContractName").Message("请输入合同名")
	valid.Required(data.ContractNo, "ContractNo").Message("请输入合同编号")
	valid.Required(data.ContractPrice, "ContractPrice").Message("请输入合同金额")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}
	//查询项目是否存在
	project, err := models.ProjectGetInfo(int64(data.ProjectId))
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	//检测企业是否和项目是对应的
	if project.Id != data.ProjectId {
		return nil, errors.New("请选择正确的项目")
	}
	if project.CompanyId != data.CompanyId {
		return nil, errors.New("非法请求 not company")
	}

	model := models.Contract{
		Model:         models.Model{},
		ContractName:  data.ContractName,
		ContractNo:    data.ContractNo,
		UseTime:       data.UseTime,
		UseAddress:    data.UseAddress,
		Price:         data.ContractPrice,
		AName:         "",
		ATel:          "",
		AEmail:        "",
		BName:         "",
		BTel:          "",
		BEmail:        "",
		ContractPrice: data.ContractPrice,
		Attachment:    utils.JsonEncode(data.Attachment),
		ContractType:  data.ContractType,
		ProjectId:     project.Id,
		//Project:               models.Project{},
		StartDate:             data.StartDate,
		EndDate:               data.EndDate,
		PayWay:                data.PayWay,
		BreachItem:            data.BreachItem,
		TotalContractTaxPrice: data.TotalContractTaxPrice,
		Remark:                data.Remark,
		ItemReceiptAmount:     0,
		InStorageAmount:       0,
		RequestAccount:        0,
		ReceiptAccount:        0,
		PayAccount:            0,
		HasR:                  0,
		CompanyId:             data.CompanyId,
		Cuid:                  data.Cuid,
		BindState:             data.BindState,
		PlatformUid:           data.PlatformUid,
		PlatformKey:           data.PlatformKey,
		PlatformId:            data.PlatformId,
		IsPlatform:            data.IsPlatform,
	}

	if err := models.ContractAdd(&model); err != nil {
		return nil, err
	}
	return &model, nil
}

// 编辑项目
func Edit(data *ContractAdd) (*models.Contract, error) {
	// 表单验证
	log.Println(utils.JsonEncode(data))
	valid := validation.Validation{}
	valid.Required(data.ContractName, "ContractName").Message("请输入合同名")
	valid.Required(data.ContractNo, "ContractNo").Message("请输入合同编号")
	valid.Required(data.ContractPrice, "ContractPrice").Message("请输入合同金额")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}

	log.Println("???")
	model, err := models.ContractInfo(data.Id)
	if err != nil {
		return nil, err
	}
	model.Id = data.Id
	model.ContractName = data.ContractName
	model.ContractNo = data.ContractNo
	model.UseTime = data.UseTime
	model.UseAddress = data.UseAddress
	model.Price = data.ContractPrice
	model.ContractPrice = data.ContractPrice
	model.Attachment = utils.JsonEncode(data.Attachment)
	model.ContractType = data.ContractType
	model.StartDate = data.StartDate
	model.EndDate = data.EndDate
	model.PayWay = data.PayWay
	model.BreachItem = data.BreachItem
	model.TotalContractTaxPrice = data.TotalContractTaxPrice
	model.Remark = data.Remark

	log.Println(model)
	if err := models.ContractEdit(model.Id, model); err != nil {
		return nil, err
	}

	return model, nil
}

// 获取Api列表
func ApiLists(page int, limit int, maps string) ([]ContractAdd, error) {
	offset := (page - 1) * limit
	list, err := models.ContractGetLists(offset, limit, maps)
	if err != nil {
		return nil, err
	}
	cb := make([]ContractAdd, len(list))
	for i := 0; i < len(list); i++ {
		cb[i] = ContractAdd{
			Id:                    list[i].Id,
			ContractName:          list[i].ContractName,
			ContractNo:            list[i].ContractNo,
			UseTime:               list[i].UseTime,
			UseAddress:            list[i].UseAddress,
			Price:                 list[i].Price,
			AName:                 list[i].AName,
			ATel:                  list[i].ATel,
			AEmail:                list[i].AEmail,
			BName:                 list[i].BName,
			BTel:                  list[i].BTel,
			BEmail:                list[i].BEmail,
			ContractPrice:         list[i].ContractPrice,
			Attachment:            nil,
			ContractType:          list[i].ContractType,
			ProjectId:             list[i].ProjectId,
			StartDate:             list[i].StartDate,
			EndDate:               list[i].EndDate,
			PayWay:                list[i].PayWay,
			BreachItem:            list[i].BreachItem,
			TotalContractTaxPrice: list[i].TotalContractTaxPrice,
			Remark:                list[i].Remark,
			ItemReceiptAmount:     list[i].ItemReceiptAmount,
			InStorageAmount:       list[i].InStorageAmount,
			RequestAccount:        list[i].RequestAccount,
			ReceiptAccount:        list[i].ReceiptAccount,
			PayAccount:            list[i].PayAccount,
			HasR:                  list[i].HasR,
			CompanyId:             list[i].CompanyId,
			Cuid:                  list[i].Cuid,
			CreatedAt:             list[i].CreatedAt,
			UpdatedAt:             list[i].UpdatedAt,
			Project:               list[i].Project,
			Company:               list[i].Company,
		}
	}
	return cb, nil
}

// 获取Select
func SelectLists(company_id int64) ([]ContractSelectData, error) {
	lists, err := models.ContractGetSelect("flag =1 AND company_id = " + strconv.Itoa(int(company_id)))
	if err != nil {
		return nil, err
	}
	cb := make([]ContractSelectData, len(lists))
	for i := 0; i < len(lists); i++ {
		cb[i] = ContractSelectData{
			Id:           lists[i].Id,
			ContractName: lists[i].ContractName,
			ContractNo:   lists[i].ContractNo,
		}
	}
	return cb, nil
}
