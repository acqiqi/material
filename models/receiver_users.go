package models

import "github.com/jinzhu/gorm"

type ReceiverUsers struct {
	Model
	CompanyId   int64    `json:"company_id"`
	Company     Contract `gorm:"ForeignKey:CompanyId" json:"company"`
	ProjectId   int64    `json:"project_id"` // 项目id
	Project     Project  `gorm:"ForeignKey:ProjectId" json:"project"`
	Contract    Contract `gorm:"ForeignKey:ContractId" json:"contract"`
	Cuid        int64    `json:"cuid"`
	Users       Users    `gorm:"ForeignKey:Cuid" json:"users"`
	Nickname    string   `json:"nickname"`
	PlatformKey string   `json:"platform_key"`
	PlatformUid string   `json:"platform_uid"`
	Mobile      string   `json:"mobile"`
}

// 收货人详情 平台查询 不检测cuid
func ReceiverUsersGetInfoByPlatform(platform_key, platform_uid string) (*ReceiverUsers, error) {
	var d ReceiverUsers
	err := db.Where("platform_key = ? AND platform_uid = ? AND flag = 1", platform_key, platform_uid).Preload("Users").First(&d).Error
	if err != nil {
		return &ReceiverUsers{}, err
	}
	return &d, nil
}

func ReceiverUsersGetInfo(id int64) (*ReceiverUsers, error) {
	var d ReceiverUsers
	err := db.Where("id = ? AND flag = 1", id).Preload("Users").First(&d).Error
	if err != nil {
		return &ReceiverUsers{}, err
	}
	return &d, nil
}

// 收货人列表
func ReceiverUsersGetLists(pageNum int, pageSize int, maps interface{}) ([]*ReceiverUsers, error) {
	var d []*ReceiverUsers
	err := db.Model(&ReceiverUsers{}).Preload("Users").Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&d).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return []*ReceiverUsers{}, err
	}
	return d, nil
}

// 查询收货人列表数量
func ReceiverUsersGetListsCount(maps interface{}) int {
	var d []*ReceiverUsers
	count := 0
	db.Preload("Users").Where(maps).Find(&d).Count(&count)
	return count
}

// 新增收货人
func ReceiverUsersAdd(d *ReceiverUsers) error {
	d.Flag = 1
	if err := db.Create(&d).Error; err != nil {
		return err
	}
	return nil
}

// 编辑收货人
func ReceiverUsersEdit(id int64, data interface{}) error {
	if err := db.Model(&ReceiverUsers{}).Where("id = ? AND flag = 1 ", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// 检测用户是否有权限
func ReceiverUsersCheckAuth(cuid, project_id int64) (*ReceiverUsers, error) {
	var d ReceiverUsers
	err := db.Where("cuid = ? AND project_id = ? AND flag = 1", cuid, project_id).Preload("Project").Preload("Users").First(&d).Error
	if err != nil {
		return &ReceiverUsers{}, err
	}
	return &d, nil
}
