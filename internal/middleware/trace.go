package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
)

type contextKey string

const traceIDKey contextKey = "trace_id"

func TraceID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b := make([]byte, 8)
			rand.Read(b)
			traceID := fmt.Sprintf("%x", b)

			ctx := context.WithValue(r.Context(), traceIDKey, traceID)

			w.Header().Set("X-Request-ID", traceID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey).(string); ok {
		return id
	}
	return ""
}
