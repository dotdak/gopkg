package promise

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/dotdak/gopkg/perror"
)

type Result struct {
	Data interface{}
	Err  error
}

type Func func(context.Context) *Result

func New(fn func(context.Context) (interface{}, error)) Func {
	return func(c context.Context) *Result {
		data, err := fn(c)
		return &Result{
			Data: data,
			Err:  err,
		}
	}
}

func AllSettled(ctx context.Context, fns ...Func) []*Result {
	var result sync.Map
	var wg sync.WaitGroup
	for i, fn := range fns {
		i := i
		fn := fn
		wg.Add(1)
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
					result.Store(i, errReturn)
				}
				wg.Done()
			}()

			result.Store(i, fn(ctx))
		}()
	}

	wg.Wait()
	output := make([]*Result, len(fns))
	for i := range fns {
		r, _ := result.Load(i)
		output[i] = r.(*Result)
	}
	return output
}
