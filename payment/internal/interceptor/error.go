package interceptor

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/romart333/my-go-microservices/payment/internal/errors"
)

func ErrorInterceptor(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	slog.ErrorContext(ctx, "обработка gRPC-запроса", "method", info.FullMethod, "error", err)

	switch {
	case errors.Is(err, errs.ErrInvalidPaymentMethod), errors.Is(err, errs.ErrInvalidOrderUUID):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
}
