package scope

import (
	"context"
	"log/slog"
	"net/http"
)

type ctxKey int

const (
	LogCtxKey ctxKey = iota
)

func CtxLog(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, LogCtxKey, log.With("api", "rest"))
}

func Log(r *http.Request) *slog.Logger {
	log := r.Context().Value(LogCtxKey).(*slog.Logger)

	return log
}
