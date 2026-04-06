// export_test.go exposes internal constructors needed only by tests.
// This file is compiled only when running `go test`.
package evm

import (
	"github.com/rs/zerolog"
)

// NewEVMStakerForTest creates an EVMStaker with nil client, signer, and
// contract bindings. This is sufficient for tests that exercise in-memory
// logic (GetTiers, ChainID) and for tests that assert ErrContractNotConnected
// is returned when no bindings are wired.
func NewEVMStakerForTest(logger zerolog.Logger) (*EVMStaker, error) {
	return &EVMStaker{
		client:  nil,
		signer:  nil,
		staking: nil,
		token:   nil,
		logger:  logger,
	}, nil
}
