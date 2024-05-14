package logger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type OtelHandler struct {
	next slog.Handler
}

func NewOtelHandler(n slog.Handler) slog.Handler {
	return &OtelHandler{next: n}
}

func (x *OtelHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return x.next.Enabled(ctx, lvl)
}

func (x *OtelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return x.next.WithAttrs(attrs)
}

func (x *OtelHandler) WithGroup(name string) slog.Handler {
	return x.next.WithGroup(name)
}

func (x *OtelHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx == nil {
		return x.next.Handle(ctx, r)
	}

	span := trace.SpanFromContext(ctx)
	if span == nil || !span.IsRecording() {
		return x.next.Handle(ctx, r)
	}

	spanContext := span.SpanContext()
	if spanContext.HasTraceID() {
		traceID := spanContext.TraceID().String()
		r.AddAttrs(slog.String("trace_id", traceID))
	}

	if spanContext.HasSpanID() {
		spanID := spanContext.SpanID().String()
		r.AddAttrs(slog.String("span_id", spanID))
	}

	if r.Level >= slog.LevelError {
		span.SetStatus(codes.Error, r.Message)
	}

	return x.next.Handle(ctx, r)
}
