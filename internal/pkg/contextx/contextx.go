package contextx

import "context"

// 定义用于上下文的键
type (
	//userIDKey 定义用户ID的上下文键
	userIDKey struct{}
	//requestIDKey 定义请求ID的上下文键
	requestIDKey struct{}
)

// WithUserID 将用户ID 存到上下文中
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// UserID从上下文中提取用户ID
func UserID(ctx context.Context) string {
	userID := ctx.Value(userIDKey{}).(string)
	return userID
}

//WithRquestID将请求ID存放到上下文中
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

//RequestID从上下文中提取请求
func RequestID(ctx context.Context) string {
	requestID := ctx.Value(requestIDKey{}).(string)
	return requestID
}
