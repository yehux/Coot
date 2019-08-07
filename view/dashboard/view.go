package dashboard

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"Coot/core/dbUtil"
	"strconv"
)

func Html(c *gin.Context) {
	sql1 := `select count(1) as sum from coot_tasks;`
	sql2 := `select count(1) as sum from coot_tasks where task_id!="";`
	sql3 := `select count(1) as sum from coot_tasks where task_id="";`
	sql4 := `select count(1) as sum from coot_tasks where last_exec_time="";`

	result1 := dbUtil.Query(sql1)
	result2 := dbUtil.Query(sql2)
	result3 := dbUtil.Query(sql3)
	result4 := dbUtil.Query(sql4)

	totalSum := strconv.FormatInt(result1[0]["sum"].(int64), 10)
	startSum := strconv.FormatInt(result2[0]["sum"].(int64), 10)
	stopSum := strconv.FormatInt(result3[0]["sum"].(int64), 10)
	noExecuSum := strconv.FormatInt(result4[0]["sum"].(int64), 10)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"totalSum":   totalSum,
		"startSum":   startSum,
		"stopSum":    stopSum,
		"noExecuSum": noExecuSum,
	})
}
