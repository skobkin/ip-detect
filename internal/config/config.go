// Package config defines and loads application configuration values.
package config

import (
	"fmt"
	"log/slog"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 10 * time.Second
	defaultLookupTimeout   = 500 * time.Millisecond
)

// Config aggregates all configuration sections.
type Config struct {
	Server   ServerConfig
	Proxy    ProxyConfig
	Resolver ResolverConfig
	Metadata MetadataConfig
	Logging  LoggingConfig
}

// ServerConfig controls HTTP server behavior.
type ServerConfig struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// ProxyConfig governs how proxy headers are trusted.
type ProxyConfig struct {
	TrustForwarded bool
	TrustedSubnets []netip.Prefix
}

// ResolverConfig tunes reverse-DNS lookups.
type ResolverConfig struct {
	EnableReverseDNS bool
	LookupTimeout    time.Duration
}

// MetadataConfig toggles extra response fields.
type MetadataConfig struct {
	IncludeUserAgent         bool
	IncludeTimestamp         bool
	IncludeConnection        bool
	IncludeTLS               bool
	IncludeClientPreferences bool
	IncludeOriginContext     bool
	IncludeClientHints       bool
	IncludeProxyDetails      bool
	IncludeRequestHeaders    bool
}

// LoggingConfig customizes application logging.
type LoggingConfig struct {
	Level  slog.Level
	Format string
}

// Default returns the default configuration.
func Default() Config {
	return Config{
		Server: ServerConfig{
			Addr:            ":8080",
			ReadTimeout:     defaultReadTimeout,
			WriteTimeout:    defaultWriteTimeout,
			ShutdownTimeout: defaultShutdownTimeout,
		},
		Proxy: ProxyConfig{
			TrustForwarded: true,
			TrustedSubnets: nil,
		},
		Resolver: ResolverConfig{
			EnableReverseDNS: true,
			LookupTimeout:    defaultLookupTimeout,
		},
		Metadata: MetadataConfig{
			IncludeUserAgent:         true,
			IncludeTimestamp:         true,
			IncludeConnection:        true,
			IncludeTLS:               true,
			IncludeClientPreferences: true,
			IncludeOriginContext:     true,
			IncludeClientHints:       true,
			IncludeProxyDetails:      false,
			IncludeRequestHeaders:    false,
		},
		Logging: LoggingConfig{
			Level:  slog.LevelInfo,
			Format: "text",
		},
	}
}

// Load reads configuration from environment variables.
func Load() (Config, error) {
	cfg := Default()

	if v := strings.TrimSpace(os.Getenv("IPD_ADDR")); v != "" {
		cfg.Server.Addr = v
	}

	if v := os.Getenv("IPD_READ_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_READ_TIMEOUT: %w", err)
		}

		cfg.Server.ReadTimeout = d
	}

	if v := os.Getenv("IPD_WRITE_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_WRITE_TIMEOUT: %w", err)
		}

		cfg.Server.WriteTimeout = d
	}

	if v := os.Getenv("IPD_SHUTDOWN_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_SHUTDOWN_TIMEOUT: %w", err)
		}

		cfg.Server.ShutdownTimeout = d
	}

	if v := os.Getenv("IPD_TRUST_FORWARDED"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_TRUST_FORWARDED: %w", err)
		}

		cfg.Proxy.TrustForwarded = b
	}

	if v := os.Getenv("IPD_TRUSTED_SUBNETS"); v != "" {
		entries := strings.Split(v, ",")
		cfg.Proxy.TrustedSubnets = nil

		for _, raw := range entries {
			trimmed := strings.TrimSpace(raw)
			if trimmed == "" {
				continue
			}

			prefix, err := netip.ParsePrefix(trimmed)
			if err != nil {
				return Config{}, fmt.Errorf("invalid CIDR in IPD_TRUSTED_SUBNETS: %w", err)
			}

			cfg.Proxy.TrustedSubnets = append(cfg.Proxy.TrustedSubnets, prefix)
		}
	}

	if v := os.Getenv("IPD_RESOLVE_PTR"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_RESOLVE_PTR: %w", err)
		}

		cfg.Resolver.EnableReverseDNS = b
	}

	if v := os.Getenv("IPD_RESOLVE_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_RESOLVE_TIMEOUT: %w", err)
		}

		cfg.Resolver.LookupTimeout = d
	}

	if v := os.Getenv("IPD_INCLUDE_UA"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_INCLUDE_UA: %w", err)
		}

		cfg.Metadata.IncludeUserAgent = b
	}

	if v := os.Getenv("IPD_INCLUDE_TS"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_INCLUDE_TS: %w", err)
		}

		cfg.Metadata.IncludeTimestamp = b
	}

	if v := os.Getenv("IPD_INCLUDE_CONNECTION"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_INCLUDE_CONNECTION: %w", err)
		}

		cfg.Metadata.IncludeConnection = b
	}

	if v := os.Getenv("IPD_INCLUDE_TLS"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_INCLUDE_TLS: %w", err)
		}

		cfg.Metadata.IncludeTLS = b
	}

	if v := os.Getenv("IPD_INCLUDE_PREFERENCES"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_INCLUDE_PREFERENCES: %w", err)
		}

		cfg.Metadata.IncludeClientPreferences = b
	}

	if v := os.Getenv("IPD_INCLUDE_ORIGIN"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_INCLUDE_ORIGIN: %w", err)
		}

		cfg.Metadata.IncludeOriginContext = b
	}

	if v := os.Getenv("IPD_INCLUDE_CLIENT_HINTS"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_INCLUDE_CLIENT_HINTS: %w", err)
		}

		cfg.Metadata.IncludeClientHints = b
	}

	if v := os.Getenv("IPD_INCLUDE_PROXY"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_INCLUDE_PROXY: %w", err)
		}

		cfg.Metadata.IncludeProxyDetails = b
	}

	if v := os.Getenv("IPD_INCLUDE_HEADERS"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid IPD_INCLUDE_HEADERS: %w", err)
		}

		cfg.Metadata.IncludeRequestHeaders = b
	}

	if v := strings.TrimSpace(os.Getenv("IPD_LOG_FORMAT")); v != "" {
		v = strings.ToLower(v)
		switch v {
		case "text", "json":
			cfg.Logging.Format = v
		default:
			return Config{}, fmt.Errorf("invalid IPD_LOG_FORMAT: %s", v)
		}
	}

	if v := strings.TrimSpace(os.Getenv("IPD_LOG_LEVEL")); v != "" {
		level, err := parseLogLevel(v)
		if err != nil {
			return Config{}, err
		}

		cfg.Logging.Level = level
	}

	return cfg, nil
}

func parseLogLevel(value string) (slog.Level, error) {
	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("invalid IPD_LOG_LEVEL: %s", value)
	}
}
