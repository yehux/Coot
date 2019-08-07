package task

import (
	"Coot/core/dbUtil"
	"Coot/core/job"
	"Coot/error"
	"Coot/utils/file"
	"Coot/utils/md5"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
	"strings"
	"time"
)

type alert struct {
	isAlertMail         string
	isAlertAlertOver    string
	isAlertPushBullet   string
	isAlertPushFangTang string
}

// Task List 页面
func Html(c *gin.Context) {
	c.HTML(http.StatusOK, "task.html", gin.H{})
}

// Task详情 页面
func HtmlDetail(c *gin.Context) {
	id, _ := c.GetQuery("id")

	sql := `select id,task_name,task_explain,task_id,task_time_type,task_time,last_exec_time,script_type,script_path,alert_type,alert_rec_mail,create_time from coot_tasks WHERE id = ?;`
	result := dbUtil.Query(sql, id)

	taskName := result[0]["task_name"].(string)
	taskExplain := result[0]["task_explain"].(string)
	taskTimeType := result[0]["task_time_type"].(string)
	taskTime := result[0]["task_time"].(string)
	scriptType := result[0]["script_type"].(string)
	scriptPath := result[0]["script_path"].(string)
	alertType := result[0]["alert_type"].(string)
	alertRecMail := result[0]["alert_rec_mail"].(string)

	code := file.ReadFile(scriptPath)

	arr := strings.Split(alertType, ",")

	a := alert{
		"0",
		"0",
		"0",
		"0",
	}

	if len(arr) > 0 {
		for _, v := range arr {
			if v == "mail" {
				a.isAlertMail = "1"
			}

			// 判断是否开启 alertOver 通知
			if v == "alertOver" {
				a.isAlertAlertOver = "1"
			}

			// 判断是否开启 pushBullet 通知
			if v == "pushBullet" {
				a.isAlertPushBullet = "1"
			}

			if v == "fangtang" {
				a.isAlertPushFangTang = "1"
			}
		}
	}

	c.HTML(http.StatusOK, "taskDetail.html", gin.H{
		"id":                  id,
		"taskName":            taskName,
		"taskExplain":         taskExplain,
		"taskTimeType":        taskTimeType,
		"taskTime":            taskTime,
		"scriptType":          scriptType,
		"code":                code,
		"isAlertMail":         a.isAlertMail,
		"isAlertAlertOver":    a.isAlertAlertOver,
		"isAlertPushBullet":   a.isAlertPushBullet,
		"isAlertPushFangTang": a.isAlertPushFangTang,
		"alertRecMail":        alertRecMail,
	})
}

// 查询任务列表
func GetTaskList(c *gin.Context) {
	sql := `select id,task_name,task_explain,task_id,task_time_type,task_time,last_exec_time,script_type,script_path,create_time from coot_tasks ORDER BY id desc;`
	result := dbUtil.Query(sql)
	c.JSON(http.StatusOK, error.ErrSuccess(result))
}

// Task Add 页面
func HtmlAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "taskAdd.html", gin.H{})
}

// 添加任务
func PostTaskAdd(c *gin.Context) {
	isPlugScript := c.PostForm("is_plug_script")
	taskName := c.PostForm("taskName")
	taskExplain := c.PostForm("taskExplain")
	taskTimeType := c.PostForm("taskTimeType")
	taskTime := c.PostForm("taskTime")
	taskLanuage := c.PostForm("taskLanuage")
	isAlert := c.PostForm("is_alert")
	mailList := c.PostForm("mail_list")
	code := c.PostForm("code")

	// 获取时间戳，生成MD5
	currTimeStr := time.Now().Format("2006-01-02 15:04")
	uid := uuid.NewV4()
	fileName := md5.Md5(currTimeStr + uid.String())

	var fileType = ""

	if taskLanuage == "Python" {
		fileType = "py"
	} else if taskLanuage == "Shell" {
		fileType = "sh"
	} else {
		c.JSON(http.StatusOK, error.ErrFailFileType())
	}

	// 写入文件
	filePath := "./scripts/" + fileName + "." + fileType
	file.Output(code, filePath)

	sql := `
		INSERT INTO coot_tasks (
			task_name,
			task_explain,
			task_id,
			task_time_type,
			task_time,
			last_exec_time,
			is_plug_script,
			script_type,
			script_path,
			alert_type,
			alert_rec_mail,
			create_time
		)
		VALUES
			(?,?,?,?,?,?,?,?,?,?,?,?);
	`

	dbUtil.Insert(sql, taskName, taskExplain, "", taskTimeType, taskTime, "", isPlugScript, taskLanuage, filePath, isAlert, mailList, currTimeStr)

	c.JSON(http.StatusOK, error.ErrSuccessNull())
}

