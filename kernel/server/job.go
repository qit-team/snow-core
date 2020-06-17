package server

import (
	"fmt"
	"github.com/qit-team/work"
	"time"
)

const __MIN_WAIT_TIMEOUT_SECONDS = 60

func waitJobStop(job *work.Job, waitTimeout int) {
	//等待结束
	WaitStop()

	//暂停新的Cron任务执行
	job.Stop()

	if waitTimeout < __MIN_WAIT_TIMEOUT_SECONDS {
		waitTimeout = __MIN_WAIT_TIMEOUT_SECONDS
	}

	err := job.WaitStop(time.Duration(waitTimeout) * time.Second)
	if err != nil {
		fmt.Println("wait stop error", err)
	}

	CloseService()
}

// Start Job Worker
func StartJob(pidFile string, registerWorker func(*work.Job), waitTimeout int) error {
	//注册Job Worker
	job := work.New()
	registerWorker(job)
	job.Start()

	//写pid文件
	WritePidFile(pidFile)

	//注册信号量
	RegisterSignal()

	//等待停止信号
	waitJobStop(job, waitTimeout)
	return nil
}
