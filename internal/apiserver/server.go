package apiserver

import (
	"net"
	"time"

	"github.com/wshadm/miniblog/internal/pkg/logger"
	v1 "github.com/wshadm/miniblog/pkg/api/apiserver/v1"
	"github.com/wshadm/miniblog/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	handler "github.com/wshadm/miniblog/internal/apiserver/handler/grpc"
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
}

// UnionServer 定义一个联合服务器，根据ServerMode 决定要启动的服务器类型。
type UnionServer struct {
	cfg *Config
	srv *grpc.Server
	lis net.Listener
}

// NewUnionServer 根据配置创建联合服务器
func (c *Config) NewUnionServer() (*UnionServer, error) {
	lis, err := net.Listen("tcp", c.GRPCOptions.Addr)
	if err != nil {
		return nil, err
	}
	//创建GRPC Server 实例
	grpcsrv := grpc.NewServer()
	v1.RegisterMiniBlogServer(grpcsrv, handler.NewHandler())
	reflection.Register(grpcsrv)
	return &UnionServer{
		cfg: c,
		srv: grpcsrv,
		lis: lis,
	}, nil
}

// Run运行应用
func (s *UnionServer) Run() error {
	// fmt.Printf("ServerMode from ServerOptions: %s\n", s.cfg.JWTKey)
	// fmt.Printf("ServerMode from Viper: %s\n", viper.GetString("jwt-key"))
	// jsonData, _ := json.MarshalIndent(s.cfg, "", " ")
	// fmt.Println(string(jsonData))
	// select {}
	logger.L().Info().Msgf("Start to listening the incoming requests on grpc address", "addr", s.cfg.GRPCOptions.Addr)
	return s.srv.Serve(s.lis)
}
