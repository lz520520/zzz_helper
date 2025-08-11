package common_model

type TaskStatus struct {
	TaskType  string         `json:"task_type"` // portscan/brute
	Status    string         `json:"status"`    // start/pause/stop
	TaskId    string         `json:"task_id"`
	QuitCh    chan bool      `json:"-"`
	PauseFlag *int32         `json:"-"`
	Params    map[string]any `json:"params"`
}

type TaskPauseReq struct {
	TaskId string `json:"task_id"`
	Status string `json:"status"`
}

type TasksResp struct {
	Status bool         `json:"status"`
	Tasks  []TaskStatus `json:"tasks"`
	Err    string       `json:"err"`
}
type TaskStatusResp struct {
	Status     bool       `json:"status"`
	TaskStatus TaskStatus `json:"task_status"`
	Err        string     `json:"err"`
}
