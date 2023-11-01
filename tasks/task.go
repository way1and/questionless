package tasks

import (
	"filler/driver"
	"filler/models"
	"filler/utils"
	"time"
)

var Tasks = make(map[int64]*models.Task)
var StatusNoTask = true

func Start(data models.SubmitData) *models.Task {
	task := new(models.Task)
	task.ID = time.Now().Unix()
	task.Secret = utils.RandomString(32)

	Tasks[task.ID] = task
	task.Info = models.TaskInfo{
		Url:          data.Url,
		Type:         data.Type,
		Num:          data.Num,
		CurrentNum:   0,
		SuccessCount: 0,
		FailedCount:  0,
		Running:      true,
	}
	go driver.Exec(task, data.Data)

	if StatusNoTask {
		go DeleteFreeTask()
	}

	return task
}

func GetTask(id int64, secret string) (*models.Task, bool) {
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

func DeleteFreeTask() {
	// 轮询 删除任务
	StatusNoTask = false
	d, _ := time.ParseDuration("60s")
	for {
		time.Sleep(d)
		for id, task := range Tasks {
			if task == nil {
				continue
			} else if Tasks[id].Info.Running == false || Tasks[id].Finished {
				delete(Tasks, id)
			}
		}
		break
	}
	// 恢复无任务状态
	StatusNoTask = true
}
