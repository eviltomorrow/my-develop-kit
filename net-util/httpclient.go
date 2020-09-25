package netutil

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// DefaultHTTPHeader default http header
var DefaultHTTPHeader = map[string]string{
	"Connection":                "keep-alive",
	"Cache-Control":             "max-age=0",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"Accept-Encoding":           "gzip, deflate",
	"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8",
}

// GetHTTP get http
func GetHTTP(url string, timeout time.Duration, header map[string]string) (string, error) {
	var client = &http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	header = setHeader(header)
	for key, val := range header {
		request.Header.Add(key, val)
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		buf, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("HTTP status: %v, Body: %v", response.StatusCode, string(buf))

	}

	var buffer []byte
	contentEncode := response.Header.Get("Content-Encoding")
	switch {
	case contentEncode == "gzip":
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return "", err
		}
		defer reader.Close()

		buf, err := ioutil.ReadAll(reader)
		if err != nil {
			return "", err
		}
		buffer = buf
	default:
		buf, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		buffer = buf
	}

	var data string
	contentType := response.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, GB18030):
		data = bytesToString(GB18030, buffer)
	default:
		data = bytesToString(UTF8, buffer)
	}

	return data, nil
}

func setHeader(header map[string]string) map[string]string {
	var data = make(map[string]string, len(DefaultHTTPHeader))
	for k, v := range DefaultHTTPHeader {
		data[k] = v
	}

	for k, v := range header {
		data[k] = v
	}
	return data
}

//
const (
	UTF8     = "UTF-8"
	GB18030  = "GB18030"
	GBK      = "GBK"
	HZGB2312 = "HZGB2312"
)

// BytesToString 字节转换为字符串
func bytesToString(charset string, buf []byte) string {
	var str string
	switch charset {
	case GB18030:
		tmp, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(buf)
		str = string(tmp)
	case GBK:
		tmp, _ := simplifiedchinese.GBK.NewDecoder().Bytes(buf)
		str = string(tmp)
	case HZGB2312:
		tmp, _ := simplifiedchinese.HZGB2312.NewDecoder().Bytes(buf)
		str = string(tmp)
	case UTF8:
		fallthrough
	default:
		str = string(buf)
	}
	return str
}
