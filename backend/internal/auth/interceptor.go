package auth

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
)

// contextKey is an unexported type for context keys in this package,
// preventing collisions with keys defined in other packages.
type contextKey struct{}

// claimsKey is the singleton key used to store Claims in a context.
var claimsKey = contextKey{}

// ClaimsFromContext retrieves Claims stored by AuthInterceptor.
// Returns false if no claims are present (unauthenticated request).
func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(claimsKey).(Claims)
	return c, ok
}

// ContextWithClaims returns a new context with the given Claims attached.
// This is intended for use in tests and middleware that need to inject claims
// without going through the full JWT verification flow.
func ContextWithClaims(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}

// AuthInterceptor returns a connect.UnaryInterceptorFunc that validates the
// Bearer JWT on every request, except GetTiers which is public.
// On success the verified Claims are stored in the request context.
func AuthInterceptor(jwtSvc *JWTService) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// GetTiers is public — no auth required.
			if strings.HasSuffix(req.Spec().Procedure, "/GetTiers") {
				return next(ctx, req)
			}

			authHeader := req.Header().Get("Authorization")
			if authHeader == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("auth: missing Authorization header"))
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("auth: Authorization header must use Bearer scheme"))
			}

			tokenStr := strings.TrimPrefix(authHeader, prefix)
			claims, err := jwtSvc.Verify(tokenStr)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			ctx = context.WithValue(ctx, claimsKey, claims)
			return next(ctx, req)
		}
	}
}

// RequirePermission returns a connect.UnaryInterceptorFunc that verifies the
// caller (whose Claims must already be in ctx) holds the given permission.
// This interceptor must be chained after AuthInterceptor.
func RequirePermission(perm Permission) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			claims, ok := ClaimsFromContext(ctx)
			if !ok {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("auth: no claims in context — AuthInterceptor must run first"))
			}
			if !HasPermission(claims.Role, perm) {
				return nil, connect.NewError(connect.CodePermissionDenied, errors.New("auth: permission denied"))
			}
			return next(ctx, req)
		}
	}
}
