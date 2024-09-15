package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

var MdwChainRequestResponseLogAuthentication = MiddlewareChain(
	MdwRequestResponseLoggerMiddleware,
	MdwRequireAuthenticationMiddleware,
)

// MdwRequireSuperUserOrHigherMiddleware returns a middleware that requires
// the user role to be either "root" or "superuser".
func MdwRequireSuperUserOrHigherMiddleware() Middleware {
	return MdwRequireAuthorizationMiddleware(RootRole, SuperUserRole)
}

// ----------------------------------------------------------------------------
/// Request/Response Logger Middleware

// MdwRequestLoggerMiddleware logs the incoming HTTP request to the standard
// logger.
func MdwRequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// MdwRequestResponseLoggerMiddleware logs the request and response status code.
func MdwRequestResponseLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("Request: %s %s", r.Method, r.RequestURI)

		// Capture the response
		rec := httptest.NewRecorder()

		next.ServeHTTP(rec, r)

		// Copy the recorded response back to the original response writer
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())

		// Log the response status code
		log.Printf("Response Status Code: %d", rec.Code)
	})
}

// ----------------------------------------------------------------------------
/// Middleware Chain

type Middleware func(http.Handler) http.HandlerFunc

// MiddlewareChain creates a single Middleware that chains together multiple
// middlewares. The order of execution is in reverse order of the arguments.
func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}

// ----------------------------------------------------------------------------
/// Authentication and Authorization Middleware

type ContextKey string

const (
	ContextUserIDKey ContextKey = "UserIDKey"
	ContextRoleKey   ContextKey = "RoleKey"
)

type Role string

const (
	RootRole            Role = "root"
	SuperUserRole       Role = "superuser"
	MerchantType001Role Role = "merchant_type_001"
	MerchantType002Role Role = "merchant_type_002"
	BuyerRole           Role = "buyer"
)

// MdwRequireAuthenticationMiddleware returns a middleware that requires a valid
// Authorization header in the format of "Bearer <user_id>". The user_id is
// validated and saved to the request context with key ContextUserIDKey. The
// user role is determined based on the user_id and saved to the request context
// with key ContextRoleKey.
func MdwRequireAuthenticationMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Printf("MdwRequireAuthenticationMiddleware ERROR: no Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// assume token is valid
		userID := strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
		if userID == "" {
			log.Printf("MdwRequireAuthenticationMiddleware ERROR: empty user ID")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userIDInt, err := strconv.Atoi(userID) // Convert userID to int
		if err != nil {
			log.Printf("MdwRequireAuthenticationMiddleware ERROR: invalid user ID: %s: %v", userID, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var role Role
		switch {
		case userIDInt == 1:
			role = RootRole
		case userIDInt >= 2 && userIDInt <= 99:
			role = SuperUserRole
		case userIDInt >= 100 && userIDInt <= 199:
			role = MerchantType001Role
		case userIDInt >= 200 && userIDInt <= 299:
			role = MerchantType002Role
		default:
			role = BuyerRole
		}

		ctx := context.WithValue(r.Context(), ContextUserIDKey, userID)
		ctx = context.WithValue(ctx, ContextRoleKey, role)

		log.Printf("MdwRequireAuthenticationMiddleware DEBUG: %s/%s", userID, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// MdwRequireAuthorizationMiddleware returns a middleware that requires the user role to be one of the roles passed in the allowedRoles argument.
// If the user role is not found in the request context, it returns Unauthorized with status code 401.
// If the user role is not in the list of allowed roles, it returns Unauthorized with status code 401.
func MdwRequireAuthorizationMiddleware(allowedRoles ...Role) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(ContextRoleKey).(Role)
			if !ok {
				log.Printf("MdwRequireAuthorizationMiddleware ERROR: ok is %v", ok)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					log.Printf("MdwRequireAuthorizationMiddleware DEBUG: user role %s matches allowed role %s", userRole, allowedRole)
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Printf("MdwRequireAuthorizationMiddleware ERROR: %s", userRole)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}
}
