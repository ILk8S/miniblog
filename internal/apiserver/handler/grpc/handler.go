package grpc

import (
	apiv1 "github.com/wshadm/miniblog/pkg/api/apiserver/v1"
)

// Handler 负责处理博客模块的请求
type Handler struct {
	apiv1.UnimplementedMiniBlogServer
}

func NewHandler() *Handler {
	return &Handler{}
}
