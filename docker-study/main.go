package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		var i = 0
		for {
			time.Sleep(2 * time.Second)
			i++
			fmt.Printf("i: %d\r\n", i)
		}
	}()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080")
}
