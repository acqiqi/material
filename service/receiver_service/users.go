package receiver_service

import (
	"errors"
	"material/lib/dd"
	"material/models"
)

type ReceiverUsers struct {
	Id         int64 `json:"id"`
	CompanyId  int   `json:"company_id"`
	ContractId int   `json:"contract_id"`
	Cuid       int   `json:"cuid"`

	PlatformKey string `json:"platform_key"`
	PlatformUid string `json:"platform_uid"`
}

type UserAdd struct {
	PlatformUid string `json:"platform_uid"`
	Mobile      string `json:"mobile"`
	Nickname    string `json:"nickname"`
}

type ReceiverUsersCallback struct {
	Id          int64  `json:"id"`
	CompanyId   int64  `json:"company_id"`
	ProjectId   int64  `json:"project_id"`
	Cuid        int64  `json:"cuid"`
	PlatformKey string `json:"platform_key"`
	PlatformUid string `json:"platform_uid"`
	Nickname    string `json:"nickname"`
	Mobile      string `json:"mobile"`
}

// 同步用户
func SyncUsers(users []UserAdd, project *models.Project, platform_key string) ([]ReceiverUsersCallback, error) {
	dd_utils := new(dd.UCUtils)
	if _, err := dd_utils.GetToken(); err != nil {
		return []ReceiverUsersCallback{}, err
	}
	cb_users := []ReceiverUsersCallback{}
	//第一层检测
	for _, v := range users {
		if v.PlatformUid != "" {
			if v.Nickname == "" {
				continue //跳过不操作
			}
			// 检测三方是否注册
			cloud_user, err := dd_utils.GetUserInfoOrMobile(v.Mobile)
			if err != nil {
				//注册
				cloud_user, err = dd_utils.UserMobileReg(v.Mobile)
				if err != nil {
					continue //直接跳出当前循环，业务不成立
				}
			}
			// 查询是否在项目中注册
			user, err := models.ReceiverUsersGetInfoByPlatform(platform_key, v.PlatformUid)
			if err != nil {
				//直接注册 奥利给
				user, err = MobileReg(models.ReceiverUsers{
					CompanyId:   project.CompanyId,
					ProjectId:   project.Id,
					Cuid:        cloud_user.Data.UserInfo.Id,
					PlatformKey: platform_key,
					PlatformUid: v.PlatformUid,
					Nickname:    v.Nickname,
					Mobile:      v.Mobile,
				})
				if err != nil {
					continue //直接跳出业务不成立
				}
			}
			//查询本地是否注册
			_, err = models.GetUsersInfoCuid(cloud_user.Data.UserInfo.Id)
			if err != nil {
				//直接注册
				user_model := models.Users{
					Cuid:     cloud_user.Data.UserInfo.Id,
					Nickname: cloud_user.Data.UserInfo.Nickname,
					Avatar:   cloud_user.Data.UserInfo.Avatar,
					MUserKey: models.GetMUserKey(),
				}
				models.AddUsers(&user_model)
			}

			cb_users = append(cb_users, ReceiverUsersCallback{
				Id:          user.Id,
				CompanyId:   user.CompanyId,
				ProjectId:   user.ProjectId,
				Cuid:        user.Cuid,
				PlatformKey: user.PlatformKey,
				PlatformUid: user.PlatformUid,
				Nickname:    user.Nickname,
				Mobile:      user.Mobile,
			})
		} // 不注册
	}
	return cb_users, nil
}

func MobileReg(m models.ReceiverUsers) (*models.ReceiverUsers, error) {
	if err := models.ReceiverUsersAdd(&m); err != nil {
		return nil, errors.New("用户注册失败")
	}
	cb, err := models.ReceiverUsersGetInfo(m.Id)
	if err != nil {
		return nil, err
	}
	return cb, nil
}
