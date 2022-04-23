package gateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/dotdak/gopkg/env"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type RegisterFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption)

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
	for _, f := range register {
		f(ctx, s, httpPort, opts)
	}

	var addr strings.Builder
	addr.WriteString(":")
	addr.WriteString(httpPort)
	return &GrpcGateway{
		server: &http.Server{
			Addr:    addr.String(),
			Handler: s,
		},
	}
}

type GrpcGateway struct {
	server *http.Server
}

func (s *GrpcGateway) Start(ctx context.Context) error {
	return s.server.ListenAndServe()
}
