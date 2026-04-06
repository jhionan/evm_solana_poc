package security_test

import (
	"math/big"
	"testing"

	"github.com/jhionan/multichain-staking/internal/security"
	"github.com/stretchr/testify/assert"
)

// --- ValidateEVMAddress ---

func TestValidateEVMAddress_Valid(t *testing.T) {
	validAddresses := []string{
		"0xAbCdEf0123456789AbCdEf0123456789AbCdEf01",
		"0x0000000000000000000000000000000000000000",
		"0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"0xabcdef0123456789abcdef0123456789abcdef01",
	}
	for _, addr := range validAddresses {
		t.Run(addr, func(t *testing.T) {
			assert.NoError(t, security.ValidateEVMAddress(addr))
		})
	}
}

func TestValidateEVMAddress_Invalid(t *testing.T) {
	invalidAddresses := []string{
		"",
		"0x",
		"0xGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG",  // invalid hex chars
		"AbCdEf0123456789AbCdEf0123456789AbCdEf01",    // no 0x prefix
		"0xAbCdEf0123456789AbCdEf0123456789AbCdEf0",   // too short (41 chars after 0x)
		"0xAbCdEf0123456789AbCdEf0123456789AbCdEf012",  // too long (43 chars after 0x)
		"0x 123456789abcdef0123456789abcdef01234567",   // space
	}
	for _, addr := range invalidAddresses {
		t.Run(addr, func(t *testing.T) {
			assert.Error(t, security.ValidateEVMAddress(addr))
		})
	}
}

// --- ValidateSolanaAddress ---

func TestValidateSolanaAddress_Valid(t *testing.T) {
	validAddresses := []string{
		"So11111111111111111111111111111111111111112",       // 44 chars
		"9n4nbM75f5Ui33ZbPYXn59EwSgE8CGsHtAeTH5YFeJ9E", // 44 chars mixed
		"11111111111111111111111111111111",                  // 32 chars (min)
	}
	for _, addr := range validAddresses {
		t.Run(addr, func(t *testing.T) {
			assert.NoError(t, security.ValidateSolanaAddress(addr))
		})
	}
}

func TestValidateSolanaAddress_Invalid(t *testing.T) {
	invalidAddresses := []string{
		"",
		"0",                                                   // too short
		"0OIl",                                                // disallowed chars
		"So1111111111111111111111111111111111111111122",        // too long (45 chars)
		"1111111111111111111111111111111",                      // 31 chars (too short)
	}
	for _, addr := range invalidAddresses {
		t.Run(addr, func(t *testing.T) {
			assert.Error(t, security.ValidateSolanaAddress(addr))
		})
	}
}

// --- ValidateStakeAmount ---

func TestValidateStakeAmount_Valid(t *testing.T) {
	cases := []*big.Int{
		big.NewInt(1),
		big.NewInt(1000000),
		new(big.Int).Mul(big.NewInt(1e9), big.NewInt(1e9)), // 1e18
	}
	for _, amt := range cases {
		t.Run(amt.String(), func(t *testing.T) {
			assert.NoError(t, security.ValidateStakeAmount(amt))
		})
	}
}

func TestValidateStakeAmount_Invalid(t *testing.T) {
	cases := []*big.Int{
		nil,
		big.NewInt(0),
		big.NewInt(-1),
		big.NewInt(-1000),
	}
	for _, amt := range cases {
		name := "<nil>"
		if amt != nil {
			name = amt.String()
		}
		t.Run(name, func(t *testing.T) {
			assert.Error(t, security.ValidateStakeAmount(amt))
		})
	}
}
