package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/jhionan/multichain-staking/internal/auth"
)

// DB is the minimal interface over db.Queries that AuditInterceptor needs.
// Using an interface keeps the interceptor decoupled from pgxpool and makes
// it straightforward to substitute a fake in tests.
type DB interface {
	GetLatestAuditLog(ctx context.Context) (AuditLogRow, error)
	InsertAuditLog(ctx context.Context, arg InsertAuditLogParams) error
}

// AuditLogRow is the data returned by GetLatestAuditLog.
// It mirrors db.AuditLog but uses plain Go types so the audit package has no
// dependency on the generated db package.
type AuditLogRow struct {
	Hash    string
	HasHash bool // false when the table is empty (no previous entry)
}

// InsertAuditLogParams carries the columns for a new audit_log row.
type InsertAuditLogParams struct {
	Action   string
	Actor    string
	ChainID  string // empty string means NULL
	Details  []byte
	PrevHash string // empty string means NULL
	Hash     string
}

// mutatingProcedures is the set of procedure suffixes that must be audit-logged.
// Read-only procedures (GetTiers, GetPosition, ListPositions) are excluded.
var mutatingProcedures = map[string]bool{
	"/Stake":        true,
	"/Unstake":      true,
	"/ClaimRewards": true,
}

// AuditInterceptor logs every successful mutating RPC call to the audit_log
// table, chaining each entry's hash to the previous one.
type AuditInterceptor struct {
	db DB
}

// NewAuditInterceptor constructs an AuditInterceptor backed by the given DB.
func NewAuditInterceptor(db DB) *AuditInterceptor {
	return &AuditInterceptor{db: db}
}

// Interceptor returns a connect.UnaryInterceptorFunc that appends to the audit
// log after every successful mutating operation.
func (a *AuditInterceptor) Interceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Execute the handler first so we only log successful operations.
			resp, err := next(ctx, req)
			if err != nil {
				return resp, err
			}

			procedure := req.Spec().Procedure
			if !isMutating(procedure) {
				return resp, nil
			}

			// Best-effort audit log; do not fail the RPC on logging errors.
			if logErr := a.record(ctx, req, procedure); logErr != nil {
				// In production you would emit a metric / alert here.
				// We intentionally do not return this error to the caller.
				_ = logErr
			}

			return resp, nil
		}
	}
}

// record writes one entry to the audit log.
func (a *AuditInterceptor) record(ctx context.Context, req connect.AnyRequest, procedure string) error {
	actor := actorFromContext(ctx)
	action := actionFromProcedure(procedure)
	chainID := chainIDFromRequest(req)

	details, _ := json.Marshal(map[string]string{
		"procedure": procedure,
		"actor":     actor,
		"chain_id":  chainID,
	})

	// Retrieve the previous hash to maintain the chain.
	prevHash := ""
	latest, err := a.db.GetLatestAuditLog(ctx)
	if err == nil && latest.HasHash {
		prevHash = latest.Hash
	}
	// pgx.ErrNoRows means the table is empty — that's fine, prevHash stays "".
	// Any other DB error is swallowed here: we prefer a broken hash chain to
	// failing the user's request.

	newHash := ComputeHash(action, actor, chainID, string(details), prevHash)

	return a.db.InsertAuditLog(ctx, InsertAuditLogParams{
		Action:   action,
		Actor:    actor,
		ChainID:  chainID,
		Details:  details,
		PrevHash: prevHash,
		Hash:     newHash,
	})
}

// isMutating reports whether the given fully-qualified procedure name is one
// of the operations that must be audit-logged.
func isMutating(procedure string) bool {
	for suffix := range mutatingProcedures {
		if strings.HasSuffix(procedure, suffix) {
			return true
		}
	}
	return false
}

// actorFromContext extracts the wallet address from JWT claims, or returns
// "anonymous" when the request is unauthenticated.
func actorFromContext(ctx context.Context) string {
	if claims, ok := auth.ClaimsFromContext(ctx); ok && claims.Wallet != "" {
		return claims.Wallet
	}
	return "anonymous"
}

// actionFromProcedure derives a short action name from the RPC procedure path.
// "/staking.v1.StakingService/Stake" → "Stake"
func actionFromProcedure(procedure string) string {
	parts := strings.Split(procedure, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return procedure
}

// chainIDFromRequest tries to read the "chain" field from the proto message
// using reflection so we don't need to import the generated types here.
// Returns an empty string when the field is absent or the message is not proto.
func chainIDFromRequest(req connect.AnyRequest) string {
	msg, ok := req.Any().(proto.Message)
	if !ok {
		return ""
	}

	fd := msg.ProtoReflect().Descriptor().Fields().ByName("chain")
	if fd == nil {
		return ""
	}

	val := msg.ProtoReflect().Get(fd)
	switch fd.Kind() {
	case protoreflect.EnumKind:
		enumVal := val.Enum()
		enumDesc := fd.Enum().Values().ByNumber(enumVal)
		if enumDesc != nil {
			name := string(enumDesc.Name())
			// Strip the "CHAIN_" prefix for readability: "CHAIN_EVM" → "evm".
			name = strings.TrimPrefix(name, "CHAIN_")
			return strings.ToLower(name)
		}
	case protoreflect.StringKind:
		return val.String()
	}

	return fmt.Sprintf("%v", val.Interface())
}

