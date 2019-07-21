package server

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/qit-team/snow-core/config"
	"strconv"
	"github.com/fvbock/endless"
	"syscall"
)

/**
 * 启动gin引擎
 * @wiki https://github.com/fvbock/endless#signals
 */
func runEngine(engine *gin.Engine, addr string, pidPath string) error {
	//设置gin调试模式
	if !GetDebug() {
		gin.SetMode(gin.ReleaseMode)
	}

	server := endless.NewServer(addr, engine)
	server.BeforeBegin = func(add string) {
		pid := syscall.Getpid()
		if gin.Mode() != gin.ReleaseMode {
			fmt.Printf("Actual pid is %d \n\r", pid)
		}
		WritePidFile(pidPath, pid)
	}
	err := server.ListenAndServe()
	return err
}

// Start proxy with config file
func StartHttp(pidFile string, apiConf config.ApiConfig, registerRoute func(*gin.Engine)) error {
	//配置路由引擎
	engine := gin.Default()
	registerRoute(engine)
	addr := apiConf.Host + ":" + strconv.Itoa(apiConf.Port)
	runEngine(engine, addr, pidFile)

	//因为信号处理由endless接管实现平滑重启和关闭，这里模拟通用的结束信号
	go func() {
		Stop()
	}()

	//等待停止信号
	WaitStop()
	return nil
}
