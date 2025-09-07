package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/wshadm/miniblog/internal/pkg/logger"
	"github.com/wshadm/miniblog/pkg/options"
)

// HTTPServer 代表一个HTTP服务器
type HTTPServer struct {
	srv *http.Server
}

// NewHTTPServer 创建一个新的HTTP服务器实例
func NewHTTPServer(httpOptions *options.HTTPOptions, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		srv: &http.Server{
			Addr:    httpOptions.Addr,
			Handler: handler,
		},
	}
}

// RunOrDie 启动HTTP服务器
func (s *HTTPServer) RunOrDie() {
	logger.L().Info().Msgf("Start to listening the incoming requests", "protocol", protocolName(s.srv), "addr", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.L().Fatal().Msgf("Failed to server HTTP(s) server", "err", err)
	}
}

// HTTP的优雅关闭
func (s *HTTPServer) GracefulStop(ctx context.Context) {
	logger.L().Info().Msg("Gracefully stop HTTP(s) server")
	if err := s.srv.Shutdown(ctx); err != nil {
		logger.L().Err(err).Msg("HTTP(s) server forced to shutdown")
	}
}
