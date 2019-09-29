package server

import (
	"fmt"
	"github.com/robfig/cron"
	"testing"
)

var consolvechan chan int

func TestStartConsole(t *testing.T) {
	pidFile := "../../.env_console_pid"
	consolvechan = make(chan int, 1)

	go func() {
		<-consolvechan
		stopServer(pidFile)
	}()
	StartConsole(pidFile, TempRegisterSchedule)
}

func TempRegisterSchedule(c *cron.Cron) {
	//c.AddFunc("0 30 * * * *", test)
	c.AddFunc("@every 1s", testConsole)
}

func testConsole() {
	fmt.Println("run test console")
	consolvechan <- 1
}
