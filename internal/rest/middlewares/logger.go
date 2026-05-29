package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/paranoideed/uni-products-svc/internal/rest/scope"
)

func (p *Provider) Logger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(scope.CtxLog(r.Context(), log)))
		})
	}
}
