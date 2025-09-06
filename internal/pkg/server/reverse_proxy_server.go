package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/wshadm/miniblog/internal/pkg/logger"
	"github.com/wshadm/miniblog/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

// GRPCGatewayServer 代表一个GRPC网关服务器
type GRPCGatewayServer struct {
	srv *http.Server
}

//NewGRPCGatewayServer 创建示例
func NewGRPCGatewayServer(httpOptions *options.HTTPOptions, grpcOptions *options.GRPCOptions,
	registerHandler func(mux *runtime.ServeMux, conn *grpc.ClientConn) error) (*GRPCGatewayServer, error) {
		dialOptions := []grpc.DialOption{grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.DefaultConfig,
			MinConnectTimeout: 10 * time.Second,
		})}
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn, err := grpc.NewClient(grpcOptions.Addr, dialOptions...)
		if err != nil {
			logger.L().Error().Err(err).Msgf("Failed to dial context: %s", err)
			return nil, err
		}
		gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseEnumNumbers: true,
			},
		}))
		if err  := registerHandler(gwmux, conn); err != nil {
			logger.L().Error().Err(err).Msgf("Failed to register handler", "err", err)
			return nil, err
		}
		return &GRPCGatewayServer{
			srv: &http.Server{
				Addr: httpOptions.Addr,
				Handler: gwmux,
			},
		}, nil
}

// RunOrDie 启动 GRPC 网关服务器并在出错时记录致命错误.
func (s *GRPCGatewayServer) RunOrDie() {
	logger.L().Info().Msgf("Start to listening the incoming requests", "protocol", protocolName(s.srv), "addr", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		logger.L().Fatal().Msgf("Failed to server HTTP(s) server", "err", err)
	}

}

// GracefulStop 优雅地关闭 GRPC 网关服务器.
func (s *GRPCGatewayServer) GracefulStop(ctx context.Context) {
	logger.L().Info().Msg("Gracefully stop HTTP(s) server")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		logger.L().Error().Err(err).Msgf("HTTP(s) server forced to shutdown", "err", err)
	}
}
