package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func startupServerHTTPS(port int, certPath, keyPath string) {
	var router = gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.GET("/host/info", GetHostMachineInfo)
	router.POST("/host/info", PostHostMachineInfo)

	err := router.RunTLS(fmt.Sprintf(":%d", port), certPath, keyPath)
	if err != nil {
		log.Fatalf("Startup Server HTTPS failure, nest error: %v", err)
	}
}
