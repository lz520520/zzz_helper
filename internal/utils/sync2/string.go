package sync2

import "sync"

type SyncString struct {
	value string
	mutex sync.RWMutex
}

func NewSyncString(value string) SyncString {
	return SyncString{
		value: value,
		mutex: sync.RWMutex{},
	}
}

func (s *SyncString) SetValue(value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.value = value
}

func (s *SyncString) Append(value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value += value
}

func (s *SyncString) AddLine(value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.value == "" {
		s.value += value
	} else {
		s.value += "\n" + value
	}
}
func (s *SyncString) GetValue() (value string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.value
}
