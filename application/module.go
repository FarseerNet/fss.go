package application

import (
	"fss/application/job"
	"fss/domain"
	"github.com/farseer-go/fs/modules"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{domain.Module{}}
}

func (module Module) PreInitialize() {
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
	job.RegisterCheckClientOfflineJob()
	job.RegisterCheckFinishStatusJob()
	job.RegisterCheckWorkStatusJob()
	job.RegisterSyncTaskGroupAvgSpeedJob()

	job.PrintSysInfoJob()
	job.InitSysTaskJob()
}

func (module Module) Shutdown() {
}
