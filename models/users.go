package models

type Users struct {
	Model
	Cuid      int    `json:"cuid"`       // ucenter用户id
	Username  string `json:"username"`   // 账号 - 暂时不用
	Password  string `json:"password"`   // 密码 - 暂时不用
	Nickname  string `json:"nickname"`   // 昵称
	Mobile    int    `json:"mobile"`     // 手机号
	MUserKey  string `json:"m_user_key"` // 用户Key
	DeletedAt string `json:"deleted_at"`
}
