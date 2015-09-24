// maparray_test
package maparray

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func TestMapArray(t *testing.T) {
	ma := NewLimitMapArray(2048, false)
	gorotines := 100000
	over := make(chan byte, 100)
	lock := &sync.Mutex{}
	for i := 0; i < gorotines; i++ {
		go func(i int) {
			//time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			lock.Lock()
			ma.Set(strconv.FormatInt(int64(i), 10), i)
			lock.Unlock()
			ma.Randoms(8, 32)
			over <- '0'
		}(i)
	}
	i := 0
	var (
		value interface{}
		ok    bool
	)
	for {
		select {
		case <-over:
			i++
			if i == gorotines {
				keys := ma.Keys()
				for _, key := range keys {
					if value, ok = ma.Get(key); ok {
						fmt.Print(value, `   `)
					}
				}
				//fmt.Println()
				fmt.Println(ma.Length())
				for key := range ma.index {
					ma.Remove(key)
				}
				fmt.Println(ma.Length())
				return
			}
		}

	}

}
