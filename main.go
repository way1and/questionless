package main

import (
	"filler/controller"
	"filler/driver"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	var port int
	var ip string
	var driverPort int
	flag.IntVar(&port, "port", 80, "指定端口号, 端口号应为数字")
	flag.StringVar(&ip, "ip", "localhost", "指定IP 以-间隔")
	flag.IntVar(&driverPort, "dport", 81, "指定端口号")
	driver.PORT = driverPort
	flag.Parse()

	server := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	server.Use(cors.New(config))
	server.Use(cors.Default())

	// 静态
	server.Static("/static", "./static")
	server.StaticFile("/", "./static/index.html")
	server.StaticFile("/index", "./static/index.html")

	// 动态
	server.GET("/scan", controller.Scan)
	server.GET("/status", controller.Status)
	server.POST("/submit", controller.Submit)
	server.POST("/close", controller.Close)

	err := server.Run(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return
	}
}
