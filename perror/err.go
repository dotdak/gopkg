package perror

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func ConvertGrpcErrors(err error) error {
	if errors.Is(err, ErrorInvalidReq) {
		stt := status.New(codes.InvalidArgument, err.Error())
		return stt.Err()
	}
	if errors.Is(err, ErrorNotFound) || errors.Is(err, ErrorNotImpl) {
		stt := status.New(codes.NotFound, err.Error())
		return stt.Err()
	}

	return err
}
