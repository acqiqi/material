package models

import (
	"material/lib/utils"
	"strconv"
)

type Users struct {
	Model
	Cuid     int64  `json:"cuid"`       // ucenter用户id
	Username string `json:"username"`   // 账号 - 暂时不用
	Password string `json:"password"`   // 密码 - 暂时不用
	Nickname string `json:"nickname"`   // 昵称
	Avatar   string `json:"avatar"`     //头像
	Mobile   int    `json:"mobile"`     // 手机号
	MUserKey string `json:"m_user_key"` // 用户Key

	// 目前要处理的三方平台的存储信息
	DUserKey string `json:"d_user_key"` //绑定三方账号信息

}

// 查询DD的UserKey
func GetUsersInfoDD(d_user_key string) (*Users, error) {
	var user Users
	err := db.Where("d_user_key = ? AND flag =1", d_user_key).First(&user).Error
	if err != nil {
		return &Users{}, err
	}
	return &user, nil
}

// 获取用户
func GetUsersInfoCuid(id int64) (*Users, error) {
	var user Users
	err := db.Where("cuid = ? AND flag =1", id).First(&user).Error
	if err != nil {
		return &Users{}, err
	}
	return &user, nil
}

func GetUsersInfo(id int64) (*Users, error) {
	var user Users
	err := db.Where("id = ? AND flag =1", id).First(&user).Error
	if err != nil {
		return &Users{}, err
	}
	return &user, nil
}

// 获取UserKey
func GetMUserKey() string {
	k := utils.RandInt64(10000000, 99999999)
	var user Users
	err := db.Where("m_user_key = ? AND flag =1", k).First(&user).Error
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
