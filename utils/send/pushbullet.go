package send

import (
	"Coot/error"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

/*发送pushBullet推送*/
func SendPushBullet(findAlertConfig []map[string]interface{}, args ...string) interface{} {
	infoArr := strings.Split(findAlertConfig[0]["info"].(string), "&&")
	var req *http.Request
	body := map[string]string{
		"title": args[0],
		"body":  args[1],
		"type":  "note", //默认消息类型
	}
	headers := map[string]string{
		"Access-Token": infoArr[0],
		"Content-Type": "application/json",
	}
	url := "https://api.pushbullet.com/v2/pushes"
	bodyJson, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		error.Check(err, "发送pushBullet创建req报错")
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	client := &http.Client{}
	result, err := client.Do(req)
	if err != nil {
		error.Check(err, "发送pushBullet返回resp报错")
	}
	return result
}
