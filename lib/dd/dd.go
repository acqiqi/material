package dd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"material/lib/utils"
	"net/http"
)

type DDUtils struct {
	PlatformKey    string `json:"platform_key"`
	PlatformSecret string `json:"platform_secret"`
	Token          string `json:"token"`
}

const (
	DD_UTILS_BASE_URL  = "http://www.luakit.com:81/"
	DD_MY_PLATFORM_KEY = "material"
)

func (d *DDUtils) Setup() {

}

func (d *DDUtils) GetToken() {

}

// Http Post Json 请求
func httpPostJson(url string, body interface{}, header map[string]string, cb interface{}) error {

	requestBody := utils.JsonEncode(body)
	var jsonStr = []byte(requestBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for i, v := range header {
		req.Header.Set(i, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New(string(b))
	}
	log.Println(string(b))
	utils.JsonDecode(string(b), &cb)
	return nil
}
