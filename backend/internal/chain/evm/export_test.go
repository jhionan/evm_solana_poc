// export_test.go exposes internal constructors needed only by tests.
// This file is compiled only when running `go test`.
package evm

import (
	"github.com/rs/zerolog"
)

// NewEVMStakerForTest creates an EVMStaker with a nil ethclient, which is
// sufficient for any test that does not exercise RPC-dependent methods
// (HealthCheck, Stake, etc.).  This avoids spinning up a real or mock node
// just to test in-memory logic like GetTiers and ChainID.
func NewEVMStakerForTest(logger zerolog.Logger) (*EVMStaker, error) {
	return &EVMStaker{
		client: nil,
		signer: nil,
		logger: logger,
	}, nil
}
