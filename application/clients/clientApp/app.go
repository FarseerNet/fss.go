package clientApp

import (
	"fss/domain/clients/client"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/mapper"
)

// ToList 取出全局客户端列表
func ToList() []DTO {
	repository := container.Resolve[client.Repository]()
	lstDO := repository.ToList()
	return mapper.Array[DTO](lstDO)
}

// GetCount 客户端数量
func GetCount() int64 {
	repository := container.Resolve[client.Repository]()
	return repository.GetCount()
}

// UpdateClient 更新客户端的使用时间
func UpdateClient(dto DTO) {
	repository := container.Resolve[client.Repository]()
	do := mapper.Single[client.DomainObject](&dto)
	repository.Update(do)
}
