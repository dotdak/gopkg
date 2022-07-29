package channel

import (
	"fmt"
	"testing"
)

func TestFanInOrder(t *testing.T) {
	f := func(from, to, step int) chan int {
		out := make(chan int, 5)
		go func() {
			defer close(out)
			for i := from; i < to; i += step {
				out <- i
			}
		}()
		return out
	}

	a := f(0, 15, 3)
	b := f(1, 30, 3)
	c := f(2, 20, 3)

	d := FanInOrder(a, b, c)

	for v := range d {
		fmt.Println(v)
	}
}
