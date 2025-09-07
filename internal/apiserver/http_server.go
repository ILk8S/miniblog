package apiserver

import (
	"context"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	handler "github.com/wshadm/miniblog/internal/apiserver/handler/http"
	"github.com/wshadm/miniblog/internal/pkg/server"
)

// ginServer定义一个使用Gin框架开发的HTTP服务器
type ginServer struct {
	srv server.Server
}

// 确保*ginServer实现了 server.Server接口
var _ server.Server = (*ginServer)(nil)

func (c *ServerConfig) NewGinServer() *ginServer {
	//创建Gin引擎
	engine := gin.New()
	//注册RESTAPI 路由
	c.InstallRESTAPI(engine)
	httpsrv := server.NewHTTPServer(c.cfg.HTTPOptions, engine)

	return &ginServer{
		srv: httpsrv,
	}
}

// 注册API路由、路径和HTTP方法，严格遵循REST规范
func (c *ServerConfig) InstallRESTAPI(engine *gin.Engine) {
	//注册业务无关的API接口
	InstallGenericAPI(engine)
	//创建核心业务处理器
	handler := handler.NewHandler()
	//注册健康检查接口
	engine.GET("/healthz", handler.Healthz)
}

// InstallGenericAPI注册业务无关的路由，例如pprof、404处理等
func InstallGenericAPI(engine *gin.Engine) {
	//注册pprof路由
	pprof.Register(engine)
	//注册404 路由处理,未匹配到的路由将返回404响应
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "Page not found.")
	})
}

func (s *ginServer) RunOrDie() {
	s.srv.RunOrDie()
}

func (s *ginServer) GracefulStop(ctx context.Context) {
	s.srv.GracefulStop(ctx)
}
