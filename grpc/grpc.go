package common

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dotdak/gopkg/logger"
	"github.com/dotdak/gopkg/pcontext"
	"github.com/dotdak/gopkg/perror"
	"github.com/golang-jwt/jwt"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	logger.DEBUG().Info("--> unary: ", "interceptor", info.FullMethod)
	return handler(ctx, req)
}
func StreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	logger.DEBUG().Info("--> stream: ", "interceptor", info.FullMethod)
	return handler(srv, ss)
}

func NewGrpcConnInSecure(dialAddr string) (*grpc.ClientConn, error) {
	return grpc.Dial(
		dialAddr,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)), // gzip compression
	)
}

func DefaultJSONMarshalOptions() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
		UseEnumNumbers:  true,
	}
}

func DefaultJSONUnmarshalOptions() protojson.UnmarshalOptions {
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
}

func JSONMarshalOmitEmptyOptions() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: false,
		UseEnumNumbers:  true,
	}
}

func DefaultGateMuxOpts() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions:   DefaultJSONMarshalOptions(),
				UnmarshalOptions: DefaultJSONUnmarshalOptions(),
			},
		}),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	}
}

func GrpcServerOption() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Minute * 5,
		}),
		grpc.ChainUnaryInterceptor(
			UnaryInterceptor,
			grpc_ctxtags.UnaryServerInterceptor(),
			UnaryServerHeaderForward(),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(TraceLog)),
		),
		grpc.ChainStreamInterceptor(
			StreamInterceptor,
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(TraceLog)),
		),
	}
}

func GrpcServerOptionWithAuth() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Minute * 5,
		}),
		grpc.ChainUnaryInterceptor(
			UnaryInterceptor,
			grpc_ctxtags.UnaryServerInterceptor(),
			UnaryServerHeaderForward(),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(TraceLog)),
			grpc_auth.UnaryServerInterceptor(authFunc),
		),
		grpc.ChainStreamInterceptor(
			StreamInterceptor,
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(TraceLog)),
			grpc_auth.StreamServerInterceptor(authFunc),
		),
	}
}

func TraceLog(r interface{}) (err error) {
	stack := debug.Stack()
	if e, ok := r.(error); ok {
		logger.LOG().Error(fmt.Errorf("%w, stack %s", e, stack), "stack error")
		return e
	} else {
		logger.LOG().Error(fmt.Errorf("%w, stack %s", perror.ErrorInternal, stack), "stack error")
		return perror.ErrorInternal
	}
}

func UnaryServerHeaderForward() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}
		v := make(map[string]string)
		for k := range md {
			vs := md.Get(k)
			if len(vs) > 0 {
				v[k] = vs[0]
			}
		}
		if len(v) > 0 {
			ctx = context.WithValue(ctx, pcontext.CtxHeaderForward, v)
		}
		return handler(ctx, req)
	}
}

func authFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	grpc_ctxtags.Extract(ctx).Set("auth.sub", tokenInfo)

	newCtx := context.WithValue(ctx, pcontext.CtxInternalUser, tokenInfo)
	return newCtx, nil
}

func parseToken(token string) (map[string]string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("%w: malformed token", perror.ErrorInvalidReq)
	}

	sDec, err := jwt.DecodeSegment(parts[1])
	if err != nil {
		return nil, err
	}

	var info map[string]string
	if err := jsoniter.Unmarshal(sDec, &info); err != nil {
		return nil, fmt.Errorf("%w: %s", perror.ErrorInvalidReq, err.Error())
	}

	return info, nil
}
