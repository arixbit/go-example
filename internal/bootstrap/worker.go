package bootstrap

import (
	"errors"
	"fmt"
	"strings"

	"go-skeleton/config"
	"go-skeleton/internal/taskqueue"
	"go-skeleton/pkg/cache"
	"go-skeleton/pkg/database"
)

// InitWorker initializes resources required by the async worker process.
func InitWorker(cfg *config.Config) (*Registry, error) {
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
		dbMgr, err = initDatabase(cfg)
		if err != nil {
			return nil, fmt.Errorf("init worker database: %w", err)
		}
	}

	queueClient := newAsynqClient(cfg)
	return &Registry{
		Cfg:         cfg,
		DB:          dbMgr,
		Cache:       cacheClient,
		Queue:       taskqueue.NewQueue(queueClient),
		queueClient: queueClient,
	}, nil
}
