package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(TraceIDHeaderKey.String())
		if traceID == "" {
			traceID = r.Header.Get("Cdn-Requestid") //sent by Bunny CDN
		}
		if traceID == "" {
			traceID = uuid.New().String()
		}

		r.Header.Set(TraceIDHeaderKey.String(), traceID)
		ctx := r.Context()
		ctx = context.WithValue(ctx, TraceIDContextKey, traceID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
