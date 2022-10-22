package gateway

import (
	"context"
	"net/http"

	"github.com/dotdak/gopkg/pgrpc"
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
	return NewGrpcGatewayWithOpts(ctx, httpAddr, grpcAddr, register, []runtime.ServeMuxOption{
		runtime.WithErrorHandler(pgrpc.CustomHTTPErrorHandler),
	})
}

func NewGrpcGatewayWithOpts(
	ctx context.Context,
	httpAddr, grpcAddr string,
	register []RegisterFunc,
	muxOpts []runtime.ServeMuxOption,
) *GrpcGateway {
	mux := runtime.NewServeMux(muxOpts...)
	opts := defaultDialOption()

	for _, f := range register {
		if err := f(ctx, mux, grpcAddr, opts); err != nil {
			panic(err)
		}
	}

	return &GrpcGateway{
		Server: &http.Server{
			Addr:    httpAddr,
			Handler: mux,
		},
		addr:     httpAddr,
		grpcAddr: grpcAddr,
		mux:      mux,
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
	return s.Server.ListenAndServe()
}

func (s *GrpcGateway) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
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

func (s *GrpcGateway) GetMux() *runtime.ServeMux {
	return s.mux
}

func WrapperHandler(h http.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		h(w, r)
	}
}

func WithOption(newOpts []grpc.DialOption, f RegisterFunc) RegisterFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
		return f(ctx, mux, endpoint, newOpts)
	}
}
