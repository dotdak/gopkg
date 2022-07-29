package channel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	arr1 = []int{0, 10, 20, 30}
	arr2 = []int{1, 11, 21, 31, 41, 51}
	arr3 = []int{2, 12, 22}
)

var f = func(in []int) chan int {
	out := make(chan int, 5)
	go func() {
		defer close(out)
		for _, i := range in {
			out <- i
		}
	}()
	return out
}

func TestFanInOrder(t *testing.T) {
	art := assert.New(t)
	aCh := f(arr1)
	bCh := f(arr2)
	cCh := f(arr3)

	d := FanInOrder(aCh, bCh, cCh)
	e := FanInOrderClassic(arr1, arr2, arr3)

	f := make([]int, 0, len(d))
	for v := range d {
		f = append(f, v)
	}
	fmt.Println(e)
	fmt.Println(f)
	art.Equal(e, f)
}

func BenchmarkWithoutChannel(b *testing.B) {
	for n := 0; n < b.N; n++ {
		FanInOrderClassic(arr1, arr2, arr3)
	}
}

func BenchmarkWithChannel(b *testing.B) {
	for n := 0; n < b.N; n++ {
		aCh := f(arr1)
		bCh := f(arr2)
		cCh := f(arr3)
		FanInOrder(aCh, bCh, cCh)
	}
}
