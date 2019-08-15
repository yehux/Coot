package report

import (
	"Coot/core/dbUtil"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/*清空日志*/
func DeleteLogsAll(c *gin.Context) {
	sql := `
		Delete from coot_logs;
		`
	dbUtil.Delete(sql)
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"success",
		"data":nil,
	})
}
/*主页面*/
func Html(c *gin.Context){
	sql1 := `select count(1) as sum from coot_tasks;`
	sql2 := `select count(1) as sum from coot_tasks where task_id!="";`
	sql3 := `select count(1) as sum from coot_tasks where task_id="";`
	sql4 := `select count(1) as sum from coot_tasks where last_exec_time="";`
	sql5 := `select id,task_id,task_name,content,cmd,time_type,status,created_at from coot_logs order by id desc limit(8)`

	result1 := dbUtil.Query(sql1)
	result2 := dbUtil.Query(sql2)
	result3 := dbUtil.Query(sql3)
	result4 := dbUtil.Query(sql4)
	result5 := dbUtil.Query(sql5)
	totalSum := strconv.FormatInt(result1[0]["sum"].(int64), 10)
	startSum := strconv.FormatInt(result2[0]["sum"].(int64), 10)
	stopSum := strconv.FormatInt(result3[0]["sum"].(int64), 10)
	noExecuSum := strconv.FormatInt(result4[0]["sum"].(int64), 10)

	c.HTML(http.StatusOK, "report.html", gin.H{
		"totalSum":   totalSum,
		"startSum":   startSum,
		"stopSum":    stopSum,
		"noExecuSum": noExecuSum,
		"tenExecuData":result5,
	})
}
/*获取数据*/
func GetLogs(c *gin.Context){
	wait:=`select task_name,count([id]) as count,created_at  from [coot_logs]  where status = 0 group by [task_id]`
	okSql:=`select task_name,count([id]) as count,created_at  from [coot_logs]  where status = 1 group by [task_id]`
	failSql:=`select task_name,count([id]) as count,created_at  from [coot_logs]  where status = -1 group by [task_id]`
	waitLogs:=dbUtil.Query(wait)
	okLogs:=dbUtil.Query(okSql)
	failLogs:=dbUtil.Query(failSql)
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"success",
		"waitLogs":waitLogs,
		"okLogs":okLogs,
		"failLogs":failLogs,
	})
}
/*获取最新日志*/
func GetNewLogs(c *gin.Context){
	id:=c.Query("id")
	sql := `select id,task_id,task_name,content,cmd,time_type,status,created_at from coot_logs where id > ? order by id `
	result:=dbUtil.Query(sql,id)
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"success",
		"data":result,
	})
}