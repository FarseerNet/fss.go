package main

import (
	"fss/infrastructure"
	"fss/interfaces"
	"github.com/farseernet/farseer.go/modules"
)

type StartupModule struct {
}

func (module StartupModule) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{interfaces.Module{}, infrastructure.Module{}}
}

func (module StartupModule) PreInitialize() {
}

func (module StartupModule) Initialize() {
}

func (module StartupModule) PostInitialize() {
}

func (module StartupModule) Shutdown() {
}
