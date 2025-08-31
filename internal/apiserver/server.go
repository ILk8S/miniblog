package apiserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	handler "github.com/wshadm/miniblog/internal/apiserver/handler/grpc"
	"github.com/wshadm/miniblog/internal/pkg/logger"
	apiv1 "github.com/wshadm/miniblog/pkg/api/apiserver/v1"
	"github.com/wshadm/miniblog/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	apiv1.RegisterMiniBlogServer(grpcsrv, handler.NewHandler())
	reflection.Register(grpcsrv)
	return &UnionServer{
		cfg: c,
		srv: grpcsrv,
		lis: lis,
	}, nil
}

// Run运行应用
func (s *UnionServer) Run() error {
	logger.L().Info().Msgf("Start to listening the incoming requests on grpc address: %s", s.cfg.GRPCOptions.Addr)
	go func() {
		if err := s.srv.Serve(s.lis); err != nil {
			logger.L().Err(err).Msg("gRPC启动失败，将退出程序")
			os.Exit(1) //退出程序
		}
	}()
	conn, err := grpc.NewClient(s.cfg.GRPCOptions.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			// 设置序列化 protobuf 数据时，枚举类型的字段以数字格式输出.
			// 否则，默认会以字符串格式输出，跟枚举类型定义不一致，带来理解成本.
			UseEnumNumbers: true,
		},
	}))
	if err := apiv1.RegisterMiniBlogHandler(context.Background(), gwmux, conn); err != nil {
		return err
	}
	logger.L().Info().Msgf("Start to listening the incoming requests", "protocol", "http", "addr", s.cfg.HTTPOptions.Addr)
	httpsrv := &http.Server{
		Addr:    s.cfg.HTTPOptions.Addr,
		Handler: gwmux,
	}
	if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
