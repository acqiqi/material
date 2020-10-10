package models

import "github.com/jinzhu/gorm"

type Packing struct {
	Model
	PackingName string   `orm:"packing_name"` // 包装名称
	SerialNo    string   `orm:"serial_no"`    // 包装编号
	Count       int      `orm:"count"`        // 产品总数
	ReturnCount int      `orm:"return_count"` // 包装下退货数量
	Remark      string   `orm:"remark"`       // 描述
	CompanyId   int64    `orm:"company_id"`
	Company     Contract `gorm:"ForeignKey:CompanyId" json:"Company"`
	ProductId   int64    `orm:"product_id"`
	MaterialId  int64    `orm:"material_id"`

	ContractId int64    `json:"contract_id"` //合同
	Contract   Contract `gorm:"ForeignKey:ContractId" json:"contract"`
}

// 产品类型详情
func PackingGetInfo(id int64) (*Packing, error) {
	var pc Packing
	err := db.Where("id = ? AND flag = 1", id).Preload("Company").First(&pc).Error
	if err != nil {
		return &Packing{}, err
	}
	return &pc, nil
}

// 产品类型列表
func PackingGetLists(pageNum int, pageSize int, maps interface{}) ([]*Packing, error) {
	var pc []*Packing
	err := db.Model(&Packing{}).Preload("Company").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&pc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return pc, nil
}

// 查询产品类型列表数量
func PackingGetListsCount(maps interface{}) int {
	var packing []*Packing
	count := 0
	db.Preload("Company").Where(maps).Find(&packing).Count(&count)
	return count
}

// 新增打包 T
func PackingAddT(packing *Packing, t *gorm.DB) error {
	packing.Flag = 1
	if err := t.Create(&packing).Error; err != nil {
		return err
	}
	return nil
}

// 新增打包类型
func PackingAdd(pc *Packing) error {
	pc.Flag = 1
	if err := db.Create(&pc).Error; err != nil {
		return err
	}
	return nil
}

// 编辑打包类型
func PackingEdit(id int64, data interface{}) error {
	if err := db.Model(&Packing{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
