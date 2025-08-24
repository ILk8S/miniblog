package grpc

import (
	"context"
	"time"

	v1 "github.com/wshadm/miniblog/pkg/api/apiserver/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Healthz 服务健康检查.
func (h *Handler) Healthz(ctx context.Context, rq *emptypb.Empty) (*v1.HealthResponse, error) {
	return &v1.HealthResponse{
		Status:    v1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil
}
