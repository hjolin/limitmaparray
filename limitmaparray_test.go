// maparray_test
package maparray

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestMapArray(t *testing.T) {
	ma := NewLimitMapArray(1001, nil, nil)
	gorotines := 100000
	over := make(chan byte, 100)
	for i := 0; i < gorotines; i++ {
		go func(i int) {
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			ma.Set(strconv.FormatInt(int64(i), 10), i)
			over <- '0'
		}(i)
	}
	i := 0
	for {
		select {
		case <-over:
			i++
			if i == gorotines {
				fmt.Println(ma.length)
				for j := 0; j < ma.length; j++ {
					fmt.Print(ma.elements[ma.busyIndexs[j]].key, `  `)
				}
				ma.Remove(strconv.FormatInt(rand.Int63n(int64(gorotines)), 10))
				fmt.Println(ma.Length(), `:`, len(ma.index))
				return
			}
		}

	}

}
