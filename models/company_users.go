package models

import "github.com/jinzhu/gorm"

type CompanyUsers struct {
	Model
	CompanyId int     `json:"company_id"` // 企业id
	Company   Company `gorm:"ForeignKey:CompanyId" json:"company"`
	Cuid      int     `json:"cuid"`
	IsMain    int     `json:"is_main"`   // 是否主用户
	RoleId    int     `json:"role_id"`   // 角色id
	RuleData  string  `json:"rule_data"` // 其他权限
	Status    int     `json:"status"`    // 状态 0停用 1正常
	DeletedAt string  `json:"deleted_at"`
	IsDefault int     `json:"is_default"` // 默认企业
}

// 获取我的企业列表
func CompanyUsersGetMyList(cuid int64) ([]CompanyUsers, error) {
	var company_users []CompanyUsers
	err := db.Where("cuid = ?", cuid).Preload("Company").Find(&company_users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return company_users, nil
}

// 获取我的企业详情
func CompanyUsersGetInfo(id int64) (*CompanyUsers, error) {
	var cu CompanyUsers
	err := db.Where("company_id = ?", id).Preload("Company").First(&cu).Error
	if err != nil {
		return &CompanyUsers{}, err
	}
	return &cu, nil
}

// 新增企业
func CompanyUsersAdd(company_users *CompanyUsers) error {
	company_users.Flag = 1
	if err := db.Create(&company_users).Error; err != nil {
		return err
	}
	return nil
}
