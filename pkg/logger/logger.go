package logger

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"os"
)

const (
	Key       = "Logger"
	RequestId = "request_id"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, Key, &Logger{logger})
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	logger, ok := ctx.Value(Key).(*Logger)
	if !ok {
		return nil
	}
	return logger
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestId) != nil {
		fields = append(fields, zap.String("request_id", ctx.Value(RequestId).(string)))
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestId) != nil {
		fields = append(fields, zap.String("request_id", ctx.Value(RequestId).(string)))
	}
	l.l.Fatal(msg, fields...)
	os.Exit(1)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestId) != nil {
		fields = append(fields, zap.String("request_id", ctx.Value(RequestId).(string)))
	}
	l.l.Error(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestId) != nil {
		fields = append(fields, zap.String("request_id", ctx.Value(RequestId).(string)))
	}
	l.l.Warn(msg, fields...)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestId) != nil {
		fields = append(fields, zap.String("request_id", ctx.Value(RequestId).(string)))
	}
	l.l.Debug(msg, fields...)
}

func Interceptor(ctx context.Context) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		guid := uuid.New().String()
		GetLoggerFromCtx(ctx).Info(ctx, "guid", zap.String("guid", guid))
		ctx = context.WithValue(ctx, RequestId, guid)
		return handler(ctx, req)
	}
}
