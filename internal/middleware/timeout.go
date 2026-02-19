package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Cinnamoon-dev/blue-gopher/internal/http/handlers"
)

func Timeout(timeout time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		r = r.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			next.ServeHTTP(w, r)
			close(done)
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				handlers.RespondJSON(w, http.StatusGatewayTimeout, map[string]string{"error": "Gateway Timeout"})
			}
		}
	})
}
