package audit_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	stakingv1 "github.com/jhionan/multichain-staking/gen/staking/v1"
	"github.com/jhionan/multichain-staking/internal/audit"
	"github.com/jhionan/multichain-staking/internal/auth"
)

// ---------------------------------------------------------------------------
// Fake DB
// ---------------------------------------------------------------------------

// fakeDB is an in-memory implementation of audit.DB for tests.
type fakeDB struct {
	rows    []audit.InsertAuditLogParams
	noRows  bool  // true → GetLatestAuditLog returns "no rows"
	getErr  error // non-nil → GetLatestAuditLog returns this error
	insErr  error // non-nil → InsertAuditLog returns this error
}

func (f *fakeDB) GetLatestAuditLog(_ context.Context) (audit.AuditLogRow, error) {
	if f.getErr != nil {
		return audit.AuditLogRow{}, f.getErr
	}
	// Rows already inserted take priority over the noRows seed flag so that
	// hash chain tests can read back entries from the same test run.
	if len(f.rows) > 0 {
		last := f.rows[len(f.rows)-1]
		return audit.AuditLogRow{Hash: last.Hash, HasHash: true}, nil
	}
	return audit.AuditLogRow{HasHash: false}, nil
}

func (f *fakeDB) InsertAuditLog(_ context.Context, arg audit.InsertAuditLogParams) error {
	if f.insErr != nil {
		return f.insErr
	}
	f.rows = append(f.rows, arg)
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// noopHandler is a connect.UnaryFunc that always returns success.
func noopAuditHandler(_ context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
	return connect.NewResponse(&stakingv1.StakeResponse{}), nil
}

// errorHandler returns an error instead of a response.
func errorHandler(_ context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
	return nil, connect.NewError(connect.CodeInternal, errors.New("handler error"))
}

// stakeRequest builds a connect.Request for the Stake procedure.
func stakeRequest(wallet string) *connect.Request[stakingv1.StakeRequest] {
	req := connect.NewRequest(&stakingv1.StakeRequest{
		Chain:  stakingv1.Chain_CHAIN_EVM,
		Wallet: wallet,
		Amount: "1000",
		Tier:   stakingv1.Tier_TIER_BRONZE,
	})
	// Manually set the procedure on the spec — connect uses request type to
	// derive Spec in production, but in unit tests we need to craft it.
	return req
}

// withProcedure wraps a proto message in a testAnyRequest with a fixed procedure.
type testAnyRequest struct {
	connect.AnyRequest
	spec connect.Spec
	msg  proto.Message
	hdr  http.Header
}

func (r *testAnyRequest) Spec() connect.Spec  { return r.spec }
func (r *testAnyRequest) Header() http.Header { return r.hdr }
func (r *testAnyRequest) Any() interface{}    { return r.msg }

func makeMutatingRequest(procedure string) connect.AnyRequest {
	return &testAnyRequest{
		spec: connect.Spec{Procedure: procedure},
		msg: &stakingv1.StakeRequest{
			Chain:  stakingv1.Chain_CHAIN_EVM,
			Wallet: "0xDave",
			Amount: "500",
		},
		hdr: http.Header{},
	}
}

func makeReadRequest(procedure string) connect.AnyRequest {
	return &testAnyRequest{
		spec: connect.Spec{Procedure: procedure},
		msg:  &stakingv1.GetTiersRequest{Chain: stakingv1.Chain_CHAIN_EVM},
		hdr:  http.Header{},
	}
}

func withWalletCtx(wallet string) context.Context {
	return auth.ContextWithClaims(context.Background(), auth.Claims{Wallet: wallet})
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// TestAuditInterceptor_MutatingOperation_IsLogged verifies that a Stake call
// produces one audit log row.
func TestAuditInterceptor_MutatingOperation_IsLogged(t *testing.T) {
	db := &fakeDB{noRows: true}
	ai := audit.NewAuditInterceptor(db)
	handler := ai.Interceptor()(noopAuditHandler)

	ctx := withWalletCtx("0xAlice")
	req := makeMutatingRequest("/staking.v1.StakingService/Stake")

	_, err := handler(ctx, req)
	require.NoError(t, err)

	require.Len(t, db.rows, 1, "expected one audit row")
	row := db.rows[0]
	assert.Equal(t, "Stake", row.Action)
	assert.Equal(t, "0xAlice", row.Actor)
	assert.NotEmpty(t, row.Hash, "hash must be populated")
	assert.Empty(t, row.PrevHash, "first entry has no prev hash")
}

// TestAuditInterceptor_ReadOperation_NotLogged verifies that GetTiers does not
// produce an audit log row.
func TestAuditInterceptor_ReadOperation_NotLogged(t *testing.T) {
	db := &fakeDB{}
	ai := audit.NewAuditInterceptor(db)
	handler := ai.Interceptor()(noopAuditHandler)

	ctx := context.Background()
	req := makeReadRequest("/staking.v1.StakingService/GetTiers")

	_, err := handler(ctx, req)
	require.NoError(t, err)

	assert.Empty(t, db.rows, "read operations must not be logged")
}

// TestAuditInterceptor_HandlerError_NotLogged verifies that a failed handler
// does not produce an audit log entry.
func TestAuditInterceptor_HandlerError_NotLogged(t *testing.T) {
	db := &fakeDB{}
	ai := audit.NewAuditInterceptor(db)
	handler := ai.Interceptor()(errorHandler)

	ctx := withWalletCtx("0xEve")
	req := makeMutatingRequest("/staking.v1.StakingService/Stake")

	_, err := handler(ctx, req)
	require.Error(t, err, "handler error must propagate")

	assert.Empty(t, db.rows, "failed operations must not be logged")
}

// TestAuditInterceptor_HashChain_SecondEntryLinksToPrevious verifies that the
// hash of the second log entry uses the hash of the first as prevHash.
func TestAuditInterceptor_HashChain_SecondEntryLinksToPrevious(t *testing.T) {
	db := &fakeDB{noRows: true}
	ai := audit.NewAuditInterceptor(db)
	handler := ai.Interceptor()(noopAuditHandler)

	ctx := withWalletCtx("0xCarol")
	stake := makeMutatingRequest("/staking.v1.StakingService/Stake")
	unstake := makeMutatingRequest("/staking.v1.StakingService/Unstake")

	_, err := handler(ctx, stake)
	require.NoError(t, err)

	_, err = handler(ctx, unstake)
	require.NoError(t, err)

	require.Len(t, db.rows, 2)

	firstHash := db.rows[0].Hash
	secondPrevHash := db.rows[1].PrevHash

	assert.Equal(t, firstHash, secondPrevHash, "second entry's prevHash must equal first entry's hash")
	assert.NotEqual(t, db.rows[0].Hash, db.rows[1].Hash, "each entry must have a unique hash")
}

// TestAuditInterceptor_AllMutatingProcedures_AreLogged verifies all three
// mutating RPCs are recognised.
func TestAuditInterceptor_AllMutatingProcedures_AreLogged(t *testing.T) {
	procedures := []string{
		"/staking.v1.StakingService/Stake",
		"/staking.v1.StakingService/Unstake",
		"/staking.v1.StakingService/ClaimRewards",
	}

	for _, proc := range procedures {
		t.Run(proc, func(t *testing.T) {
			db := &fakeDB{noRows: true}
			ai := audit.NewAuditInterceptor(db)
			handler := ai.Interceptor()(noopAuditHandler)

			ctx := withWalletCtx("0xFrank")
			req := makeMutatingRequest(proc)

			_, err := handler(ctx, req)
			require.NoError(t, err)

			assert.Len(t, db.rows, 1, "procedure %s should produce an audit row", proc)
		})
	}
}

// TestAuditInterceptor_HashChain_IsCorrect verifies that ComputeHash is called
// with the right inputs and produces the expected deterministic hash.
func TestAuditInterceptor_HashChain_IsCorrect(t *testing.T) {
	db := &fakeDB{noRows: true}
	ai := audit.NewAuditInterceptor(db)
	handler := ai.Interceptor()(noopAuditHandler)

	ctx := withWalletCtx("0xGrace")
	req := makeMutatingRequest("/staking.v1.StakingService/Stake")

	_, err := handler(ctx, req)
	require.NoError(t, err)

	require.Len(t, db.rows, 1)
	row := db.rows[0]

	// Recompute the expected hash and verify it matches.
	expected := audit.ComputeHash(row.Action, row.Actor, row.ChainID, string(row.Details), row.PrevHash)
	assert.Equal(t, expected, row.Hash, "stored hash must equal recomputed hash")
}
