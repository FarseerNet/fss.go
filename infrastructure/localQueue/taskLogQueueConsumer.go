package localQueue

import (
	"fss/infrastructure/repository"
	"fss/infrastructure/repository/model"
	"github.com/farseer-go/linq"
	"github.com/farseer-go/queue"
)

func SubscribeTaskLogQueue() {
	queue.Subscribe("TaskLogQueue", "", 1000, taskLogQueueConsumer)
}

// 将日志指写入
func taskLogQueueConsumer(subscribeName string, message []any, remainingCount int) {
	// 转成BuildLogVO数组
	var lstPO []model.TaskLogPO
	linq.From(message).Select(&lstPO, func(item any) any {
		return item.(model.TaskLogPO)
	})
	repository.NewTaskLogRepository().AddBatch(lstPO)
}
