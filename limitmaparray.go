// maparray
package maparray

import (
	"errors"
	`math/rand`
)

func NewLimitMapArray(capacity int, coverMaxTry int, sRule SelectRuler, cRule CoverRuler) *LimitMapArray {
	if capacity > 0 {
		array := &LimitMapArray{
			index:       make(map[string]int),
			capacity:    capacity,
			length:      -1,
			sRule:       sRule,
			cRule:       cRule,
			coverMaxTry: coverMaxTry,
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
	sRule       SelectRuler
	cRule       CoverRuler
	elements    []*element
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
	if _, ok := ra.index[key]; !ok {
		if !ra.IsFull() {
			ra.length++
			ra.index[key] = ra.length
			ra.elements[ra.length] = &element{key, value}
			return nil
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
		delete(ra.index, key)
		if idx == ra.length {
			ra.length--
		} else {
			last := ra.elements[ra.length]
			ra.index[last.key] = idx
			ra.elements[idx] = last
			ra.length--
		}
		return value
	}
	return nil
}

func (ra *LimitMapArray) IsFull() bool {
	return ra.length == ra.capacity-1
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
		values := make([]interface{}, ra.length+1)
		for i := 0; i <= ra.length; i++ {
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

var (
	SetFullAMapArrayErr = errors.New(`cannot set full map array`)
)
