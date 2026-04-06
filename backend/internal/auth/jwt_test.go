package auth_test

import (
	"testing"
	"time"

	"github.com/jhionan/multichain-staking/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTService_SignVerify_Roundtrip(t *testing.T) {
	svc := auth.NewJWTService("super-secret-key")

	original := auth.Claims{
		Wallet: "0xDeAdBeEf",
		Role:   auth.RoleUser,
	}

	token, err := svc.Sign(original, time.Hour)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	got, err := svc.Verify(token)
	require.NoError(t, err)

	assert.Equal(t, original.Wallet, got.Wallet)
	assert.Equal(t, original.Role, got.Role)
	assert.NotNil(t, got.IssuedAt)
	assert.NotNil(t, got.ExpiresAt)
	assert.True(t, got.ExpiresAt.After(got.IssuedAt.Time))
}

func TestJWTService_Verify_ExpiredToken_Fails(t *testing.T) {
	svc := auth.NewJWTService("super-secret-key")

	claims := auth.Claims{
		Wallet: "0xExpired",
		Role:   auth.RoleUser,
	}

	// Sign with a negative duration so the token is immediately expired.
	token, err := svc.Sign(claims, -time.Second)
	require.NoError(t, err)

	_, err = svc.Verify(token)
	assert.Error(t, err, "expected error for expired token")
}

func TestJWTService_Verify_WrongSecret_Fails(t *testing.T) {
	signer := auth.NewJWTService("correct-secret")
	verifier := auth.NewJWTService("wrong-secret")

	claims := auth.Claims{
		Wallet: "0xSomeWallet",
		Role:   auth.RoleAdmin,
	}

	token, err := signer.Sign(claims, time.Hour)
	require.NoError(t, err)

	_, err = verifier.Verify(token)
	assert.Error(t, err, "expected error for wrong secret")
}

func TestJWTService_AdminRole_Roundtrip(t *testing.T) {
	svc := auth.NewJWTService("admin-secret")

	claims := auth.Claims{
		Wallet: "0xAdminWallet",
		Role:   auth.RoleAdmin,
	}

	token, err := svc.Sign(claims, time.Minute*15)
	require.NoError(t, err)

	got, err := svc.Verify(token)
	require.NoError(t, err)
	assert.Equal(t, auth.RoleAdmin, got.Role)
}
