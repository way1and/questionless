package models

type Task struct {
	ID       int64    `json:"id,omitempty"`
	Secret   string   `json:"-"`
	Running  bool     `json:"running,omitempty"`
	Finished bool     `json:"finished"`
	Info     TaskInfo `json:"info"`
}

func (task *Task) GetInfo() TaskInfo {
	return task.Info
}

func (task *Task) Stop() bool {
	task.Running = false
	return true
}
