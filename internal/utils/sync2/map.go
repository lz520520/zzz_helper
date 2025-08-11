package sync2

import (
	"reflect"
	"sync"
)

type SyncMap struct {
	sync.Map
}

func (this *SyncMap) HasValue(checkValue any) bool {
	isExist := false
	this.Range(func(key, value any) bool {
		if reflect.TypeOf(value).String() == reflect.TypeOf(checkValue).String() {
			if checkValue == value {
				isExist = true
				return false
			}
		}

		return true
	})
	return isExist
}

func (this *SyncMap) Length() int {
	count := 0
	this.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}
