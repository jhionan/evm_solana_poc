// Command server is the entry point for the multichain-staking API server.
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
	"github.com/rs/zerolog"

	"github.com/jhionan/multichain-staking/gen/staking/v1/stakingv1connect"
	"github.com/jhionan/multichain-staking/internal/api"
	"github.com/jhionan/multichain-staking/internal/auth"
	"github.com/jhionan/multichain-staking/internal/chain"
	"github.com/jhionan/multichain-staking/internal/config"
	"github.com/jhionan/multichain-staking/internal/staking"
	"github.com/jhionan/multichain-staking/pkg/middleware"
)

func main() {
	// 1. Logger setup — human-readable console output.
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Logger()

	// 2. Load configuration from environment.
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}
	log.Info().Str("env", cfg.AppEnv).Int("port", cfg.ServerPort).Msg("configuration loaded")

	// 3. Auth service.
	jwtSvc := auth.NewJWTService(cfg.JWTSecret)

	// 4. Chain adapters — empty until contracts are deployed.
	adapters := []chain.ChainStaker{}

	// 5. Staking service.
	stakingSvc := staking.NewService(adapters)

	// 6. ConnectRPC handler.
	handler := api.NewHandler(stakingSvc)

	// 7. Mount handler with AuthInterceptor.
	mux := http.NewServeMux()
	path, connectHandler := stakingv1connect.NewStakingServiceHandler(
		handler,
		connect.WithInterceptors(auth.AuthInterceptor(jwtSvc)),
	)
	mux.Handle(path, connectHandler)

	// 8. HTTP server with security middleware and timeouts.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      middleware.SecurityHeaders(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 9. Graceful shutdown via OS signal.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 10. Start server in background.
	go func() {
		log.Info().Str("addr", srv.Addr).Msg("server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	// Wait for shutdown signal.
	<-ctx.Done()
	log.Info().Msg("shutdown signal received")

	// 11. Graceful shutdown with 10-second deadline.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
		os.Exit(1)
	}
	log.Info().Msg("server stopped gracefully")
}
