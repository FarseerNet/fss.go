package infrastructure

import (
	"fss/infrastructure/domainEvent"
	_ "fss/infrastructure/domainEvent"
	"fss/infrastructure/localQueue"
	_ "fss/infrastructure/localQueue"
	"fss/infrastructure/repository"
	_ "fss/infrastructure/repository"
	"github.com/farseer-go/data"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/queue"
	"github.com/farseer-go/redis"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{data.Module{}, redis.Module{}, eventBus.Module{}, queue.Module{}}
}

func (module Module) PreInitialize() {
	domainEvent.SubscribeTaskFinishEvent()

	localQueue.SubscribeTaskLogQueue()

	repository.RegisterClientRepository()
	repository.RegisterTaskGroupRepository()
	repository.RegisterTaskLogRepository()
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
}

func (module Module) Shutdown() {
}
