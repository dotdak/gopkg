package gateway

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

type RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

func defaultDialOption() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}
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
	*http.Server
	addr     string
	grpcAddr string
	opts     []grpc.DialOption
	mux      *runtime.ServeMux
}

func (s *GrpcGateway) Start(ctx context.Context) error {
	s.Server = &http.Server{
		Addr:    s.addr,
		Handler: s.mux,
	}
	return s.ListenAndServe()
}

func (s *GrpcGateway) Stop(ctx context.Context) {
	s.Shutdown(ctx)
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

func WithOption(newOpts []grpc.DialOption, f RegisterFunc) RegisterFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
		return f(ctx, mux, endpoint, newOpts)
	}
}
