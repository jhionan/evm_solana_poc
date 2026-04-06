package audit_test

import (
	"encoding/hex"
	"testing"

	"github.com/jhionan/multichain-staking/internal/audit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeHash_ProducesValidHexSHA256(t *testing.T) {
	h := audit.ComputeHash("stake", "0xWallet", "evm", "amount=100", "")

	// A SHA-256 hex string is always 64 hex characters.
	assert.Len(t, h, 64, "expected 64-char hex SHA-256")

	// Must be valid hex.
	_, err := hex.DecodeString(h)
	require.NoError(t, err, "hash must be valid hexadecimal")
}

func TestComputeHash_ChainIntegrity_Hash2UsesPrevHash(t *testing.T) {
	hash1 := audit.ComputeHash("stake", "0xAlice", "evm", "amount=50", "")
	assert.NotEmpty(t, hash1)

	hash2 := audit.ComputeHash("unstake", "0xAlice", "evm", "amount=50", hash1)
	assert.NotEmpty(t, hash2)

	// hash1 and hash2 must be different because prevHash differs.
	assert.NotEqual(t, hash1, hash2)

	// Re-computing hash2 with the same inputs must be deterministic.
	hash2Again := audit.ComputeHash("unstake", "0xAlice", "evm", "amount=50", hash1)
	assert.Equal(t, hash2, hash2Again)
}

func TestComputeHash_TamperingDetection(t *testing.T) {
	hash1 := audit.ComputeHash("stake", "0xBob", "solana", "amount=200", "")

	// Simulate an honest hash2 that references hash1.
	hash2 := audit.ComputeHash("claim", "0xBob", "solana", "rewards=10", hash1)

	// Now a verifier re-computes hash2 using the same inputs + hash1.
	// If hash1 were tampered (different string), the result would differ.
	tamperedHash1 := hash1[:len(hash1)-1] + "X" // flip last char
	recomputedWithTampered := audit.ComputeHash("claim", "0xBob", "solana", "rewards=10", tamperedHash1)

	assert.NotEqual(t, hash2, recomputedWithTampered, "tampered prevHash must produce a different hash")
}

func TestComputeHash_Deterministic(t *testing.T) {
	args := [5]string{"stake", "0xCarol", "evm", "amount=999", "abc123"}
	h1 := audit.ComputeHash(args[0], args[1], args[2], args[3], args[4])
	h2 := audit.ComputeHash(args[0], args[1], args[2], args[3], args[4])
	assert.Equal(t, h1, h2, "same inputs must always produce the same hash")
}
