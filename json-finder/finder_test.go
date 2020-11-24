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
			key, err := BuildKey(parentKey, nil, valueKey)
			So(err, ShouldBeNil)
			t.Logf("Key: %v\r\n", key)
		})

		Convey("Case 2", func() {
			var parentKey = []string{
				"instance.cpu",
				"instance.disk.type",
			}
			var valueKey = "instance.disk.value"
			key, err := BuildKey(parentKey, nil, valueKey)
			So(err, ShouldBeNil)
			t.Logf("Key: %v\r\n", key)
		})

		Convey("Case 3", func() {
			var parentKey = []string{
				"instance.cpu",
				"instance.mem.type",
			}
			var valueKey = "instance.disk.value"
			key, err := BuildKey(parentKey, nil, valueKey)
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
				"instanceList.disk.type",
			}
			var brotherKey = []string{
				"instanceList.disk.type",
			}
			var valueKey = "instanceList.disk.size"
			key, err := BuildKey(parentKey, brotherKey, valueKey)
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
			var brotherKey = []string{
				"instance.disk.type",
			}
			var valueKey = "instance.disk.type"
			key, err := BuildKey(parentKey, brotherKey, valueKey)
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
			key, err := BuildKey(parentKey, nil, valueKey)
			So(err, ShouldBeNil)
			t.Logf("Key: %s", key)
		})
	})
}

func BenchmarkFindKey(b *testing.B) {
	buf, _ := ioutil.ReadFile("data.json")

	var parentKey = []string{
		"instanceList.ip",
		"instanceList.disk.type",
	}
	var valueKey = "instanceList.disk.1"
	key, _ := BuildKey(parentKey, nil, valueKey)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FindKey(string(buf), key)
	}

}
