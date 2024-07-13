package config

import (
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
)

type Config struct {
    DB  *gorm.DB
    RDB *redis.Client
}
