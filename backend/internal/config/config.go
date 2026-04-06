// Package config provides Viper-based configuration loading and validation.
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration values.
type Config struct {
	AppEnv     string `mapstructure:"app_env"`
	ServerPort int    `mapstructure:"server_port"`

	DatabaseURL    string `mapstructure:"database_url"`
	ValkeyURL      string `mapstructure:"valkey_url"`
	ValkeyPassword string `mapstructure:"valkey_password"`

	JWTSecret string `mapstructure:"jwt_secret"`

	EVMRpcURL          string `mapstructure:"evm_rpc_url"`
	EVMWsURL           string `mapstructure:"evm_ws_url"`
	EVMPrivateKey      string `mapstructure:"evm_private_key"`
	EVMStakingContract string `mapstructure:"evm_staking_contract"`
	EVMTokenContract   string `mapstructure:"evm_token_contract"`
	EVMDeployBlock     int64  `mapstructure:"evm_deploy_block"`

	SolanaRpcURL     string `mapstructure:"solana_rpc_url"`
	SolanaWsURL      string `mapstructure:"solana_ws_url"`
	SolanaPrivateKey string `mapstructure:"solana_private_key"`
	SolanaProgramID  string `mapstructure:"solana_program_id"`
}

// allKeys maps viper keys (lowercase) to their environment variable names.
var allKeys = []struct {
	key    string
	envVar string
}{
	{"app_env", "APP_ENV"},
	{"server_port", "SERVER_PORT"},
	{"database_url", "DATABASE_URL"},
	{"valkey_url", "VALKEY_URL"},
	{"valkey_password", "VALKEY_PASSWORD"},
	{"jwt_secret", "JWT_SECRET"},
	{"evm_rpc_url", "EVM_RPC_URL"},
	{"evm_ws_url", "EVM_WS_URL"},
	{"evm_private_key", "EVM_PRIVATE_KEY"},
	{"evm_staking_contract", "EVM_STAKING_CONTRACT"},
	{"evm_token_contract", "EVM_TOKEN_CONTRACT"},
	{"evm_deploy_block", "EVM_DEPLOY_BLOCK"},
	{"solana_rpc_url", "SOLANA_RPC_URL"},
	{"solana_ws_url", "SOLANA_WS_URL"},
	{"solana_private_key", "SOLANA_PRIVATE_KEY"},
	{"solana_program_id", "SOLANA_PROGRAM_ID"},
}

// Load reads configuration from environment variables, applies defaults,
// unmarshals into Config, and validates required fields.
func Load() (*Config, error) {
	v := viper.New()

	// Defaults.
	v.SetDefault("app_env", "local")
	v.SetDefault("server_port", 8080)
	v.SetDefault("valkey_url", "localhost:6379")
	v.SetDefault("valkey_password", "")

	// Explicitly bind each env var so AutomaticEnv resolves correctly.
	for _, k := range allKeys {
		if err := v.BindEnv(k.key, k.envVar); err != nil {
			return nil, fmt.Errorf("config: binding env var %s: %w", k.envVar, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal failed: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

var validEnvs = map[string]bool{
	"local":      true,
	"staging":    true,
	"production": true,
}

// validate checks all required fields and business rules.
func validate(cfg *Config) error {
	var errs []string

	if cfg.DatabaseURL == "" {
		errs = append(errs, "DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		errs = append(errs, "JWT_SECRET is required")
	} else if len(cfg.JWTSecret) < 32 {
		errs = append(errs, "JWT_SECRET must be at least 32 characters")
	}

	if cfg.EVMRpcURL == "" {
		errs = append(errs, "EVM_RPC_URL is required")
	}

	if cfg.EVMPrivateKey == "" {
		errs = append(errs, "EVM_PRIVATE_KEY is required")
	}

	if !validEnvs[cfg.AppEnv] {
		errs = append(errs, fmt.Sprintf("APP_ENV must be one of local, staging, production (got %q)", cfg.AppEnv))
	}

	if len(errs) > 0 {
		return fmt.Errorf("config validation failed: %s", strings.Join(errs, "; "))
	}

	return nil
}
