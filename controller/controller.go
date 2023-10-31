package controller

import (
	"filler/driver"
	"filler/models"
	"filler/tasks"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func Scan(c *gin.Context) {
	url := c.Query("url")
	li := driver.NewPage(url).GetStructure()
	c.JSON(http.StatusOK, gin.H{"data": li})
	return
}

func Submit(c *gin.Context) {
	var data models.SubmitData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "提交数据格式有误"})
		return
	}

	t := tasks.Start(data)
	c.JSON(http.StatusOK, gin.H{"msg": "启动成功", "data": t, "secret": t.Secret})
	return
}

func Close(c *gin.Context) {
	var data map[string]any
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请求参数错误"})
		return
	}
	taskID := cast.ToInt64(data["task_id"])
	secret := cast.ToString(data["task_secret"])
	task, exist := tasks.GetTask(taskID, secret)
	if !exist {
		c.JSON(http.StatusPreconditionFailed, gin.H{"msg": "任务不存在"})
		return
	}

	if exist && task == nil {
		c.JSON(http.StatusForbidden, gin.H{"msg": "task_id 与 task_secret 不匹配"})
		return
	}

	success := task.Stop()
	if !success {
		c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "服务器错误, 中断失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "中断成功"})
	return
}

func Status(c *gin.Context) {
	fmt.Println(c.Query("task_id"))
	taskID := cast.ToInt64(c.Query("task_id"))
	secret := c.Query("task_secret")
	task, exist := tasks.GetTask(taskID, secret)
	if !exist {
		c.JSON(http.StatusPreconditionFailed, gin.H{"msg": "任务不存在"})
		return
	}

	if exist && task == nil {
		c.JSON(http.StatusForbidden, gin.H{"msg": "task_id 与 task_secret 不匹配"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "", "data": task.GetInfo()})
	return
}