// 更新任务
func PostTaskUpdate(c *gin.Context) {
	id := c.PostForm("id")
	taskName := c.PostForm("taskName")
	taskExplain := c.PostForm("taskExplain")
	taskTimeType := c.PostForm("taskTimeType")
	taskTime := c.PostForm("taskTime")
	taskLanuage := c.PostForm("taskLanuage")
	isAlert := c.PostForm("is_alert")
	mailList := c.PostForm("mail_list")
	code := c.PostForm("code")

	// 获取时间戳，生成MD5
	currTimeStr := time.Now().Format("2006-01-02 15:04")
	uid := uuid.NewV4()
	fileName := md5.Md5(currTimeStr + uid.String())

	var fileType = ""

	if taskLanuage == "Python" {
		fileType = "py"
	} else if taskLanuage == "Shell" {
		fileType = "sh"
	} else {
		c.JSON(http.StatusOK, error.ErrFailFileType())
	}

	// 写入文件
	filePath := "./scripts/" + fileName + "." + fileType
	file.Output(code, filePath)

	sql := `
		UPDATE coot_tasks
		SET task_name = ?,
		 task_explain = ?,
		 task_time_type = ?,
		 task_time = ?,
		 script_type = ?,
		 script_path = ?,
		 alert_type = ?,
		 alert_rec_mail = ?
		WHERE
			id = ?;
	`

	dbUtil.Update(sql, taskName, taskExplain, taskTimeType, taskTime, taskLanuage, filePath, isAlert, mailList, id)

	stop(id)

	c.JSON(http.StatusOK, error.ErrSuccessNull())
}

// 删除任务
func PostTaskDel(c *gin.Context) {
	id := c.PostForm("id")

	sql := `select id,task_name,task_explain,task_id,task_time_type,task_time,last_exec_time,script_type,script_path,create_time from coot_tasks WHERE id = ?;`
	result := dbUtil.Query(sql, id)

	taskId := result[0]["task_id"]

	job.StopJob(taskId.(string))

	sqlDel := `delete from coot_tasks where id=?;`
	dbUtil.Delete(sqlDel, id)

	c.JSON(http.StatusOK, error.ErrSuccessNull())
}

func updateTaskId(task_id string, id string) {
	sql := `
		UPDATE coot_tasks
		SET task_id = ?
		WHERE
			id = ?;
		`
	dbUtil.Update(sql, task_id, id)
}

func UpdateTaskAll() {
	sql := `
		UPDATE coot_tasks
		SET task_id = "";
		`
	dbUtil.Update(sql)
}

func start(id string) {
	sql := `select id,task_name,task_explain,task_id,task_time_type,task_time,last_exec_time,script_type,script_path,alert_type,alert_rec_mail,create_time from coot_tasks WHERE id = ?;`
	result := dbUtil.Query(sql, id)

	taskName := result[0]["task_name"]
	taskTimeType := result[0]["task_time_type"]
	taskTime := result[0]["task_time"]
	scriptType := result[0]["script_type"]
	scriptPath := result[0]["script_path"]
	alertType := result[0]["alert_type"]
	alertRecMail := result[0]["alert_rec_mail"]

	// 启动任务
	taskId := job.AddJob(&job.Task{
		id,
		taskName.(string),
		"",
		taskTimeType.(string),
		taskTime.(string),
		scriptType.(string),
		scriptPath.(string),
		alertType.(string),
		alertRecMail.(string),
	})

	// 更新数据库
	updateTaskId(taskId, id)
}

// 启动任务
func TaskStart(c *gin.Context) {
	id := c.PostForm("id")
	start(id)
	c.JSON(http.StatusOK, error.ErrSuccessNull())
}

func stop(id string) {
	sql := `select task_id from coot_tasks WHERE id = ?;`
	result := dbUtil.Query(sql, id)

	taskId := result[0]["task_id"]

	// 停止任务
	job.StopJob(taskId.(string))

	// 更新数据库
	updateTaskId("", id)

}

// 关闭任务
func TaskStop(c *gin.Context) {
	id := c.PostForm("id")
	stop(id)
	c.JSON(http.StatusOK, error.ErrSuccessNull())
}

// 重启任务
func taskRestart(id string) {
	stop(id)
	start(id)
}
