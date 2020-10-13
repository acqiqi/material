package pr_service

import (
	"errors"
	"github.com/astaxie/beego/validation"
	uuid "github.com/satori/go.uuid"
	"log"
	"material/lib/app"
	"material/lib/utils"
	"material/models"
)

type PrAdd struct {
	CompanyId int64          `json:"company_id"`
	Company   models.Company `gorm:"ForeignKey:CompanyId" json:"company"`

	ProjectId int64          `json:"project_id"`
	Project   models.Project `gorm:"ForeignKey:ProjectId" json:"project"`

	Price        float64    `json:"price"`     // 实际金额
	TplPrice     float64    `json:"tpl_price"` // 模板金额 输入的
	Type         int        `json:"type"`      // 0进度款 1预付款 2结算款 3质保金
	Desc         string     `json:"desc"`      // 描述
	Count        float64    `json:"count"`     // 数量
	Cuid         int64      `json:"cuid"`
	Status       int        `json:"status"` // 0正常 1接收 2审批通过 -1未通过
	PlatformKey  string     `json:"platform_key"`
	PlatformUid  string     `json:"platform_uid"`
	PlatformId   string     `json:"platform_id"`
	IsPlatform   int        `json:"is_platform"`   // 是否平台审核
	PlatformMsg  string     `json:"platform_msg"`  // 平台审核消息
	PrNo         string     `json:"pr_no"`         // 请款唯一编号
	ApprovalTime utils.Time `json:"approval_time"` // 审批时间
}

const (
	// 付款类型
	PR_TYPE_PRORESS   = 0 // 进度款
	PR_TYPE_PRE       = 1 // 预付款
	PR_TYPE_ACCOUNT   = 2 // 结算款
	PR_TYPE_RETENTION = 3 // 质保金
)

// 请款类型
type PRType struct {
	Key int    `json:"key"`
	Msg string `json:"msg"`
}

// 请款数据
var PRTypeData = []PRType{
	{Key: PR_TYPE_PRORESS, Msg: "进度款"},
	{Key: PR_TYPE_PRE, Msg: "预付款"},
	{Key: PR_TYPE_ACCOUNT, Msg: "结算款"},
	{Key: PR_TYPE_RETENTION, Msg: "质保金"},
}

// 检测type是否存在
func CheckType(t int) bool {
	for _, v := range PRTypeData {
		if v.Key == t {
			return true
		}
	}
	return false
}

//获取支付金额 合同id 请款类型
func CheckPrice(contract_id int64, pr_type int) (price float64, sw bool, err error) {
	// 查询合同是否存在
	cc, err := models.ContractConfigGetInfo(contract_id)
	if err != nil {
		return 0, false, errors.New("合同不存在")
	}
	switch pr_type {
	case PR_TYPE_PRORESS: // 进度款

		break
	case PR_TYPE_PRE: // 预付款
		if cc.PrePayFlag == 0 {
			return 0, false, nil
		}
		if cc.PrePayPercent > 0 {
			price = cc.Contract.Price * (cc.PrePayPercent / 100)
			return price, true, nil
		} else {
			return cc.PrePayAmount, true, nil
		}
		break
	case PR_TYPE_ACCOUNT: // 结算款
		break
	case PR_TYPE_RETENTION: // 质保金
		break
	default:
		return 0, false, errors.New("付款类型有误")
	}
	return 0, false, errors.New("付款类型有误 ！")
}

//新增请款
func Add(data *PrAdd) (*models.Pr, error) {
	// 表单验证
	valid := validation.Validation{}
	valid.Required(data.TplPrice, "TplPrice").Message("请输入金额")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		log.Println(valid.Errors)
		return nil, app.ErrorsGetOne(valid.Errors)
	}

	//检测type是否正确
	if !CheckType(data.Type) {
		return nil, errors.New("请选择正确的请款类型")
	}

	log.Println("???")
	model := models.Pr{}
	model.CompanyId = data.CompanyId
	model.ProjectId = data.ProjectId
	model.Type = data.Type
	model.Desc = data.Desc
	model.Cuid = data.Cuid
	model.Status = 0
	model.PlatformKey = data.PlatformKey
	model.PlatformUid = data.PlatformUid
	model.PlatformId = data.PlatformId
	model.IsPlatform = data.IsPlatform
	model.PlatformMsg = ""
	model.ApprovalTime = data.ApprovalTime

	model.PrNo = uuid.NewV4().String() //编号
	model.Count = 0                    // !!! 这里稍后处理
	model.Price = data.Price
	model.Price = data.TplPrice

	if err := models.PrAdd(&model); err != nil {
		return nil, err
	}
	return &model, nil
}

// 获取Api列表
func ApiLists(page int, limit int, maps string) ([]*models.Pr, error) {
	offset := (page - 1) * limit
	return models.PrGetLists(offset, limit, maps)
}
