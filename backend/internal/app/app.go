package app

import (
	"context"
	"errors"
	"log/slog"

	"github.com/beatfraps/wa-dashboard/backend/db/sqlc"
	analyticsmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/analytics"
	authmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/auth"
	broadcastmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/broadcast"
	inboxmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/inbox"
	platformadminmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/platformadmin"
	templatemod "github.com/beatfraps/wa-dashboard/backend/internal/modules/template"
	tenantmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/tenant"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/config"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/db"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/httpx"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/obs"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/queue"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	Template      *templatemod.Module
	PlatformAdmin *platformadminmod.Module
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

	queries := sqlc.NewFromPool(pool.Pool)
	srv.Echo.Use(auth.TenantActiveGuard(func(ctx context.Context, tenantID uuid.UUID) (string, error) {
		tenant, err := queries.GetTenantByID(ctx, tenantID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", apperrors.Unauthorized("")
			}
			return "", err
		}
		return tenant.Status, nil
	}))

	authModule := authmod.NewModule(pool.Pool, tokens, cfg.RefreshTokenTTL, cfg.AccessTokenExpiresIn(), cfg.PublicRegistrationEnabled)
	platformAdminModule := platformadminmod.NewModule(pool.Pool, tokens, cfg.RefreshTokenTTL, cfg.AccessTokenExpiresIn())
	tenantModule := tenantmod.NewModule(pool.Pool)
	broadcastModule := broadcastmod.NewModule(pool.Pool, qClient)
	inboxModule := inboxmod.NewModule(pool.Pool)
	analyticsModule := analyticsmod.NewModule()
	templateModule := templatemod.NewModule(pool.Pool)

	if cfg.PlatformAdminSeedEmail != "" && cfg.PlatformAdminSeedPassword != "" {
		if err := platformAdminModule.Service.SeedPlatformAdmin(
			contextBackground(),
			cfg.PlatformAdminSeedEmail,
			cfg.PlatformAdminSeedPassword,
			cfg.PlatformAdminSeedFullName,
		); err != nil {
			pool.Close()
			qClient.Close()
			return nil, err
		}
	}

	v1 := srv.Echo.Group("/api/v1")
	authModule.RegisterRoutes(v1.Group("/auth"))
	platformAdminModule.RegisterRoutes(v1.Group("/admin"))
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
		Template:      templateModule,
		PlatformAdmin: platformAdminModule,
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
