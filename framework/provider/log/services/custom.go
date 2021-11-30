package services

import (
	"github.com/go-ddh/nice/framework"
	"github.com/go-ddh/nice/framework/contract"
	"io"
)

type NiceCustomLog struct {
	NiceLog
}

func NewNiceCustomLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	output := params[4].(io.Writer)

	log := &NiceConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	log.SetOutput(output)
	log.c = c
	return log, nil
}
