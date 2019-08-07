package send

import (
	"Coot/error"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type AlertOverModel struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func SendAlertOver(findAlertConfig []map[string]interface{}, args ...string) interface{} {
	infoArr := strings.Split(findAlertConfig[0]["info"].(string), "&&")
	resp, err := http.PostForm("https://api.alertover.com/v1/alert",
		url.Values{
			"source":   {infoArr[0]},
			"receiver": {infoArr[1]},
			"title":    {args[0]},
			"content":  {args[1]},
		})

	if err != nil {
		error.Check(err, "发送alertOver通知失败1")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		error.Check(err, "发送alertOver通知失败2")
	}

	_alertOverModel := AlertOverModel{}
	result := json.Unmarshal(body, &_alertOverModel)
	return result
}
