package config

import (
	"context"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	GlobalConfig AppConfig
	DB           *gorm.DB
	Rdb          *redis.Client
	Ctx          = context.Background()
)
