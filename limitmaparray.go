// maparray
package maparray

import (
	"errors"
	`math/rand`
)

func NewLimitMapArray(capacity int, sRule SelectRuler, cRule CoverRuler) *LimitMapArray {
	if capacity > 0 {
		array := &LimitMapArray{
			index:      make(map[string]int),
			capacity:   capacity,
			freeIndexs: make([]int, capacity),
			busyIndexs: make([]int, capacity),
			elements:   make([]*element, capacity),
			length:     0,
			sRule:      sRule,
			cRule:      cRule,
		}
		for i := 0; i < capacity; i++ {
			array.freeIndexs[i] = i
		}
		return array
	}
	return nil
}

type LimitMapArray struct {
	index      map[string]int
	capacity   int
	length     int
	freeIndexs []int
	busyIndexs []int
	sRule      SelectRuler
	cRule      CoverRuler
	elements   []*element
}

type element struct {
	key   string
	value interface{}
}

func (ra *LimitMapArray) Length() int {
	// TODO bug
	return len(ra.index)
}

func (ra *LimitMapArray) Set(key string, value interface{}) error {
	if !ra.IsFull() {
		idx := ra.freeIndexs[ra.capacity-ra.length-1]
		ra.index[key] = idx
		ra.busyIndexs[ra.length] = idx
		ra.elements[idx] = &element{key, value}
		ra.length++
		return nil
	} else {
		var (
			idx      int
			tmpKey   string
			tmpValue interface{}
			try      int
		)
		for {
			idx = rand.Intn(ra.length)
			try++
			tmpKey = ra.elements[idx].key
			tmpValue = ra.elements[ra.index[tmpKey]].value
			if try >= maxTry || ra.cRule == nil || ra.cRule.ShouldCover(tmpValue) {
				ra.index[key] = idx
				ra.elements[idx] = &element{key, value}
				delete(ra.index, key)
				return nil
			}
		}
	}
	return SetFullAMapArrayErr
}

func (ra *LimitMapArray) Get(key string) interface{} {
	if idx, ok := ra.index[key]; ok {
		return ra.elements[idx].value
	}
	return nil
}

func (ra *LimitMapArray) Remove(key string) interface{} {
	if idx, ok := ra.index[key]; ok {
		value := ra.elements[idx].value
		ra.freeIndexs[ra.capacity-ra.length] = idx
		ra.length--
		delete(ra.index, key)
		return value
	}
	return nil
}

func (ra *LimitMapArray) IsFull() bool {
	return ra.length == ra.capacity
}

func (ra *LimitMapArray) IsEmpty() bool {
	return ra.length == 0
}

func (ra *LimitMapArray) RandomOne() interface{} {
	idx := rand.Int31n(int32(ra.length))
	return ra.elements[ra.busyIndexs[idx]].value
}

func (ra *LimitMapArray) Randoms(limit int, maxTry int) []interface{} {
	if ra.length <= limit {
		values := make([]interface{}, ra.length)
		for i := 0; i < ra.length; i++ {
			values[i] = ra.elements[ra.busyIndexs[i]].value
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
			idx = rand.Intn(ra.length)
			try++
			if _, ok = idxs[idx]; !ok && (ra.sRule == nil || ra.sRule.Check(ra.elements[ra.busyIndexs[idx]])) {
				idxs[idx] = '0'
				values[i] = ra.elements[ra.busyIndexs[idx]]
				i++
			}
			if try >= maxTry {
				return values[:i]
			}
		}
		return values
	}
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

const (
	maxTry = 128
)
