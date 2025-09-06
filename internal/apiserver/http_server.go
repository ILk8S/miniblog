package apiserver

import (
	"context"

	"github.com/wshadm/miniblog/internal/pkg/server"
)

//ginServer定义一个使用Gin框架开发的HTTP服务器
type ginServer struct {

}

//确保*ginServer实现了 server.Server接口
var _ server.Server = (*ginServer)(nil)

func (c *ServerConfig) NewGinServer() *ginServer {
	return &ginServer{}
}

func (s *ginServer) RunOrDie() {
	select {}
}

func (s *ginServer) GracefulStop(ctx context.Context) {
	
}