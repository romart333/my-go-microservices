package interceptor

import (
	"context"
	"log/slog"
	"path"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	method := path.Base(info.FullMethod)
	slog.Info("начало grpc метода", "method", method)

	resp, err := handler(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		slog.Error("error", "error", err, "method", info.FullMethod, "status", st.Code(), "message", st.Message())
	} else {
		slog.Info("успешное завершение grpc метода", "method", method)
	}

	return resp, err
}
