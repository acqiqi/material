package models

type Contract struct {
	Model
	ContractName          string  `json:"contract_name"`  // 合同名
	ContractNo            string  `json:"contract_no"`    // 合同编号
	UseTime               int     `json:"use_time"`       // 签订时间
	UseAddress            string  `json:"use_address"`    // 签约地点
	Price                 float64 `json:"price"`          // 全部总金额
	AName                 string  `json:"a_name"`         // 甲方名
	ATel                  string  `json:"a_tel"`          // 甲方电话
	AEmail                string  `json:"a_email"`        // 甲方email
	BName                 string  `json:"b_name"`         // 乙方名
	BTel                  string  `json:"b_tel"`          // 乙方电话
	BEmail                string  `json:"b_email"`        // 乙方email
	ContractPrice         float64 `json:"contract_price"` // 合同金额
	Attachment            string  `json:"attachment"`     // 合同附件
	ContractType          string  `json:"contract_type"`  // 合同类型 供应商合同 框架协议
	ProjectId             int     `json:"project_id"`     // 项目id
	Project               Company `gorm:"ForeignKey:ProjectId" json:"project"`
	StartDate             int     `json:"start_date"`               // 合同开始时间
	EndDate               int     `json:"end_date"`                 // 合同结束时间
	PayWay                string  `json:"pay_way"`                  // 付款方式
	BreachItem            string  `json:"breach_item"`              // 违约条款
	TotalContractTaxPrice float64 `json:"total_contract_tax_price"` // 合同含税总价
	Remark                string  `json:"remark"`                   // 备注
	ItemReceiptAmount     float64 `json:"item_receipt_amount"`      // 已开进项发票总额
	InStorageAmount       float64 `json:"in_storage_amount"`        // 合同入库材料总金额
	RequestAccount        float64 `json:"request_account"`          // 总请款金额
	ReceiptAccount        float64 `json:"receipt_account"`          // 已收发票金额
	PayAccount            float64 `json:"pay_account"`              // 付款总金额
	HasR                  float64 `json:"has_r"`                    // 已请总金额
	CompanyId             int     `json:"company_id"`               // 公司id
	Company               Company `gorm:"ForeignKey:CompanyId" json:"company"`
	Cuid                  int     `json:"cuid"`
	DeletedAt             string  `json:"deleted_at"`
}

// 新增产品
func ContractAdd(contract *Contract) error {
	contract.Flag = 1
	if err := db.Create(&contract).Error; err != nil {
		return err
	}
	return nil
}

// 编辑产品
func ContractEdit(id int64, data interface{}) error {
	if err := db.Model(&Contract{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// 获取产品详情
func ContractInfo(id int64) (*Contract, error) {
	var contract Contract
	err := db.Where("id = ?", id).Preload("Company").First(&contract).Error
	if err != nil {
		return &Contract{}, err
	}
	return &contract, nil
}
