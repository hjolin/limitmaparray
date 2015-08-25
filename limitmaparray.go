// maparray
package maparray

import (
	"errors"
	`math/rand`
)

func NewLimitMapArray(capacity int, scalable bool) *LimitMapArray {
	if capacity > 0 {
		array := &LimitMapArray{
			index:       make(map[string]int),
			capacity:    capacity,
			length:      0,
			scalable:    scalable,
			coverMaxTry: CoverMaxTry,
			elements:    make([]*element, capacity),
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
}

type element struct {
	key   string
	value interface{}
}

func (ra *LimitMapArray) SetSelectRuler(sRule SelectRuler) {
	ra.sRule = sRule
}

func (ra *LimitMapArray) SetCoverRuler(cRule CoverRuler) {
	ra.cRule = cRule
}

func (ra *LimitMapArray) SetCoverMaxTry(coverMaxTry int) {
	ra.coverMaxTry = coverMaxTry
}

func (ra *LimitMapArray) Length() int {
	return ra.length
}

func (ra *LimitMapArray) Set(key string, value interface{}) error {
	if _, ok := ra.index[key]; !ok {
		if !ra.IsFull() {
			ra.index[key] = ra.length
			ra.elements[ra.length] = &element{key, value}
			ra.length++
			return nil
		} else if ra.scalable {
			ra.capacity = ra.capacity << 1
			ra.elements = append(ra.elements, ra.elements...)
			ra.index[key] = ra.length
			ra.elements[ra.length] = &element{key, value}
			ra.length++
		} else {
			var (
				idx int
				ele *element
				try int
			)
			for {
				idx = rand.Intn(ra.capacity)
				try++
				ele = ra.elements[idx]
				if try >= ra.coverMaxTry || ra.cRule == nil || ra.cRule.ShouldCover(ele.value) {
					delete(ra.index, ele.key)
					ra.index[key] = idx
					ra.elements[idx] = &element{key, value}
					return nil
				}
			}
		}
	} else {
		ra.elements[ra.index[key]].value = value
	}
	return SetFullAMapArrayErr
}

func (ra *LimitMapArray) Get(key string) (interface{}, bool) {
	if idx, ok := ra.index[key]; ok {
		return ra.elements[idx].value, true
	}
	return nil, false
}

func (ra *LimitMapArray) Contains(key string) bool {
	_, ok := ra.index[key]
	return ok
}

func (ra *LimitMapArray) Remove(key string) interface{} {
	if idx, ok := ra.index[key]; ok {
		value := ra.elements[idx].value
		delete(ra.index, key)
		if idx == ra.length-1 {
			ra.length--
		} else {
			last := ra.elements[ra.length-1]
			ra.index[last.key] = idx
			ra.elements[idx] = last
			ra.length--
		}
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
	return ra.elements[idx].value
}

func (ra *LimitMapArray) Randoms(limit int, maxTry int) []interface{} {
	if ra.length < limit {
		values := make([]interface{}, ra.length)
		for i := 0; i < ra.length; i++ {
			values[i] = ra.elements[i].value
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
			if _, ok = idxs[idx]; !ok && (ra.sRule == nil || ra.sRule.Check(ra.elements[idx].value)) {
				idxs[idx] = '0'
				values[i] = ra.elements[idx].value
				i++
			}
			if try >= maxTry {
				return values[:i]
			}
		}
		return values
	}
}

func (ra *LimitMapArray) Iterate() *Iterate {
	return &Iterate{
		cur:   0,
		array: ra,
	}
}

type SelectRuler interface {
	Check(interface{}) bool
}

type CoverRuler interface {
	ShouldCover(interface{}) bool
}

type Iterate struct {
	cur   int
	array *LimitMapArray
}

func (it *Iterate) Next() interface{} {
	if it.cur <= it.array.length {
		value := it.array.elements[it.cur].value
		it.cur++
		return value
	}
	return nil
}

const (
	CoverMaxTry = 16
)

var (
	SetFullAMapArrayErr = errors.New(`cannot set full map array`)
)
