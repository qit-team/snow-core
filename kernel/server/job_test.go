package server

import (
	"fmt"
	"github.com/qit-team/snow-core/config"
	"github.com/qit-team/snow-core/queue"
	"github.com/qit-team/snow-core/queue/redisqueue"
	"github.com/qit-team/snow-core/redis"
	"github.com/qit-team/work"
	"sync"
	"testing"
	"time"
)

var (
	jb       *work.Job
	register func(job *work.Job)
	mu       sync.RWMutex
)

var q queue.Queue

func init() {

	redisConf := config.RedisConfig{
		Master: config.RedisBaseConfig{
			Host: "127.0.0.1",
			Port: 6379,
		},
	}

	// 注册redis类
	err := redis.Pr.Register("redis", redisConf, true)
	if err != nil {
		fmt.Println(err)
	}

	// 为了让redisqueue的driver先进行注册
	redisqueue.GetRedisQueue("redis")
	q = queue.GetQueue("redis", queue.DriverTypeRedis)
}

func TestStartJob(t *testing.T) {
	pidFile := "../../.env_pid"

	StartJob(pidFile, TempRegisterWorker, 60)
}

func TempRegisterWorker(job *work.Job) {
	TempSetJob(job)

	//设置worker的任务投递回调函数
	job.AddFunc("topic-test", test)
	//设置worker的任务投递回调函数，和并发数
	job.AddFunc("topic-test1", test, 2)
	//使用worker结构进行注册
	job.AddWorker("topic-test2", &work.Worker{Call: work.MyWorkerFunc(test), MaxConcurrency: 1})

	TempRegisterQueueDriver(job)
}

func TempSetJob(job *work.Job) {
	if jb == nil {
		jb = job
	}
}

func TempRegisterQueueDriver(job *work.Job) {
	q := queue.GetQueue(redis.SingletonMain, queue.DriverTypeRedis)
	job.AddQueue(q, "topic-test1", "topic-test2")
	job.AddQueue(q)
}

func test(task work.Task) work.TaskResult {
	time.Sleep(time.Millisecond * 5)
	s, err := work.JsonEncode(task)
	if err != nil {

		return work.TaskResult{Id: task.Id, State: work.StateFailedWithAck}
	} else {
		//work.StateSucceed 会进行ack确认
		fmt.Println("do task", s)
		return work.TaskResult{Id: task.Id, State: work.StateSucceed}
	}
}
