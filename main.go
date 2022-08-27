package main

import (
	"fss/interfaces"
	"github.com/beego/beego/v2/server/web"
	"github.com/farseer-go/fs"
)

func main() { // main函数，程序执行的入口
	fs.Initialize[StartupModule]("fss")

	web.AutoRouter(&interfaces.MetaController{})
	web.AutoRouter(&interfaces.TaskController{})
	web.Run()
}
