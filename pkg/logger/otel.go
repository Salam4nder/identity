package logger

import (
	"context"
	"log/slog"
)

type OtelLogger struct {
	l *slog.Logger
}

func NewOtelLogger(l *slog.Logger) slog.Handler {
	return &OtelLogger{l: l}
}

func (x *OtelLogger) Enabled(ctx context.Context, lvl slog.Level) bool {
	return x.l.Enabled(ctx, lvl)
}

func (x *OtelLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return x.l.Handler().WithAttrs(attrs)
}

func (x *OtelLogger) WithGroup(name string) slog.Handler {
	return x.l.Handler().WithGroup(name)
}

func (x *OtelLogger) Handle(ctx context.Context, r slog.Record) error {
	return nil
}
