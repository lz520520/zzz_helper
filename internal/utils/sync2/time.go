package sync2

import (
	"sync"
	"time"
)

type SyncTime struct {
	t     time.Time
	mutex sync.RWMutex
}

func NewSyncTime() *SyncTime {
	return &SyncTime{t: time.Now()}
}
func (this *SyncTime) GetTime() time.Time {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.t
}

func (this *SyncTime) RefreshTime() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.t = time.Now()
}
func (this *SyncTime) SetTime(t time.Time) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.t = t
}
func (this *SyncTime) GetInterval() time.Duration {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return time.Now().Sub(this.t)
}
