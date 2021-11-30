package main

import (
	"github.com/go-ddh/nice/app/console"
	"github.com/go-ddh/nice/app/http"
	"github.com/go-ddh/nice/framework"
	"github.com/go-ddh/nice/framework/provider/app"
	"github.com/go-ddh/nice/framework/provider/cache"
	"github.com/go-ddh/nice/framework/provider/config"
	"github.com/go-ddh/nice/framework/provider/distributed"
	"github.com/go-ddh/nice/framework/provider/env"
	"github.com/go-ddh/nice/framework/provider/id"
	"github.com/go-ddh/nice/framework/provider/kernel"
	"github.com/go-ddh/nice/framework/provider/log"
	"github.com/go-ddh/nice/framework/provider/orm"
	"github.com/go-ddh/nice/framework/provider/redis"
	"github.com/go-ddh/nice/framework/provider/ssh"
	"github.com/go-ddh/nice/framework/provider/trace"
)

func main() {
	// 初始化服务容器
	container := framework.NewNiceContainer()
	// 绑定App服务提供者
	//container.Bind(&devops.OpsProvider{})
	// 绑定App服务提供者
	container.Bind(&app.NiceAppProvider{})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.NiceEnvProvider{})
	container.Bind(&distributed.LocalDistributedProvider{})
	container.Bind(&config.NiceConfigProvider{})
	container.Bind(&id.NiceIDProvider{})
	container.Bind(&trace.NiceTraceProvider{})
	container.Bind(&log.NiceLogServiceProvider{})
	container.Bind(&orm.GormProvider{})
	container.Bind(&redis.RedisProvider{})
	container.Bind(&cache.NiceCacheProvider{})
	container.Bind(&ssh.SSHProvider{})
	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(container); err == nil {
		container.Bind(&kernel.NiceKernelProvider{HttpEngine: engine})
	}

	// 运行root命令
	console.RunCommand(container)
}
