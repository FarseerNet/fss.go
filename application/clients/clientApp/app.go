package clientApp

import (
	"fss/domain/clients/client"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/mapper"
)

// ToList 取出全局客户端列表
func ToList() collections.List[client.DomainObject] {
	repository := container.Resolve[client.Repository]()
	return repository.ToList()
}

// GetCount 客户端数量
func GetCount() int64 {
	repository := container.Resolve[client.Repository]()
	return repository.GetCount()
}

// UpdateClient 更新客户端的使用时间
func UpdateClient(dto DTO) {
	if dto.Id > 0 {
		repository := container.Resolve[client.Repository]()
		do := mapper.Single[client.DomainObject](dto)
		repository.Update(do)
	}
}
