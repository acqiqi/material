package models

import (
	"github.com/jinzhu/gorm"
	"log"
)

type PlatformCompany struct {
	Model
	CompanyId   int64   `json:"company_id"` // 企业id
	Company     Company `gorm:"ForeignKey:CompanyId" json:"company"`
	PlatformKey string  `json:"platform_key"` // 平台key
	PlatformUid string  `json:"platform_uid"` // 平台用户id
	CompanyKey  string  `json:"company_key"`
	DataOrigin  string  `json:"data_origin"`
	Opt         string  `json:"opt"`
	SupplierId  string  `json:"supplier_id"`
}

func PlatformCompanyGetLists(maps string) ([]*PlatformCompany, error) {
	var pc []*PlatformCompany
	err := db.Where(maps).Order("id desc").Preload("Company").Find(&pc).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return pc, nil
}

func PlatformCompanyGetInfo(id int64) (*PlatformCompany, error) {
	var pc PlatformCompany
	err := db.Where("id = ? AND flag =1", id).Preload("Company").First(&pc).Error
	if err != nil {
		return &PlatformCompany{}, err
	}
	return &pc, nil
}

func PlatformCompanyGetInfoCheck(company_id int64, platform_key string) (*PlatformCompany, error) {
	var pc PlatformCompany
	err := db.Where("company_id = ? AND platform_key= ?  AND flag =1", company_id, platform_key).Preload("Company").First(&pc).Error
	if err != nil {
		return &PlatformCompany{}, err
	}
	return &pc, nil
}

//检测平台是否存在
func PlatformCompanyCheck(company_id int64, platform_key string, platform_uid string) (*PlatformCompany, error) {
	log.Println(company_id, platform_key, platform_uid)
	var pc PlatformCompany
	err := db.Where("company_id = ? AND platform_key = ? AND platform_uid = ? AND flag =1",
		company_id, platform_key, platform_uid).Preload("Company").First(&pc).Error
	if err != nil {
		return &PlatformCompany{}, err
	}
	return &pc, nil
}

// 新增平台公司连接
func PlatformCompanyAdd(pc *PlatformCompany) error {
	pc.Flag = 1
	if err := db.Create(&pc).Error; err != nil {
		return err
	}
	return nil
}

// 新增平台公司连接
func PlatformCompanyDelete(pc *PlatformCompany) error {
	pc.Flag = -1
	return PlatformCompanyEdit(pc.Id, &pc)
}

// 编辑平台公司连接
func PlatformCompanyEdit(id int64, data interface{}) error {
	if err := db.Model(&PlatformCompany{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
