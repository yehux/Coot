package setting

import (
	"Coot/core/exec"
	"Coot/utils/color"
	"Coot/view"
	"Coot/view/setting"
	"Coot/view/task"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func RunWeb(addr string) {
	gin.DisableConsoleColor()
	f, _ := os.Create("./logs/coot.log")
	gin.DefaultWriter = io.MultiWriter(f)
	// 引入gin
	r := gin.Default()
	// 引入html资源
	r.LoadHTMLGlob("web/*")

	// 引入静态资源
	r.Static("/static", "./static")

	// 加载路由
	view.LoadUrl(r)

	// 清除 task_id
	task.UpdateTaskAll()
	// 清除 日志
	//report.DeleteLogsAll()
	//初始日志开关
	setting.InitLogsOff()
	r.Run(addr)
}

func Init() {
	fmt.Println("test")
}

func Help() {
	exec.Execute("clear")
	logo := `    ______            __
   / ____/___  ____  / /_
  / /   / __ \/ __ \/ __/
 / /___/ /_/ / /_/ / /_
 \____/\____/\____/\__/`
	fmt.Print(color.Yellow(logo))
	fmt.Println(color.White("   v0.1\n"))
	fmt.Println(color.White(" Play IFTTT, Experience Geek Life, Internet Automation."))
	fmt.Println("")
	fmt.Println(color.Yellow(" + [ ABOUT ]--------------------------------------------------- +"))
	fmt.Println("")
	fmt.Println(color.Green("   + Home : "), color.White("https://coot.io"), color.Green("    + Team : "), color.White("https://yehu.io"))
	fmt.Println("")
	fmt.Println(color.Yellow(" + [ ARGUMENTS ]----------------------------------------------- +"))
	fmt.Println("")
	fmt.Println(color.Cyan("   run,--run"), color.White("	       Start up service, Default localhost:9000"))
	//fmt.Println(color.Cyan("   init,--init"), color.White("		   Initialization, Wipe data"))
	fmt.Println(color.Cyan("   version,--version"), color.White("  Coot Version"))
	fmt.Println(color.Cyan("   help,--help"), color.White("	       Help"))
	fmt.Println("")
	fmt.Println(color.Yellow(" + ------------------------------------------------------------ +"))
	fmt.Println("")
}
