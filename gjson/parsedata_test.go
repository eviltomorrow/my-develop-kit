package httpclient

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testJSON = `{
	"name": {"first": "Tom", "last": "Anderson"},
	"age":37,
	"children": ["Sara","Alex","Jack"],
	"fav.movie": "Deer Hunter",
	"friends": [
	  {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
	  {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
	  {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
	]
}`

var reg = regexp.MustCompile("\\s+")

func TestGetJSONValue(t *testing.T) {
	var data = reg.ReplaceAllString(testJSON, "")
	Convey("Test GetJSONValue", t, func() {
		Convey("JSON: ", func() {
			result, err := GetJSON(data, "friends.age")
			So(err, ShouldBeNil)
			t.Logf("%v", err)
			t.Logf("%v", result)
		})
	})
}
