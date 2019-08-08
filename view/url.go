package view

import (
	"Coot/view/dashboard"
	"Coot/view/login"
	"Coot/view/plug"
	"Coot/view/report"
	"Coot/view/setting"
	"Coot/view/task"
	"github.com/gin-gonic/gin"
)

func LoadUrl(r *gin.Engine) {
	r.GET("/login", login.Html)
	r.POST("/login", login.Login)

	// 注销
	r.GET("/logout", login.Logout)

	// 仪表盘
	r.GET("/", login.Jump, dashboard.Html)
	r.GET("/dashboard", login.Jump, dashboard.Html)

	// 任务
	r.GET("/task", login.Jump, task.Html)
	r.GET("/task/detail", login.Jump, task.HtmlDetail)
	r.GET("/task/add", login.Jump, task.HtmlAdd)
	r.GET("/get/task/list", login.Jump, task.GetTaskList)
	r.POST("/post/task/add", login.Jump, task.PostTaskAdd)
	r.POST("/post/task/update", login.Jump, task.PostTaskUpdate)
	r.POST("/post/task/del", login.Jump, task.PostTaskDel)
	r.POST("/task/start", login.Jump, task.TaskStart)

	r.POST("/task/stop", login.Jump, task.TaskStop)

	// 插件
	r.GET("/plugs", login.Jump, plug.Html)

	//日志报告
	r.GET("/report",login.Jump,report.Html)
	r.GET("/get/report/data",login.Jump,report.GetLogs)
	r.GET("/get/report/getNewLogs",login.Jump,report.GetNewLogs)
	r.POST("/post/report/delete",login.Jump,report.DeleteLogsAll)

	// 设置
	r.GET("/setting", login.Jump, setting.Html)
	r.GET("/get/setting/info", login.Jump, setting.GetSettingInfo)
	r.POST("/post/setting/update", login.Jump, setting.UpdateEmailInfo)
	r.POST("/post/setting/login", login.Jump, setting.UpdateLoginInfo)
	r.POST("/post/setting/alertOver", login.Jump, setting.UpdateAlertOverInfo)
	r.POST("/post/setting/pushBullet", login.Jump, setting.UpdatePushBulletInfo)
	r.POST("/post/setting/pushFangTang", login.Jump, setting.UpdatePushFangTangInfo)
	r.POST("/post/setting/checkSetting", login.Jump, setting.UpdateStatusSetting)
	r.POST("/post/setting/checkLogSetting", login.Jump, setting.UpdateLogStatusSetting)
}
