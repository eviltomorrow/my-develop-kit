package zlog

import (
	"fmt"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	global, prop, err := InitLogger(&Config{
		Level:            "info",
		Format:           "text",
		DisableTimestamp: false,
		File: FileLogConfig{
			Filename: "/tmp/log/zlog/data.log",
			MaxSize:  300,
		},
	})
	if err != nil {
		fmt.Printf("配置日志信息错误，nest error: %v\r\n", err)
		os.Exit(1)
	}
	ReplaceGlobals(global, prop)

	Info("this is shepard")
}
