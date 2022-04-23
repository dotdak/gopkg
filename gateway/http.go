package gateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/dotdak/gopkg/env"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

func defaultDialOption() []grpc.DialOption {
	return []grpc.DialOption{}
}

func NewGrpcGateway(
	ctx context.Context,
	register ...RegisterFunc,
) *GrpcGateway {
	s := runtime.NewServeMux()
	httpPort := env.EnvString("HTTP_PORT", "8080")
	opts := defaultDialOption()

	var addrBuilder strings.Builder
	addrBuilder.WriteString(":")
	addrBuilder.WriteString(httpPort)

	addr := addrBuilder.String()
	for _, f := range register {
		if err := f(ctx, s, addr, opts); err != nil {
			panic(err)
		}
	}
	return &GrpcGateway{
		addr: addr,
		mux:  s,
		opts: opts,
	}
}

type GrpcGateway struct {
	addr string
	opts []grpc.DialOption
	mux  *runtime.ServeMux
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
