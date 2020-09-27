package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

// HostMachine host machine
type HostMachine struct {
	Desc   string `json:"desc"`
	OS     string `json:"os"`
	ARCH   string `json:"arch"`
	CPU    int    `json:"cpu"`
	Memory int64  `json:"memory"`
	IP     string `json:"ip"`
	Disk   string `json:"disk"`
}

func (hm *HostMachine) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(hm)
	return string(buf)
}

var instance = &HostMachine{
	Desc:   "HTTP Plugin Test Server.",
	OS:     "CentOS 7.2 x64",
	ARCH:   "x64",
	CPU:    4,
	Memory: 1024 * 1024 * 1024 * 16,
	IP:     "192.168.11.157",
	Disk:   `{"Type":"SSD","Size":"250GB"}`,
}

// GetHostMachineInfo get host machine info
func GetHostMachineInfo(ctx *gin.Context) {
	for k, v := range ctx.Request.Header {
		log.Printf("%v: %v\r\n", k, v)
	}

	log.Printf("Method: %v\r\n", ctx.Request.Method)
	log.Printf("Protocol: %v\r\n", ctx.Request.Proto)
	ctx.String(http.StatusOK, instance.String())
}

// PostHostMachineInfo post host machine info
func PostHostMachineInfo(ctx *gin.Context) {
	log.Printf("Method: %v\r\n", ctx.Request.Method)
	ctx.Request.ParseForm()
	fmt.Println(ctx.Request.Form)
	var cpu = ctx.DefaultPostForm("cpu", "4")
	var ip = ctx.DefaultPostForm("ip", "localhost")

	cpuInt, err := strconv.Atoi(cpu)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	instance.CPU = cpuInt
	instance.IP = ip
	ctx.String(http.StatusOK, instance.String())
}
