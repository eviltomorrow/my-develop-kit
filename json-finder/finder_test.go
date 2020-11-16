package finder

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildKey(t *testing.T) {
	Convey("Test Build Key", t, func() {
		Convey("Case 1", func() {
			var parentKey = []string{
				"instance.cpu",
				"instance.ip",
			}
			var valueKey = "instance.disk.type"
			key, err := BuildKey(parentKey, valueKey)
			So(err, ShouldBeNil)
			t.Logf("Key: %v\r\n", key)
		})

		Convey("Case 2", func() {
			var parentKey = []string{
				"instance.cpu",
				"instance.disk.type",
			}
			var valueKey = "instance.disk.value"
			key, err := BuildKey(parentKey, valueKey)
			So(err, ShouldBeNil)
			t.Logf("Key: %v\r\n", key)
		})

		Convey("Case 3", func() {
			var parentKey = []string{
				"instance.cpu",
				"instance.mem.type",
			}
			var valueKey = "instance.disk.value"
			key, err := BuildKey(parentKey, valueKey)
			So(err, ShouldNotBeNil)
			t.Logf("Key: %v\r\n", key)
		})
	})
}
