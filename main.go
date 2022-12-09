package main

import (
	"fss/interfaces"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
)

func main() { // main函数，程序执行的入口知
	fs.Initialize[StartupModule]("fss")

	webapi.RegisterRoutes(routeMeta)
	webapi.RegisterController(&interfaces.TaskController{})

	webapi.UseApiResponse()
	webapi.Run()
}
