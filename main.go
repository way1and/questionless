package main

import (
	"filler/controller"
	"filler/driver"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"strings"
	"time"
)

func main() {

	var port int
	var ip string
	var driverPORTS string
	flag.IntVar(&port, "port", 80, "指定端口号, 端口号应为数字")
	flag.StringVar(&ip, "ip", "localhost", "指定IP 以-间隔")
	flag.StringVar(&driverPORTS, "dports", "9000-9010-9020-9030-9040-9050-9060-9070-9080-9090", "指定driver口号")

	ports := strings.Split(strings.TrimSpace(driverPORTS), "-")
	for _, p := range ports {
		port := cast.ToInt(p)
		if port == 0 {
			fmt.Println("dport指定有误, 不能包含数字外参数, 多个dport以`-`分割 你的输入为", driverPORTS)
			return
		}
		if !driver.CheckPort(port) {
			fmt.Println("dPorts 中包含被占用的端口:", p)
			return
		}
		driver.AddPORT(port)
	}

	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
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

	go PrintInfo(ip, port)
	err := server.Run(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println(err)
		return
	}

}

func PrintInfo(ip, port any) {
	duration, _ := time.ParseDuration("1s")
	time.Sleep(duration)
	fmt.Println()
	fmt.Println(`   ____                  _   _             _               
  / __ \                | | (_)           | |              
 | |  | |_   _  ___  ___| |_ _  ___  _ __ | | ___  ___ ___ 
 | |  | | | | |/ _ \/ __| __| |/ _ \| '_ \| |/ _ \/ __/ __|
 | |__| | |_| |  __/\__ \ |_| | (_) | | | | |  __/\__ \__ \
  \___\_\\__,_|\___||___/\__|_|\___/|_| |_|_|\___||___/___/
                                                            `)
	fmt.Println(`Version: v1.0.0 stable`)
	fmt.Println("Author: Way1and Yang")
	fmt.Println("Github: https://github.com/Way1and/questionless")
	fmt.Println("Gin server running in:", fmt.Sprintf("http://%v:%v", ip, port))
	fmt.Printf("\033[1;32;40m%s\033[0m\n", "All service ready.")
	fmt.Println()
	fmt.Println("Application URL:", fmt.Sprintf("http://%v:%v/index", ip, port), " Click open in browser.")
	return
}
