package signer

import "context"

// TxSigner abstracts transaction signing across different chain implementations.
// Each concrete adapter (EVM, Solana, …) must implement this interface.
type TxSigner interface {
	// Address returns the public address controlled by this signer.
	Address() string

	// SignAndSend signs the raw transaction payload and broadcasts it.
	// Returns the transaction hash on success.
	SignAndSend(ctx context.Context, payload []byte) (txHash string, err error)

	// Nonce returns the current nonce/sequence number for the signer's address.
	Nonce(ctx context.Context) (uint64, error)
}
