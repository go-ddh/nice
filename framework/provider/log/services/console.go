package services

import (
	"os"

	"github.com/go-ddh/nice/framework"
	"github.com/go-ddh/nice/framework/contract"
)

// NiceConsoleLog 代表控制台输出
type NiceConsoleLog struct {
	NiceLog
}

// NewNiceConsoleLog 实例化NiceConsoleLog
func NewNiceConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &NiceConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	// 最重要的将内容输出到控制台
	log.SetOutput(os.Stdout)
	log.c = c
	return log, nil
}
