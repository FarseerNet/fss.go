package main

import (
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fss"
)

type StartupModule struct {
}

func (module StartupModule) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{fss.Module{}}
}

func (module StartupModule) PreInitialize() {
}

func (module StartupModule) Initialize() {
}

func (module StartupModule) PostInitialize() {
	RegisterDemo1()
	RegisterDemo2()
	RegisterDemo3()
}

func (module StartupModule) Shutdown() {
}
