package models

type PlatformUsers struct {
	Model
	Cuid        int    `json:"cuid"`         // 用户中心用户id
	PlatformKey string `json:"platform_key"` // 平台key
	PlatformUid string `json:"platform_uid"` // 平台用户id
}
