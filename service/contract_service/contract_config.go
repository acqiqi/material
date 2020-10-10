package contract_service

import (
	"material/models"
)

// 合同付款方式
type ContractConfigAdd struct {
	Id                int64   `json:"id"`
	OtherAmount       float64 `json:"other_amount"`        // 合同其他费用
	OtherAmountRemark string  `json:"other_amount_remark"` // 合同其他费用说明
	AutoFlag          int     `json:"auto_flag"`           // 是否默认按照自动生成材料总额计算
	PrePayFlag        int     `json:"pre_pay_flag"`        // 是否有预付款
	PrePayPercent     float64 `json:"pre_pay_percent"`     // 预付款百分比
	PrePayAmount      float64 `json:"pre_pay_amount"`      // 按照预付款金额
	PrePayRemark      string  `json:"pre_pay_remark"`      // 预付款说明
	ProgressByMonth   int     `json:"progress_by_month"`   // 按月结,否则按批结
	ProgressPercent   float64 `json:"progress_percent"`    // 进度款百分比
	ProgressRemark    string  `json:"progress_remark"`     // 进度款说明
	AccountFlag       int     `json:"account_flag"`        // 是否有结算款
	AccountPercent    float64 `json:"account_percent"`     // 结算款比例
	AccountRemark     string  `json:"account_remark"`      // 结算款备注
	RetentionFlag     int     `json:"retention_flag"`      // 是否有质保金
	RetentionPercent  float64 `json:"retention_percent"`   // 质保金比例
	RetentionRemark   string  `json:"retention_remark"`    // 质保金备注
	ContractId        int64   `json:"contract_id"`
}

//新增
func AddConfig(data *ContractConfigAdd) (*models.ContractConfig, error) {

	model := models.ContractConfig{
		Model:             models.Model{},
		OtherAmount:       data.OtherAmount,
		OtherAmountRemark: data.OtherAmountRemark,
		AutoFlag:          data.AutoFlag,
		PrePayFlag:        data.PrePayFlag,
		PrePayPercent:     data.PrePayPercent,
		PrePayAmount:      data.PrePayAmount,
		PrePayRemark:      data.PrePayRemark,
		ProgressByMonth:   data.ProgressByMonth,
		ProgressPercent:   data.ProgressPercent,
		ProgressRemark:    data.ProgressRemark,
		AccountFlag:       data.AccountFlag,
		AccountPercent:    data.AccountPercent,
		AccountRemark:     data.AccountRemark,
		RetentionFlag:     data.RetentionFlag,
		RetentionPercent:  data.RetentionPercent,
		RetentionRemark:   data.RetentionRemark,
		ContractId:        data.ContractId,
	}

	if err := models.ContractConfigAdd(&model); err != nil {
		return nil, err
	}
	return &model, nil
}
