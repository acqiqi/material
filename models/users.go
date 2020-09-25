package models

import (
	"material/lib/utils"
	"strconv"
)

type Users struct {
	Model
	Cuid      int    `json:"cuid"`       // ucenter用户id
	Username  string `json:"username"`   // 账号 - 暂时不用
	Password  string `json:"password"`   // 密码 - 暂时不用
	Nickname  string `json:"nickname"`   // 昵称
	Avatar    string `json:"avatar"`     //头像
	Mobile    int    `json:"mobile"`     // 手机号
	MUserKey  string `json:"m_user_key"` // 用户Key
	DeletedAt string `json:"deleted_at"`
}

// 获取用户
func GetUsersInfoCuid(id int64) (*Users, error) {
	var user Users
	err := db.Where("cuid = ?", id).First(&user).Error
	if err != nil {
		return &Users{}, err
	}
	return &user, nil
}

func GetUsersInfo(id int64) (*Users, error) {
	var user Users
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return &Users{}, err
	}
	return &user, nil
}

// 获取UserKey
func GetMUserKey() string {
	k := utils.RandInt64(10000000, 99999999)
	var user Users
	err := db.Where("m_user_key = ?", k).First(&user).Error
	if err != nil {
		return strconv.FormatInt(k, 10)
	}
	return GetMUserKey()
}

func AddUsers(user *Users) error {
	user.Avatar = "https://cdn.ddgongjiang.com/041194a5705e6ff65287cfc0188b019f.png"
	user.Flag = 1
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
