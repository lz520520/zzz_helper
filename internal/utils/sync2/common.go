package sync2

import "sync"

func NewSyncCommon[T any](value T) *SyncCommon[T] {
	return &SyncCommon[T]{
		value: value,
		mutex: sync.RWMutex{},
	}
}

type SyncCommon[T any] struct {
	value T
	mutex sync.RWMutex
}

func (s *SyncCommon[T]) SetValue(value T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value = value
}

func (s *SyncCommon[T]) GetValue() T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.value
}
