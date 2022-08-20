package client

import (
	"fss/domain/clients/client"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/mapper"
	"time"
)

type app struct {
	repository client.Repository
}

func NewApp() *app {
	return &app{repository: container.Resolve[client.Repository]()}
}

// ToList 取出全局客户端列表
func (r *app) ToList() []DTO {
	lstDO := r.repository.ToList()
	return mapper.Array[DTO](lstDO)
}

// GetCount 客户端数量
func (r *app) GetCount() int64 {
	return r.repository.GetCount()
}

// UpdateClient 更新客户端的使用时间
func (r *app) UpdateClient(dto DTO) {
	do := mapper.Single[client.DomainObject](dto)
	r.repository.Update(do)
}

func (r *app) GetClient() DTO {
	do := client.DomainObject{
		//Id:         Headers["ClientId"].ToString().ConvertType(0),
		//Ip:         Headers["ClientIp"].ToString().Split(',')[0].Trim(),
		//Name:       Headers["ClientName"],
		//Jobs:       Headers["ClientJobs"].ToString().Split(','),
		ActivateAt: time.Now(),
	}

	// 更新客户端的使用时间
	r.repository.Update(do)

	return mapper.Single[DTO](do)
}
