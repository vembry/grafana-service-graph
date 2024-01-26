package main

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func setupLogger() *Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.DisableStacktrace = true
	config.DisableCaller = true
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// To log correctly line & function we add caller skip (1)
	// Because we wrap zap-uber by 1 layer of logging
	// https://github.com/uber-go/zap/blob/master/options.go#L108
	logger, _ := config.Build(zap.AddCallerSkip(1), zap.AddCaller())

	return &Logger{
		logger.Sugar(),
	}
}

func (l *Logger) InfofC(ctx context.Context, msg string, fields ...interface{}) {
	l.SugaredLogger.With("trace_id", trace.SpanFromContext(ctx).SpanContext().TraceID().String()).Errorf(msg, fields...)
}
func (l *Logger) ErrorfC(ctx context.Context, msg string, fields ...interface{}) {
	l.SugaredLogger.With("trace_id", trace.SpanFromContext(ctx).SpanContext().TraceID().String()).Infof(msg, fields...)
}
