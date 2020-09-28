package e

import (
	"material/lib/utils"
	"material/models"
)

const (
	PLATFORM_ACTION_PROJECT_RECEIVE = "PROJECT_RECEIVE" //项目接收
)

// 项目接收回调结构体
type PlatformProjectReceiveCallback struct {
	Id          int64          `json:"id"`
	ProjectName string         `json:"project_name"` // 项目名称
	State       int            `json:"state"`        // 0 在建中 1已完成
	CompanyId   int64          `json:"company_id"`   //关联企业
	Company     models.Company `json:"company"`      //关联企业
	BindState   int            `json:"bind_state"`   // 绑定工程账号的状态,手动新建的项目默认已绑定.未绑定,待处理,已绑定
	PlatformKey string         `json:"platform_key"` // 平台key
	PlatformUid string         `json:"platform_uid"` // 平台用户id
	PlatformId  string         `json:"platform_id"`  // 平台用户id
	CreatedAt   utils.Time     `json:"created_at"`
	Status      int            `json:"status"` // 0未同步 1已同步
}
