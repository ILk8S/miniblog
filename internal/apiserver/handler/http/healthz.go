package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wshadm/miniblog/internal/pkg/log"
	apiv1 "github.com/wshadm/miniblog/pkg/api/apiserver/v1"
)

// Healthz服务健康检查
func (h *Handler) Healthz(c *gin.Context) {
	log.W(c.Request.Context()).Infow("Healthz handler is called", "method", "Healthz", "status", "healthy")
	//返回Json响应
	c.JSON(200, &apiv1.HealthResponse{
		Status:    apiv1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	})
}
