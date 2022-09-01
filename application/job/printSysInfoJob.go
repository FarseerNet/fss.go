package job

import (
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/flog"
)

func PrintSysInfoJob() {
	reservedTaskCount := configure.GetInt("FSS.ReservedTaskCount")
	flog.Infof("当前系统设置至少保留：%d条 历史任务", reservedTaskCount)
}
