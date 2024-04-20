package middlewares

import (
	"log"
	"net/http"
	"time"
)

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[API] %s %s (Latency: %v)", r.Method, r.URL.Path, time.Since(started).Milliseconds())
	})
}
