package server

import (
	"fmt"
	"github.com/qit-team/work"
	"time"
)

func waitJobStop(job *work.Job) {
	//等待结束
	WaitStop()

	//暂停新的Cron任务执行
	job.Stop()

	err := job.WaitStop(60 * time.Second)
	if err != nil {
		fmt.Println("wait stop error", err)
	}

	CloseService()
}

// Start Job Worker
func StartJob(pidFile string, registerWorker func(*work.Job)) error {
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
