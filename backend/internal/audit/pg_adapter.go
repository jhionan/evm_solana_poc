package audit

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/jhionan/multichain-staking/db/sqlc"
)

// PGAuditDB adapts the sqlc-generated Queries to the audit.DB interface.
type PGAuditDB struct {
	q *db.Queries
}

func NewPGAuditDB(pool *pgxpool.Pool) *PGAuditDB {
	return &PGAuditDB{q: db.New(pool)}
}

func (a *PGAuditDB) GetLatestAuditLog(ctx context.Context) (AuditLogRow, error) {
	row, err := a.q.GetLatestAuditLog(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		return AuditLogRow{HasHash: false}, nil
	}
	if err != nil {
		return AuditLogRow{}, err
	}
	return AuditLogRow{Hash: row.Hash, HasHash: true}, nil
}

func (a *PGAuditDB) InsertAuditLog(ctx context.Context, arg InsertAuditLogParams) error {
	_, err := a.q.InsertAuditLog(ctx, db.InsertAuditLogParams{
		Action:   arg.Action,
		Actor:    arg.Actor,
		ChainID:  pgtype.Text{String: arg.ChainID, Valid: arg.ChainID != ""},
		Details:  arg.Details,
		PrevHash: pgtype.Text{String: arg.PrevHash, Valid: arg.PrevHash != ""},
		Hash:     arg.Hash,
	})
	return err
}
