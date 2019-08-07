package send

import (
	"Coot/error"
	"net/http"
	"strings"
)

/*发送 方糖 Server酱 推送*/
func SendFangTang(findAlertConfig []map[string]interface{}, args ...string) interface{} {
	infoArr := strings.Split(findAlertConfig[0]["info"].(string), "&&")
	var req *http.Request
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	url := "https://sc.ftqq.com/" + strings.TrimSpace(infoArr[0]) + ".send"
	body := "text=" + args[0] + "&desp=" + args[1]
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		error.Check(err, "发送方糖推送创建req报错")
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
		error.Check(err, "发送方糖推送返回resp报错")
	}
	return result
}
