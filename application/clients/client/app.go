package client

type app struct {
	//repository containerLog.Repository
}

func NewApp() *app {
	//return &app{repository: container.Resolve[containerLog.Repository]()}
	return &app{}
}
