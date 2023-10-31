package models

type TaskInfo struct {
	Url          string `json:"url,omitempty"`
	Type         string `json:"type,omitempty"`
	Num          int    `json:"num,omitempty"`
	CurrentNum   int    `json:"current_num,omitempty"`
	SuccessCount int    `json:"success_count,omitempty"`
	FailedCount  int    `json:"failed_count,omitempty"`
}
