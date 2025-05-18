package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"vatsim-auth-service/internal/monitoring"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Враппер для подсчёта статуса
		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", rec.status).
			Dur("duration", duration).
			Msg("Request handled")

		monitoring.HTTPRequests.WithLabelValues(r.URL.Path, r.Method, http.StatusText(rec.status)).Inc()
		monitoring.HTTPRequestDuration.WithLabelValues(r.URL.Path).Observe(duration.Seconds())
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
