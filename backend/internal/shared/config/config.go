package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL      string        `env:"DATABASE_URL,required"`
	RedisURL         string        `env:"REDIS_URL,required"`
	JWTSecret        string        `env:"JWT_SECRET,required"`
	Port             int           `env:"PORT" envDefault:"8080"`
	LogLevel         string        `env:"LOG_LEVEL" envDefault:"info"`
	CORSOrigins      string        `env:"CORS_ORIGINS" envDefault:"http://localhost:3000"`
	JWTAccessTTL     time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`
	RefreshTokenTTL  time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"720h"`
	MigrationsPath   string        `env:"MIGRATIONS_PATH" envDefault:"db/migrations"`
	Environment      string        `env:"ENVIRONMENT" envDefault:"development"`
}

func (c Config) CORSOriginList() []string {
	parts := strings.Split(c.CORSOrigins, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func (c Config) AccessTokenExpiresIn() int {
	return int(c.JWTAccessTTL.Seconds())
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}
	return cfg, nil
}
