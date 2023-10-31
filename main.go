package main

import (
	"filler/controller"
	"flag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	var port = flag.String("port", "80", "指定IP")
	var ip = flag.String("ip", "localhost", "指定IP")

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

	err := server.Run(*ip + ":" + *port)
	if err != nil {
		return
	}
}
