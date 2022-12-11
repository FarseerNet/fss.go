package clientApp

import "time"

type DTO struct {
	Id         int64     `webapi:"Clientid"`   // 客户端ID
	Ip         string    `webapi:"Clientip"`   // 客户端IP
	Name       string    `webapi:"Clientname"` // 客户端名称
	Jobs       []string  `webapi:"Clientjobs"` // 客户端能执行的任务
	ActivateAt time.Time // 活动时间
}
