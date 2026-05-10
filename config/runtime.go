package config

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	applog "go-skeleton/pkg/log"
)

// LoadEnv preloads optional dotenv files without overriding existing environment variables.
func LoadEnv(paths ...string) {
	files := append([]string{}, paths...)
	files = append(files, ".env")

	for _, file := range files {
		_ = godotenv.Load(file)
	}
}

// InitRuntime initializes process-wide runtime settings.
func InitRuntime(cfg *Config, service ...string) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	gin.SetMode(cfg.Server.GinMode)
	serviceName := ""
	if len(service) > 0 {
		serviceName = service[0]
	}
	if _, err := applog.Init(applog.Config{
		Level:           cfg.Log.Level,
		Format:          cfg.Log.Format,
		StacktraceLevel: cfg.Log.StacktraceLevel,
		Service:         serviceName,
	}); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	return nil
}
