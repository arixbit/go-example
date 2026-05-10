package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hibiken/asynq"

	"go-skeleton/internal/taskqueue"
	"go-skeleton/pkg/cache"
	"go-skeleton/pkg/database"
	"go-skeleton/pkg/validator"
)

// Registry holds shared dependencies initialized at process startup.
type Registry struct {
	Cfg         *Config
	DB          *database.DBManager
	Cache       *cache.Client
	QueueClient *asynq.Client
	Queue       *taskqueue.Queue
}

// InitAPI initializes dependencies required by the HTTP API process.
func InitAPI(cfg *Config) (*Registry, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}

	dbMgr, err := database.Init(database.Config{
		DSN:             cfg.Postgres.DSN,
		LogLevel:        cfg.Postgres.LogLevel,
		MaxIdleConns:    cfg.Postgres.MaxIdleConns,
		MaxOpenConns:    cfg.Postgres.MaxOpenConns,
		ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Postgres.ConnMaxIdleTime,
	})
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}
	if dbMgr.DB() == nil {
		return nil, errors.New("postgres dsn is required for api")
	}

	cacheClient, err := initCache(cfg)
	if err != nil {
		return nil, fmt.Errorf("init cache: %w", err)
	}

	validator.InitValidator()
	queueClient := newAsynqClient(cfg)

	return &Registry{
		Cfg:         cfg,
		DB:          dbMgr,
		Cache:       cacheClient,
		QueueClient: queueClient,
		Queue:       taskqueue.NewQueue(queueClient),
	}, nil
}

// InitWorker initializes dependencies required by the async worker process.
func InitWorker(cfg *Config) (*Registry, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}
	if strings.TrimSpace(cfg.Redis.Addr) == "" {
		return nil, errors.New("redis address is required for worker")
	}

	cacheClient, err := cache.NewClient(cache.RedisConfig{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.CacheDB,
	})
	if err != nil {
		return nil, fmt.Errorf("init worker cache: %w", err)
	}

	var dbMgr *database.DBManager
	if strings.TrimSpace(cfg.Postgres.DSN) != "" {
		dbMgr, err = database.Init(database.Config{
			DSN:             cfg.Postgres.DSN,
			LogLevel:        cfg.Postgres.LogLevel,
			MaxIdleConns:    cfg.Postgres.MaxIdleConns,
			MaxOpenConns:    cfg.Postgres.MaxOpenConns,
			ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime,
			ConnMaxIdleTime: cfg.Postgres.ConnMaxIdleTime,
		})
		if err != nil {
			return nil, fmt.Errorf("init worker database: %w", err)
		}
	}

	queueClient := newAsynqClient(cfg)
	return &Registry{
		Cfg:         cfg,
		DB:          dbMgr,
		Cache:       cacheClient,
		QueueClient: queueClient,
		Queue:       taskqueue.NewQueue(queueClient),
	}, nil
}

// Close releases resources owned by the registry.
func (r *Registry) Close() error {
	if r == nil {
		return nil
	}

	var errs []error
	if r.QueueClient != nil {
		if err := r.QueueClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close queue client: %w", err))
		}
	}
	if r.Cache != nil {
		if err := r.Cache.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close cache: %w", err))
		}
	}
	if r.DB != nil {
		if err := r.DB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close database: %w", err))
		}
	}
	return errors.Join(errs...)
}

// Load reads all configuration from environment variables. Call after LoadEnv.
func Load() *Config {
	port := getEnvOrDefault("SERVER_PORT", ":3000")
	ginMode := getEnvOrDefault("GIN_MODE", "release")
	logLevel := getEnvOrDefault("LOG_LEVEL", "info")
	logFormat := getEnvOrDefault("LOG_FORMAT", "json")
	stacktraceLevel := getEnvOrDefault("LOG_STACKTRACE_LEVEL", "error")

	return &Config{
		Server: ServerConfig{
			Port:           port,
			GinMode:        ginMode,
			TrustedProxies: parseCSV(os.Getenv("TRUSTED_PROXIES")),
			RequestTimeout: durationEnv("REQUEST_TIMEOUT", 30*time.Second),
		},
		Postgres: PostgresConfig{
			DSN:             os.Getenv("POSTGRES"),
			LogLevel:        os.Getenv("GORM_LOG_LEVEL"),
			MaxIdleConns:    intEnv("DB_MAX_IDLE_CONNS", 15),
			MaxOpenConns:    intEnv("DB_MAX_OPEN_CONNS", 30),
			ConnMaxLifetime: durationEnv("DB_CONN_MAX_LIFETIME", 30*time.Minute),
			ConnMaxIdleTime: durationEnv("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			CacheDB:  intEnv("REDIS_CACHE_DB", 0),
			QueueDB:  intEnv("REDIS_QUEUE_DB", 6),
		},
		Cors: CorsConfig{
			AllowOrigins: parseCSV(os.Getenv("CORS_ALLOW_ORIGINS")),
		},
		Log: LogConfig{
			Level:           logLevel,
			Format:          logFormat,
			StacktraceLevel: stacktraceLevel,
			AuditEnabled:    boolEnv("AUDIT_LOG_ENABLED", true),
			AuditExcludes:   parseCSV(os.Getenv("AUDIT_LOG_EXCLUDE_PATHS")),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: intEnv("RATE_LIMIT_PER_MINUTE", 0),
		},
	}
}

func initCache(cfg *Config) (*cache.Client, error) {
	if strings.TrimSpace(cfg.Redis.Addr) == "" {
		return nil, nil
	}
	return cache.NewClient(cache.RedisConfig{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.CacheDB,
	})
}

func newAsynqClient(cfg *Config) *asynq.Client {
	if cfg == nil || strings.TrimSpace(cfg.Redis.Addr) == "" {
		return nil
	}
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.QueueDB,
	})
}

func getEnvOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func parseCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func intEnv(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func boolEnv(key string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return parsed
}
