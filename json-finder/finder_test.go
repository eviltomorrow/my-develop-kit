package finder

import (
	"io/ioutil"
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

func TestFindKey(t *testing.T) {
	Convey("Test Find Key", t, func() {
		Convey("Case 1", func() {
			buf, err := ioutil.ReadFile("data.json")
			So(err, ShouldBeNil)

			var parentKey = []string{
				"instanceList.ip",
			}
			var valueKey = "instanceList.disk.size"
			key, err := BuildKey(parentKey, valueKey)
			So(err, ShouldBeNil)

			FindKey(string(buf), key)
		})
	})
}

func TestDeepCloneKey(t *testing.T) {
	Convey("Test Deep Clone Key", t, func() {
		Convey("Case1", func() {
			var parentKey = []string{
				"instance.ip",
			}
			var valueKey = "instance.disk.type"
			key, err := BuildKey(parentKey, valueKey)
			So(err, ShouldBeNil)
			key.PrintPtr()

			var newKey = deepCloneKey(key)
			newKey.PrintPtr()
		})
	})
}

func TestSortParentKey(t *testing.T) {
	Convey("Test Sort Parent Key", t, func() {
		Convey("Case1", func() {
			var parentKey = []string{
				"instance.disk.type",
				"instance.cpu",
				"instance",
				"instance.c",
			}
			var valueKey = "instance.disk.type"
			key, err := BuildKey(parentKey, valueKey)
			So(err, ShouldBeNil)
			t.Logf("Key: %s", key)
		})
	})
}
