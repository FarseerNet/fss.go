package main

import (
	"fss/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"strconv"
	"time"
)

func initData() {
	for i := 1; i <= 5; i++ {
		data.NewContext[taskGroupRepository]("default").TaskGroup.Insert(&model.TaskGroupPO{
			Caption:     "demo.job1-" + strconv.Itoa(i),
			JobName:     "demo.job1",
			Data:        collections.NewDictionary[string, string](),
			StartAt:     time.Now(),
			NextAt:      time.Now(),
			IntervalMs:  1000,
			Cron:        "",
			ActivateAt:  time.Now(),
			LastRunAt:   time.Now(),
			RunSpeedAvg: 0,
			RunCount:    0,
			IsEnable:    true,
		})
	}

	for i := 6; i <= 10; i++ {
		data.NewContext[taskGroupRepository]("default").TaskGroup.Insert(&model.TaskGroupPO{
			Caption:     "demo.job2-" + strconv.Itoa(i),
			JobName:     "demo.job2",
			Data:        collections.NewDictionary[string, string](),
			StartAt:     time.Now(),
			NextAt:      time.Now(),
			IntervalMs:  1000,
			Cron:        "",
			ActivateAt:  time.Now(),
			LastRunAt:   time.Now(),
			RunSpeedAvg: 0,
			RunCount:    0,
			IsEnable:    true,
		})
	}

	for i := 11; i <= 15; i++ {
		data.NewContext[taskGroupRepository]("default").TaskGroup.Insert(&model.TaskGroupPO{
			Caption:     "demo.job3-" + strconv.Itoa(i),
			JobName:     "demo.job3",
			Data:        collections.NewDictionary[string, string](),
			StartAt:     time.Now(),
			NextAt:      time.Now(),
			IntervalMs:  1000,
			Cron:        "",
			ActivateAt:  time.Now(),
			LastRunAt:   time.Now(),
			RunSpeedAvg: 0,
			RunCount:    0,
			IsEnable:    true,
		})
	}
}

type taskGroupRepository struct {
	TaskGroup data.TableSet[model.TaskGroupPO] `data:"name=task_group"`
}
