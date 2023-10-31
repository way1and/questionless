package models

type Task struct {
	ID       int64    `json:"id,omitempty"`
	Secret   string   `json:"-"`
	Finished bool     `json:"finished"`
	Info     TaskInfo `json:"info"`
}

func (task *Task) GetInfo() TaskInfo {
	return task.Info
}

func (task *Task) Stop() bool {
	task.Info.Running = false
	return true
}
