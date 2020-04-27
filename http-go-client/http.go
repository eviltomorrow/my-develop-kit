package client

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPGetURLWithCookie 发送get请求
func HTTPGetURLWithCookie(addr string, params map[string]string, header map[string]string, cookie map[string]string) ([]byte, error) {
	var client = &http.Client{}
	request, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range header {
		request.Header.Add(key, val)
	}

	for key, val := range cookie {
		c := &http.Cookie{
			Name:     key,
			Value:    val,
			HttpOnly: true,
		}
		request.AddCookie(c)
	}

	if params != nil {
		var values = url.Values{}
		for key, val := range params {
			values.Add(key, val)
		}
		data := values.Encode()
		request.URL.RawQuery = data
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		encode := response.Header.Get("Content-Encoding")
		switch {
		case encode == "gzip":
			reader, err := gzip.NewReader(response.Body)
			if err != nil {
				return nil, err
			}
			defer reader.Close()
			return ioutil.ReadAll(reader)

		default:
			return ioutil.ReadAll(response.Body)
		}
	} else {
		buf, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("解析报文失败[%v]，nest erro r: cde = %v", err, response.StatusCode)
		}
		return nil, fmt.Errorf("解析报文成功[%s]，nest erro r: cde = %v", buf, response.StatusCode)
	}
}

// HTTPGetURLForCookie reload cookie
func HTTPGetURLForCookie(addr string, params map[string]string, header map[string]string) (map[string]string, error) {
	var cookie = make(map[string]string, 16)
	var client = &http.Client{}

	request, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range header {
		request.Header.Add(key, val)
	}

	if params != nil {
		var values = url.Values{}
		for key, val := range params {
			values.Add(key, val)
		}
		data := values.Encode()
		request.URL.RawQuery = data
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		for _, c := range response.Cookies() {
			cookie[c.Name] = c.Value
		}
		return cookie, nil
	}

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("解析报文失败[%v]，nest erro r: cde = %v", err, response.StatusCode)
	}
	return nil, fmt.Errorf("解析报文成功[%s]，nest erro r: cde = %v", buf, response.StatusCode)

}
