package server

import (
	"github.com/robfig/cron"
	"fmt"
	"time"
)

func waitConsoleStop(c *cron.Cron) {
	//等待结束
	WaitStop()

	//暂停新的Cron任务执行
	c.Stop()

	//等待执行中的cron任务结束，目前简单实现等待5s后结束
	if GetDebug() {
		fmt.Println("wait 5 sencods")
	}
	time.Sleep(time.Second * 5)

	CloseService()
}

// Start Cron Schedule
func StartConsole(pidFile string, registerSchedule func(*cron.Cron)) error {
	//注册Cron执行计划
	cronEngine := cron.New()
	registerSchedule(cronEngine)
	cronEngine.Start()

	//写pid文件
	WritePidFile(pidFile)

	//注册信号量
	RegisterSignal()

	//等待停止信号
	waitConsoleStop(cronEngine)
	return nil
}
