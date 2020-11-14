package models

// 平台
type Platform struct {
	Model
	PlatformName       string `json:"platform_name"` // 平台名称
	PlatformKey        string `json:"platform_key"`  // 平台key 完全标识
	PlatformSk         string `json:"platform_sk"`
	PlatformSecret     string `json:"platform_secret"` // 平台secret
	Status             int    `json:"status"`          // 状态 1运行 0暂停维护 -1禁用
	PayName            string `json:"pay_name"`
	PayAk              string `json:"pay_ak"`
	PaySk              string `json:"pay_sk"`
	Ak                 string `json:"ak"`
	Sk                 string `json:"sk"`
	PayNotifyUrl       string `json:"pay_notify_url"` // 支付回调
	PayNotifyFunc      string `json:"pay_notify_func"`
	MessageCallbackUrl string `json:"message_callback_url"` // 消息中心回调地址
}

// 获取平台详情
func PlatformGetInfoOrKey(key string) (*Platform, error) {
	var platform Platform
	err := db.Where("platform_key = ? AND flag =1", key).First(&platform).Error
	if err != nil {
		return &Platform{}, err
	}
	return &platform, nil
}

// 利用Ak 获取平台详情
func PlatformGetInfoOrAk(key string) (*Platform, error) {
	var platform Platform
	err := db.Where("ak = ? AND flag =1", key).First(&platform).Error
	if err != nil {
		return &Platform{}, err
	}
	return &platform, nil
}
