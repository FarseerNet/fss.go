package client

import "time"

type DTO struct {
	// 客户端ID
	Id int64
	// 客户端IP
	Ip string
	// 客户端名称
	Name string
	// 客户端能执行的任务
	Jobs []string
	// 活动时间
	ActivateAt time.Time
}
