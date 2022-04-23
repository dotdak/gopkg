package promise

import (
	"context"
	"fmt"
)

func Retry(ctx context.Context, n int, f Func) (interface{}, error) {
	var err error
	for i := 0; i < n; i++ {
		out := f(ctx)
		if out.Err == nil {
			return out.Data, nil
		}

		err = out.Err
	}

	return nil, fmt.Errorf("%w: retried %v times", err, n)
}
