package middleware

import (
	"context"
	"github.com/make-money-fast/plugins/logger"
	"google.golang.org/grpc"
)

type GRpcLogConfigure struct {
	EnableRequestBody  bool
	EnableResponseBody bool
	SkipLoggerFunc     func(ctx context.Context, method string) bool
}

// UnaryServerInterceptor 日志注入器.
func UnaryServerInterceptor(
	log *logger.Entry,
	config GRpcLogConfigure,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !config.SkipLoggerFunc(ctx, info.FullMethod) {
			resp, err := handler(ctx, req)
			return resp, err
		}
		//TODO:: logger me.
		resp, err := handler(ctx, req)
		return resp, err
	}
}
