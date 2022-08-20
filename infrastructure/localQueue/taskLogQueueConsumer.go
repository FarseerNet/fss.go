package localQueue

import (
	"fss/infrastructure/repository"
	"fss/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/queue"
)

func SubscribeTaskLogQueue() {
	queue.Subscribe("TaskLogQueue", "", 1000, taskLogQueueConsumer)
}

// 将日志指写入
func taskLogQueueConsumer(subscribeName string, message []any, remainingCount int) {
	// 转成BuildLogVO数组
	var lstPO collections.List[model.TaskLogPO]
	collections.NewListAny(message...).Select(&lstPO, func(item any) any {
		return item.(model.TaskLogPO)
	})
	repository.NewTaskLogRepository().AddBatch(lstPO)
}
