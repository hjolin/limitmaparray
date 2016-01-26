// maparray
package maparray

import (
	"errors"
	"sync"
	//"fmt"
	"math/rand"
)

func NewLimitMapArray(capacity int, scalable bool) *LimitMapArray {
	if capacity > 0 {
		array := &LimitMapArray{
			index:    make(map[string]int),
			capacity: capacity,
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
	index       map[string]int
	capacity    int
	length      int
	coverMaxTry int
	scalable    bool
	sRule       SelectRuler
	cRule       CoverRuler
	elements    []*element
	lock        *sync.RWMutex
}

type element struct {
	key   string
	value interface{}
}

func (this *LimitMapArray) SetSelectRuler(sRule SelectRuler) {
	this.sRule = sRule
}

func (this *LimitMapArray) SetCoverRuler(cRule CoverRuler) {
	this.cRule = cRule
}

func (this *LimitMapArray) SetCoverMaxTry(coverMaxTry int) {
	this.coverMaxTry = coverMaxTry
}

func (this *LimitMapArray) Length() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.length
}

func (this *LimitMapArray) Set(key string, value interface{}) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.index[key]; !ok {
		if this.length != this.capacity {
			this.index[key] = this.length
			this.elements[this.length] = &element{key, value}
			this.length++
			return nil
		} else if this.scalable {
			this.capacity = this.capacity << 1
			this.elements = append(this.elements, this.elements...)
			this.index[key] = this.length
			this.elements[this.length] = &element{key, value}
			this.length++
		} else {
			var (
				idx int
				ele *element
				try int
			)
			for {
				idx = rand.Intn(this.capacity)
				try++
				ele = this.elements[idx]
				if try >= this.coverMaxTry || this.cRule == nil || this.cRule.ShouldCover(ele.value) {
					delete(this.index, ele.key)
					this.index[key] = idx
					this.elements[idx] = &element{key, value}
					return nil
				}
			}
		}
	} else {
		this.elements[this.index[key]].value = value
	}
	return SetFullAMapArrayErr
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

func (this *LimitMapArray) IsFull() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.length == this.capacity
}

func (this *LimitMapArray) IsEmpty() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.length == 0
}

func (this *LimitMapArray) RandomOne() interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	idx := rand.Int31n(int32(this.length))
	return this.elements[idx].value
}

func (this *LimitMapArray) Randoms(limit int, maxTry int) []interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if this.length < limit {
		values := make([]interface{}, this.length)
		for i := 0; i < this.length; i++ {
			values[i] = this.elements[i].value
		}
		return values
	} else {
		idxs := make(map[int]byte, limit)
		values := make([]interface{}, limit)
		var (
			idx int
			ok  bool
			i   int
			try int
		)
		for i < limit {
			idx = rand.Intn(this.length)
			try++
			if _, ok = idxs[idx]; !ok && (try >= maxTry || this.sRule == nil || this.sRule.Check(this.elements[idx].value)) {
				idxs[idx] = '0'
				values[i] = this.elements[idx].value
				i++
			}
		}
		return values
	}
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

type SelectRuler interface {
	Check(interface{}) bool
}

type CoverRuler interface {
	ShouldCover(interface{}) bool
}

var (
	SetFullAMapArrayErr = errors.New(`cannot set full map array`)
)
