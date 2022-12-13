package application

import (
	"context"
	"fss/application/job"
	"fss/domain"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fss"
	"github.com/farseer-go/tasks"
	"time"
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
	fss.RegisterJob("FSS.CheckClientOffline", job.CheckClientOfflineJob)
	tasks.Run("FSS.CheckClientOffline", 30*time.Second, job.CheckFinishStatusJob, context.Background())
	tasks.Run("FSS.CheckWorkStatus", 30*time.Second, job.CheckWorkStatusJob, context.Background())
	fss.RegisterJob("FSS.ClearHisTask", job.ClearHisTaskJob)
	fss.RegisterJob("FSS.SyncTaskGroupAvgSpeed", job.SyncTaskGroupAvgSpeedJob)
	fss.RegisterJob("FSS.SyncTaskGroup", job.SyncTaskGroupJob)

	fs.AddInitCallback(job.PrintSysInfoJob)
	fs.AddInitCallback(job.InitSysTaskJob)
}

func (module Module) Shutdown() {
}
