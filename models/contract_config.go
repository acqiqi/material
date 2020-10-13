package models

// 合同付款方式
type ContractConfig struct {
	Model
	OtherAmount       float64  `json:"other_amount"`        // 合同其他费用
	OtherAmountRemark string   `json:"other_amount_remark"` // 合同其他费用说明
	AutoFlag          int      `json:"auto_flag"`           // 是否默认按照自动生成材料总额计算
	PrePayFlag        int      `json:"pre_pay_flag"`        // 是否有预付款
	PrePayPercent     float64  `json:"pre_pay_percent"`     // 预付款百分比
	PrePayAmount      float64  `json:"pre_pay_amount"`      // 按照预付款金额
	PrePayRemark      string   `json:"pre_pay_remark"`      // 预付款说明
	ProgressByMonth   int      `json:"progress_by_month"`   // 按月结,否则按批结
	ProgressPercent   float64  `json:"progress_percent"`    // 进度款百分比
	ProgressRemark    string   `json:"progress_remark"`     // 进度款说明
	AccountFlag       int      `json:"account_flag"`        // 是否有结算款
	AccountPercent    float64  `json:"account_percent"`     // 结算款比例
	AccountRemark     string   `json:"account_remark"`      // 结算款备注
	RetentionFlag     int      `json:"retention_flag"`      // 是否有质保金
	RetentionPercent  float64  `json:"retention_percent"`   // 质保金比例
	RetentionRemark   string   `json:"retention_remark"`    // 质保金备注
	ContractId        int64    `json:"contract_id"`
	Contract          Contract `gorm:"ForeignKey:ContractId" json:"contract"`
}

// 新增合同付款方式
func ContractConfigAdd(cc *ContractConfig) error {
	cc.Flag = 1
	if err := db.Create(&cc).Error; err != nil {
		return err
	}
	return nil
}

// 编辑合同付款方式
func ContractConfigEdit(id int64, data interface{}) error {
	if err := db.Model(&ContractConfig{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// 获取合同付款方式详情
func ContractConfigGetInfo(id int64) (*ContractConfig, error) {
	var cc ContractConfig
	err := db.Where("company_id = ? AND flag =1", id).Preload("Company").First(&cc).Error
	if err != nil {
		return &ContractConfig{}, err
	}
	return &cc, nil
}
