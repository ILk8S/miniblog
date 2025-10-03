package apiserver

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	handler "github.com/wshadm/miniblog/internal/apiserver/handler/grpc"
	mw "github.com/wshadm/miniblog/internal/pkg/middleware/grpc"
	"github.com/wshadm/miniblog/internal/pkg/server"
	apiv1 "github.com/wshadm/miniblog/pkg/api/apiserver/v1"
	"google.golang.org/grpc"
)

// grpcServer 定义一个 gRPC 服务器.
type grpcServer struct {
	srv server.Server
	// stop 为优雅关停函数.
	stop func(context.Context)
}

// 确保*grpcServer 实现了server.Server接口
var _ server.Server = (*grpcServer)(nil)

// NewGRPCServerOr 创建并初始化 gRPC 或者 gRPC +  gRPC-Gateway 服务器.
// 一般函数命名中有Or，表示“或者”的含义，暗示该函数有多种选择
func (c *ServerConfig) NewGRPCServerOr() (server.Server, error) {
	//配置gRPC服务器选项，包括拦截器链
	serverOptions := []grpc.ServerOption{
		//注意拦截器顺序
		grpc.ChainUnaryInterceptor(
			//请求ID拦截器
			mw.RequestIDInterceptor(),
		),
	}
	//创建grpc服务器
	grpcsrv, err := server.NewGRPCServer(
		c.cfg.GRPCOptions,
		serverOptions,
		func(s grpc.ServiceRegistrar) {
			apiv1.RegisterMiniBlogServer(s, handler.NewHandler())
		})

	if err != nil {
		return nil, err
	}
	if c.cfg.ServerMode == GRPCServerMode {
		return &grpcServer{
			srv: grpcsrv,
			stop: func(ctx context.Context) {
				grpcsrv.GracefulStop(ctx)
			},
		}, nil
	}
	// 先启动 gRPC 服务器，因为 HTTP 服务器依赖 gRPC 服务器.
	go grpcsrv.RunOrDie()
	httpsrv, err := server.NewGRPCGatewayServer(c.cfg.HTTPOptions, c.cfg.GRPCOptions,
		func(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
			return apiv1.RegisterMiniBlogHandler(context.Background(), mux, conn)
		})
	if err != nil {
		return nil, err
	}
	return &grpcServer{
		srv: httpsrv,
		stop: func(ctx context.Context) {
			grpcsrv.GracefulStop(ctx)
			httpsrv.GracefulStop(ctx)
		},
	}, nil
}

func (s *grpcServer) RunOrDie() {
	s.srv.RunOrDie()
}

func (s *grpcServer) GracefulStop(ctx context.Context) {
	s.stop(ctx)
}
