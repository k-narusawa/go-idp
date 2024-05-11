package main

import (
	"context"
	"log/slog"
)

type ContextHandler struct {
	slog.Handler
}

func (ch ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return ch.Handler.Enabled(ctx, level)
}

// Handle backend for api, this will be used to configure how the logs will be structured
func (ch ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(ch.addRequestId(ctx)...)
	return ch.Handler.Handle(ctx, r)
}

// WithAttrs overriding default implementation otherwise it will call the starting JSON Handler
func (ch ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ContextHandler{ch.Handler.WithAttrs(attrs)}
}

// WithGroup overriding default implementation otherwise it will call the starting JSON Handler
func (ch ContextHandler) WithGroup(name string) slog.Handler {
	return ContextHandler{ch.Handler.WithGroup(name)}
}

func (ch ContextHandler) addRequestId(ctx context.Context) []slog.Attr {
	var as []slog.Attr
	correlation := getDefaultValueFromContext(ctx, "correlation_id")
	method := getDefaultValueFromContext(ctx, "request_method")
	path := getDefaultValueFromContext(ctx, "request_path")
	agent := getDefaultValueFromContext(ctx, "request_user_agent")

	group := slog.Group("meta_information", slog.String("correlation_id", correlation),
		slog.String("request_method", method),
		slog.String("request_path", path),
		slog.String("request_user_agent", agent))
	as = append(as, group)
	return as
}

// getDefaultValueFromContext get default value from context
func getDefaultValueFromContext(ctx context.Context, key string) string {
	value := ""
	ctxValue := ctx.Value(key)
	if ctxValue != nil {
		value = ctxValue.(string)
	}
	return value
}
