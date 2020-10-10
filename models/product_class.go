package models

import "github.com/jinzhu/gorm"

//材料类型表
type ProductClass struct {
	Model
	ClassName   string      `json:"class_name"` // 材料类型名称
	Desc        string      `json:"desc"`       // 描述
	CatsId      int         `json:"cats_id"`
	ProductCats ProductCats `gorm:"ForeignKey:CatsId" json:"product_cats"`
	CompanyId   int64       `json:"company_id"`
	Company     Company     `gorm:"ForeignKey:CompanyId" json:"company"`
	Contract    Contract    `gorm:"ForeignKey:ContractId" json:"contract"`
	Cuid        int         `json:"cuid"`
}

// 获取材料大类select
func ProductClassGetSelect(maps string) ([]*ProductClass, error) {
	var pc []*ProductClass
	err := db.Where(maps).Order("id asc").Find(&pc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return pc, nil
}

// 产品类型详情
func ProductClassGetInfo(id int64) (*ProductClass, error) {
	var pc ProductClass
	err := db.Where("id = ? AND flag = 1", id).Preload("Company").First(&pc).Error
	if err != nil {
		return &ProductClass{}, err
	}
	return &pc, nil
}

// 产品类型列表
func ProductClassGetLists(pageNum int, pageSize int, maps interface{}) ([]*ProductClass, error) {
	var pc []*ProductClass
	err := db.Model(&ProductClass{}).Preload("Company").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&pc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return pc, nil
}

// 查询产品类型列表数量
func ProductClassGetListsCount(maps interface{}) int {
	var pc []*ProductClass
	count := 0
	db.Preload("Company").Where(maps).Find(&pc).Count(&count)
	return count
}

// 新增产品类型
func ProductClassAdd(pc *ProductClass) error {
	pc.Flag = 1
	if err := db.Create(&pc).Error; err != nil {
		return err
	}
	return nil
}

// 编辑产品类型
func ProductClassEdit(id int64, data interface{}) error {
	if err := db.Model(&ProductClass{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
