package logger

import (
	"context"
	"fmt"
	"log"
	"os"

	"log/slog"

	"github.com/pkg/errors"
)

type Logger struct {
	logger         *slog.Logger
	onError        func(l *Logger, msg string, err error, arg ...any)
	loggerContexts []LoggerContext
}

type Opts struct {
	Level   slog.Level
	OnError func(l *Logger, msg string, err error, arg ...any)
}

type LoggerContext struct {
	Key   string
	Value string
}

func New() *Logger {
	options := slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, &options)
	logger := slog.New(jsonHandler)

	return &Logger{
		logger: logger,
		onError: func(l *Logger, msg string, err error, arg ...any) {
			traceIDContext, ok := l.LoggerContext("traceID")
			if !ok {
				log.Println(msg)
				return
			}

			log.Printf("%s \n", traceIDContext.Value)
		},
	}
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		logger:  l.logger.With(args...),
		onError: l.onError,
	}
}

func (l *Logger) LoggerContext(k string) (*LoggerContext, bool) {
	for _, v := range l.loggerContexts {
		if v.Key == k {
			return &v, true
		}
	}

	return nil, false
}

func (l *Logger) SetLoggerContexts(args ...LoggerContext) {
	l.loggerContexts = append(l.loggerContexts, args...)
}

func (l *Logger) Debug(msg string, arg ...any) {
	l.logger.Debug(msg, arg...)
}

func (l *Logger) Info(msg string, arg ...any) {
	l.logger.Info(msg, arg...)
}

func (l *Logger) Warn(msg string, arg ...any) {
	l.logger.Warn(msg, arg...)
}

func (l *Logger) Error(msg string, err error, arg ...any) {
	arg = append(arg, slog.String("stack", fmt.Sprintf("%+v", errors.WithStack(err))))
	l.logger.Error(msg, append([]any{err}, arg...)...)

	go func() {
		l.onError(l, msg, err, arg...)
	}()
}

type traceLoggerCtxKey struct{}

// context に logger を詰める
func TraceLoggerWith(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, traceLoggerCtxKey{}, logger)
}

// context から Logger を抜き出す
func TraceLoggerFrom(ctx context.Context) (*Logger, bool) {
	traceLogger, ok := ctx.Value(traceLoggerCtxKey{}).(*Logger)

	return traceLogger, ok
}
