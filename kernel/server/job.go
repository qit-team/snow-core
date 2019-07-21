package server

import (
	"github.com/qit-team/snow-core/config"
	"fmt"
	"github.com/qit-team/work"
)

func waitJobStop(job *work.Job) {
	WaitStop()

	//暂停新的Cron任务执行
	job.Stop()
	err := job.WaitStop(0)
	if err != nil {
		fmt.Println("wait stop error", err)
	}
}

// Start Job Worker
func StartJob(confFile, pidFile string, boot func(*config.Config) error, registerWorker func(*work.Job)) error {
	//加载配置文件
	conf, err := config.Load(confFile)
	if err != nil {
		return err
	}

	//容器初始化
	err = boot(conf)
	if err != nil {
		return fmt.Errorf("container ini failed %s", err.Error())
	}

	//注册Job Worker
	job := work.New()
	registerWorker(job)
	job.Start()

	//写pid文件
	WritePidFile(pidFile)

	//注册信号量
	RegisterSignal()

	//等待停止信号
	waitJobStop(job)
	return nil
}
