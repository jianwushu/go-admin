package global

import (
	"go-admin/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Redis  *redis.Client
	Log    *zap.Logger
	Config config.Config
)
