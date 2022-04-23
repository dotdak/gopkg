package perror

import "errors"

var (
	ErrorNotFound   = errors.New("not found")
	ErrorInvalidReq = errors.New("invalid request")
	ErrorInternal   = errors.New("internal error")
	ErrorTimeout    = errors.New("timeout")
	ErrorNotImpl    = errors.New("not implemented")
)

func IfErr(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}

	return nil
}
