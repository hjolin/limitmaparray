// maparray_test
package maparray

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func TestMapArray(t *testing.T) {
	ma := NewLimitMapArray(5, false)
	gorotines := 100000
	over := make(chan byte, 100)
	lock := &sync.Mutex{}
	for i := 0; i < gorotines; i++ {
		go func(i int) {
			//time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			lock.Lock()
			ma.Set(strconv.FormatInt(int64(i), 10), i)
			lock.Unlock()
			over <- '0'
		}(i)
	}
	i := 0
	for {
		select {
		case <-over:
			i++
			if i == gorotines {
				//				it := ma.Iterate()
				//for key := range ma.index {
				//fmt.Print(key, `  `)
				//	ma.Remove(key)
				//}
				//fmt.Println()
				fmt.Println(ma.Length(), ma.length, len(ma.elements), len(ma.index))
				for key := range ma.index {
					fmt.Print(key, `  `)
				}
				return
			}
		}

	}

}
