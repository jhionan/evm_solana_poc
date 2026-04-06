package config_test

import (
	"os"
	"testing"

	"github.com/jhionan/multichain-staking/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helpers to set/unset env vars and defer cleanup
func setEnv(t *testing.T, pairs ...string) {
	t.Helper()
	for i := 0; i < len(pairs); i += 2 {
		key, val := pairs[i], pairs[i+1]
		prev, exists := os.LookupEnv(key)
		if exists {
			t.Cleanup(func() { os.Setenv(key, prev) })
		} else {
			t.Cleanup(func() { os.Unsetenv(key) })
		}
		require.NoError(t, os.Setenv(key, val))
	}
}

func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		prev, exists := os.LookupEnv(key)
		if exists {
			t.Cleanup(func() { os.Setenv(key, prev) })
		}
		require.NoError(t, os.Unsetenv(key))
	}
}

// requiredVars returns the minimal set of required env vars.
func setRequiredEnv(t *testing.T) {
	t.Helper()
	setEnv(t,
		"DATABASE_URL", "postgres://user:pass@localhost:5432/db",
		"JWT_SECRET", "this-is-a-32-char-secret-minimum!",
		"EVM_RPC_URL", "https://mainnet.infura.io/v3/key",
		"EVM_PRIVATE_KEY", "0xdeadbeefdeadbeef",
	)
}

func TestLoad_Defaults(t *testing.T) {
	setRequiredEnv(t)
	// Make sure optional vars are not set so defaults kick in.
	unsetEnv(t, "APP_ENV", "SERVER_PORT", "VALKEY_URL", "VALKEY_PASSWORD")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "local", cfg.AppEnv)
	assert.Equal(t, 8080, cfg.ServerPort)
	assert.Equal(t, "localhost:6379", cfg.ValkeyURL)
	assert.Equal(t, "", cfg.ValkeyPassword)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
	assert.Equal(t, "this-is-a-32-char-secret-minimum!", cfg.JWTSecret)
	assert.Equal(t, "https://mainnet.infura.io/v3/key", cfg.EVMRpcURL)
	assert.Equal(t, "0xdeadbeefdeadbeef", cfg.EVMPrivateKey)
}

func TestLoad_MissingRequired(t *testing.T) {
	// Unset all relevant env vars so validation catches missing DATABASE_URL.
	unsetEnv(t,
		"DATABASE_URL", "JWT_SECRET", "EVM_RPC_URL", "EVM_PRIVATE_KEY",
		"APP_ENV", "SERVER_PORT", "VALKEY_URL",
	)

	_, err := config.Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DATABASE_URL")
}

func TestLoad_JWTSecretTooShort(t *testing.T) {
	setRequiredEnv(t)
	// Override JWT_SECRET with a short value.
	setEnv(t, "JWT_SECRET", "tooshort")

	_, err := config.Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET")
}

func TestLoad_InvalidAppEnv(t *testing.T) {
	setRequiredEnv(t)
	setEnv(t, "APP_ENV", "unknown")

	_, err := config.Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "APP_ENV")
}

func TestLoad_ValidAppEnvStaging(t *testing.T) {
	setRequiredEnv(t)
	setEnv(t, "APP_ENV", "staging")

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "staging", cfg.AppEnv)
}
