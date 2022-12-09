package main

import (
	"fss/infrastructure"
	"fss/interfaces"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/webapi"
	"strings"
)

type StartupModule struct {
}

func (module StartupModule) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{webapi.Module{}, infrastructure.Module{}, interfaces.Module{}}
}

func (module StartupModule) PreInitialize() {
}

func (module StartupModule) Initialize() {
}

func (module StartupModule) PostInitialize() {
	// 服务端也使用了fss客户端，自动设置客户端的Server地址
	webApiUrl := configure.GetString("WebApi.Url")
	if strings.HasPrefix(webApiUrl, ":") {
		webApiUrl = "http://127.0.0.1" + webApiUrl
	} else if !strings.HasPrefix(webApiUrl, "http://") {
		webApiUrl = "http://" + webApiUrl
	}
	configure.SetDefault("FSS.Server", webApiUrl)
}

func (module StartupModule) Shutdown() {
}
