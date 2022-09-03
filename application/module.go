package application

import (
	"fss/application/job"
	"fss/domain"
	"github.com/farseer-go/fs"
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

	fs.AddInitCallback(job.PrintSysInfoJob)
	fs.AddInitCallback(job.InitSysTaskJob)
}

func (module Module) Shutdown() {
}
