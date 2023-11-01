package controller

import (
	"filler/driver"
	"filler/models"
	"filler/tasks"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
)

func Scan(c *gin.Context) {
	url := c.Query("url")
	li := driver.NewPage(url).GetStructure()
	if li == nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "本机无法连接到问卷星或该问卷已关闭"})
		return
	}

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
