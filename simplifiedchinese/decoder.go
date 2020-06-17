package simplifiedchinese

import (
	"golang.org/x/text/encoding/simplifiedchinese"
)

//
const (
	UTF8     = "UTF-8"
	GB18030  = "GB18030"
	GBK      = "GBK"
	HZGB2312 = "HZGB2312"
)

// BytesToString 字节转换为字符串
func BytesToString(charset string, buf []byte) string {

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
