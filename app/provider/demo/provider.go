package demo

import (
	"github.com/go-ddh/nice/framework"
)

type OpsProvider struct {
	framework.ServiceProvider
	c framework.Container
}

func (op *OpsProvider) Name() string {
	return DemoKey
}

func (op *OpsProvider) Register(c framework.Container) framework.NewInstance {
	return NewService
}

func (op *OpsProvider) IsDefer() bool {
	return false
}

func (op *OpsProvider) Params(c framework.Container) []interface{} {
	return []interface{}{op.c}
}

func (op *OpsProvider) Boot(c framework.Container) error {
	op.c = c
	return nil
}
