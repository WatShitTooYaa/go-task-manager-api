package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/auth"
	"github.com/WatShitTooYaa/go-task-manager-api/internal/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, "Missing authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Unauthorized(w, "Invalid authorization format")
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(w, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		reqID := middleware.GetReqID(r.Context())

		// Log request
		fmt.Printf("[%s] %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(ww, r)

		// Log duration
		duration := time.Since(start)

		logEvent := log.Info()
		if ww.Status() >= 400 {
			logEvent = log.Warn()
		}
		if ww.Status() >= 400 {
			logEvent = log.Warn()
		}
		if ww.Status() >= 500 {
			logEvent = log.Error()
		}

		logEvent.
			Str("request_id", reqID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.Status()).
			Int("bytes", ww.BytesWritten()).
			Dur("duration", duration).
			Str("remote_addr", r.RemoteAddr).
			Msg("HTTP request")
		// fmt.Printf("Completed in %v\n", duration)
	})
}
