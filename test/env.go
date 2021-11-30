package test

import (
	"github.com/go-ddh/nice/framework"
	"github.com/go-ddh/nice/framework/provider/app"
	"github.com/go-ddh/nice/framework/provider/env"
)

const (
	BasePath = "/Users/yejianfeng/Documents/UGit/coredemo/"
)

func InitBaseContainer() framework.Container {
	// 初始化服务容器
	container := framework.NewNiceContainer()
	// 绑定App服务提供者
	container.Bind(&app.NiceAppProvider{BaseFolder: BasePath})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.NiceTestingEnvProvider{})
	return container
}
