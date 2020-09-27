package httpclient

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetHTTP(t *testing.T) {
	Convey("Test GetHTTP", t, func() {
		client, err := CreateClientHTTP(10*time.Second, nil, nil, nil)
		So(err, ShouldBeNil)

		var case1URL = "http://localhost:9090/host/info"
		Convey(fmt.Sprintf("Normal URL: %v", case1URL), func() {
			result, err := GetHTTP(client, case1URL, nil)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeEmpty)
			t.Logf("Data: %v", result)
		})

		var case2URL = "http://localhost:1010/host/info"
		Convey(fmt.Sprintf("Wrong URL(wrong port): %v", case2URL), func() {
			result, err := GetHTTP(client, case2URL, nil)
			So(err, ShouldNotBeNil)
			So(result, ShouldBeEmpty)
			t.Logf("Error: %v", err)
		})

		var case3URL = "http://192.168.180.67:9090/host/info"
		Convey(fmt.Sprintf("Wrong URL(wrong host): %v", case3URL), func() {
			result, err := GetHTTP(client, case3URL, nil)
			So(err, ShouldNotBeNil)
			So(result, ShouldBeEmpty)
			t.Logf("Error: %v", err)
		})

		var case4URL = "http://192.168.180.67:1010/host/info"
		Convey(fmt.Sprintf("Wrong URL(wrong host/port): %v", case4URL), func() {
			result, err := GetHTTP(client, case4URL, nil)
			So(err, ShouldNotBeNil)
			So(result, ShouldBeEmpty)
			t.Logf("Error: %v", err)
		})

		var case5URL = "http://localhost:9090/machine/info"
		Convey(fmt.Sprintf("Wrong URL(wrong path): %v", case5URL), func() {
			result, err := GetHTTP(client, case5URL, nil)
			So(err, ShouldNotBeNil)
			So(result, ShouldBeEmpty)
			t.Logf("Error: %v", err)
		})

		var header = SetHeader(map[string]string{
			"token": "123456",
		})
		Convey(fmt.Sprintf("Wrong URL(with header): %v", case1URL), func() {
			t.Logf("Header: %v", header)
			result, err := GetHTTP(client, case1URL, header)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeEmpty)
			t.Logf("Error: %v", err)
		})

		Convey("No client", func() {
			result, err := GetHTTP(nil, case1URL, nil)
			So(err, ShouldNotBeNil)
			So(result, ShouldBeEmpty)
			t.Logf("Data: %v", result)
		})

		Convey("With gzip", func() {
			result, err := GetHTTP(client, case1URL, nil)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeEmpty)
			t.Logf("Data: %v", result)
		})

		var case6URL = "https://localhost:9091/host/info"
		Convey(fmt.Sprintf("With https: %v", case6URL), func() {
			caCert, err := ioutil.ReadFile("example-server/certs/ca.crt")
			So(err, ShouldBeNil)
			clientCert, err := ioutil.ReadFile("example-server/certs/client.crt")
			So(err, ShouldBeNil)
			clientKey, err := ioutil.ReadFile("example-server/certs/client.key")
			So(err, ShouldBeNil)

			client, err = CreateClientHTTP(10*time.Second, caCert, clientCert, clientKey)
			So(err, ShouldBeNil)

			result, err := GetHTTP(client, case6URL, nil)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeEmpty)
			t.Logf("Data: %v", result)
		})
	})
}

func TestPostHTTP(t *testing.T) {
	Convey("Test PostHTTP", t, func() {

		client, err := CreateClientHTTP(10*time.Second, nil, nil, nil)
		So(err, ShouldBeNil)

		var case1URL = "http://localhost:9090/host/info"

		Convey("Case correct: ", func() {
			var form = map[string]string{
				"cpu": "1",
				"ip":  "192.168.180.67",
			}

			result, err := PostHTTP(client, case1URL, SetHeader(nil), form)
			So(err, ShouldBeNil)
			t.Log(result)
		})
	})
}
