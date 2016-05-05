package maparray

import (
	"errors"
	"math/rand"
	"sync"
)

func NewLimitMapArray(capacity int, scalable bool) *LimitMapArray {
	if capacity > 0 {
		array := &LimitMapArray{
			index:    make(map[string]int),
			length:   0,
			scalable: scalable,
			elements: make([]*element, capacity),
			lock:     &sync.RWMutex{},
		}
		return array
	}
	return nil
}

type LimitMapArray struct {
	index    map[string]int
	length   int
	scalable bool
	elements []*element
	lock     *sync.RWMutex
}

type element struct {
	key   string
	value interface{}
}

func (this *LimitMapArray) Length() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.length
}

func (this *LimitMapArray) Capcity() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return cap(this.elements)
}

func (this *LimitMapArray) Set(key string, value interface{}) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.index[key]; !ok {
		if this.length >= cap(this.elements) {
			if this.scalable {
				this.elements = append(this.elements, this.elements...)
			} else {
				return SetFullAMapArrayErr
			}
		}
		this.index[key] = this.length
		this.elements[this.length] = &element{key, value}
		this.length++
		return nil
	} else {
		this.elements[this.index[key]].value = value
		return nil
	}
}

func (this *LimitMapArray) Get(key string) (interface{}, bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if idx, ok := this.index[key]; ok {
		return this.elements[idx].value, true
	}
	return nil, false
}

func (this *LimitMapArray) Contains(key string) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	_, ok := this.index[key]
	return ok
}

func (this *LimitMapArray) Remove(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	if idx, ok := this.index[key]; ok {
		value := this.elements[idx].value
		delete(this.index, key)
		if idx == this.length-1 {
			this.length--
		} else {
			last := this.elements[this.length-1]
			this.index[last.key] = idx
			this.elements[idx] = last
			this.length--
		}
		return value
	}
	return nil
}

func (this *LimitMapArray) Full() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return !this.scalable && this.length == cap(this.elements)
}

func (this *LimitMapArray) Empty() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.length == 0
}

func (this *LimitMapArray) Random() (string, interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	idx := rand.Int31n(int32(this.length))
	return this.elements[idx].key, this.elements[idx].value
}

func (this *LimitMapArray) Keys() []string {
	i, keys := 0, make([]string, this.length)
	this.lock.RLock()
	defer this.lock.RUnlock()
	for k, _ := range this.index {
		keys[i] = k
		i++
	}
	return keys
}

var (
	SetFullAMapArrayErr = errors.New(`cannot set full map array`)
)
