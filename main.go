package main

import (
	"fss/interfaces"
	"github.com/beego/beego/v2/server/web"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
)

func main() { // main函数，程序执行的入口
	fs.Initialize[StartupModule]("fss")

	webapi.Run()

	web.AutoRouter(&interfaces.MetaController{})
	web.AutoRouter(&interfaces.TaskController{})
	web.Run()
}
