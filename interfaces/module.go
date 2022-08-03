package interfaces

import (
	"fss/application"
	"github.com/farseernet/farseer.go/modules"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{application.Module{}}
}

func (module Module) PreInitialize() {
	//TODO implement me
}

func (module Module) Initialize() {
	//TODO implement me
}

func (module Module) PostInitialize() {
	//TODO implement me
}

func (module Module) Shutdown() {
	//TODO implement me
}
