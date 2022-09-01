package main

import (
	"fss/interfaces"
	"github.com/beego/beego/v2/server/web"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/webapi"
)

func main() { // main函数，程序执行的入口
	fs.Initialize[StartupModule]("fss")

	try := exception.Try(func() {
		exception.ThrowRefuseException("test is throw")
	})
	try.CatchStringException(func(exp string) {
		flog.Info(exp)
	})
	try.CatchRefuseException(func(exp *exception.RefuseException) {
		flog.Warning(exp.Message)
		exception.ThrowRefuseException(exp.Message)
	})
	try.CatchStringException(func(exp string) {
		flog.Error(exp)
	})

	//try.ThrowUnCatch()
	webapi.Run()

	web.AutoRouter(&interfaces.MetaController{})
	web.AutoRouter(&interfaces.TaskController{})
	//web.Run()
}

func test() {
	exception.ThrowRefuseException("test is throw")
}
