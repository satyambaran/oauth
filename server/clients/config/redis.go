package config

import (
    "log"
    "os"

    "github.com/go-redis/redis/v8"
    "github.com/joho/godotenv"
)

func InitRedis() *redis.Client {
    err := godotenv.Load()
    if err != nil {
        panic("unable to load env file")
    }
    redisUrl := os.Getenv("REDIS_URL")
    rdb := redis.NewClient(&redis.Options{
        Addr:         redisUrl,
        Password:     "",
        DB:           0,
        PoolSize:     10, // maximum number of socket connections
        MinIdleConns: 3,  // minimum number of idle connections which is useful when establishing new connection is slow
    })
    _, err = rdb.Ping(rdb.Context()).Result()
    if err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }
    return rdb
}
