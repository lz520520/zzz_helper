package lzhttp

import "sync"

func NewLzHeader(header Header) *LzHeader {
	return &LzHeader{header: header}
}

type LzHeader struct {
	header Header
	mutex  sync.RWMutex
}

func (this *LzHeader) Set(key, value string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.header.Set(key, value)
}

func (this *LzHeader) Add(key, value string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.header.Add(key, value)
}

func (this *LzHeader) Get(key string) string {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.header.Get(key)
}

func (this *LzHeader) GetOriginHeader() Header {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	cloneMap := make(Header)
	for k, v := range this.header {
		cloneSlice := make([]string, 0)
		for _, vv := range v {
			cloneSlice = append(cloneSlice, vv)
		}
		cloneMap[k] = cloneSlice
	}
	return cloneMap
}
