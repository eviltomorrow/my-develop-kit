package httpclient

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testXML = `<?xml version="1.0" encoding="UTF-8"?>
<nodeList>
  <clusterList>
    <id>1</id>
    <hostPoolId>2</hostPoolId>
    <name>提供虚拟机CVK</name>
    <description/>
    <enableHA>0</enableHA>
    <priority>3</priority>
    <enableLB>0</enableLB>
    <persistTime>0</persistTime>
    <checkInterval>0</checkInterval>
    <enableIPM>0</enableIPM>
    <persistTimeIPM>0</persistTimeIPM>
    <checkIntervalIPM>0</checkIntervalIPM>
    <operatorGroupId>0</operatorGroupId>
    <operatorGroupCode/>
    <childNum>5</childNum>
  </clusterList>
  <hostList>
    <id>107</id>
    <user>root</user>
    <pwd>1q2w3e</pwd>
    <hostPoolId>2</hostPoolId>
    <name>cvm80</name>
    <ip>192.168.10.80</ip>
    <status>1</status>
    <haEnable>1</haEnable>
  </hostList>
</nodeList>`

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

func TestGetJSON(t *testing.T) {
	var data = reg.ReplaceAllString(testJSON, "")
	Convey("Test GetJSONValue", t, func() {
		Convey("JSON: ", func() {
			result, err := GetJSON(data, "name")
			So(err, ShouldBeNil)
			t.Logf("%v", err)
			t.Logf("%v", result)
		})
	})
}

func TestGetXML(t *testing.T) {
	var data = reg.ReplaceAllString(testXML, "")
	Convey("Test GetXMLValue", t, func() {
		Convey("XML: ", func() {
			result, err := GetXML(data, "nodeList.hostList")
			So(err, ShouldBeNil)
			t.Logf("%v", err)
			t.Logf("%v", result)
		})
	})
}

func BenchmarkGetJSON(b *testing.B) {
	var data = reg.ReplaceAllString(testJSON, "")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetJSON(data, "name.first")
	}
}

func BenchmarkGetXML(b *testing.B) {
	var data = reg.ReplaceAllString(testXML, "")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetXML(data, "nodeList.hostList.ip")
	}
}
