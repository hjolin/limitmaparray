package maparray

import (
	"errors"
	"math/rand"
	"sync"
)

func NewLimitMapArray(capacity int) *LimitMapArray {
	array := &LimitMapArray{
		index:    make(map[string]int),
		length:   0,
		capcity:  capacity,
		elements: make([]*element, 1),
		lock:     &sync.RWMutex{},
	}
	return array
}

type LimitMapArray struct {
	index    map[string]int
	elements []*element
	length   int
	capcity  int
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

func (this *LimitMapArray) Capacity() int {
	return this.capcity
}

func (this *LimitMapArray) RealCapacity() int {
	return cap(this.elements)
}

func (this *LimitMapArray) Set(key string, value interface{}) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.index[key]; !ok {
		if this.capcity > 0 && this.length >= this.capcity {
			idx := rand.Intn(this.length)
			delete(this.index, this.elements[idx].key)
			this.index[key] = idx
			this.elements[idx].key = key
			this.elements[idx].value = value
			return LimitMapArrayFullErr
		}
		if this.length >= cap(this.elements) {
			var size int = this.length << 1
			if this.capcity > 0 && size > this.capcity {
				size = this.capcity
			}
			this.resize(size)
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

func (this *LimitMapArray) GetByKey(key string) interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if idx, ok := this.index[key]; ok {
		return this.elements[idx].value
	}
	return nil
}

func (this *LimitMapArray) GetByIndex(idx int) (string, interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if idx >= 0 && idx < this.length {
		return this.elements[idx].key, this.elements[idx].value
	}
	return ``, nil
}

func (this *LimitMapArray) Contains(key string) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	_, ok := this.index[key]
	return ok
}

func (this *LimitMapArray) RemoveByKey(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	if idx, ok := this.index[key]; ok {
		element := this.remove(idx)
		return element.value
	}
	return nil
}

func (this *LimitMapArray) RemoveByIndex(idx int) (string, interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if element := this.remove(idx); element != nil {
		return element.key, element.value
	}
	return ``, nil
}

func (this *LimitMapArray) remove(idx int) *element {
	if idx >= 0 && idx < this.length {
		element := this.elements[idx]
		delete(this.index, element.key)
		if idx != this.length-1 {
			last := this.elements[this.length-1]
			this.index[last.key] = idx
			this.elements[idx] = last
		}
		this.length--
		if this.length < cap(this.elements)>>1 {
			this.resize(cap(this.elements) >> 1)
		}
		return element
	}
	return nil
}

func (this *LimitMapArray) resize(capcity int) {
	newElements := make([]*element, capcity)
	copy(newElements, this.elements)
	this.elements = newElements
}

func (this *LimitMapArray) Full() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.capcity > 0 && this.length == this.capcity
}

func (this *LimitMapArray) Empty() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.length == 0
}

func (this *LimitMapArray) Keys() []string {
	this.lock.RLock()
	defer this.lock.RUnlock()
	values := make([]string, this.length)
	for i := 0; i < this.length; i++ {
		values[i] = this.elements[i].key
	}
	return values
}

func (this *LimitMapArray) Values() []interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	values := make([]interface{}, this.length)
	for i := 0; i < this.length; i++ {
		values[i] = this.elements[i].value
	}
	return values
}

var LimitMapArrayFullErr = errors.New(`cannot set full map array`)
