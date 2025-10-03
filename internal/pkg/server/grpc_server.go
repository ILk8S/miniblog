package server

import (
	"context"
	"net"

	"github.com/wshadm/miniblog/internal/pkg/logger"
	"github.com/wshadm/miniblog/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// GRPCServer代表一个GRPC服务器
type GRPCServer struct {
	srv *grpc.Server
	lis net.Listener
}

func NewGRPCServer(
	grpcOptions *options.GRPCOptions,
	serverOptions []grpc.ServerOption,
	registerServer func(grpc.ServiceRegistrar)) (*GRPCServer, error) {
	lis, err := net.Listen("tcp", grpcOptions.Addr)
	if err != nil {
		logger.L().Error().Msgf("Failed to listen", "err", err)
		return nil, err
	}
	grpcsrv := grpc.NewServer(serverOptions...)
	registerServer(grpcsrv)
	registerHealthServer(grpcsrv)
	reflection.Register(grpcsrv)
	return &GRPCServer{
		srv: grpcsrv,
		lis: lis,
	}, nil
}
func (s *GRPCServer) RunOrDie() {
	logger.L().Info().Msgf("Start to listening the incoming requests", "protocol", "grpc", "addr", s.lis.Addr().String())
	if err := s.srv.Serve(s.lis); err != nil {
		logger.L().Fatal().Msgf("Failed to serve grpc server", "err", err)
	}

}

func (s *GRPCServer) GracefulStop(ctx context.Context) {
	logger.L().Info().Msg("Gracefully stop grpc server")
	s.srv.GracefulStop()
}

// registerHealthServer 注册健康检查服务.
func registerHealthServer(grpcsrv *grpc.Server) {
	healthServer := health.NewServer()
	// 设定服务的健康状态
	healthServer.SetServingStatus("MiniBlog", grpc_health_v1.HealthCheckResponse_SERVING)

	//注册健康检查服务,将health注册给grpcsrv
	grpc_health_v1.RegisterHealthServer(grpcsrv, healthServer)
}
