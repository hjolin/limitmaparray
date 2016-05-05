// maparray_test
package maparray

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMapArray(t *testing.T) {
	ma := NewLimitMapArray(2, true)
	gorotines := 17
	over := make(chan byte, 100)
	for i := 0; i < gorotines; i++ {
		go func(i int) {
			//time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			fmt.Println(ma.Set(strconv.FormatInt(int64(i), 10), i))
			ma.Random()
			over <- '0'
		}(i)
	}
	i := 0
	var (
		//	value interface{}
		ok bool
	)
	for {
		select {
		case <-over:
			i++
			if i == gorotines {
				keys := ma.Keys()
				for _, key := range keys {
					if _, ok = ma.Get(key); ok {
						//fmt.Print(value, `   `)
					}
				}
				fmt.Println(ma.Length())
				//				for key := range ma.index {
				//					ma.Remove(key)
				//				}

				fmt.Println(ma.Keys())
				fmt.Println(ma.Length(), ma.Capcity())
				return
			}
		}

	}

}
