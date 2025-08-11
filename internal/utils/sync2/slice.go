package sync2

import (
	"fmt"
	"sync"
)

type SliceType interface {
	GetId() string
	SetId(id string)
}

func NewSliceString(meta string) *SliceString {
	return &SliceString{meta: meta}
}

type SliceString struct {
	meta string
}

func (this *SliceString) GetId() string {
	return this.meta
}
func (this *SliceString) SetId(id string) {
	this.meta = id
}

type SyncSlice[T SliceType] struct {
	values []T
	mutex  sync.RWMutex
}

func NewSyncSlice[T SliceType]() *SyncSlice[T] {
	return &SyncSlice[T]{
		values: make([]T, 0),
		mutex:  sync.RWMutex{},
	}
}

func (s *SyncSlice[T]) SetValue(value T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	isExist := false
	for i, v := range s.values {
		if v.GetId() == value.GetId() {
			s.values[i] = value
			isExist = true
			break
		}
	}
	if !isExist {
		s.values = append(s.values, value)
	}

}

func (s *SyncSlice[T]) SwitchValue(from, to string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 获取所有id
	ids := make([]string, 0)
	for _, value := range s.values {
		ids = append(ids, value.GetId())
	}

	var fromValue T
	var ok bool
	// 将原数据提取并在slice中删除
	for i, v := range s.values {
		if v.GetId() == from {
			s.values = append(s.values[:i], s.values[i+1:]...)
			fromValue = v
			ok = true
			break
		}
	}

	if !ok {
		err = fmt.Errorf("%s is not found", from)
		return
	}

	// 将fromValue插入到新slice中
	ok = false
	for i, v := range s.values {
		if v.GetId() == to {
			pre := s.values[:i]
			suf := s.values[i:]
			tmp := append([]T{}, pre...)
			tmp2 := append(tmp, []T{fromValue}...)
			tmp3 := append(tmp2, suf...)
			s.values = tmp3
			ok = true
			break
		}
	}

	if !ok {
		err = fmt.Errorf("%s is not found", to)
		return
	}

	// 更新id
	for i, v := range s.values {
		v.SetId(ids[i])
	}
	return
}

// 插入新值的同时，移除重复的值
func (s *SyncSlice[T]) SetValueAndRemoveDup(value T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newValues := make([]T, 0)
	for _, v := range s.values {
		if v.GetId() != value.GetId() {
			newValues = append(newValues, v)
		}
	}
	newValues = append(newValues, value)
	s.values = newValues

}

func (s *SyncSlice[T]) DelValue(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, v := range s.values {
		if v.GetId() == id {
			s.values = append(s.values[:i], s.values[i+1:]...)
			break
		}
	}
}

func (s *SyncSlice[T]) AddValue(value T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.values = append(s.values, value)
}

func (s *SyncSlice[T]) GetAllValueAndDelete() []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	result := make([]T, 0)
	for _, v := range s.values {
		exist := false
		for i, t := range result {
			if v.GetId() == t.GetId() {
				result[i] = v
				exist = true
			}
		}
		if !exist {
			result = append(result, v)
		}
	}
	s.values = s.values[0:0]
	return result
}

func (s *SyncSlice[T]) GetValue(id string) (value T, ok bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, v := range s.values {
		if v.GetId() == id {
			value = v
			ok = true
			break
		}
	}
	return
}

// 获取指定值并移除
func (s *SyncSlice[T]) GetValueAndDelete(id string) (value T, ok bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, v := range s.values {
		if v.GetId() == id {
			s.values = append(s.values[:i], s.values[i+1:]...)
			value = v
			ok = true
			break
		}
	}
	return
}

// 获取所有值，这里可以限制获取数量，是否倒序
func (s *SyncSlice[T]) GetAllValue(count int, reverse bool) (value []T) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	tmp := make([]T, 0)
	value = make([]T, 0)
	// 这里做了限制，只获取100个
	if count > 0 {
		if len(s.values) < count {
			tmp = s.values
		} else {
			tmp = s.values[len(s.values)-count:]
		}
	} else {
		tmp = s.values
	}
	if reverse {
		for i := len(tmp) - 1; i >= 0; i-- {
			value = append(value, tmp[i])
		}
	} else {
		value = tmp
	}

	return
}

func (s *SyncSlice[T]) SetAllValue(value []T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.values = value
}

func (s *SyncSlice[T]) GetAllId(count int, reverse bool) (ids []string) {
	ids = make([]string, 0)
	values := s.GetAllValue(count, reverse)
	for _, value := range values {
		ids = append(ids, value.GetId())
	}
	return
}
