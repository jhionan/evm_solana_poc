// Command server is the entry point for the multichain-staking API server.
// It wires all components together and runs until an OS signal is received.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/jhionan/multichain-staking/gen/staking/v1/stakingv1connect"
	"github.com/jhionan/multichain-staking/internal/api"
	"github.com/jhionan/multichain-staking/internal/audit"
	"github.com/jhionan/multichain-staking/internal/auth"
	"github.com/jhionan/multichain-staking/internal/chain"
	"github.com/jhionan/multichain-staking/internal/chain/evm"
	solanachain "github.com/jhionan/multichain-staking/internal/chain/solana"
	"github.com/jhionan/multichain-staking/internal/config"
	"github.com/jhionan/multichain-staking/internal/indexer"
	"github.com/jhionan/multichain-staking/internal/security"
	"github.com/jhionan/multichain-staking/internal/signer"
	"github.com/jhionan/multichain-staking/internal/staking"
	"github.com/jhionan/multichain-staking/pkg/middleware"
)

func main() {
	// -------------------------------------------------------------------------
	// 1. Logger: ConsoleWriter for local/dev, JSON for staging/production.
	// -------------------------------------------------------------------------
	var baseLogger zerolog.Logger

	appEnvHint := os.Getenv("APP_ENV")
	if appEnvHint == "" || appEnvHint == "local" {
		baseLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			Logger()
	} else {
		baseLogger = zerolog.New(os.Stderr).
			With().
			Timestamp().
			Logger()
	}

	log := baseLogger

	// -------------------------------------------------------------------------
	// 2. Config: load and validate all environment variables.
	// -------------------------------------------------------------------------
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}
	log.Info().
		Str("env", cfg.AppEnv).
		Int("port", cfg.ServerPort).
		Msg("configuration loaded")

	// -------------------------------------------------------------------------
	// 3. Root context — cancelled on OS signal for graceful shutdown.
	// -------------------------------------------------------------------------
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// -------------------------------------------------------------------------
	// 4. Database: pgxpool connection pool.
	// -------------------------------------------------------------------------
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database connection pool")
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("database ping failed")
	}
	log.Info().Msg("database connected")

	// -------------------------------------------------------------------------
	// 5. Auth service.
	// -------------------------------------------------------------------------
	jwtSvc := auth.NewJWTService(cfg.JWTSecret)
	log.Info().Msg("JWT service initialised")

	// -------------------------------------------------------------------------
	// 6. Chain adapters (optional — skipped gracefully if env vars not set).
	// -------------------------------------------------------------------------
	var adapters []chain.ChainStaker
	var evmAdapter *evm.EVMStaker // retained for indexer wiring

	// -- 6a. EVM chain --
	if cfg.EVMStakingContract != "" && cfg.EVMTokenContract != "" {
		log.Info().
			Str("rpc_url", cfg.EVMRpcURL).
			Msg("wiring EVM chain adapter")

		ethClient, err := ethclient.DialContext(ctx, cfg.EVMRpcURL)
		if err != nil {
			log.Fatal().Err(err).Str("url", cfg.EVMRpcURL).Msg("failed to connect to EVM RPC")
		}

		chainID, err := ethClient.ChainID(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to fetch EVM chain ID")
		}
		log.Info().Str("chain_id", chainID.String()).Msg("EVM chain ID resolved")

		evmSigner, err := signer.NewEVMSigner(ethClient, cfg.EVMPrivateKey, chainID)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create EVM signer")
		}
		log.Info().Str("address", evmSigner.Address()).Msg("EVM signer ready")

		stakingAddr := common.HexToAddress(cfg.EVMStakingContract)
		tokenAddr := common.HexToAddress(cfg.EVMTokenContract)

		evmStaker, err := evm.NewEVMStaker(ethClient, stakingAddr, tokenAddr, evmSigner, log)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create EVM staker")
		}

		adapters = append(adapters, evmStaker)
		evmAdapter = evmStaker
		log.Info().
			Str("staking_contract", cfg.EVMStakingContract).
			Str("token_contract", cfg.EVMTokenContract).
			Msg("EVM chain adapter registered")
	} else {
		log.Warn().
			Msg("EVM_STAKING_CONTRACT or EVM_TOKEN_CONTRACT not set — EVM chain adapter skipped")
	}

	// -- 6b. Solana chain --
	if cfg.SolanaProgramID != "" {
		log.Info().
			Str("rpc_url", cfg.SolanaRpcURL).
			Msg("wiring Solana chain adapter")

		solClient := rpc.New(cfg.SolanaRpcURL)

		programID, err := solanago.PublicKeyFromBase58(cfg.SolanaProgramID)
		if err != nil {
			log.Fatal().Err(err).Str("program_id", cfg.SolanaProgramID).Msg("invalid Solana program ID")
		}

		var authority solanago.PrivateKey
		if cfg.SolanaPrivateKey != "" {
			authority, err = solanago.PrivateKeyFromBase58(cfg.SolanaPrivateKey)
			if err != nil {
				log.Fatal().Err(err).Msg("invalid Solana private key")
			}
		}

		solanaStaker, err := solanachain.NewSolanaStaker(solClient, programID, authority, log)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create Solana staker")
		}

		adapters = append(adapters, solanaStaker)
		log.Info().
			Str("program_id", cfg.SolanaProgramID).
			Msg("Solana chain adapter registered")
	} else {
		log.Warn().Msg("SOLANA_PROGRAM_ID not set — Solana chain adapter skipped")
	}

	// -------------------------------------------------------------------------
	// 7. Staking service.
	// -------------------------------------------------------------------------
	stakingSvc := staking.NewService(adapters)
	log.Info().Int("adapters", len(adapters)).Msg("staking service initialised")

	// -------------------------------------------------------------------------
	// 8. Indexer (only if at least one EVM adapter is registered).
	// -------------------------------------------------------------------------
	if evmAdapter != nil {
		log.Info().Msg("wiring EVM event indexer")

		pgStore := indexer.NewPGStore(pool)

		// Fetch the chain ID string for the EVMEventSource.
		var evmChainIDStr string
		{
			ethClient, dialErr := ethclient.DialContext(ctx, cfg.EVMRpcURL)
			if dialErr != nil {
				log.Fatal().Err(dialErr).Msg("indexer: failed to connect to EVM RPC for chain ID")
			}
			cid, cidErr := ethClient.ChainID(ctx)
			if cidErr != nil {
				log.Fatal().Err(cidErr).Msg("indexer: failed to fetch chain ID")
			}
			evmChainIDStr = cid.String()
			// We use a fresh dial here; the source manages its own clients.
			ethClient.Close()
		}

		// HTTP client for historical FilterLogs calls.
		httpClient, err := ethclient.DialContext(ctx, cfg.EVMRpcURL)
		if err != nil {
			log.Fatal().Err(err).Msg("indexer: failed to connect HTTP EVM client")
		}

		// WebSocket client for live subscriptions (optional).
		var wsClient *ethclient.Client
		if cfg.EVMWsURL != "" {
			wsClient, err = ethclient.DialContext(ctx, cfg.EVMWsURL)
			if err != nil {
				log.Warn().Err(err).Str("ws_url", cfg.EVMWsURL).
					Msg("indexer: failed to connect WS client — live subscription disabled")
				wsClient = nil
			} else {
				log.Info().Str("ws_url", cfg.EVMWsURL).Msg("indexer: WebSocket client connected")
			}
		} else {
			log.Warn().Msg("EVM_WS_URL not set — indexer will run catch-up only, no live subscription")
		}

		contractAddr := common.HexToAddress(cfg.EVMStakingContract)

		evmSource := indexer.NewEVMEventSource(
			evmChainIDStr,
			contractAddr,
			httpClient,
			wsClient,
		)

		idx := indexer.NewIndexer(evmSource, pgStore, 0)

		go func() {
			log.Info().Str("chain_id", evmChainIDStr).Msg("indexer: starting")
			if runErr := idx.Run(ctx); runErr != nil {
				if ctx.Err() != nil {
					// Context cancelled — normal shutdown path.
					log.Info().Msg("indexer: stopped (context cancelled)")
					return
				}
				log.Error().Err(runErr).Msg("indexer: exited with error")
			}
		}()
	} else {
		log.Warn().Msg("no EVM adapter registered — indexer not started")
	}

	// -------------------------------------------------------------------------
	// 9. Rate limiter (Valkey-backed, optional).
	// -------------------------------------------------------------------------
	var interceptors []connect.Interceptor
	interceptors = append(interceptors, auth.AuthInterceptor(jwtSvc))

	if cfg.ValkeyURL != "" {
		rl, rlErr := security.NewRateLimiter(cfg.ValkeyURL, cfg.ValkeyPassword, 60, time.Minute)
		if rlErr != nil {
			log.Warn().Err(rlErr).Msg("rate limiter unavailable — skipping")
		} else {
			interceptors = append(interceptors, rl.Interceptor())
			log.Info().Str("addr", cfg.ValkeyURL).Int("limit", 60).Msg("rate limiter enabled")
		}
	}

	// -------------------------------------------------------------------------
	// 10. Audit interceptor (PostgreSQL-backed).
	// -------------------------------------------------------------------------
	auditDB := audit.NewPGAuditDB(pool)
	auditInterceptor := audit.NewAuditInterceptor(auditDB)
	interceptors = append(interceptors, auditInterceptor.Interceptor())
	log.Info().Msg("audit interceptor enabled")

	// -------------------------------------------------------------------------
	// 11. ConnectRPC handler and route registration.
	// -------------------------------------------------------------------------
	handler := api.NewHandler(stakingSvc)

	mux := http.NewServeMux()
	path, connectHandler := stakingv1connect.NewStakingServiceHandler(
		handler,
		connect.WithInterceptors(interceptors...),
	)
	mux.Handle(path, connectHandler)
	log.Info().Str("path", path).Msg("ConnectRPC handler mounted")

	// -------------------------------------------------------------------------
	// 10. HTTP server with security middleware and timeouts.
	// -------------------------------------------------------------------------
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      middleware.SecurityHeaders(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("server starting")
		if listenErr := srv.ListenAndServe(); listenErr != nil && listenErr != http.ErrServerClosed {
			log.Fatal().Err(listenErr).Msg("server error")
		}
	}()

	// -------------------------------------------------------------------------
	// 11. Block until shutdown signal, then drain gracefully.
	// -------------------------------------------------------------------------
	<-ctx.Done()
	log.Info().Msg("shutdown signal received — draining")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
		os.Exit(1)
	}

	// pool.Close() is called by defer above.
	log.Info().Msg("server stopped gracefully")

}
