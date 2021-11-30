package id

import (
	tests "github.com/go-ddh/nice/test"
	"testing"

	"github.com/go-ddh/nice/framework/contract"
	"github.com/go-ddh/nice/framework/provider/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConsoleLog_Normal(t *testing.T) {
	Convey("test hade console log normal case", t, func() {
		c := tests.InitBaseContainer()
		c.Bind(&config.NiceConfigProvider{})

		err := c.Bind(&NiceIDProvider{})
		So(err, ShouldBeNil)

		idService := c.MustMake(contract.IDKey).(contract.IDService)
		xid := idService.NewID()
		So(xid, ShouldNotBeEmpty)
	})
}
