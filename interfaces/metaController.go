package interfaces

import "github.com/beego/beego/v2/server/web"

type MetaController struct {
	web.Controller
}

// GetClientList 取出全局客户端列表
func (r *MetaController) GetClientList() {
	//name := r.GetString("name")
	r.Ctx.WriteString("Hello, world")
}
