package server

import (
	"context"
	"net/http"
)

// 抽象出一个Server接口，使用统一的方法来启动服务。也可以实现代码复用性
type Server interface {
	RunOrDie()
	GracefulStop(ctx context.Context)
}

// protocolName 从 http.Server 中获取协议名.
func protocolName(server *http.Server) string {
	if server.TLSConfig != nil {
		return "https"
	}
	return "http"
}
