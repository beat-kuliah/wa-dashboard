package app

import (
	"context"
	"log/slog"

	analyticsmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/analytics"
	authmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/auth"
	broadcastmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/broadcast"
	inboxmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/inbox"
	templatemod "github.com/beatfraps/wa-dashboard/backend/internal/modules/template"
	tenantmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/tenant"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/config"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/db"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/httpx"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/obs"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/queue"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Config   config.Config
	Logger   *slog.Logger
	Pool     *db.Pool
	Server   *httpx.Server
	Queue    *queue.Client
	Auth     *authmod.Module
	Tenant   *tenantmod.Module
	Broadcast *broadcastmod.Module
	Inbox    *inboxmod.Module
	Analytics *analyticsmod.Module
	Template *templatemod.Module
}

func New(cfg config.Config, logger *slog.Logger) (*Application, error) {
	pool, err := db.NewPool(contextBackground(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	tokens := auth.NewTokenService(cfg.JWTSecret, cfg.JWTAccessTTL)
	qClient, err := queue.NewClient(cfg.RedisURL)
	if err != nil {
		pool.Close()
		return nil, err
	}

	srv := httpx.NewServer(httpx.ServerConfig{
		CORSOrigins: cfg.CORSOriginList(),
		Logger:      logger,
		Tokens:      tokens,
	})

	health := obs.NewHealth(pool, cfg.RedisURL)
	health.Register(srv.Echo)

	authModule := authmod.NewModule(pool.Pool, tokens, cfg.RefreshTokenTTL, cfg.AccessTokenExpiresIn())
	tenantModule := tenantmod.NewModule(pool.Pool)
	broadcastModule := broadcastmod.NewModule(pool.Pool, qClient)
	inboxModule := inboxmod.NewModule(pool.Pool)
	analyticsModule := analyticsmod.NewModule()
	templateModule := templatemod.NewModule(pool.Pool)

	v1 := srv.Echo.Group("/api/v1")
	authModule.RegisterRoutes(v1.Group("/auth"))
	tenantModule.RegisterRoutes(v1.Group("/tenant"))
	broadcastModule.RegisterRoutes(v1.Group("/broadcasts"))
	inboxModule.RegisterRoutes(v1.Group("/inbox"))
	analyticsModule.RegisterRoutes(v1.Group("/analytics"))
	templateModule.RegisterRoutes(v1.Group("/templates"))

	return &Application{
		Config:    cfg,
		Logger:    logger,
		Pool:      pool,
		Server:    srv,
		Queue:     qClient,
		Auth:      authModule,
		Tenant:    tenantModule,
		Broadcast: broadcastModule,
		Inbox:     inboxModule,
		Analytics: analyticsModule,
		Template:  templateModule,
	}, nil
}

func contextBackground() context.Context {
	return context.Background()
}

func (a *Application) Close() {
	if a.Queue != nil {
		a.Queue.Close()
	}
	if a.Pool != nil {
		a.Pool.Close()
	}
}

func (a *Application) PgxPool() *pgxpool.Pool {
	return a.Pool.Pool
}
