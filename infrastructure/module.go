package infrastructure

import (
	_ "fss/infrastructure/domainEvent"
	_ "fss/infrastructure/localQueue"
	_ "fss/infrastructure/repository"
	"github.com/farseernet/farseer.go/data"
	"github.com/farseernet/farseer.go/modules"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{data.Module{}}
}

func (module Module) PreInitialize() {
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
}

func (module Module) Shutdown() {
}
