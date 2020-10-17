package small

import (
	"github.com/gin-gonic/gin"
	"log"
	"material/lib/dd"
	"material/lib/gredis"
	"material/models"
	"strconv"

	"material/lib/e"
)

// JWT is jwt middleware
func Small() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			//e.ApiOpt(c, e.INVALID_PARAMS, e.GetMsg(e.INVALID_PARAMS), e.GetEmptyStruct())
			e.ApiOpt(c, e.API_NOT_AUTH_CODE, e.GetMsg(e.API_NOT_AUTH_CODE), e.GetEmptyStruct())
			return
		} else {
			id := gredis.GetCacheString(token)
			if id == "" { //没有获取到token
				//如果token 没有存储则查询一次三方
				dd := new(dd.UCUtils)
				cb, err := dd.UserGetInfo(token)
				if err != nil {
					e.ApiOpt(c, cb.Code, cb.Msg, e.GetEmptyStruct())
					return
				}
				if cb.Code != 0 {
					e.ApiOpt(c, cb.Code, cb.Msg, e.GetEmptyStruct())
					return
				}
				id = strconv.FormatInt(cb.Data.UserInfo.Id, 10)
				//查询是否注册
				_, err = models.GetUsersInfoCuid(cb.Data.UserInfo.Id)
				if err != nil {
					//产生注册
					user_model := models.Users{
						Cuid:     cb.Data.UserInfo.Id,
						Nickname: cb.Data.UserInfo.Nickname,
						Avatar:   cb.Data.UserInfo.Avatar,
						MUserKey: models.GetMUserKey(),
					}
					models.AddUsers(&user_model)
				}
				gredis.SetCacheString(token, id, 60*60*24*30)
			}

			user_id, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				e.ApiOpt(c, e.API_NOT_AUTH_CODE, e.GetMsg(e.API_NOT_AUTH_CODE), e.GetEmptyStruct())
				return
			}
			user, err := models.GetUsersInfoCuid(user_id)
			if err != nil {
				log.Println(err)
				//e.ApiOpt(c, e.ERROR_AUTH, e.GetMsg(e.ERROR_AUTH), e.GetEmptyStruct())
				e.ApiOpt(c, e.API_NOT_AUTH_CODE, e.GetMsg(e.API_NOT_AUTH_CODE), e.GetEmptyStruct())
				return
			}

			c.Set("user_info", *user)
		}
		c.Next()
	}
}
