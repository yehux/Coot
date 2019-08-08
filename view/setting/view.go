package setting

import (
	"Coot/core/dbUtil"
	"Coot/core/job"
	"Coot/error"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func Html(c *gin.Context) {
	settingData := getSetting(1) //设置配置等
	alertData:=getSetting(0)   //通知配置等
	c.HTML(http.StatusOK, "setting.html", gin.H{
		"settingData": settingData,
		"alertData":alertData,
	})
}

/*获取配置*/
func getSetting(settingType int) []map[string]interface{} {
	sql := "select id,type,info,setting_name,setting_dis,update_time,status from coot_setting where setting_type=?"
	result := dbUtil.Query(sql,settingType)
	return result
}

/*检查是否配置信息*/
func checkInfo(id string) bool {
	sql := "select id,info,type from coot_setting where id = ?"
	result := dbUtil.Query(sql, id)
	info := result[0]["info"].(string)
	typeStr := result[0]["type"].(string)
	infoArr := strings.Split(info, "&&")
	num := len(infoArr)

	if num == 4 && typeStr == "mail" {
		return true
	}
	if num == 2 && typeStr == "login" {
		return true
	}
	if num == 2 && typeStr == "alertOver" {
		return true
	}
	if num == 1 && typeStr == "pushBullet" {
		return true
	}
	if num == 1 && typeStr == "fangTang" {
		return true
	}
	if num == 1 && typeStr == "logsOff" {
		return true
	}
	return false
}
func joinInfo(args ...string) string {
	and := "&&"
	info := ""
	for _, value := range args {
		if value == "" {
			return ""
		}
		info += value + and
	}
	info = info[:len(info)-2]
	return info
}

/*更新邮件通知*/
func UpdateEmailInfo(c *gin.Context) {
	email := c.PostForm("email")
	id := c.PostForm("id")
	pass := c.PostForm("pass")
	host := c.PostForm("host")
	port := c.PostForm("port")
	//subType := c.PostForm("type")
	info := joinInfo(host, port, email, pass)
	sql := `
		UPDATE  coot_setting 
		set	info = ?,
			status = ?,
			update_time = ?
		where id = ?;`
	dbUtil.Update(sql, info, 0, time.Now().Format("2006-01-02 15:04"), id)
	c.JSON(http.StatusOK, error.ErrSuccessNull())
}

/*更新登录用户信息*/
func UpdateLoginInfo(c *gin.Context) {
	loginName := c.PostForm("loginName")
	loginPwd := c.PostForm("loginPwd")
	id := c.PostForm("id")
	info := joinInfo(loginName, loginPwd)
	sql := `
		UPDATE  coot_setting 
		set	info = ?,
			status = ?,
			update_time = ?
		where id = ?;`
	dbUtil.Update(sql, info, 0, time.Now().Format("2006-01-02 15:04"), id)
	c.JSON(http.StatusOK, error.ErrSuccessNull())
}

/*更新alertOver推送配置*/
func UpdateAlertOverInfo(c *gin.Context) {
	alertSource := c.PostForm("alertSource")
	alertReceiver := c.PostForm("alertReceiver")
	id := c.PostForm("id")
	info := joinInfo(alertSource, alertReceiver)
	sql := `
		UPDATE  coot_setting 
		set	info = ?,
			update_time = ?
		where id = ?;`
	dbUtil.Update(sql, info, time.Now().Format("2006-01-02 15:04"), id)
	c.JSON(http.StatusOK, error.ErrSuccessNull())
}

/*更新pushBullet推送配置*/
func UpdatePushBulletInfo(c *gin.Context) {
	pushBulletToken := c.PostForm("pushBulletToken")
	id := c.PostForm("id")
	info := joinInfo(pushBulletToken)
	sql := `
		UPDATE  coot_setting 
		set	info = ?,
			status = ?,
			update_time = ?
		where id = ?;`
	dbUtil.Update(sql, info, 0, time.Now().Format("2006-01-02 15:04"), id)
	c.JSON(http.StatusOK, error.ErrSuccessNull())
}

/*更新pushFangTang推送配置*/
func UpdatePushFangTangInfo(c *gin.Context) {
	pushBulletToken := c.PostForm("pushBulletToken")
	id := c.PostForm("id")
	info := joinInfo(pushBulletToken)
	sql := `
		UPDATE  coot_setting 
		set	info = ?,
			status = ?,
			update_time = ?
		where id = ?;`
	dbUtil.Update(sql, info, 0, time.Now().Format("2006-01-02 15:04"), id)
	c.JSON(http.StatusOK, error.ErrSuccessNull())
}

/*更新设置状态*/
func UpdateStatusSetting(c *gin.Context) {
	id := c.PostForm("id")
	status := c.PostForm("status")

	if !checkInfo(id) && status == "1" {
		c.JSON(http.StatusOK, gin.H{"code": 10003, "msg": "请配置后在启用", "data": nil})
		return
	}
	sql := `update coot_setting
		set status = ?,
			update_time=?
		where id = ?`
	dbUtil.Update(sql, status, time.Now().Format("2006-01-02 15:04"), id)
	c.JSON(http.StatusOK, error.ErrSuccessNull())
}

/*是否启用日志功能*/
func UpdateLogStatusSetting(c *gin.Context) {
	id := c.PostForm("id")
	status := c.PostForm("status")
	if !checkInfo(id) && status == "1" {
		c.JSON(http.StatusOK, gin.H{"code": 10003, "msg": "请配置后在启用", "data": nil})
		return
	}
	sql := `update coot_setting
		set status = ?,
			update_time=?
		where id = ?`
	dbUtil.Update(sql, status, time.Now().Format("2006-01-02 15:04"), id)
	job.LogOff=!job.LogOff
	c.JSON(http.StatusOK, error.ErrSuccessNull())
}
/*初始日志开关*/
func InitLogsOff(){
	sql := `select id,status from coot_setting where [type] = 'logsOff'`
	result := dbUtil.Query(sql)
	status := result[0]["status"].(int64)
	if status==1{
		job.LogOff=true
	}
}
/*根据id获取设置详情*/
func GetSettingInfo(c *gin.Context) {
	id, _ := c.GetQuery("id")
	sql := `select id,type,info,status from coot_setting where id = ?`
	result := dbUtil.Query(sql, id)
	c.JSON(http.StatusOK, error.ErrSuccess(result))
}
