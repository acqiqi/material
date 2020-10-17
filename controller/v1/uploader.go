package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"material/lib/e"
	"material/lib/setting"
)

func Uploader(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")
	fileSize := fileHeader.Size

	putPolicy := storage.PutPolicy{
		Scope: setting.QiniuSetting.Bucket,
	}
	mac := qbox.NewMac(setting.QiniuSetting.AccessKey, setting.QiniuSetting.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		e.ApiErr(c, "上传失败"+err.Error())
		return
	}

	fmt.Println(ret.Key, ret.Hash)
	url := setting.QiniuSetting.Domain + ret.Key
	e.ApiOk(c, "上传成功", struct {
		Url string `json:"url"`
	}{
		Url: url,
	})
}
