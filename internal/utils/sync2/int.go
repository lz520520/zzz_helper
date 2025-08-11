package sync2

import (
	"sync"
	"sync/atomic"
)

type SyncUint16 struct {
	value uint16
	mutex sync.RWMutex
}

func NewSyncUint16(value uint16) SyncUint16 {
	return SyncUint16{
		value: value,
		mutex: sync.RWMutex{},
	}
}

func (s *SyncUint16) SetValue(value uint16) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value = value
}

func (s *SyncUint16) GetValue() (value uint16) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.value
}

func (s *SyncUint16) Add(delta uint16) uint16 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value += delta
	return s.value
}

type SyncInt32 struct {
	value int32
	mutex sync.RWMutex
}

func NewSyncInt32(value int32) SyncInt32 {
	return SyncInt32{
		value: value,
		mutex: sync.RWMutex{},
	}
}

func (s *SyncInt32) SetValue(value int32) {
	atomic.StoreInt32(&s.value, value)
}

func (s *SyncInt32) GetValue() (value int32) {
	return atomic.LoadInt32(&s.value)
}

func (s *SyncInt32) Add(delta int32) {
	atomic.AddInt32(&s.value, delta)
}

type SyncUint64 struct {
	value uint64
}

func NewSyncUint64(value uint64) *SyncUint64 {
	return &SyncUint64{
		value: value,
	}
}

func (s *SyncUint64) SetValue(value uint64) {
	atomic.StoreUint64(&s.value, value)
}

func (s *SyncUint64) GetValue() (value uint64) {
	return atomic.LoadUint64(&s.value)
}

func (s *SyncUint64) Add(delta uint64) uint64 {
	return atomic.AddUint64(&s.value, delta)
}

//
//type SyncUint32 struct {
//	value uint32
//	mutex sync.RWMutex
//}
//
//func NewSyncUint32(value uint32) SyncUint32 {
//	return SyncUint32{
//		value: value,
//		mutex: sync.RWMutex{},
//	}
//}
//
//func (s *SyncUint32) SetValue(value uint32) {
//	atomic(&s.value, value)
//}
//
//func (s *SyncUint32) GetValue() (value uint32) {
//	return atomic.LoadUint32(&s.value)
//}
//
//func (s *SyncUint32) Add(delta uint32) {
//	atomic.AddUint32(&s.value, delta)
//}
