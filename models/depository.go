package models

import "github.com/jinzhu/gorm"

// 仓库
type Depository struct {
	Model
	Name      string  `json:"name"`       // 仓库名称
	Desc      string  `json:"desc"`       // 描述
	Address   string  `json:"address"`    // 仓库地址
	CompanyId int64   `json:"company_id"` // 企业id
	Company   Company `gorm:"ForeignKey:CompanyId" json:"company"`
	Status    int     `json:"status"` // 状态 0停用 1正常
}

// 仓库详情
func DepositoryGetInfo(id int64) (*Depository, error) {
	var d Depository
	err := db.Where("id = ? AND flag = 1", id).Preload("Company").First(&d).Error
	if err != nil {
		return &Depository{}, err
	}
	return &d, nil
}

// 仓库列表
func DepositoryGetLists(pageNum int, pageSize int, maps interface{}) ([]*Depository, error) {
	var d []*Depository
	err := db.Model(&Depository{}).Preload("Company").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&d).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return d, nil
}

// 查询仓库列表数量
func DepositoryGetListsCount(maps interface{}) int {
	var d []*Depository
	count := 0
	db.Preload("Company").Where(maps).Find(&d).Count(&count)
	return count
}

// 新增仓库
func DepositoryAdd(d *Depository) error {
	d.Flag = 1
	if err := db.Create(&d).Error; err != nil {
		return err
	}
	return nil
}

// 编辑仓库
func DepositoryEdit(id int64, data interface{}) error {
	if err := db.Model(&Depository{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func DepositoryGetSelect(maps string) ([]*Depository, error) {
	var d []*Depository
	err := db.Where(maps).Order("id desc").Find(&d).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return d, nil
}
