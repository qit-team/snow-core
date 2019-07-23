package server

import (
	"github.com/qit-team/snow-core/command"
)

// Execute one-time command
func ExecuteCommand(name string, registerCommand func(*command.Command)) error {
	//注册并执行某个name对应的脚本
	c := command.New()
	registerCommand(c)
	err := c.Execute(name)
	return err
}
