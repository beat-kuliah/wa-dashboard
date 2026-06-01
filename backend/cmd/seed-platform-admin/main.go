package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	authmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/auth"
	platformadminmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/platformadmin"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/config"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/db"
	"github.com/jackc/pgx/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	email := strings.TrimSpace(cfg.PlatformAdminSeedEmail)
	password := cfg.PlatformAdminSeedPassword
	fullName := strings.TrimSpace(cfg.PlatformAdminSeedFullName)
	if email == "" || password == "" {
		fmt.Fprintln(os.Stderr, "PLATFORM_ADMIN_EMAIL and PLATFORM_ADMIN_PASSWORD are required")
		os.Exit(1)
	}
	if fullName == "" {
		fullName = "Platform Admin"
	}

	pool, err := db.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "database error: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	ctx := context.Background()
	repo := platformadminmod.NewRepository(pool.Pool)

	if _, err := repo.GetPlatformAdminByEmail(ctx, email); err == nil {
		fmt.Printf("platform admin already exists: %s\n", email)
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		fmt.Fprintf(os.Stderr, "lookup error: %v\n", err)
		os.Exit(1)
	}

	hash, err := authmod.HashPassword(password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hash error: %v\n", err)
		os.Exit(1)
	}

	admin, err := repo.CreatePlatformAdmin(ctx, email, hash, fullName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("created platform admin %s (%s)\n", admin.Email, admin.ID)
}
