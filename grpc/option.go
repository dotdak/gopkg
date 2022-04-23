package pgrpc

import (
	"context"

	"github.com/dotdak/gopkg/pcontext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CallWithAuth(ctx context.Context) grpc.CallOption {
	userMeta := metadata.New(pcontext.GetInternalUser(ctx))
	return grpc.Header(&userMeta)
}
