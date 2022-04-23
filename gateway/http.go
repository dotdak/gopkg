package gateway

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

func defaultDialOption() []grpc.DialOption {
	return []grpc.DialOption{}
}

func NewGrpcGateway(
	ctx context.Context,
	httpAddr, grpcAddr string,
	register ...RegisterFunc,
) *GrpcGateway {
	s := runtime.NewServeMux()
	opts := defaultDialOption()

	for _, f := range register {
		if err := f(ctx, s, grpcAddr, opts); err != nil {
			panic(err)
		}
	}
	return &GrpcGateway{
		addr:     httpAddr,
		grpcAddr: grpcAddr,
		mux:      s,
		opts:     opts,
	}
}

type GrpcGateway struct {
	addr     string
	grpcAddr string
	opts     []grpc.DialOption
	mux      *runtime.ServeMux
}

func (s *GrpcGateway) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.addr,
		Handler: s.mux,
	}
	return server.ListenAndServe()
}

func (s *GrpcGateway) Register(
	ctx context.Context,
	register ...RegisterFunc,
) error {
	for _, f := range register {
		if err := f(ctx, s.mux, s.addr, s.opts); err != nil {
			return err
		}
	}
	return nil
}
