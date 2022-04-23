package promise

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/dotdak/gopkg/perror"
)

func Future(ctx context.Context, f Func) <-chan *Result {
	c := make(chan *Result, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				errReturn := &Result{}
				if e, ok := r.(error); ok {
					errReturn.Err = fmt.Errorf("%w, stack %s", e, stack)
				} else {
					errReturn.Err = fmt.Errorf("%v: %w, stack: %s", r, perror.ErrorInternal, stack)
				}

				c <- errReturn
			}
		}()
		c <- f(ctx)
	}()

	return c
}
