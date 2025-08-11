package main_control

import (
	"sync/atomic"
	"zzz_helper/modules/module_common/common_model"
)

func (this *Control) TaskPause(req common_model.TaskPauseReq) (resp common_model.CommonResp) {
	if v, ok := this.taskCaches.Load(req.TaskId); ok {
		if req.Status == "暂停" {
			atomic.StoreInt32(v.PauseFlag, 1)
		} else if req.Status == "继续" {
			atomic.StoreInt32(v.PauseFlag, 0)
		}
		resp.Status = true
	} else {
		resp.Err = "task is not found"
	}
	return
}

func (this *Control) TaskStop(req common_model.CommonReq) (resp common_model.CommonResp) {
	if v, ok := this.taskCaches.Load(req.Msg); ok {
		select {
		case <-v.QuitCh:
		default:
			close(v.QuitCh)
		}
		resp.Status = true
	} else {
		resp.Err = "task is not found"
	}
	return
}

func (this *Control) GetTasks(req common_model.CommonReq) (resp common_model.TasksResp) {
	resp.Tasks = make([]common_model.TaskStatus, 0)
	this.taskCaches.Range(func(key string, value common_model.TaskStatus) bool {
		if value.TaskType == req.Msg {
			resp.Tasks = append(resp.Tasks, value)
		}
		return true
	})
	resp.Status = true
	return
}

func (this *Control) GetTaskStatus(req common_model.CommonReq) (resp common_model.TaskStatusResp) {
	if v, ok := this.taskCaches.Load(req.Msg); ok {
		resp.TaskStatus = v
		resp.Status = true
	} else {
		resp.Err = "task is not found"
	}

	return
}
