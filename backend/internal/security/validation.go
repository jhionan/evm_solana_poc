// Package security provides input validation helpers for blockchain addresses
// and financial amounts.
package security

import (
	"fmt"
	"math/big"
	"regexp"
)

var (
	// evmAddressRe matches a 0x-prefixed 40-hex-character EVM address.
	evmAddressRe = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

	// solanaAddressRe matches a base58-encoded Solana address (32–44 chars,
	// excluding the ambiguous characters 0, O, I, l).
	solanaAddressRe = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`)
)

// ValidateEVMAddress returns an error if addr is not a valid EVM address.
// A valid address must match the regex ^0x[0-9a-fA-F]{40}$.
func ValidateEVMAddress(addr string) error {
	if !evmAddressRe.MatchString(addr) {
		return fmt.Errorf("invalid EVM address %q: must match ^0x[0-9a-fA-F]{40}$", addr)
	}
	return nil
}

// ValidateSolanaAddress returns an error if addr is not a valid Solana address.
// A valid address must match the regex ^[1-9A-HJ-NP-Za-km-z]{32,44}$.
func ValidateSolanaAddress(addr string) error {
	if !solanaAddressRe.MatchString(addr) {
		return fmt.Errorf("invalid Solana address %q: must match ^[1-9A-HJ-NP-Za-km-z]{32,44}$", addr)
	}
	return nil
}

// ValidateStakeAmount returns an error if amount is nil or not a positive integer.
func ValidateStakeAmount(amount *big.Int) error {
	if amount == nil {
		return fmt.Errorf("stake amount must not be nil")
	}
	if amount.Sign() <= 0 {
		return fmt.Errorf("stake amount must be positive, got %s", amount.String())
	}
	return nil
}
