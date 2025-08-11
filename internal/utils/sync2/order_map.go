package sync2

import "sync"

type order[K comparable, V any] struct {
	key   K
	value V
}

func NewOrderMap[K comparable, V any]() *OrderMap[K, V] {
	m := &OrderMap[K, V]{}
	m.caches = make([]*order[K, V], 0)
	return m
}

type OrderMap[K comparable, V any] struct {
	caches  []*order[K, V]
	mutex   sync.RWMutex
	maxSize int // 添加最大长度
}

func (this *OrderMap[K, V]) SetMaxSize(size int) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.maxSize = size

}

func (this *OrderMap[K, V]) Load(key K) (value V, ok bool) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	for _, cache := range this.caches {
		if cache.key == key {
			value = cache.value
			ok = true
			break
		}
	}
	return
}

func (this *OrderMap[K, V]) Range(f func(key K, value V) bool) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	for _, cache := range this.caches {
		status := f(cache.key, cache.value)
		if !status {
			break
		}
	}
}

func (this *OrderMap[K, V]) Store(key K, value V) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	isExist := false
	for _, cache := range this.caches {
		if cache.key == key {
			cache.value = value
			isExist = true
			break
		}
	}
	if !isExist {
		if this.maxSize > 0 {
			overflow := len(this.caches) + 1 - this.maxSize
			if overflow > 0 {
				this.caches = this.caches[overflow:] // 删除第一个元素
			}
		}
		cache := new(order[K, V])
		cache.key = key
		cache.value = value
		this.caches = append(this.caches, cache)
	}
	return
}

func (this *OrderMap[K, V]) Delete(key K) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for i, cache := range this.caches {
		if cache.key == key {
			this.caches = append(this.caches[:i], this.caches[i+1:]...)
			break
		}
	}
}

func (this *OrderMap[K, V]) Clear() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.caches = make([]*order[K, V], 0)
}
