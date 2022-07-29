package channel

// Maxsize of in is 63
func FanInOrder[V any](in ...chan V) chan V {
	out := make(chan V, len(in))
	go func() {
		defer close(out)
		var complete int64 = 1<<len(in) - 1
		var current int64 = 0
		inx := 0
		for current < complete {
			for k, ch := range in {
				v, ok := <-ch
				if !ok {
					current |= 1 << k
					continue
				}
				out <- v
			}
			inx++
			inx %= len(in)
		}
	}()

	return out
}
