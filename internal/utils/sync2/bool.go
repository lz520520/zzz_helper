package sync2

import "sync"

type SyncBool struct {
	value bool
	mutex sync.RWMutex
}

func NewSyncBool(value bool) SyncBool {
	return SyncBool{
		value: value,
		mutex: sync.RWMutex{},
	}
}

func (s *SyncBool) SetValue(value bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.value = value
}

func (s *SyncBool) GetValue() (value bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.value
}

func (s *SyncBool) GetOrSetValue(new bool) (old bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	old = s.value
	if new != old {
		s.value = new
	}
	return old
}
