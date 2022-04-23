package pgrpc

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewGrpcHealthCheckServer() grpc_health_v1.HealthServer {
	return &GrpcHealthCheckServerImpl{}
}

type GrpcHealthCheckServerImpl struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *GrpcHealthCheckServerImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}
