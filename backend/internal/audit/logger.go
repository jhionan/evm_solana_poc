package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ComputeHash produces a hex-encoded SHA-256 hash that chains audit entries.
// The input is the pipe-delimited concatenation:
//
//	"action|actor|chainID|details|prevHash"
//
// Pass an empty string for prevHash on the first entry in a chain.
func ComputeHash(action, actor, chainID, details, prevHash string) string {
	raw := fmt.Sprintf("%s|%s|%s|%s|%s", action, actor, chainID, details, prevHash)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
