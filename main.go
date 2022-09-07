package main

import (
	"fss/interfaces"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
)

func main() { // main函数，程序执行的入口
	fs.Initialize[StartupModule]("fss")

	webapi.AutoRouter(&interfaces.MetaController{})
	webapi.AutoRouter(&interfaces.TaskController{})
	webapi.Run()
}

func test() {
	//try := exception.Try(func() {
	//	exception.ThrowRefuseException("test is throw")
	//})
	//try.CatchStringException(func(exp string) {
	//	flog.Info(exp)
	//})
	//try.CatchRefuseException(func(exp *exception.RefuseException) {
	//	flog.Warning(exp.Message)
	//	exception.ThrowRefuseException(exp.Message)
	//})
	//try.CatchStringException(func(exp string) {
	//	flog.Error(exp)
	//})
	//try.ThrowUnCatch()
}
