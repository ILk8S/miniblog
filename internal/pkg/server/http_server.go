package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/wshadm/miniblog/internal/pkg/logger"
)

// HTTPServer 代表一个HTTP服务器
type HTTPServer struct {
	srv *http.Server
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
