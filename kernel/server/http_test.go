package server

import (
	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/config"

	"fmt"
	"go/build"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestStartHttp(t *testing.T) {
	pidFile := "../../.env_pid"
	var apiConf config.ApiConfig
	apiConf.Host = "127.0.0.1"
	apiConf.Port = 9000
	go func() {
		for i := 1; i < 100; i++ {
			// 进程启动后http服务没法自动停掉，需要借助os.exec执行自动停止
			stopServer(pidFile)
		}
	}()
	StartHttp(pidFile, apiConf, RegisterRoute)
}

//api路由配置
func RegisterRoute(router *gin.Engine) {
	router.GET("/hello", HandleHello)
}

func HandleHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":        200,
		"message":     "ok",
		"request_uri": c.Request.URL.Path,
		"data":        "test",
	})
	c.Abort()
	return
}

func stopServer(pidPath string) error {
	pid, _ := ReadPidFile(pidPath)
	pidStr := strconv.Itoa(pid)
	cmdName, cmdPath, command := "stop http", gopath(), "kill -TERM "+pidStr

	cmds := strings.Split(command, " ")
	err := runTool(cmdName, cmdPath, cmds[0], cmds[1:])
	return err
}

// 封装os.exec
func runTool(name, dir, cmd string, args []string) (err error) {
	toolCmd := &exec.Cmd{
		Path:   cmd,
		Args:   append([]string{cmd}, args...),
		Dir:    dir,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
	}

	if filepath.Base(cmd) == cmd {
		var lp string
		if lp, err = exec.LookPath(cmd); err == nil {
			toolCmd.Path = lp
		}
	}
	if err = toolCmd.Run(); err != nil {
		if e, ok := err.(*exec.ExitError); !ok || !e.Exited() {
			fmt.Fprintf(os.Stderr, "运行 %s 出错: %v\n", name, err)
		}
	}
	return
}

// 获取gopath路径
func gopath() (gp string) {
	gopaths := strings.Split(os.Getenv("GOPATH"), ":")
	if len(gopaths) == 1 {
		return gopaths[0]
	}
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	abspwd, err := filepath.Abs(pwd)
	if err != nil {
		return
	}
	for _, gopath := range gopaths {
		absgp, err := filepath.Abs(gopath)
		if err != nil {
			return
		}
		if strings.HasPrefix(abspwd, absgp) {
			return absgp
		}
	}
	return build.Default.GOPATH
}
