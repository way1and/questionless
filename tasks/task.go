package tasks

import (
	"filler/driver"
	"filler/models"
	"filler/utils"
	"fmt"
	"time"
)

var Tasks = make(map[int64]*models.Task)

func Start(data models.SubmitData) *models.Task {
	task := new(models.Task)
	task.ID = time.Now().Unix()
	task.Running = true
	task.Secret = utils.RandomString(32)

	Tasks[task.ID] = task
	task.Info = models.TaskInfo{
		Url:          data.Url,
		Type:         data.Type,
		Num:          data.Num,
		CurrentNum:   0,
		SuccessCount: 0,
		FailedCount:  0,
	}
	go driver.Exec(task, data.Data)

	return task
}

func GetTask(id int64, secret string) (*models.Task, bool) {
	fmt.Println("Get INFO 任务", id, "密钥", secret)
	task := Tasks[id]
	if task == nil {
		return nil, false
	}

	if task.Secret == secret {
		return task, true
	} else {
		// 密钥错误
		return nil, true
	}
}

func DeleteTask(id int64) {
	delete(Tasks, id)
}
