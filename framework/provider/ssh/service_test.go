package ssh

import (
	"github.com/go-ddh/nice/framework/provider/config"
	"github.com/go-ddh/nice/framework/provider/log"
	tests "github.com/go-ddh/nice/test"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNiceSSHService_Load(t *testing.T) {
	container := tests.InitBaseContainer()
	container.Bind(&config.NiceConfigProvider{})
	container.Bind(&log.NiceLogServiceProvider{})

	Convey("test get client", t, func() {
		hadeRedis, err := NewNiceSSH(container)
		So(err, ShouldBeNil)
		service, ok := hadeRedis.(*NiceSSH)
		So(ok, ShouldBeTrue)
		client, err := service.GetClient(WithConfigPath("ssh.web-01"))
		So(err, ShouldBeNil)
		So(client, ShouldNotBeNil)
		session, err := client.NewSession()
		So(err, ShouldBeNil)
		out, err := session.Output("pwd")
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
		session.Close()
	})
}
