package models

import (
	"github.com/jinzhu/gorm"
	"material/lib/utils"
	"strconv"
)

type Company struct {
	Model
	Cuid       int    `json:"cuid"`         // 注册人cuid
	Name       string `json:"name"`         // 企业名称
	Mobile     string `json:"mobile"`       // 企业手机号
	Tel        string `json:"tel"`          // 电话号
	Address    string `json:"address"`      // 企业地址
	Desc       string `json:"desc"`         // 描述
	AuthPics   string `json:"auth_pics"`    // 资质 多图
	VipLevel   int    `json:"vip_level"`    // 企业购买等级
	VipEndTime int    `json:"vip_end_time"` // 到期时间
	Status     int    `json:"status"`       // 状态 0 停业  1营业 -1停用
	CompanyKey string `json:"company_key"`  //企业Key

}

func CompanyGetInfo(id int64) (*Company, error) {
	var user Company
	err := db.Where("id = ? AND flag =1", id).First(&user).Error
	if err != nil {
		return &Company{}, err
	}
	return &user, nil
}

func CompanyGetInfoOrKey(key string) (*Company, error) {
	var user Company
	err := db.Where("company_key = ? AND flag =1", key).First(&user).Error
	if err != nil {
		return &Company{}, err
	}
	return &user, nil
}

//获取我的主企业
func CompanyGetUserList(cuid int64) ([]Company, error) {
	var companys []Company
	err := db.Where("cuid = ? AND flag =1", cuid).Find(&companys).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return companys, nil
}

// 新增企业
func CompanyAdd(company *Company) error {
	company.Flag = 1
	if err := db.Create(&company).Error; err != nil {
		return err
	}
	return nil
}

// 获取企业Key
func CompanyGetKey() string {
	k := utils.RandInt64(10000000, 99999999)
	var company Company
	err := db.Where("company_key = ? AND flag =1", k).First(&company).Error
	if err != nil {
		return strconv.FormatInt(k, 10)
	}
	return GetMUserKey()
}
