package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"material/lib/utils"
)

type AdminUsers struct {
	Model
	Username string `json:"username"` // 账号
	Password string `json:"password"` // 密码
	Mobile   string `json:"mobile"`   // 手机号
	Nickname string `json:"nickname"` // 昵称
	Email    string `json:"email"`    // 邮箱
	Gender   string `json:"gender"`   // 性别
	Status   int    `json:"status"`   // 状态 0停用 1启用
}

// GetArticle Get a single article based on ID
func GetAdminUsers(id int64) (*AdminUsers, error) {
	var admin_users AdminUsers
	err := db.Where("id = ?", id).First(&admin_users).Error
	if err != nil {
		return &AdminUsers{}, err
	}
	return &admin_users, nil
}

//
func GetAdminUsersList(pageNum int, pageSize int, maps interface{}) ([]*AdminUsers, error) {
	var admin_users []*AdminUsers
	err := db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&admin_users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return admin_users, nil
}

//检测密码
func CheckAuth(username, password string) (int64, error) {
	var auth AdminUsers
	password = utils.PasswordEncode(password)
	err := db.Select("id").Where(AdminUsers{Username: username, Password: password}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	if auth.Id > 0 {
		return auth.Id, nil
	}

	return 0, errors.New("Users Or Password Error")
}
