package maparray

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

const (
	MaxElementsNum = 500
)

type indices struct {
	peerIndex   int
	ispIndex    uint8
	provinIndex uint8
	clsIndex    int
}

func NewLimitMapArray(capacity int) *LimitMapArray {
	return &LimitMapArray{
		indices:    make(map[string]indices),
		classIndex: make(map[uint8]map[uint8][]int),
		length:     0,
		capacity:   capacity,
		elements:   make([]*element, MaxElementsNum),
		lock:       new(sync.RWMutex),
	}
}

type LimitMapArray struct {
	indices    map[string]indices
	classIndex map[uint8]map[uint8][]int
	elements   []*element
	length     int
	capacity   int
	lock       *sync.RWMutex
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
	return this.capacity
}

func (this *LimitMapArray) RealCapacity() int {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return cap(this.elements)
}

func (this *LimitMapArray) Set(key string, value interface{}, isp uint8, province uint8) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.indices[key]; !ok {
		if this.length == this.capacity {
			now := time.Now()
			rand.Seed(int64(now.Nanosecond()))
			peerindex := rand.Intn(this.length)
			key := this.elements[peerindex].key
			// get its indices by key
			idx, ok := this.indices[key]
			if !ok {
				return
			}

			if idx.peerIndex < this.length {
				// first remove its reference in this.classIndex.
				cindex := this.classIndex[idx.ispIndex][idx.provinIndex]
				if idx.clsIndex < len(cindex)-1 {
					last := cindex[len(cindex)-1]
					cindex[idx.clsIndex] = last

					e := this.elements[last]
					i := this.indices[e.key]
					i.clsIndex = idx.clsIndex
					this.indices[e.key] = i

				}
				this.classIndex[idx.ispIndex][idx.provinIndex] = this.classIndex[idx.ispIndex][idx.provinIndex][:len(cindex)-1]

				// then remove itself in this.elements.
				if idx.peerIndex < this.length-1 {
					this.elements[idx.peerIndex] = this.elements[this.length-1]
					i := this.indices[this.elements[this.length-1].key]
					i.peerIndex = idx.peerIndex
					this.indices[this.elements[this.length-1].key] = i

					this.classIndex[i.ispIndex][i.provinIndex][i.clsIndex] = i.peerIndex
				}

				this.length--
				// remove its indice in this.indices.
				delete(this.indices, key)
			} else {
				return
			}

		}

		if this.length == cap(this.elements) {
			var size int = this.length << 1
			if size > this.capacity {
				size = this.capacity
			}
			this.resize(size)
		}

		var i = indices{
			peerIndex:   this.length,
			ispIndex:    isp,
			provinIndex: province,
		}

		this.elements[this.length] = &element{key, value}

		if _, ok := this.classIndex[isp]; !ok {
			this.classIndex[isp] = make(map[uint8][]int)
		}
		if _, ok := this.classIndex[isp][province]; !ok {
			this.classIndex[isp][province] = make([]int, 0)
		}

		this.classIndex[isp][province] = append(this.classIndex[isp][province], i.peerIndex)
		i.clsIndex = len(this.classIndex[isp][province]) - 1
		this.indices[key] = i
		this.length++
		return
	}
}

func (this *LimitMapArray) GetByKey(key string) interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if idx, ok := this.indices[key]; ok {
		if idx.peerIndex >= 0 && idx.peerIndex < this.length {
			return this.elements[idx.peerIndex].value
		}
	}
	return nil
}

func (this *LimitMapArray) GetByIndex(idx int) (string, interface{}) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if idx >= 0 && idx < this.length {
		return this.elements[idx].key, this.elements[idx].value
	}
	return "", nil
}

func shuffle(s []int) {
	t := time.Now()
	rand.Seed(int64(t.Nanosecond()))

	for i := len(s) - 1; i > 0; i-- {
		j := rand.Intn(i)
		s[i], s[j] = s[j], s[i]
	}
}

func (this *LimitMapArray) Randoms(key string, num int) []interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if idx, ok := this.indices[key]; ok {
		cindex := this.classIndex[idx.ispIndex][idx.provinIndex]
		peerindex := make([]int, len(cindex))
		copied := copy(peerindex, cindex)
		shuffle(peerindex)

		var result = make([]interface{}, 0, num)
		var numClass int
		var numOthers int

		if num > copied {
			numClass = copied
			numOthers = num - numClass
		} else {
			numClass = num
			numOthers = 0
		}

		for i := 0; i < numClass; i++ {
			result = append(result, this.elements[peerindex[i]].value)
		}

		now := time.Now().Nanosecond()
		rand.Seed(int64(now))
		appended, offset, i := make([]int, numOthers), 0, 0
		for {
			if i == numOthers || i == this.length-numClass {
				break
			}

			j := rand.Intn(this.length)
			if intContains(peerindex[:numClass], j) || intContains(appended[:offset], j) {
				continue
			}
			result = append(result, this.elements[j].value)
			appended[offset] = j
			offset++
			i++
		}
		return result
	}
	return nil

}

func intContains(s []int, i int) bool {
	for _, v := range s {
		if v == i {
			return true
		}
	}

	return false
}

func (this *LimitMapArray) ContainPeer(key string) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	_, ok := this.indices[key]
	return ok
}

func (this *LimitMapArray) RemoveByKey(key string) (err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	// get its indices by key
	idx, ok := this.indices[key]
	if !ok {
		return errors.New("key not found")
	}

	if idx.peerIndex < this.length {
		// first remove its reference in this.classIndex.
		cindex := this.classIndex[idx.ispIndex][idx.provinIndex]
		if idx.clsIndex < len(cindex)-1 {
			last := cindex[len(cindex)-1]
			cindex[idx.clsIndex] = last

			e := this.elements[last]
			i := this.indices[e.key]
			i.clsIndex = idx.clsIndex
			this.indices[e.key] = i

		}
		this.classIndex[idx.ispIndex][idx.provinIndex] = this.classIndex[idx.ispIndex][idx.provinIndex][:len(cindex)-1]

		// then remove itself in this.elements.
		if idx.peerIndex < this.length-1 {
			this.elements[idx.peerIndex] = this.elements[this.length-1]
			i := this.indices[this.elements[this.length-1].key]
			i.peerIndex = idx.peerIndex
			this.indices[this.elements[this.length-1].key] = i

			this.classIndex[i.ispIndex][i.provinIndex][i.clsIndex] = i.peerIndex
		}

		this.length--
		// remove its indice in this.indices.
		delete(this.indices, key)
		return nil
	} else {
		return errors.New("peerIndex out of range")
	}

}

// TODO: fix bug
func (this *LimitMapArray) resize(size int) {
	newElements := make([]*element, size)
	copy(newElements, this.elements)
	this.elements = newElements
}

func (this *LimitMapArray) Full() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.capacity > 0 && this.length == this.capacity
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
