package test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdd(t *testing.T) {
	Convey("测试 Add 函数", t, func() {
		So(Add(1, 2), ShouldEqual, 34)
	})
}
