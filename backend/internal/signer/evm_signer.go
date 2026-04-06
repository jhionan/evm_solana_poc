package signer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EVMSigner implements TxSigner for EVM-compatible chains.
// It holds the private key in memory and serialises nonce access with a mutex
// to prevent accidental nonce reuse under concurrent calls.
type EVMSigner struct {
	client  *ethclient.Client
	privKey *ecdsa.PrivateKey
	address common.Address
	chainID *big.Int
	mu      sync.Mutex
}

// NewEVMSigner creates an EVMSigner from a hex-encoded private key string.
// The "0x" prefix is stripped automatically if present.
func NewEVMSigner(client *ethclient.Client, privateKeyHex string, chainID *big.Int) (*EVMSigner, error) {
	if client == nil {
		return nil, fmt.Errorf("signer: ethclient must not be nil")
	}
	if chainID == nil {
		return nil, fmt.Errorf("signer: chainID must not be nil")
	}

	hex := strings.TrimPrefix(privateKeyHex, "0x")
	privKey, err := crypto.HexToECDSA(hex)
	if err != nil {
		return nil, fmt.Errorf("signer: invalid private key: %w", err)
	}

	pubKey, ok := privKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("signer: could not cast public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*pubKey)

	return &EVMSigner{
		client:  client,
		privKey: privKey,
		address: address,
		chainID: chainID,
	}, nil
}

// Address returns the checksummed hex address derived from the private key.
func (s *EVMSigner) Address() string {
	return s.address.Hex()
}

// Nonce fetches the pending nonce for the signer's address from the node.
// The mutex ensures only one goroutine reads+uses the nonce at a time.
func (s *EVMSigner) Nonce(ctx context.Context) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nonce, err := s.client.PendingNonceAt(ctx, s.address)
	if err != nil {
		return 0, fmt.Errorf("signer: failed to fetch nonce: %w", err)
	}
	return nonce, nil
}

// SignAndSend signs the raw RLP-encoded transaction bytes and broadcasts them.
// payload must be a fully-populated but unsigned *types.Transaction encoded
// with rlp.EncodeToBytes. Returns the transaction hash on success.
func (s *EVMSigner) SignAndSend(ctx context.Context, payload []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(payload); err != nil {
		return "", fmt.Errorf("signer: failed to decode transaction payload: %w", err)
	}

	signer := types.LatestSignerForChainID(s.chainID)
	signed, err := types.SignTx(tx, signer, s.privKey)
	if err != nil {
		return "", fmt.Errorf("signer: failed to sign transaction: %w", err)
	}

	if err := s.client.SendTransaction(ctx, signed); err != nil {
		return "", fmt.Errorf("signer: failed to broadcast transaction: %w", err)
	}

	return signed.Hash().Hex(), nil
}

// SignTx signs a pre-built *types.Transaction and returns the signed copy.
// This is a lower-level helper for adapters that construct transactions
// themselves before broadcasting separately. Protected by mutex.
func (s *EVMSigner) SignTx(tx *types.Transaction) (*types.Transaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	signer := types.LatestSignerForChainID(s.chainID)
	signed, err := types.SignTx(tx, signer, s.privKey)
	if err != nil {
		return nil, fmt.Errorf("signer: failed to sign transaction: %w", err)
	}
	return signed, nil
}

// compile-time interface assertion
var _ TxSigner = (*EVMSigner)(nil)
