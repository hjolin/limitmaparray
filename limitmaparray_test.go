package maparray

import (
	//	"fmt"
	//	"math/rand"
	//	"strconv"
	"testing"
)

func TestMapArray(t *testing.T) {
	//	ma := NewLimitMapArray(100)
	gorotines := 100000
	over := make(chan byte, 2048)
	for i := 0; i < gorotines; i++ {
		go func(i int) {
			//			//time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			//			ma.Set(strconv.FormatInt(int64(i), 10), i)
			//			//			for ma.Set(strconv.FormatInt(int64(i), 10), i) != nil {
			//			//				ma.RemoveByIndex(rand.Intn(ma.Length()))
			//			//			}
			over <- '0'
		}(i)
	}
	var j int
	for {
		select {
		case <-over:
			j++
			if j == gorotines {
				//fmt.Println(ma.Keys())
				//				fmt.Println(ma.Length(), ma.realCapacity())
				//				for j, count := 0, ma.length>>1+1; j < count; j++ {
				//					ma.RemoveByIndex(0)
				//				fmt.Println(ma.Length())
				//				for key := range ma.index {
				//					ma.RemoveByKey(key)
				//				}
				//				fmt.Println(ma.Keys())
				//				fmt.Println(ma.Length(), ma.realCapacity())
				return
			}
		}
	}
}
