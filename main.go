package main

import (
	"context"
	"fmt"
	"github.com/erikdubbelboer/gspt"
	"github.com/go-ddh/nice/app/devops"
	"github.com/go-ddh/nice/app/provider/demo"
	"github.com/go-ddh/nice/framework"
	"github.com/go-ddh/nice/framework/contract"
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
	"github.com/go-ddh/nice/framework/util"
	"github.com/sevlyar/go-daemon"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// 初始化服务容器
	container := framework.NewNiceContainer()
	// 绑定App服务提供者
	container.Bind(&demo.OpsProvider{})
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
	if engine, err := devops.NewHttpEngine(container); err == nil {
		container.Bind(&kernel.NiceKernelProvider{HttpEngine: engine})
	}
	initService(container)
}

func initService(container *framework.NiceContainer) error {
	// 从服务容器中获取kernel的服务实例
	kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
	// 从kernel服务实例中获取引擎
	core := kernelService.HttpEngine()
	appAddress := ""
	appDaemon := false
	if appAddress == "" {
		envService := container.MustMake(contract.EnvKey).(contract.Env)
		if envService.Get("ADDRESS") != "" {
			appAddress = envService.Get("ADDRESS")
		} else {
			configService := container.MustMake(contract.ConfigKey).(contract.Config)
			if configService.IsExist("app.address") {
				appAddress = configService.GetString("app.address")
			} else {
				appAddress = ":8888"
			}
		}
	}
	// 创建一个Server服务
	server := &http.Server{
		Handler: core,
		Addr:    appAddress,
	}

	appService := container.MustMake(contract.AppKey).(contract.App)

	pidFolder := appService.RuntimeFolder()
	if !util.Exists(pidFolder) {
		if err := os.MkdirAll(pidFolder, os.ModePerm); err != nil {
			return err
		}
	}
	serverPidFile := filepath.Join(pidFolder, "app.pid")
	logFolder := appService.LogFolder()
	if !util.Exists(logFolder) {
		if err := os.MkdirAll(logFolder, os.ModePerm); err != nil {
			return err
		}
	}
	// 应用日志
	serverLogFile := filepath.Join(logFolder, "app.log")
	currentFolder := util.GetExecDirectory()
	// daemon 模式
	if appDaemon {
		// 创建一个Context
		ctxt := &daemon.Context{
			// 设置pid文件
			PidFileName: serverPidFile,
			PidFilePerm: 0664,
			// 设置日志文件
			LogFileName: serverLogFile,
			LogFilePerm: 0640,
			// 设置工作路径
			WorkDir: currentFolder,
			// 设置所有设置文件的mask，默认为750
			Umask: 027,
			// 子进程的参数，按照这个参数设置，子进程的命令为 ./nice app start --daemon=true
			Args: []string{"", "app", "start", "--daemon=true"},
		}
		// 启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
		d, err := ctxt.Reborn()
		if err != nil {
			return err
		}
		if d != nil {
			// 父进程直接打印启动成功信息，不做任何操作
			fmt.Println("app启动成功，pid:", d.Pid)
			fmt.Println("日志文件:", serverLogFile)
			return nil
		}
		defer ctxt.Release()
		// 子进程执行真正的app启动操作
		fmt.Println("demon started")
		//spot.SetProcTitle("nice app")
		if err := startAppServe(server, container); err != nil {
			fmt.Println(err)
		}
		return nil
	}

	// 非demon模式，直接执行
	content := strconv.Itoa(os.Getpid())
	fmt.Println("[PID]", content)
	err := ioutil.WriteFile(serverPidFile, []byte(content), 0644)
	if err != nil {
		return err
	}
	gspt.SetProcTitle("nice app")

	fmt.Println("app serve url:", appAddress)
	if err := startAppServe(server, container); err != nil {
		fmt.Println(err)
	}
	return nil
}

// 启动AppServer, 这个函数会将当前goroutine阻塞
func startAppServe(server *http.Server, c framework.Container) error {
	// 这个goroutine是启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	closeWait := 5
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	if configService.IsExist("app.close_wait") {
		closeWait = configService.GetInt("app.close_wait")
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(closeWait)*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		return err
	}
	return nil
}
