package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GroVlAn/auth-base/crypto"
	"github.com/GroVlAn/auth-user/internal/config"
	grpcHandler "github.com/GroVlAn/auth-user/internal/handler/grpc-handler"
	httpHandler "github.com/GroVlAn/auth-user/internal/handler/http-handler"
	"github.com/GroVlAn/auth-user/internal/infrastructure/database"
	grpcClient "github.com/GroVlAn/auth-user/internal/infrastructure/grpc-client"
	"github.com/GroVlAn/auth-user/internal/infrastructure/secrets"
	vaultClient "github.com/GroVlAn/auth-user/internal/infrastructure/vault-client"
	"github.com/GroVlAn/auth-user/internal/repository"
	grpcServer "github.com/GroVlAn/auth-user/internal/server/grpc-server"
	httpserver "github.com/GroVlAn/auth-user/internal/server/http-server"
	"github.com/GroVlAn/auth-user/internal/service"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	localConfigPath = "configs/config-local.yml"
)

func main() {
	timeStart := time.Now()

	l := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"})

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	configPath := flag.String("config", localConfigPath, "Path to the configuration file")
	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to load configuration")
	}

	vc, err := vaultClient.New(vaultClient.Conf{
		SecretToken: cfg.Vault.SecretToken,
		Address:     cfg.Vault.Address,
		Mount:       cfg.Vault.Mount,
	})
	if err != nil {
		l.Fatal().Err(err).Msg("failed to load vault client")
	}

	provider := secrets.New(vc, secrets.Paths{
		Postgres: cfg.VaultPaths.Postgres,
		Hasher:   cfg.VaultPaths.Hasher,
	})

	scrt, err := provider.Load(ctx)
	if err != nil {
		l.Fatal().Err(err).Msg("failed load secrets")
	}

	db, err := database.NewPostgresqlDB(database.PostgresSettings{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: scrt.Postgres.Username,
		Password: scrt.Postgres.Password,
		DBName:   scrt.Postgres.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		l.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			l.Error().Err(err).Msg("failed to close postgresql db connection")
		}
	}()

	r := repository.New(db)

	hasher := crypto.New(crypto.Deps{
		Time:    scrt.Hasher.Time,
		Memory:  scrt.Hasher.Memory,
		Threads: scrt.Hasher.Threads,
		KeyLen:  scrt.Hasher.KeyLen,
		SaltLen: scrt.Hasher.SaltLen,
	})

	conn, err := grpc.NewClient(
		cfg.GRPC.AccessApiHost+":"+cfg.GRPC.AccessApiPort,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to grpc user service client")
	}
	defer func() {
		if err := conn.Close(); err != nil {
			l.Error().Err(err).Msg("failed close grpc connection")
		}
	}()

	grpcClient := grpcClient.New(conn)

	s := service.New(r, hasher, grpcClient)

	h := httpHandler.New(
		l,
		s,
		httpHandler.Deps{
			BasePath:       cfg.HTTP.BaseHTTPPath,
			DefaultTimeout: cfg.Settings.DefaultTimeout,
		},
	)

	gh := grpcHandler.New(
		l,
		s,
		cfg.Settings.DefaultTimeout,
	)

	hServer := httpserver.New(
		h.Handler(),
		httpserver.Settings{
			Port:              cfg.HTTP.Port,
			MaxHeaderBytes:    cfg.HTTP.MaxHeaderBytes,
			ReadHeaderTimeout: time.Duration(cfg.HTTP.ReadHeaderTimeout) * time.Second,
			WriteTimeout:      time.Duration(cfg.HTTP.WriteTimeout) * time.Second,
		},
	)

	gServer := grpcServer.New(
		gh,
	)

	errCh := make(chan error, 2)

	go func() {
		l.Info().Msgf("starting http server on port: %s", cfg.HTTP.Port)

		err := hServer.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	go func() {
		l.Info().Msgf("starting grpc server on port: %s", cfg.GRPC.Port)

		if err := gServer.ListenAndServe(cfg.GRPC.Port); err != nil {
			errCh <- err
		}

	}()

	l.Info().
		Dur("startup_time", time.Since(timeStart)).
		Str("http_port", cfg.HTTP.Port).
		Str("grpc_port", cfg.GRPC.Port).
		Msg("server started")

	shutdown := func() {
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := hServer.Shutdown(shutdownCtx); err != nil {
			l.Error().Err(err).Msg("failed to shutdown server")
		} else {
			l.Info().Msg("server shutdown gracefully")
		}
		gServer.Stop()
	}

	select {
	case <-ctx.Done():
		shutdown()
	case err := <-errCh:
		if err != nil {
			l.Error().Err(err).Msg("server exited with error")

			shutdown()
		}
	}

}
