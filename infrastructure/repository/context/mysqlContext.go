package context

import (
	"fss/infrastructure/repository/model"
	"github.com/farseernet/farseer.go/data"
)

type MysqlContext struct {
	Build   data.TableSet[model.TaskLogPO]   `data:"name=run_log"`
	Admin   data.TableSet[model.TaskGroupPO] `data:"name=task_group"`
	Cluster data.TableSet[model.TaskPO]      `data:"name=task"`
}
