package abczap

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	abclog "go.smartfba.io/abcs/pkg/log"
)

var Module = fx.Module("abczap", fx.Provide(New))

type Params struct {
	fx.In
}

type Result struct {
	fx.Out

	Z *zap.Logger
	L Logger
}

func New(p Params) (Result, error) {
	z, err := zap.NewDevelopment()
	if err != nil {
		return Result{}, err
	}

	return Result{
		Z: z,
		L: &zapLogger{
			Logger: z,
		},
	}, nil
}

type Logger interface {
	With(fields ...zapcore.Field) Logger

	Sugar() abclog.Logger

	Log(lvl zapcore.Level, msg string, fields ...zapcore.Field)
}

type contextKey int

const (
	contextKeyValue contextKey = iota
)

func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(contextKeyValue).(Logger); ok {
		return l
	}

	return nil
}

func WithLogger(parent context.Context, l Logger) context.Context {
	return context.WithValue(parent, contextKeyValue, l)
}

func Log(ctx context.Context, lvl zapcore.Level, msg string, fields ...zapcore.Field) {
	l := FromContext(ctx)
	if l == nil {
		return
	}

	l.Log(lvl, msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	Log(ctx, zapcore.DebugLevel, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	Log(ctx, zapcore.InfoLevel, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	Log(ctx, zapcore.WarnLevel, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	Log(ctx, zapcore.ErrorLevel, msg, fields...)
}

func DPanic(ctx context.Context, msg string, fields ...zapcore.Field) {
	Log(ctx, zapcore.DPanicLevel, msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...zapcore.Field) {
	Log(ctx, zapcore.DPanicLevel, msg, fields...)
}

type zapLogger struct {
	*zap.Logger
}

func (l *zapLogger) With(fields ...zapcore.Field) Logger {
	return &zapLogger{
		Logger: l.Logger.With(fields...),
	}
}

func (l *zapLogger) Sugar() abclog.Logger {
	return &zapSugaredLogger{
		SugaredLogger: l.Logger.Sugar(),
	}
}

func (l *zapLogger) Log(lvl zapcore.Level, msg string, fields ...zapcore.Field) {
	l.Logger.Log(lvl, msg, fields...)
}

type zapSugaredLogger struct {
	*zap.SugaredLogger
}

func (l *zapSugaredLogger) With(keysAndValues ...any) abclog.Logger {
	return &zapSugaredLogger{
		SugaredLogger: l.SugaredLogger.With(keysAndValues...),
	}
}

func (l *zapSugaredLogger) Debug(msg string, keysAndValues ...any) {
	l.SugaredLogger.Debugw(msg, keysAndValues)
}

func (l *zapSugaredLogger) Info(msg string, keysAndValues ...any) {
	l.SugaredLogger.Infow(msg, keysAndValues)
}

func (l *zapSugaredLogger) Warn(msg string, keysAndValues ...any) {
	l.SugaredLogger.Warnw(msg, keysAndValues)
}

func (l *zapSugaredLogger) Error(msg string, keysAndValues ...any) {
	l.SugaredLogger.Errorw(msg, keysAndValues)
}

func (l *zapSugaredLogger) DPanic(msg string, keysAndValues ...any) {
	l.SugaredLogger.DPanicw(msg, keysAndValues)
}

func (l *zapSugaredLogger) Panic(msg string, keysAndValues ...any) {
	l.SugaredLogger.Panicw(msg, keysAndValues)
}

func (l *zapSugaredLogger) Fatal(msg string, keysAndValues ...any) {
	l.SugaredLogger.Fatalw(msg, keysAndValues)
}
