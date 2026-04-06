package security_test

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/jhionan/multichain-staking/internal/auth"
	"github.com/jhionan/multichain-staking/internal/security"
)

// startMiniRedis starts an in-memory Redis server and returns the client and a
// teardown function.
func startMiniRedis(t *testing.T) (*redis.Client, func()) {
	t.Helper()

	mr, err := miniredis.Run()
	require.NoError(t, err, "start miniredis")

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	return client, func() {
		_ = client.Close()
		mr.Close()
	}
}

// buildRateLimiter creates a RateLimiter backed by miniredis for testing.
func buildRateLimiter(t *testing.T, limit int, window time.Duration) *security.RateLimiter {
	t.Helper()
	client, teardown := startMiniRedis(t)
	t.Cleanup(teardown)
	return security.NewRateLimiterFromClient(client, limit, window)
}

// noopHandler is a minimal connect.UnaryFunc that always succeeds.
func noopHandler(_ context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
	return connect.NewResponse(&emptypb.Empty{}), nil
}

// withWallet returns a context carrying auth.Claims with the given wallet.
func withWallet(wallet string) context.Context {
	return auth.ContextWithClaims(context.Background(), auth.Claims{Wallet: wallet})
}

// TestRateLimiter_WithinLimit_PassesThrough verifies that requests below the
// threshold are forwarded to the next handler without error.
func TestRateLimiter_WithinLimit_PassesThrough(t *testing.T) {
	rl := buildRateLimiter(t, 5, time.Minute)
	handler := rl.Interceptor()(noopHandler)

	req := connect.NewRequest(&emptypb.Empty{})
	ctx := withWallet("0xAlice")

	for i := range 5 {
		_, err := handler(ctx, req)
		assert.NoError(t, err, "request %d should pass", i+1)
	}
}

// TestRateLimiter_OverLimit_ReturnsResourceExhausted verifies that the (limit+1)th
// request receives a CodeResourceExhausted error.
func TestRateLimiter_OverLimit_ReturnsResourceExhausted(t *testing.T) {
	const limit = 3
	rl := buildRateLimiter(t, limit, time.Minute)
	handler := rl.Interceptor()(noopHandler)

	req := connect.NewRequest(&emptypb.Empty{})
	ctx := withWallet("0xBob")

	// Exhaust the limit.
	for i := range limit {
		_, err := handler(ctx, req)
		assert.NoError(t, err, "request %d should pass", i+1)
	}

	// Next request must be rejected.
	_, err := handler(ctx, req)
	require.Error(t, err, "request over limit should fail")

	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, connect.CodeResourceExhausted, connectErr.Code())
}

// TestRateLimiter_DifferentWallets_IndependentCounters verifies that two distinct
// wallets each have their own counter and do not interfere with each other.
func TestRateLimiter_DifferentWallets_IndependentCounters(t *testing.T) {
	const limit = 2
	rl := buildRateLimiter(t, limit, time.Minute)
	handler := rl.Interceptor()(noopHandler)

	req := connect.NewRequest(&emptypb.Empty{})
	ctxAlice := withWallet("0xAlice")
	ctxBob := withWallet("0xBob")

	// Exhaust Bob's limit.
	for range limit {
		_, err := handler(ctxBob, req)
		require.NoError(t, err)
	}

	// Bob is now over limit.
	_, err := handler(ctxBob, req)
	require.Error(t, err)
	var connectErr *connect.Error
	require.ErrorAs(t, err, &connectErr)
	assert.Equal(t, connect.CodeResourceExhausted, connectErr.Code())

	// Alice still has her full quota.
	for range limit {
		_, err := handler(ctxAlice, req)
		assert.NoError(t, err, "Alice's requests should still pass")
	}
}

// TestRateLimiter_NoClaimsUsesIP_PassesThrough verifies that unauthenticated
// requests (no claims in context) are keyed on IP and not rejected when within
// the limit.
func TestRateLimiter_NoClaimsUsesIP_PassesThrough(t *testing.T) {
	const limit = 5
	rl := buildRateLimiter(t, limit, time.Minute)
	handler := rl.Interceptor()(noopHandler)

	req := connect.NewRequest(&emptypb.Empty{})
	req.Header().Set("X-Forwarded-For", "203.0.113.1")

	ctx := context.Background() // no claims

	for i := range limit {
		_, err := handler(ctx, req)
		assert.NoError(t, err, "request %d should pass", i+1)
	}
}
