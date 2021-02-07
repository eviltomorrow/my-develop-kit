package httpclient

import (
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	nurl "net/url"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

//
const (
	UTF8     = "UTF-8"
	GB18030  = "GB18030"
	GBK      = "GBK"
	HZGB2312 = "HZGB2312"
)

// DefaultHeader default header
var DefaultHeader = map[string]string{
	"Connection":                "keep-alive",
	"Cache-Control":             "max-age=0",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36",
	"Accept":                    "application/json,text/html,application/xml",
	"Accept-Encoding":           "gzip, deflate",
	"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8",
}

// CreateClientHTTP create http client
func CreateClientHTTP(timeout time.Duration, certCA []byte, certClient []byte, keyClient []byte) (*http.Client, error) {
	if len(certCA) == 0 {
		return &http.Client{Timeout: timeout}, nil
	}

	if len(certClient) == 0 || len(keyClient) == 0 {
		return nil, fmt.Errorf("Invalid client cert/key")
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(certCA); !ok {
		return nil, fmt.Errorf("Append CA cert failure, please check CA cert")
	}

	cert, err := tls.X509KeyPair(certClient, keyClient)
	if err != nil {
		return nil, fmt.Errorf("Load client cert/key failure, nest error: %v", err)
	}

	var client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
		Timeout: timeout,
	}
	return client, nil
}

// SetHeader set header
func SetHeader(header map[string]string) map[string]string {
	var data = clone(DefaultHeader)
	for key, val := range header {
		data[key] = val
	}
	return data
}

// GetHTTP get http for request
func GetHTTP(client *http.Client, url string, header map[string]string) (string, error) {
	return do(client, "GET", url, header, nil)
}

// PostHTTP post http for request
func PostHTTP(client *http.Client, url string, header map[string]string, form map[string]string) (string, error) {
	var val = nurl.Values{}
	for k, v := range form {
		val.Set(k, v)
	}

	if header != nil {
		_, ok := header["Content-Type"]
		if !ok {
			header["Content-Type"] = "application/x-www-form-urlencoded; charset=utf8"
		}
	}
	return do(client, "POST", url, header, ioutil.NopCloser(strings.NewReader(val.Encode())))
}

// PutHTTP put http for request
func PutHTTP(client *http.Client, url string, header map[string]string, form map[string]string) (string, error) {
	var val = nurl.Values{}
	for k, v := range form {
		val.Set(k, v)
	}
	if header != nil {
		_, ok := header["Content-Type"]
		if !ok {
			header["Content-Type"] = "application/x-www-form-urlencoded; charset=utf8"
		}
	}
	return do(client, "PUT", url, header, ioutil.NopCloser(strings.NewReader(val.Encode())))
}

// DeleteHTTP delete http for request
func DeleteHTTP(client *http.Client, url string, header map[string]string, form map[string]string) (string, error) {
	var val = nurl.Values{}
	for k, v := range form {
		val.Set(k, v)
	}
	if header != nil {
		_, ok := header["Content-Type"]
		if !ok {
			header["Content-Type"] = "application/x-www-form-urlencoded; charset=utf8"
		}
	}
	return do(client, "DELETE", url, header, ioutil.NopCloser(strings.NewReader(val.Encode())))
}

// do
func do(client *http.Client, method string, url string, header map[string]string, body io.Reader) (string, error) {
	if client == nil {
		return "", fmt.Errorf("Invalid http client, nest error: client is nil")
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", fmt.Errorf("New request failure, nest error: %v", err)
	}

	for key, val := range header {
		request.Header.Add(key, val)
	}

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("Do http request failure, nest error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		buf, _ := ioutil.ReadAll(response.Body)
		return "", fmt.Errorf("HTTP status code: %d, message: %s", response.StatusCode, buf)
	}

	var buffer []byte
	contentEncode := response.Header.Get("Content-Encoding")
	switch {
	case contentEncode == "gzip":
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return "", fmt.Errorf("Panic: gzip NewReader failure, nest error: %v", err)
		}
		defer reader.Close()

		buf, err := ioutil.ReadAll(reader)
		if err != nil {
			return "", fmt.Errorf("Panic: ReadAll failure, nest error: %v", err)
		}
		buffer = buf
	default:
		buf, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("ReadAll response body failure: %v", err)
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

func clone(data map[string]string) map[string]string {
	var newObj = make(map[string]string, len(data))
	for key, val := range data {
		newObj[key] = val
	}
	return newObj
}

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
