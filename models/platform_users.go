package models

type PlatformUsers struct {
	Model
	Cuid        int    `json:"cuid"`         // 用户中心用户id
	PlatformKey string `json:"platform_key"` // 平台key
	PlatformUid string `json:"platform_uid"` // 平台用户id
}

// 检测平台用户是否存在
func PlatformUsersCheckUser(platform_uid, platform_key string) (*Users, error) {
	var user Users
	err := db.Where("platform_uid = ? AND platform_key = ? AND flag =1", platform_uid, platform_key).First(&user).Error
	if err != nil {
		return &Users{}, err
	}
	return &user, nil
}
