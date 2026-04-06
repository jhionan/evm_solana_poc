// Package security provides security primitives for the multi-chain staking API.
package security

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/redis/go-redis/v9"

	"github.com/jhionan/multichain-staking/internal/auth"
)

const (
	// DefaultLimit is the default maximum number of requests per window.
	DefaultLimit = 60
	// DefaultWindow is the default sliding window duration.
	DefaultWindow = time.Minute
)

// RateLimiter is a ConnectRPC interceptor that rate-limits callers per wallet
// address (when authenticated) or per client IP (for public endpoints).
// It uses a Redis INCR + EXPIRE counter strategy.
//
// Production backend: Valkey (Redis-compatible OSS drop-in).
// Test backend: miniredis (see ratelimit_test.go).
type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

// NewRateLimiter connects to a Redis/Valkey instance and returns a RateLimiter.
// addr is the "host:port" of the server. password may be empty.
// limit is the maximum number of allowed calls per window duration.
func NewRateLimiter(addr, password string, limit int, window time.Duration) (*RateLimiter, error) {
	if addr == "" {
		return nil, errors.New("ratelimit: addr must not be empty")
	}
	if limit <= 0 {
		return nil, errors.New("ratelimit: limit must be > 0")
	}
	if window <= 0 {
		return nil, errors.New("ratelimit: window must be > 0")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("ratelimit: redis ping: %w", err)
	}

	return &RateLimiter{client: rdb, limit: limit, window: window}, nil
}

// NewRateLimiterFromClient creates a RateLimiter from an existing *redis.Client.
// This is useful in tests (e.g. with miniredis) or when the caller manages the
// connection pool externally.
func NewRateLimiterFromClient(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{client: client, limit: limit, window: window}
}

// Close releases the underlying Redis connection pool.
func (rl *RateLimiter) Close() error {
	return rl.client.Close()
}

// Interceptor returns a connect.UnaryInterceptorFunc that enforces the rate limit.
func (rl *RateLimiter) Interceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			key := rl.resolveKey(ctx, req)

			count, err := rl.increment(ctx, key)
			if err != nil {
				// On Redis errors, fail open — do not block legitimate traffic.
				return next(ctx, req)
			}

			if count > int64(rl.limit) {
				return nil, connect.NewError(
					connect.CodeResourceExhausted,
					fmt.Errorf("ratelimit: rate limit exceeded (%d/%d per %s)", count, rl.limit, rl.window),
				)
			}

			return next(ctx, req)
		}
	}
}

// resolveKey returns the Redis key for the current caller.
// Authenticated requests use "rl:{wallet}"; unauthenticated use "rl:{ip}".
func (rl *RateLimiter) resolveKey(ctx context.Context, req connect.AnyRequest) string {
	if claims, ok := auth.ClaimsFromContext(ctx); ok && claims.Wallet != "" {
		return "rl:" + claims.Wallet
	}
	return "rl:" + extractClientIP(req.Header())
}

// increment atomically increments the counter for key and sets the expiry on
// first use. Returns the new counter value.
func (rl *RateLimiter) increment(ctx context.Context, key string) (int64, error) {
	pipe := rl.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, rl.window)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, fmt.Errorf("ratelimit: pipeline exec: %w", err)
	}
	return incr.Val(), nil
}

// extractClientIP returns the best-effort client IP from standard HTTP headers,
// falling back to RemoteAddr when no forwarding headers are present.
func extractClientIP(h http.Header) string {
	// Respect standard forwarding headers in priority order.
	for _, hdr := range []string{"X-Forwarded-For", "X-Real-Ip", "CF-Connecting-IP"} {
		if val := h.Get(hdr); val != "" {
			// X-Forwarded-For can be "client, proxy1, proxy2" — take the first.
			ip := strings.TrimSpace(strings.SplitN(val, ",", 2)[0])
			if ip != "" {
				return ip
			}
		}
	}

	// RemoteAddr from the raw connection — strip the port.
	if addr := h.Get("RemoteAddr"); addr != "" {
		if host, _, err := net.SplitHostPort(addr); err == nil {
			return host
		}
		return addr
	}

	return "unknown"
}
