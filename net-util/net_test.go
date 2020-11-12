package netutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetExternalIP(t *testing.T) {
	Convey("GetExternalIP", t, func() {
		Convey("Case1", func() {
			ip, err := GetExternalIP()
			So(err, ShouldBeNil)
			t.Logf("IP: %v\r\n", ip)
		})
	})
}

func TestGetLocalIP(t *testing.T) {
	Convey("GetLocalIP", t, func() {
		Convey("Case1", func() {
			ip, err := GetLocalIP()
			So(err, ShouldBeNil)
			t.Logf("IP: %v\r\n", ip)
		})
	})
}

func TestGetLocalIPv4s(t *testing.T) {
	Convey("GetLocalIPv4s", t, func() {
		Convey("Case1", func() {
			ips, err := GetLocalIPv4s()
			So(err, ShouldBeNil)
			t.Logf("IP: %v\r\n", ips)
		})
	})
}
