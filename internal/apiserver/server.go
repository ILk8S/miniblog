package apiserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wshadm/miniblog/internal/pkg/logger"
	"github.com/wshadm/miniblog/internal/pkg/server"
	"github.com/wshadm/miniblog/pkg/options"
)

const (
	// GRPCServerMode 定义 gRPC 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器.
	GRPCServerMode = "grpc"
	// GRPCServerMode 定义 gRPC + HTTP 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器 + HTTP 反向代理服务器.
	GRPCGatewayServerMode = "grpc-gateway"
	// GinServerMode 定义 Gin 服务模式.
	// 使用 Gin Web 框架启动一个 HTTP 服务器.
	GinServerMode = "gin"
)

// Config 配置结构体，用于存储应用相关的配置
// 不用viper.Get 是因为这种方式能更加清晰的知道应用提供了哪些配置项
type Config struct {
	ServerMode  string
	JWTKey      string
	Expiration  time.Duration
	GRPCOptions *options.GRPCOptions
	HTTPOptions *options.HTTPOptions
}

// UnionServer 定义一个联合服务器，根据ServerMode 决定要启动的服务器类型。
type UnionServer struct {
	srv server.Server
}

// ServerConfig
type ServerConfig struct {
	cfg *Config
}

// NewUnionServer 根据配置创建联合服务器
func (c *Config) NewUnionServer() (*UnionServer, error) {
	serverConfig, err := c.NewServerConfig()
	if err != nil {
		return nil, err
	}
	logger.L().Info().Msgf("Initializing federation server", "server-mode", c.ServerMode)
	// 根据服务模式创建对应的服务实例
	// 实际企业开发中，可以根据需要只选择一种服务器模式.
	// 这里为了方便给你展示，通过 cfg.ServerMode 同时支持了 Gin 和 GRPC 2 种服务器模式.
	// 默认为 gRPC 服务器模式.
	var srv server.Server
	switch c.ServerMode {
	case GinServerMode:
		srv, err = serverConfig.NewGinServer(), nil
	default:
		srv, err = serverConfig.NewGRPCServerOr()
	}
	if err != nil {
		return nil, err
	}
	return &UnionServer{srv: srv}, nil
}

// Run运行应用
func (s *UnionServer) Run() error {
	go s.srv.RunOrDie()
	//创建一个os.Signal类型的channel，用于接收系统信号
	quit := make(chan os.Signal, 1)
	// 当执行 kill 命令时（不带参数），默认会发送 syscall.SIGTERM 信号
	// 使用 kill -2 命令会发送 syscall.SIGINT 信号（例如按 CTRL+C 触发）
	// 使用 kill -9 命令会发送 syscall.SIGKILL 信号，但 SIGKILL 信号无法被捕获，因此无需监听和处理
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//阻塞程序， 等待从quit channel中接收信号
	<-quit
	logger.L().Info().Msg("Shutting down server ...")
	//优雅关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//先关闭依赖的服务，再光比被依赖的服务
	s.srv.GracefulStop(ctx)
	return nil
}

func (c *Config) NewServerConfig() (*ServerConfig, error) {
	return &ServerConfig{
		cfg: c,
	}, nil
}
