package config

import (
	"path/filepath"
	"testing"

	"github.com/go-ddh/nice/framework/contract"
	tests "github.com/go-ddh/nice/test"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNiceConfig_GetInt(t *testing.T) {
	container := tests.InitBaseContainer()

	Convey("test nice env normal case", t, func() {
		appService := container.MustMake(contract.AppKey).(contract.App)
		envService := container.MustMake(contract.EnvKey).(contract.Env)
		folder := filepath.Join(appService.ConfigFolder(), envService.AppEnv())

		serv, err := NewNiceConfig(container, folder, map[string]string{})
		So(err, ShouldBeNil)
		conf := serv.(*NiceConfig)
		timeout := conf.GetString("database.default.timeout")
		So(timeout, ShouldEqual, "10s")
	})
}
