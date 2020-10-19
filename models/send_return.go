package models

import "github.com/jinzhu/gorm"

type SendReturn struct {
	Model
	MaterialName     string  `json:"material_name"` // 材料名称
	SendId           int64   `json:"send_id"`
	PackingId        int64   `json:"packing_id"`
	ProductId        int64   `json:"product_id"`
	Product          Product `gorm:"ForeignKey:ProductId" json:"product"`
	PackingProductId int64   `json:"packing_product_id"`
	ProjectId        int64   `json:"project_id"`
	Cuid             int64   `json:"cuid"`
	Count            float64 `json:"count"`       // 退货数量
	Reason           string  `json:"reason"`      // 退货原因
	Remark           string  `json:"remark"`      // 其他描述
	Attachment       string  `json:"attachment"`  // 照片
	RepairId         int64   `json:"repair_id"`   // 补货对应id 暂时不用
	Status           int     `json:"status"`      // 发起中 1接收
	UseCount         float64 `json:"use_count"`   // 接收数量
	ReturnType       string  `json:"return_type"` //退货类型

	CompanyId int64    `json:"company_id"`
	Company   Contract `gorm:"ForeignKey:CompanyId" json:"company"`

	ReplenishId int64 `json:"replenish_id"` //补货id productid
	IsReplenish int64 `json:"is_replenish"` //是否补货
}

// 退货详情
func SendReturnGetInfo(id int64) (*SendReturn, error) {
	var d SendReturn
	err := db.Where("id = ? AND flag = 1", id).Preload("Product").Preload("Company").First(&d).Error
	if err != nil {
		return &SendReturn{}, err
	}
	return &d, nil
}

// 退货列表
func SendReturnGetLists(pageNum int, pageSize int, maps interface{}) ([]*SendReturn, error) {
	var d []*SendReturn
	err := db.Model(&SendReturn{}).Preload("Company").Preload("Product").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&d).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return d, nil
}

// 查询退货列表数量
func SendReturnGetListsCount(maps interface{}) int {
	var d []*SendReturn
	count := 0
	db.Preload("Company").Where(maps).Find(&d).Count(&count)
	return count
}

// 退货
func SendReturnAdd(d *SendReturn) error {
	d.Flag = 1
	if err := db.Create(&d).Error; err != nil {
		return err
	}
	return nil
}

// 编辑退货
func SendReturnEdit(id int64, data interface{}) error {
	if err := db.Model(&SendReturn{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func SendReturnGetSelect(maps string) ([]*SendReturn, error) {
	var d []*SendReturn
	err := db.Where(maps).Order("id desc").Find(&d).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return d, nil
}
