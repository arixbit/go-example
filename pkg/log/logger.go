package log

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config controls logger initialization.
type Config struct {
	Level           string
	Format          string
	StacktraceLevel string
	Service         string
}

type ctxKey struct{}

var defaultLogger atomic.Pointer[zap.Logger]

func init() {
	defaultLogger.Store(zap.NewNop())
}

// Init initializes the global zap logger.
func Init(cfg Config) (*zap.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	stacktraceLevel, err := parseLevel(cfg.StacktraceLevel)
	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	development := cfg.Format == "console"
	if development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	switch cfg.Format {
	case "", "json":
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	default:
		return nil, fmt.Errorf("unsupported log format %q", cfg.Format)
	}

	options := []zap.Option{
		zap.AddStacktrace(stacktraceLevel),
		zap.AddCaller(),
	}
	if development {
		options = append(options, zap.Development())
	}

	logger := zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
		options...,
	)
	if cfg.Service != "" {
		logger = logger.With(zap.String("service", cfg.Service))
	}
	defaultLogger.Store(logger)
	zap.ReplaceGlobals(logger)

	return logger, nil
}

// L returns the initialized global logger.
func L() *zap.Logger {
	return defaultLogger.Load()
}

// SetLogger replaces the global logger and returns a restore function.
func SetLogger(logger *zap.Logger) func() {
	if logger == nil {
		logger = zap.NewNop()
	}

	previous := defaultLogger.Load()
	defaultLogger.Store(logger)
	zap.ReplaceGlobals(logger)

	return func() {
		defaultLogger.Store(previous)
		zap.ReplaceGlobals(previous)
	}
}

// Sync flushes any buffered log entries.
func Sync() error {
	return defaultLogger.Load().Sync()
}

// WithTraceID returns a new context carrying the given trace ID.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ctxKey{}, traceID)
}

// TraceIDFrom extracts the trace ID from context.
func TraceIDFrom(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(ctxKey{}).(string); ok {
		return v
	}
	return ""
}

// EnsureTraceID attaches traceID only when ctx does not already have one.
func EnsureTraceID(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if TraceIDFrom(ctx) != "" {
		return ctx
	}
	traceID = strings.TrimSpace(traceID)
	if traceID == "" {
		return ctx
	}
	return WithTraceID(ctx, traceID)
}

// NewTraceID joins stable trace parts with a colon.
func NewTraceID(parts ...string) string {
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			clean = append(clean, part)
		}
	}
	return strings.Join(clean, ":")
}

// FromContext returns a zap logger pre-filled with trace_id from context.
func FromContext(ctx context.Context) *zap.Logger {
	entry := defaultLogger.Load()
	if traceID := TraceIDFrom(ctx); traceID != "" {
		entry = entry.With(zap.String("trace_id", traceID))
	}
	return entry
}

// Error is a thin re-export for concise call sites.
func Error(err error) zap.Field {
	return zap.Error(err)
}

func parseLevel(level string) (zapcore.Level, error) {
	var parsed zapcore.Level
	if level == "" {
		level = "info"
	}
	if err := parsed.Set(level); err != nil {
		return zapcore.InfoLevel, fmt.Errorf("parse log level %q: %w", level, err)
	}
	return parsed, nil
}
