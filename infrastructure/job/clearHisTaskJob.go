package job

import "github.com/farseer-go/fss"

func RegisterClearHisTaskJob() {
	fss.RegisterJob("FSS.ClearHisTask", clearHisTaskJob)
}

// 自动清除历史任务记录
func clearHisTaskJob(context fss.IFssContext) {

}
