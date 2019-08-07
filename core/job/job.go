package job

import (
	"Coot/core/dbUtil"
	"Coot/core/exec"
	"Coot/error"
	"Coot/utils/send"
	"fmt"
	"github.com/domgoer/gotask"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	/*
	 * Id         	数据库ID
	 * Name		  	任务名称
	 * TaskId     	任务ID  添加的时候为空
	 * TimeType   	执行类型 1 秒执行，2 分钟执行，3 小时执行 ，4 每天指定时间执行，5 每月指定天和时间执行，6 年执行
	 * Time   	  	周期时间
	 * ScriptType 	脚本语言
	 * ScriptPath 	脚本路径
	 */
	Id           string
	Name         string
	TaskId       string
	TimeType     string
	Time         string
	ScriptType   string
	ScriptPath   string
	AlertType    string
	AlertRecMail string
}

func updateExecTime(id string) {
	sql := `
		UPDATE coot_tasks
		SET last_exec_time = ?
		WHERE
			id = ?;
		`

	currTimeStr := time.Now().Format("2006-01-02 15:04:05")
	dbUtil.Update(sql, currTimeStr, id)
}

// 执行任务
func execute(t *Task) {
	var id = t.Id
	var cmd string

	// 拼接命令
	if t.ScriptType == "Python" {
		cmd = "python " + t.ScriptPath
	} else if t.ScriptType == "Shell" {
		cmd = "sh " + t.ScriptPath
	}

	// 开始执行任务
	fmt.Fprintln(gin.DefaultWriter, time.Now().Format("2006-01-02 15:04:05")+" 开始执行 id:"+id+" 任务名称:"+t.Name+" 执行命令:"+cmd)
	result, err := exec.Execute(cmd)

	if err != nil {
		fmt.Fprintln(gin.DefaultWriter, time.Now().Format("2006-01-02 15:04:05")+" 执行失败 id:"+id+" 任务名称:"+t.Name+" 脚本结果:", err)
	} else {
		fmt.Fprintln(gin.DefaultWriter, time.Now().Format("2006-01-02 15:04:05")+" 执行成功 id:"+id+" 任务名称:"+t.Name+" 脚本结果:"+result)
		//执行通知
		go notice(t, result)
	}
	// 更新任务执行时间
	updateExecTime(id)
}

// 消息通知
func notice(t *Task, result string) {

	// AlertType 格式，mail,pushBullet,alertOver
	arr := strings.Split(t.AlertType, ",")

	if len(arr) > 0 {
		for _, v := range arr {
			// 判断是否开启邮箱通知
			if v == "mail" {
				sql := `select status,info from coot_setting where type="mail";`
				isAlertStatus := dbUtil.Query(sql)

				status := strconv.FormatInt(isAlertStatus[0]["status"].(int64), 10)

				// 判断总开关是否开启
				if status == "1" {
					r := strings.Split(result, "&&")

					// 判断脚本 code 是否 为 0
					if r[0] == "0" {
						recList := strings.Split(t.AlertRecMail, ",")
						send.SendMail(recList, "Coot["+t.Name+"]提醒你", r[1], isAlertStatus)
					}
				}
			}

			// 判断是否开启 alertOver 通知
			if v == "alertOver" {
				sql := `select status,info from coot_setting where type="alertOver";`
				isAlertStatus := dbUtil.Query(sql)
				status := strconv.FormatInt(isAlertStatus[0]["status"].(int64), 10)

				// 判断总开关是否开启
				if status == "1" {
					r := strings.Split(result, "&&")

					// 判断脚本 code 是否 为 0
					if r[0] == "0" {
						send.SendAlertOver(isAlertStatus, "Coot["+t.Name+"]提醒你", r[1])
					}
				}
			}

			// 判断是否开启 pushBullet 通知
			if v == "pushBullet" {
				sql := `select status,info from coot_setting where type="pushBullet";`
				isAlertStatus := dbUtil.Query(sql)

				status := strconv.FormatInt(isAlertStatus[0]["status"].(int64), 10)

				// 判断总开关是否开启
				if status == "1" {
					r := strings.Split(result, "&&")

					// 判断脚本 code 是否 为 0
					if r[0] == "0" {
						send.SendPushBullet(isAlertStatus, "Coot["+t.Name+"]提醒你", r[1])
					}
				}
			}

			// 判断是否开启 方糖 通知
			if v == "fangTang" {
				sql := `select status,info from coot_setting where type="fangtang";`
				isAlertStatus := dbUtil.Query(sql)

				status := strconv.FormatInt(isAlertStatus[0]["status"].(int64), 10)

				// 判断总开关是否开启
				if status == "1" {
					r := strings.Split(result, "&&")

					// 判断脚本 code 是否 为 0
					if r[0] == "0" {
						send.SendFangTang(isAlertStatus, "Coot["+t.Name+"]提醒你", r[1])
					}
				}
			}
		}
	}
}

func mTask(t *Task, typs string) string {
	var taskId string

	// 创建任务
	switch t.TimeType {
	case "1":
		// 秒执行
		number, err := strconv.Atoi(t.Time)
		error.Check(err, "秒时间格式化失败")

		if typs == "add" {
			task := gotask.NewTask(time.Second*time.Duration(number), func() { execute(t) })
			gotask.AddToTaskList(task)
			taskId = task.ID()
		} else if typs == "update" {
			gotask.ChangeInterval(t.TaskId, time.Second*time.Duration(number))
		}
	case "2":
		// 分钟执行
		number, err := strconv.Atoi(t.Time)
		error.Check(err, "分钟时间格式化失败")

		if typs == "add" {
			task := gotask.NewTask(time.Minute*time.Duration(number), func() { execute(t) })
			gotask.AddToTaskList(task)
			taskId = task.ID()
		} else if typs == "update" {
			gotask.ChangeInterval(t.TaskId, time.Minute*time.Duration(number))
		}
	case "3":
		// 小时执行
		number, err := strconv.Atoi(t.Time)
		error.Check(err, "小时时间格式化失败")

		if typs == "add" {
			task := gotask.NewTask(time.Hour*time.Duration(number), func() { execute(t) })
			gotask.AddToTaskList(task)
			taskId = task.ID()
		} else if typs == "update" {
			gotask.ChangeInterval(t.TaskId, time.Hour*time.Duration(number))
		}
	case "4":
		// 天执行
		task, err := gotask.NewDayTask(t.Time, func() { execute(t) })
		error.Check(err, "")
		gotask.AddToTaskList(task)
		taskId = task.ID()
	case "5":
		// 月执行
		task, err := gotask.NewMonthTask(t.Time, func() { execute(t) })
		error.Check(err, "")
		taskId = task.ID()
	case "6":
		// 年执行
		task := gotask.NewTask(time.Second*2, func() { execute(t) })
		gotask.AddToTaskList(task)
		taskId = task.ID()
	}
	return taskId
}

// 创建定时任务
func AddJob(t *Task) string {
	// 创建任务
	taskId := mTask(t, "add")

	// 返回 任务id
	return taskId
}

// 停止定时任务
func StopJob(taskId string) {
	// 停止任务
	gotask.Stop(taskId)
}

// 更新任务运行时间
func UpdateJobTime(t *Task) {
	mTask(t, "update")
}
