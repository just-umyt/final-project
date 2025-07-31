package grpc

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	myLog "cart/internal/observability/log"
	"cart/internal/observability/metrics"
)

func LoggingInterceptor(logger myLog.Logger, metric metrics.Metrics, tracer trace.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		latency := time.Now()

		defer func() {
			duration := time.Since(latency).Seconds()
			metric.ObserveLatency(info.FullMethod, duration)
		}()

		ctx, span := tracer.Start(ctx, info.FullMethod)
		defer span.End()

		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error("request failed",
				myLog.String("error", err.Error()),
				myLog.String("method", info.FullMethod),
				myLog.String("trace_id", span.SpanContext().TraceID().String()),
			)

			metric.IncError(info.FullMethod)

			return nil, err
		}

		logger.Info("request success",
			myLog.String("method", info.FullMethod),
			myLog.String("trace_id", span.SpanContext().TraceID().String()),
		)

		metric.IncRequest(info.FullMethod)

		return resp, err
	}
}
