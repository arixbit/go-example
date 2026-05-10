package database

import (
	"gorm.io/gorm/logger"
)

func newGormLogger(level string) logger.Interface {
	return logger.Default.LogMode(parseLogLevel(level))
}
