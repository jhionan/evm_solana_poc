// export_test.go exposes internal constructors needed only by tests.
// This file is compiled only when running `go test`.
package solana

import (
	solanago "github.com/gagliardetto/solana-go"
	"github.com/rs/zerolog"
)

// NewSolanaStakerForTest creates a SolanaStaker with a nil rpc.Client, which
// is sufficient for tests that only exercise in-memory logic (GetTiers,
// ChainID).
func NewSolanaStakerForTest(logger zerolog.Logger) (*SolanaStaker, error) {
	return &SolanaStaker{
		client:    nil,
		programID: solanago.PublicKey{},
		authority: solanago.PrivateKey{},
		logger:    logger,
	}, nil
}
