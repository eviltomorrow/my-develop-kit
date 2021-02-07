package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func startupServerHTTP(port int) {
	var router = gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.GET("/host/info", GetHostMachineInfo)
	router.POST("/host/info", PostHostMachineInfo)

	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Startup Server HTTP failure, nest error: %v", err)
	}
}
