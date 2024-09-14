package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

type Middleware func(http.Handler) http.HandlerFunc

var MdwChainRequestResponseLogAuth = MiddlewareChain(
	MdwRequestResponseLoggerMiddleware,
	MdwRequireAuthMiddleware,
)

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}

func MdwRequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func MdwRequestResponseLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("Request: %s %s", r.Method, r.RequestURI)

		// Create a recorder to capture the response status code
		rec := httptest.NewRecorder()

		// Serve the request to capture the response
		next.ServeHTTP(rec, r)

		// Log the response status code
		log.Printf("Response Status Code: %d", rec.Code)

		// Copy the recorded response back to the original response writer
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())
	}
}

type ContextKey string

const ContextUserIDKey ContextKey = "UserIDKey"

func MdwRequireAuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// assume token is valid
		userID := strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
		ctx := context.WithValue(r.Context(), ContextUserIDKey, userID)

		log.Printf("MdwRequireAuthMiddleware DEBUG: %s", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MdwRequireSuperUserMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(ContextUserIDKey).(string)
		if !ok || userID != "1" {
			log.Printf("RequireSuperUserMiddleware ERROR: %v %s", ok, userID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
