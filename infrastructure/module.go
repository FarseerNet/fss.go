package infrastructure

import (
	"fss/domain/tasks/taskGroup"
	"fss/infrastructure/domainEvent"
	"fss/infrastructure/localQueue"
	"fss/infrastructure/repository"
	"github.com/farseer-go/cache"
	"github.com/farseer-go/data"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fss"
	"github.com/farseer-go/queue"
	"github.com/farseer-go/redis"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{data.Module{}, redis.Module{}, eventBus.Module{}, queue.Module{}, fss.Module{}}
}

func (module Module) PreInitialize() {
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
	cache.SetProfilesInRedis[taskGroup.DomainObject]("FSS_TaskGroup", "default", "Id", 0)

	repository.RegisterClientRepository()
	repository.RegisterTaskGroupRepository()
	repository.RegisterTaskLogRepository()

	domainEvent.SubscribeTaskFinishEvent()
	localQueue.SubscribeTaskLogQueue()
}

func (module Module) Shutdown() {
}
