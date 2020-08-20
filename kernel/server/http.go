package server

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/config"
	"strconv"
	"syscall"
)

/**
 * 启动gin引擎
 * @wiki https://github.com/fvbock/endless#signals
 */
func runEngine(engine *gin.Engine, addr string, pidPath string) error {
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
	//设置gin调试模式
	if !GetDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	//配置路由引擎
	engine := gin.New()
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
